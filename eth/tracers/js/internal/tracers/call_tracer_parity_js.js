// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// callTracer is a full blown transaction tracer that extracts and reports all
// the internal calls made by a transaction, along with any useful information.
{
	// callstack is the current recursive call stack of the EVM execution.
	callstack: [{}],

	parityErrorMapping: {
		"contract creation code storage out of gas": "Out of gas",
		"out of gas": "Out of gas",
		"gas uint64 overflow": "Out of gas",
		"max code size exceeded": "Out of gas",
		"invalid jump destination": "Bad jump destination",
		"execution reverted": "Reverted",
		"return data out of bounds": "Out of bounds",
		"stack limit reached 1024 (1023)": "Out of stack",
		"precompiled failed": "Built-in failed",
		"invalid input length": "Built-in failed",
	},

	parityErrorMappingStartingWith: {
		"invalid opcode:": "Bad instruction",
		"stack underflow": "Stack underflow",
	},

	paritySkipTracesForErrors: [
		"insufficient balance for transfer"
	],

	isObjectEmpty: function(obj) {
		for (var x in obj) { return false; }
		return true;
	},

	enter: function(frame) {
		var type = frame.getType()
		var to = frame.getTo()

		var call = {
			type: type,
			from: toHex(frame.getFrom()),
			to: toHex(frame.getTo()),
			input: toHex(frame.getInput()),
			gas: '0x' + bigInt(frame.getGas()).toString('16'),
		}
		if (frame.getValue() !== undefined){
			call.value='0x' + bigInt(frame.getValue()).toString(16)
		}
		this.callstack.push(call)
	},

	exit: function(frameResult) {
		var len = this.callstack.length
		if (len > 1) {
			var call = this.callstack.pop()

			// Skip any pre-compile invocations, those are just fancy opcodes
			// NOTE: let them captured on `enter` method so as we handle internal txs state correctly
			//			 and drop them here, as pop() has removed them from the stack
			if (isPrecompiled(call.to) && (call.type == "CALL" || call.type == "STATICCALL")) {
				return;
			}

			call.gasUsed = '0x' + bigInt(frameResult.getGasUsed()).toString('16')
			var error = frameResult.getError()
			if (error === undefined) {
				call.output = toHex(frameResult.getOutput())
			} else {
				call.error = error
				if (call.type === 'CREATE' || call.type === 'CREATE2') {
					delete call.to
				}
			}
			len -= 1
			if (this.callstack[len-1].calls === undefined) {
				this.callstack[len-1].calls = []
			}
			this.callstack[len-1].calls.push(call)
		}
	},

	// fault is invoked when the actual execution of an opcode fails.
	fault: function(log, db) {
		// If the topmost call already reverted, don't handle the additional fault again
		if (this.callstack[this.callstack.length - 1].error !== undefined) {
			return;
		}
		// Pop off the just failed call
		var call = this.callstack.pop();
		call.error = log.getError();

		var opError = log.getCallError();
		if (opError !== undefined) {
			if (this.paritySkipTracesForErrors.indexOf(opError) > -1) {
				return;
			}
			call.error = opError;
		}

		// Consume all available gas and clean any leftovers
		if (call.gas !== undefined) {
			call.gas = '0x' + bigInt(call.gas).toString(16);
			call.gasUsed = call.gas
		} else {
			// Retrieve gas true allowance from the inner call.
			// We need to extract if from within the call as there may be funky gas dynamics
			// with regard to requested and actually given gas (2300 stipend, 63/64 rule).
			call.gas = '0x' + bigInt(log.getGas()).toString(16);
		}

		if (call.error === "out of gas" && call.gas === undefined) {
			call.gas = "0x0";
		}
		delete call.gasIn; delete call.gasCost;
		delete call.outOff; delete call.outLen;

		// Flatten the failed call into its parent
		var left = this.callstack.length;
		if (left > 0) {
			if (this.callstack[left-1].calls === undefined) {
				this.callstack[left-1].calls = [];
			}
			this.callstack[left-1].calls.push(call);
			return;
		}
		// Last call failed too, leave it in the stack
		this.callstack.push(call);
	},

	// result is invoked when all the opcodes have been iterated over and returns
	// the final result of the tracing.
	result: function(ctx, db) {
		var result = {
			block:   ctx.block,
			type:    ctx.type,
			from:    toHex(ctx.from),
			to:      toHex(ctx.to),
			value:   '0x' + ctx.value.toString(16),
			gas:     '0x' + bigInt(ctx.gas).toString(16),
			gasUsed: '0x' + bigInt(ctx.gasUsed).toString(16),
			input:   toHex(ctx.input),
			output:  toHex(ctx.output),
			time:    ctx.time,
		};
		var extraCtx = {
			blockHash: ctx.blockHash,
			blockNumber: ctx.block,
			transactionHash: ctx.transactionHash,
			transactionPosition: ctx.transactionPosition,
		};
		// when this.descended remains true and first item in callstack is an empty object
		// drop the first item, in order to handle edge cases in the step() loop.
		// example edge case: contract init code "0x605a600053600160006001f0ff00", search in testdata
		if (this.descended && this.callstack.length > 1 && this.isObjectEmpty(this.callstack[0])) {
			this.callstack.shift();
		}
		if (this.callstack[0].calls !== undefined) {
			result.calls = this.callstack[0].calls;
		}
		if (this.callstack[0].error !== undefined) {
			result.error = this.callstack[0].error;
		} else if (ctx.error !== undefined) {
			result.error = ctx.error;
		}
		if (result.error !== undefined && (result.error !== "execution reverted" || result.output ==="0x")) {
			delete result.output;
		}
		return this.finalize(result, extraCtx);
	},

	// finalize recreates a call object using the final desired field order for json
	// serialization. This is a nicety feature to pass meaningfully ordered results
	// to users who don't interpret it, just display it.
	finalize: function(call, extraCtx, traceAddress) {
		var data;
		if (call.type == "CREATE" || call.type == "CREATE2") {
			data = this.createResult(call);

			// update after callResult so as it affects only the root type
			call.type = "CREATE";
		} else if (call.type == "SELFDESTRUCT") {
			call.type = "SUICIDE";
			data = this.suicideResult(call);
		} else {
			data = this.callResult(call);

			// update after callResult so as it affects only the root type
			if (call.type == "CALLCODE" || call.type == "DELEGATECALL" || call.type == "STATICCALL") {
				call.type = "CALL";
			}
		}

		traceAddress = traceAddress || [];
		var sorted = {
			type: call.type.toLowerCase(),
			action: data.action,
			result: data.result,
			error: call.error,
			traceAddress: traceAddress,
			subtraces: 0,
			transactionPosition: extraCtx.transactionPosition,
			transactionHash: extraCtx.transactionHash,
			blockNumber: extraCtx.blockNumber,
			blockHash: extraCtx.blockHash,
			time: call.time,
		}

		if (sorted.error !== undefined) {
			if (this.parityErrorMapping.hasOwnProperty(sorted.error)) {
				sorted.error = this.parityErrorMapping[sorted.error];
				delete sorted.result;
			} else {
				for (var searchKey in this.parityErrorMappingStartingWith) {
					if (this.parityErrorMappingStartingWith.hasOwnProperty(searchKey) && sorted.error.indexOf(searchKey) > -1) {
						sorted.error = this.parityErrorMappingStartingWith[searchKey];
						delete sorted.result;
					}
				}
			}
		}

		for (var key in sorted) {
			if (typeof sorted[key] === "object") {
				for (var nested_key in sorted[key]) {
					if (sorted[key][nested_key] === undefined) {
						delete sorted[key][nested_key];
					}
				}
			} else if (sorted[key] === undefined) {
				delete sorted[key];
			}
		}

		var calls = call.calls;
		if (calls !== undefined) {
			sorted["subtraces"] = calls.length;
		}

		var results = [sorted];

		if (calls !== undefined) {
			for (var i=0; i<calls.length; i++) {
				var childCall = calls[i];

				// Delegatecall uses the value from parent
				if ((childCall.type == "DELEGATECALL" || childCall.type == "STATICCALL") && typeof childCall.value === "undefined") {
					childCall.value = call.value;
				}

				results = results.concat(this.finalize(childCall, extraCtx, traceAddress.concat([i])));
			}
		}
		return results;
	},

	createResult: function(call) {
		return {
			action: {
				from:           call.from,                // Sender
				value:          call.value,               // Value
				gas:            call.gas,                 // Gas
				init:           call.input,               // Initialization code
				creationMethod: call.type.toLowerCase(),  // Create Type
			},
			result: {
				gasUsed:  call.gasUsed,  // Gas used
				code:     call.output,   // Code
				address:  call.to,       // Assigned address
			}
		}
	},

	callResult: function(call) {
		return {
			action: {
				from:      call.from,               // Sender
				to:        call.to,                 // Recipient
				value:     call.value,              // Transfered Value
				gas:       call.gas,                // Gas
				input:     call.input,              // Input data
				callType:  call.type.toLowerCase(), // The type of the call
			},
			result: {
				gasUsed: call.gasUsed,  // Gas used
				output:  call.output,   // Output bytes
			}
		}
	},

	suicideResult: function(call) {
		return {
			action: {
				address:        call.from,   // Address
				refundAddress:  call.to,     // Refund address
				balance:        call.value,  // Balance
			},
			result: null
		}
	}
}
