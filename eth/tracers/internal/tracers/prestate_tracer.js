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

// prestateTracer outputs sufficient information to create a local execution of
// the transaction from a custom assembled genesisT block.
{
	// prestate is the genesisT that we're building.
	prestate: {},
	hasStepInited: false,

	// lookupAccount injects the specified account into the prestate object.
	lookupAccount: function(addr, db){
		var acc = toHex(addr);
		if (this.prestate[acc] === undefined) {
			this.prestate[acc] = {
				balance: '0x' + db.getBalance(addr).toString(16),
				nonce:   db.getNonce(addr),
				code:    toHex(db.getCode(addr)),
				storage: {}
			};
		}
	},

	// lookupStorage injects the specified storage entry of the given account into
	// the prestate object.
	lookupStorage: function(addr, key, db){
		var acc = toHex(addr);
		var idx = toHex(key);

		if (this.prestate[acc].storage[idx] === undefined) {
			this.prestate[acc].storage[idx] = toHex(db.getState(addr, key));
		}
	},

	// init is invoked before any VM execution.
	// ctx has to|msgTo|coinbase set and additional context based on each trace method.
	init: function(ctx, db) {
		// get actual "from" values for from|to|coinbase accounts.
		this.lookupAccount(ctx.from, db);
		this.lookupAccount(ctx.coinbase, db);

		// msgTo is set for the init method and it is the actual "to" value of the Tx.
		// ctx.to on the other hand is always set by the EVM and on type=CREATE
		// it is the newly created contract address
		if (ctx.msgTo !== undefined) {
			this.lookupAccount(ctx.msgTo, db);
		}
	},

	// result is invoked when all the opcodes have been iterated over and returns
	// the final result of the tracing.
	result: function(ctx, db) {
		if (ctx.type == 'CREATE') {
			// We can blibdly delete the contract prestate, as any existing state would
			// have caused the transaction to be rejected as invalid in the first place.
			delete this.prestate[toHex(ctx.to)];
		}
		// Return the assembled allocations (prestate)
		return this.prestate;
	},

	// step is invoked for every opcode that the VM executes.
	step: function(log, db) {
		if (!this.hasStepInited) {
			this.hasStepInited = true;

			this.lookupAccount(log.contract.getAddress(), db);
		}

		// Whenever new state is accessed, add it to the prestate
		switch (log.op.toString()) {
			case "EXTCODECOPY": case "EXTCODESIZE": case "BALANCE":
				this.lookupAccount(toAddress(log.stack.peek(0).toString(16)), db);
				break;
			case "CREATE":
				var from = log.contract.getAddress();
				this.lookupAccount(toContract(from, db.getNonce(from)), db);
				break;
			case "CREATE2":
				var from = log.contract.getAddress();
				// stack: salt, size, offset, endowment
				var offset = log.stack.peek(1).valueOf()
				var size = log.stack.peek(2).valueOf()
				var end = offset + size
				this.lookupAccount(toContract2(from, log.stack.peek(3).toString(16), log.memory.slice(offset, end)), db);
				break;
			case "CALL": case "CALLCODE": case "DELEGATECALL": case "STATICCALL":
				this.lookupAccount(toAddress(log.stack.peek(1).toString(16)), db);
				break;
			case 'SSTORE':case 'SLOAD':
				this.lookupStorage(log.contract.getAddress(), toWord(log.stack.peek(0).toString(16)), db);
				break;
		}
	},

	// fault is invoked when the actual execution of an opcode fails.
	fault: function(log, db) {}
}
