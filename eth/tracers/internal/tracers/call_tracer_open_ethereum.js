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

	// descended tracks whether we've just descended from an outer transaction into
	// an inner call.
	descended: false,

	oeErrorMapping: {
		"contract creation code storage out of gas": "Out of gas",
		"out of gas": "Out of gas",
		"gas uint64 overflow": "Out of gas",
		"max code size exceeded": "Out of gas",
		"invalid jump destination": "Bad jump destination",
		"execution reverted": "Reverted",
		"return data out of bounds": "Out of bounds",
	},

	oeErrorMappingStartingWith: {
		"invalid opcode:": "Bad instruction",
		"stack underflow": "Stack underflow",
	},

	// step is invoked for every opcode that the VM executes.
	step: function(log, db) {
		// Capture any errors immediately
		var error = log.getError();
		if (typeof error !== "undefined") {
			this.fault(log, db);
			return;
		}
		// We only care about system opcodes, faster if we pre-check once
		var syscall = (log.op.toNumber() & 0xf0) == 0xf0;
		if (syscall) {
			var op = log.op.toString();
		}
		// If a new contract is being created, add to the call stack
		if (syscall && (op == "CREATE" || op == "CREATE2")) {
			var inOff = log.stack.peek(1).valueOf();
			var inEnd = inOff + log.stack.peek(2).valueOf();

			// Assemble the internal call report and store for completion
			var call = {
				type:    op,
				from:    toHex(log.contract.getAddress()),
				input:   toHex(log.memory.slice(inOff, inEnd)),
				gasIn:   log.getGas(),
				gasCost: log.getCost(),
				value:   "0x" + log.stack.peek(0).toString(16)
			};
			this.callstack.push(call);
			this.descended = true;
			return;
		}
		// If a contract is being self destructed, gather that as a subcall too
		if (syscall && op == "SELFDESTRUCT") {
			var left = this.callstack.length;
			if (typeof this.callstack[left-1].calls === "undefined") {
				this.callstack[left-1].calls = [];
			}
			this.callstack[left-1].calls.push({
				type:    op,
				from:    toHex(log.contract.getAddress()),
				to:      toHex(toAddress(log.stack.peek(0).toString(16))),
				gasIn:   log.getGas(),
				gasCost: log.getCost(),
				value:   "0x" + db.getBalance(log.contract.getAddress()).toString(16)
			});
			return;
		}
		// If a new method invocation is being done, add to the call stack
		if (syscall && (op == "CALL" || op == "CALLCODE" || op == "DELEGATECALL" || op == "STATICCALL")) {
			// Skip any pre-compile invocations, those are just fancy opcodes
			var to = toAddress(log.stack.peek(1).toString(16));
			if (isPrecompiled(to)) {
				return
			}
			var off = (op == "DELEGATECALL" || op == "STATICCALL" ? 0 : 1);

			var inOff = log.stack.peek(2 + off).valueOf();
			var inEnd = inOff + log.stack.peek(3 + off).valueOf();

			// Assemble the internal call report and store for completion
			var call = {
				type:    op,
				from:    toHex(log.contract.getAddress()),
				to:      toHex(to),
				input:   toHex(log.memory.slice(inOff, inEnd)),
				gasIn:   log.getGas(),
				gasCost: log.getCost(),
				outOff:  log.stack.peek(4 + off).valueOf(),
				outLen:  log.stack.peek(5 + off).valueOf()
			};

			if (op != "DELEGATECALL" && op != "STATICCALL") {
				call.value = "0x" + log.stack.peek(2).toString(16);

			} else if (op == "STATICCALL") {
				call.value = "0x0";
			}

			this.callstack.push(call);
			this.descended = true;
			return;
		}
		// If we've just descended into an inner call, retrieve it's true allowance. We
		// need to extract if from within the call as there may be funky gas dynamics
		// with regard to requested and actually given gas (2300 stipend, 63/64 rule).
		if (this.descended) {
			if (log.getDepth() >= this.callstack.length) {
				this.callstack[this.callstack.length - 1].gas = log.getGas();
			} else {
				// TODO(karalabe): The call was made to a plain account. We currently don't
				// have access to the true gas amount inside the call and so any amount will
				// mostly be wrong since it depends on a lot of input args. Skip gas for now.
			}
			this.descended = false;
		}
		// If an existing call is returning, pop off the call stack
		if (syscall && op == "REVERT") {
			this.callstack[this.callstack.length - 1].error = "execution reverted";
			return;
		}
		if (syscall && op == "RETURN") {
			var outOff = log.stack.peek(0).valueOf();
			var outLen = log.stack.peek(1).valueOf();

			this.callstack[this.callstack.length - 1].output = toHex(log.memory.slice(outOff, outOff + outLen));
			return;
		}
		if (log.getDepth() == this.callstack.length - 1) {
			// Pop off the last call and get the execution results
			var call = this.callstack.pop();

			if (call.type == 'CREATE' || call.type == "CREATE2") {
				// If the call was a CREATE, retrieve the contract address and output code
				call.gasUsed = '0x' + bigInt(call.gasIn - call.gasCost - log.getGas()).toString(16);
				delete call.gasIn; delete call.gasCost;

				var ret = log.stack.peek(0);
				if (!ret.equals(0)) {
					call.to     = toHex(toAddress(ret.toString(16)));
				} else if (typeof call.error === "undefined") {
					call.error = "internal failure"; // TODO(karalabe): surface these faults somehow

					if (typeof call.gas !== "undefined" && call.gasUsed === '0x' + bigInt(call.gas).toString(16)) {
						call.error = "out of gas";
					} else {
						return;
					}
				}
			} else {

				// If the call was a contract call, retrieve the gas usage and output
				if (typeof call.gas !== "undefined") {
					var gasUsed = bigInt(call.gasIn - call.gasCost + call.gas - log.getGas());
					if (gasUsed.compare(0) >= 0) {
						call.gasUsed = '0x' + gasUsed.toString(16);
					}
				} else {
					call.gas = bigInt(call.gasIn - call.gasCost - log.getGas()).abs();
					call.gasUsed = '0x0';
				}

				var ret = log.stack.peek(0);
				if (ret.equals(0) && typeof call.error === "undefined") {
					call.error = "internal failure"; // TODO(karalabe): surface these faults somehow

					if (typeof call.gas !== "undefined" && call.gasUsed === '0x' + bigInt(call.gas).toString(16)) {
						call.error = "out of gas";
					} else {
						return;
					}
				}

				if (typeof call.output === "undefined") {
					call.output = '0x';
				}

				delete call.gasIn; delete call.gasCost;
				delete call.outOff; delete call.outLen;
			}
			if (typeof call.gas !== "undefined") {
				call.gas = '0x' + bigInt(call.gas).toString(16);
			}

			// Inject the call into the previous one
			var left = this.callstack.length;
			if (typeof this.callstack[left-1].calls === "undefined") {
				this.callstack[left-1].calls = [];
			}
			this.callstack[left-1].calls.push(call);
		}
	},

	// fault is invoked when the actual execution of an opcode fails.
	fault: function(log, db) {
		// If the topmost call already reverted, don't handle the additional fault again
		if (typeof this.callstack[this.callstack.length - 1].error !== "undefined") {
			return;
		}
		// Pop off the just failed call
		var call = this.callstack.pop();
		call.error = log.getError();

		// Consume all available gas and clean any leftovers
		if (typeof call.gas !== "undefined") {
			call.gas = '0x' + bigInt(call.gas).toString(16);
			call.gasUsed = call.gas
		} else {
			// Retrieve gas true allowance from the inner call.
			// We need to extract if from within the call as there may be funky gas dynamics
			// with regard to requested and actually given gas (2300 stipend, 63/64 rule).
			call.gas = '0x' + bigInt(log.getGas()).toString(16);
		}

		if (call.error === "out of gas" && typeof call.gas === "undefined") {
			call.gas = "0x0";
		}
		delete call.gasIn; delete call.gasCost;
		delete call.outOff; delete call.outLen;

		// Flatten the failed call into its parent
		var left = this.callstack.length;
		if (left > 0) {
			if (typeof this.callstack[left-1].calls === "undefined") {
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
			blockNumber: ctx.blockNumber,
			transactionHash: ctx.transactionHash,
			transactionPosition: ctx.transactionPosition,
		};
		if (typeof this.callstack[0].calls !== "undefined") {
			result.calls = this.callstack[0].calls;
		}
		if (typeof this.callstack[0].error !== "undefined") {
			result.error = this.callstack[0].error;
		} else if (typeof ctx.error !== "undefined") {
			result.error = ctx.error;
		}
		if (typeof result.error !== "undefined" && (result.error !== "execution reverted" || result.output ==="0x")) {
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
			blockNumber: call.block || extraCtx.blockNumber,
			blockHash: extraCtx.blockHash,
			time: call.time,
		}

		if (typeof sorted.error !== "undefined") {
			if (this.oeErrorMapping.hasOwnProperty(sorted.error)) {
				sorted.error = this.oeErrorMapping[sorted.error];
				delete sorted.result;
			} else {
				for (var searchKey in this.oeErrorMappingStartingWith) {
					if (this.oeErrorMappingStartingWith.hasOwnProperty(searchKey) && sorted.error.indexOf(searchKey) > -1) {
						sorted.error = this.oeErrorMappingStartingWith[searchKey];
						delete sorted.result;
					}
				}
			}
		}

		for (var key in sorted) {
			if (typeof sorted[key] === "object") {
				for (var nested_key in sorted[key]) {
					if (typeof sorted[key][nested_key] === "undefined") {
						delete sorted[key][nested_key];
					}
				}
			} else if (typeof sorted[key] === "undefined") {
				delete sorted[key];
			}
		}

		var calls = call.calls;
		if (typeof calls !== "undefined") {
			sorted["subtraces"] = calls.length;
		}

		var results = [sorted];

		if (typeof calls !== "undefined") {
			for (var i=0; i<calls.length; i++) {
				var childCall = calls[i];

				// Delegatecall uses the value from parent
				if (childCall.type == 'DELEGATECALL' && typeof childCall.value === "undefined") {
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
