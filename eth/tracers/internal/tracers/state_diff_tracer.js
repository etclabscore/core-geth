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

// stateDiffTracer outputs sufficient information to create a local execution of
// the transaction from a custom assembled genesisT block.
{
	// stateDiff is the genesisT that we're building.
	stateDiff: {},

	hasInitCalled: false,
	lastAccessedAccount: null,
	lastRefund: 0,
	hasError: false,

	diffMarkers: {
		Memory: "_",	// temp state used while running the tracer, will never be returned to the user
		Born: "+",
		Died: "-",
		Changed: "*",
		Same: "=",
	},

	isObjectEmpty: function(obj) {
		for (var x in obj) { return false; }
		return true;
	},

	toHexJs: function(val) {
		return typeof val !== "string" || val.indexOf("0x") !== 0 ? "0x" + val.toString(16) : val;
	},

	accountInit: function(acc, type) {
		if (this.stateDiff[acc] === undefined) {
			var memoryMarker = this.diffMarkers.Memory;

			this.stateDiff[acc] = {
				_type: type || this.diffMarkers.Changed, // temp storage of account's initial type
				_removed: false, // removed from state
				_error: false, // error returned by the VM
				_final: false, // stop updating state if account's state marked as final
				balance: {
					[memoryMarker]: {}
				},
				nonce: {
					[memoryMarker]: {}
				},
				code: {
					[memoryMarker]: {}
				},
				storage: {}
			};
		}
	},

	// lookupAccount injects the specified account into the stateDiff object
	lookupAccount: function(addr, db, type) {
		var acc = toHex(addr);
		this.lastAccessedAccount = acc;

		// no need to fetch updates, as the account is marked as final
		// which means that some manual calculations have been performed, usually in results()
		// NOTE: it moved at the top of fuction in order we don't read from statedb when not needed
		if (this.stateDiff[acc] !== undefined && this.stateDiff[acc]._final) {
			return;
		}

		var balance = db.getBalance(addr);
		var code = toHex(db.getCode(addr));
		var nonce = db.getNonce(addr);

		var memoryMarker = this.diffMarkers.Memory;

		this.accountInit(acc, type);

		var accountData = this.stateDiff[acc];

		// if (balance|nonce|code).from is not filled, then this is the first time lookupAccount is called
		if (this.stateDiff[acc].balance[memoryMarker].from === undefined) {
			accountData.balance[memoryMarker].from = balance;
			accountData.nonce[memoryMarker].from = nonce;
			accountData.code[memoryMarker].from = code;
		}

		// force type change
		if (type !== undefined) {
			// if an account has been Born within this run and also Died,
			// this means it will never be persisted to the state
			if (accountData._type === this.diffMarkers.Born && type === this.diffMarkers.Died) {
				accountData._removed = true;
			}

			accountData._type = type;
		}

		if (balance) {
			accountData.balance[memoryMarker].to = balance;
		}

		// if state doesn't have the account, most probably because of EIP-161, then remove it
		if (typeof nonce === "number" && nonce < accountData.nonce[memoryMarker].from) {
			delete this.stateDiff[acc];
			return;
		} else if (nonce) {
			accountData.nonce[memoryMarker].to = nonce;
		}

		if (code) {
			accountData.code[memoryMarker].to = code;
		}
	},

	// lookupStorage injects the specified storage entry of the given account into
	// the stateDiff object
	lookupStorage: function(addr, key, val, db){
		var acc = toHex(addr);
		this.lastAccessedAccount = acc;

		this.accountInit(acc);

		var accountData = this.stateDiff[acc];

		var memoryMarker = this.diffMarkers.Memory;
		var idx = toHex(key);

		if (accountData.storage[idx] === undefined) {
			var rval = toHex(db.getState(addr, key));
			accountData.storage[idx] = {
				[memoryMarker]: {
					// marker keeping if storage idx had data at least once
					_changed: !(/^(0x)?0*$/.test(rval)),
					from: rval
				}
			};
		}

		if (val) {
			var sval = toHex(val);
			accountData.storage[idx][memoryMarker].to = sval;

			if (!accountData.storage[idx][memoryMarker]._changed && !(/^(0x)?0*$/.test(sval))) {
				accountData.storage[idx][memoryMarker]._changed = true;
			}
		} else {
			var rval = toHex(db.getState(addr, key));
			accountData.storage[idx][memoryMarker].to = rval;

			if (!accountData.storage[idx][memoryMarker]._changed && !(/^(0x)?0*$/.test(rval))) {
				accountData.storage[idx][memoryMarker]._changed = true;
			}
		}
	},

	// formatSingle used to format output for Born and Died markers
	formatSingle: function(data, type) {
		type = type || this.diffMarkers.Changed;

		var memoryMarker = this.diffMarkers.Memory;
		var val = data[memoryMarker].to || data[memoryMarker].from;

		// for Died markers, we want to output the "from" value
		if (type === this.diffMarkers.Died) {
			val = data[memoryMarker].from
		}

		// this is mostly for balances and nonces, where we can't have negative values
		// it would have been considered in consensus rules, though for trace_call
		// we break the logic on VM's CanTransfer & Transfer
		if (bigInt.isInstance(val) && val.isNegative()) {
			val = bigInt.zero;
		}

		return {
			[type]: this.toHexJs(val),
		}
	},

	// formatChanged used to format output for Changed marker,
	// which handles the Same market internally
	formatChanged: function(data, type) {
		type = type || this.diffMarkers.Changed;

		var memoryMarker = this.diffMarkers.Memory;
		var from = data[memoryMarker].from;
		var to = data[memoryMarker].to;

		// this is mostly for balances and nonces, where we can't have negative values
		// it would have been considered in consensus rules, though for trace_call
		// we break the logic on VM's CanTransfer & Transfer
		if (bigInt.isInstance(from) && from.isNegative()) {
			from = bigInt.zero;
			data[memoryMarker].from = from; // used for checkIfSame
		}

		if (bigInt.isInstance(to) && to.isNegative()) {
			to = bigInt.zero;
			data[memoryMarker].to = to; // used for checkIfSame
		}

		if (this.checkIfSame(data, type)) {
			return this.diffMarkers.Same;
		}

		return {
			[this.diffMarkers.Changed]: {
				from: this.toHexJs(from),
				to: this.toHexJs(to),
			}
		}
	},

	// checkIfSame check if both "from" and "to" have the same values
	checkIfSame: function(data, type) {
		type = type || this.diffMarkers.Changed;

		var memoryMarker = this.diffMarkers.Memory;
		var from = data[memoryMarker].from;
		var to = data[memoryMarker].to;

		return (to === undefined ||
			(bigInt.isInstance(from) && from.compare(to) === 0) ||
			from === to);
	},

	// hasAccountChanges checks if the account has any changes
	// and has to be added to stateDiff's output
	hasAccountChanges: function(data) {
		var sameMarker = this.diffMarkers.Same;

		return !(data.balance === sameMarker
				&& data.nonce === sameMarker
				&& data.code === sameMarker
				&& this.isObjectEmpty(data.storage));
	},

	format: function(db) {
		for (var acc in this.stateDiff) {
			var accountAddress = toAddress(acc);
			// fetch latest balance
			this.lookupAccount(accountAddress, db);

			var memoryMarker = this.diffMarkers.Memory;
			var changedMarker = this.diffMarkers.Changed;

			var accountData = this.stateDiff[acc];

			// remove accounts with errors or marked for removal. Do this check last (not in another method),
			// otherwise account might be re-added for ctx.[from|to|coinbase] from within this.result()
			if (accountData._error || accountData._removed) {
				delete this.stateDiff[acc];
				continue;
			}

			// check if it is a new borned account
			// happens on mordor's tx: 0xc468750bd7d73f53ff3fdc74201245910d84d84bfc5c40d97e4c5a8928c92187
			if (accountData._type === changedMarker
					&& accountData.balance[memoryMarker].from == 0
					&& accountData.code[memoryMarker].from === "0x"
					&& accountData.nonce[memoryMarker].from == 0) {
				accountData._type = this.diffMarkers.Born;
			}

			var type = accountData._type;
			delete accountData._type;
			delete accountData._removed;
			delete accountData._error;
			delete accountData._final;

			if (type === changedMarker) {
				accountData.balance = this.formatChanged(accountData.balance, type);
				accountData.nonce = this.formatChanged(accountData.nonce, type);
				accountData.code = this.formatChanged(accountData.code, type);
			} else {
				accountData.balance = this.formatSingle(accountData.balance, type);
				accountData.nonce = this.formatSingle(accountData.nonce, type);
				accountData.code = this.formatSingle(accountData.code, type);
			}

			if (db.isEmpty(toAddress(acc))) {
				delete this.stateDiff[acc];
				continue;
			}

			// optimisation: pre-check if we have changes before parsing storage state
			if (!this.hasAccountChanges(accountData)) {
				delete this.stateDiff[acc];
				continue;
			}

			// handle storage entries
			for (var idx in accountData.storage) {
				// fetch latest value in storage
				this.lookupStorage(accountAddress, toWord(idx), null, db);

				var sti = accountData.storage[idx];
				if (sti[memoryMarker] === undefined
						|| sti[memoryMarker].to === undefined
						|| !sti[memoryMarker]._changed) {
					delete this.stateDiff[acc].storage[idx];
					continue;
				}

				if (type === changedMarker) {
					var res = this.formatChanged(sti, type);
					if (res === this.diffMarkers.Same) {
						delete this.stateDiff[acc].storage[idx];
					} else {
						accountData.storage[idx] = res;
					}
				} else {
					// when Died and from|to is same, then remove the storage entry from output
					// happens on mordor's tx: 0x0c59ddf8ebbaa64140db6214bbad641fff6bb066847dbef3433d434bd1fb6270
					if (type === this.diffMarkers.Died && this.checkIfSame(sti, type)) {
						delete this.stateDiff[acc].storage[idx];
					} else {
						accountData.storage[idx] = this.formatSingle(sti, type);
					}
				}
			}

			// remove unchanged accounts
			if (!this.hasAccountChanges(accountData)) {
				delete this.stateDiff[acc];
				continue;
			}
		}
	},

	// includeOpError checks for specific VM OP errors
	includeOpError: function(err) {
		return err
			&& err.indexOf('contract address collision') > -1;
			// && (err.indexOf('insufficient balance for transfer') > -1
	},

	// init is invoked before any VM execution.
	// ctx has to|msgTo|coinbase set and additional context based on each trace method.
	init: function(ctx, db) {
		this.hasInitCalled = true;

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

	// step is invoked for every opcode that the VM executes
	step: function(log, db) {
		// Capture any errors immediately
		var error = log.getError();
		var opError = log.getCallError();
		if (!this.hasError && (error !== undefined || opError !== undefined)) {
			this.hasError = true;
		}

		if ((error !== undefined || this.includeOpError(opError))
				&& this.lastAccessedAccount !== null
				&& this.stateDiff[this.lastAccessedAccount] !== undefined) {
			this.stateDiff[this.lastAccessedAccount]._error = true;  // mark account that had an error
			this.lastAccessedAccount = null;
			return;
		}

		this.lastRefund = log.getRefund();

		// whenever new state is accessed, add it to the stateDiff
		switch (log.op.toString()) {
			case "EXTCODECOPY": case "EXTCODESIZE": case "BALANCE":
				this.lookupAccount(toAddress(log.stack.peek(0).toString(16)), db);
				break;
			case "CREATE":
				var address = log.contract.getAddress();
				this.lookupAccount(toContract(address, db.getNonce(address)), db, this.diffMarkers.Born);
				break;
			case "CREATE2":
				// stack: salt, size, offset, endowment
				var offset = log.stack.peek(1).valueOf()
				var size = log.stack.peek(2).valueOf()
				var end = offset + size
				this.lookupAccount(toContract2(log.contract.getAddress(), log.stack.peek(3).toString(16), log.memory.slice(offset, end)), db, this.diffMarkers.Born);
				break;
			case "CALL": case "CALLCODE": case "DELEGATECALL": case "STATICCALL":
				var address = toAddress(log.stack.peek(1).toString(16));

				// no need to handle anything in pre-compiles
				// while also helping, to maintain the lastAccessedAccount logic
				if (isPrecompiled(address)) {
					break;
				}

				this.lookupAccount(address, db);
				break;
			case "SLOAD":
				this.lookupStorage(log.contract.getAddress(), toWord(log.stack.peek(0).toString(16)), null, db);
				break;
			case "SSTORE":
				this.lookupStorage(log.contract.getAddress(), toWord(log.stack.peek(0).toString(16)), toWord(log.stack.peek(1).toString(16)), db);
				break;
			case "SELFDESTRUCT":
				this.lookupAccount(log.contract.getAddress(), db, this.diffMarkers.Died);
				break;
		}
	},

	// fault is invoked when the actual execution of an opcode fails
	fault: function(log, db) {
		var error = log.getError();
		var opError = log.getCallError();
		if (!this.hasError && (error !== undefined || opError !== undefined)) {
			this.hasError = true;
		}

		if ((error !== undefined || this.includeOpError(opError))
				&& this.lastAccessedAccount !== null
				&& this.stateDiff[this.lastAccessedAccount] !== undefined) {
			this.stateDiff[this.lastAccessedAccount]._error = true;  // mark account that had an error
			this.lastAccessedAccount = null;
		}
	},

	// result is invoked when all the opcodes have been iterated over and returns
	// the final result of the tracing
	result: function(ctx, db) {

		// KEEP for tests:
		// 0x9965b02962ecd7aa6867fb8ea3357a9fd26f07f62c768ae4c9f902130babe97d
		// 0x9d1a0f214ebc5d727fbc9c0dd299a2e76c8321fb5fc552fab87eb6028ab2239d?
		// 0x870e57c81ae99c0bdc24351af834bfc571e9596c57202d55d77ff6a633854f5d
		// 0x0c59ddf8ebbaa64140db6214bbad641fff6bb066847dbef3433d434bd1fb6270 // Died marker (keep)
		// 0xf18306dcc1badc05c32be8b91d31d6fcd1c8003e71c770fd69bbf77623cbbbdc // Died marker (remove)
		// 0x00671034509a65920422f3f5060039183c9a04b3692c89c1bc7d92e27bd1fb83 // Slow TX (in general after this block)

		var memoryMarker = this.diffMarkers.Memory;

		// get actual "to" values for from|to|coinbase accounts
		this.lookupAccount(ctx.from, db);
		this.lookupAccount(ctx.to, db);
		this.lookupAccount(ctx.coinbase, db);

		var fromAccountHex = toHex(ctx.from);
		var toAccountHex = toHex(ctx.to);
		var coinbaseHex = toHex(ctx.coinbase);

		var convertCtxKeysToBigInt = ['gasLimit', 'gas', 'gasPrice', 'gasUsed', 'value'];
		for (var i in convertCtxKeysToBigInt) {
			var key = convertCtxKeysToBigInt[i];
			ctx[key] = bigInt(ctx[key]);
		}

		var refund = !this.hasError && this.lastRefund > 0 ? this.lastRefund : 0;

		var gasUsed = ctx.gasLimit.subtract(ctx.gas).add(ctx.gasUsed).add(refund);
		var gasLeft = ctx.gasLimit.subtract(gasUsed);

		var refundValue = gasLeft.multiply(ctx.gasPrice);
		var feesValue = gasUsed.multiply(ctx.gasPrice);
		var fullGasCost = ctx.gasLimit.multiply(ctx.gasPrice);

		// in case from balance is negative because the tracer has disabled the CanTransfer check,
		// and the Transfer happened before the CaptureStart and the interpreter execution
		var hasFromSufficientBalanceForValueAndGasCost = ctx.hasFromSufficientBalanceForValueAndGasCost || false;
		var hasFromSufficientBalanceForGasCost = ctx.hasFromSufficientBalanceForGasCost || false;

		var isCreateType = ctx.type == "CREATE" || ctx.type == "CREATE2";
		var isCallTypeOnNonExistingAccount = ctx.type == "CALL" && ctx.value.isZero() && !db.exists(ctx.to) && !isPrecompiled(ctx.to)

		// a transaction with "value" set/positive,
		// while "from" account has not enough balance to pay even for the gas cost,
		// will not be run at all and will not change the state
		if (!hasFromSufficientBalanceForGasCost && !hasFromSufficientBalanceForValueAndGasCost && ctx.value.isPositive()) {
			return {};
		}

		if (this.stateDiff[fromAccountHex] !== undefined) {
			var fromAcc = this.stateDiff[fromAccountHex];

			// remove any errors marked on the from account, as it has to be included on output
			// happens on mordor tx: 0x8f26c1acfce0178a2b037d85feeea99bb961bb46f541ad8c01c6668455952221
			fromAcc._error = false;

			// don't update from state anymore, as we applied our customs calcs
			fromAcc._final = true;
		}

		if (this.stateDiff[toAccountHex] !== undefined) {
			var toAcc = this.stateDiff[toAccountHex];

			// TODO: check if it needs to be added back
			// if (!hasFromSufficientBalanceForValueAndGasCost) {
			// 	// in case from account doesn't have enough balance, the transfer won't happen
			// 	toAcc.balance[memoryMarker].to = toAcc.balance[memoryMarker].from;
			// }

			if (isCreateType) {
				if (toAcc._type !== this.diffMarkers.Died) {
					// mark new created contracts address
					toAcc._type = this.diffMarkers.Born;
				} else {
					// if the new contract has not be persisted to state then remove it from stateDiff output
					toAcc._removed = true;
				}
			} else {
				// remove any errors marked on the to account (if not type = CREATE*), as it has to be included on output
				// happens on mordor tx: 0x89fd95d97374ccb9cdac249c74efdc57907c53beecb3e6ebce03b4ca31b0df2f
				toAcc._error = false;
			}

			// don't update from state anymore, as we applied our customs calcs
			toAcc._final = true;
		}

		if (this.stateDiff[coinbaseHex] !== undefined) {
			var coinbaseAcc = this.stateDiff[coinbaseHex];

			// remove any errors marked on the coinbase account, as it has to be included on output
			// happens on mordor tx: 0xbfca41d82781ba1888c10d96de84ff68799e328c658b34964d382eba019b3752
			coinbaseAcc._error = false;

			// don't update from state anymore, as we applied our customs calcs
			coinbaseAcc._final = true;
		}

		this.format(db);

		// Return the assembled allocations (stateDiff)
		return this.stateDiff;
	}
}
