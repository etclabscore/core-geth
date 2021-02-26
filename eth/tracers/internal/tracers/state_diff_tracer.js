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
	lastAccessedAccount: null, // TODO: check as we might can remove it at the end
	lastGasIn: null,	// TODO: this can be removed, keep it until we are sure
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

	// lookupAccount injects the specified account into the stateDiff object.
	lookupAccount: function(addr, db, type) {
		var memoryMarker = this.diffMarkers.Memory;

		var acc = toHex(addr);
		var balance = db.getBalance(addr);
		var code = toHex(db.getCode(addr));
		var nonce = db.getNonce(addr);

		this.lastAccessedAccount = acc;

		if (this.stateDiff[acc] === undefined) {
			this.stateDiff[acc] = {
				_type: type || this.diffMarkers.Changed,  // temp storage of account's initial type
				_removed: false, // removed from state
				_error: false, // evm returned an error
				_final: false, // stop updating state if account's state marked as final
				balance: {
					[memoryMarker]: {
						from: balance,
					}
				},
				nonce: {
					[memoryMarker]: {
						from: nonce,
					}
				},
				code: {
					[memoryMarker]: {
						from: code,
					}
				},
				storage: {}
			};
		}

		var accountData = this.stateDiff[acc];

		if (accountData._final) {
			return;
		}

		// console.log('nonce', type, acc, accountData.nonce, '=>',nonce)


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

		// var latestKnownNonce = accountData.nonce[memoryMarker].to || accountData.nonce[memoryMarker].from
		// if (nonce && type === this.diffMarkers.Born) {
		// 	accountData.nonce[memoryMarker].to = latestKnownNonce + 1;
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
	// the stateDiff object.
	lookupStorage: function(addr, key, val, db){
		var acc = toHex(addr);
		var idx = toHex(key);

		this.lastAccessedAccount = acc;

		var memoryMarker = this.diffMarkers.Memory;

		if (this.stateDiff[acc] === undefined) {
			return;
		}

		var accountData = this.stateDiff[acc];

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
			var tval = toHex(val);
			accountData.storage[idx][memoryMarker].to = tval;

			if (!accountData.storage[idx][memoryMarker]._changed && !(/^(0x)?0*$/.test(tval))) {
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

	formatSingle: function(data, type) {
		type = type || this.diffMarkers.Changed;

		var memoryMarker = this.diffMarkers.Memory;
		var val = data[memoryMarker].to || data[memoryMarker].from;

		if (type === this.diffMarkers.Died) {
			val = data[memoryMarker].from
		}

		if (bigInt.isInstance(val) && val.isNegative()) {
			val = bigInt.zero;
		}

		return {
			[type]: this.toHexJs(val),
		}
	},

	formatChanged: function(data, type) {
		type = type || this.diffMarkers.Changed;

		var memoryMarker = this.diffMarkers.Memory;
		var changedMarker = this.diffMarkers.Changed;
		var sameMarker = this.diffMarkers.Same;

		var from = data[memoryMarker].from;
		var to = data[memoryMarker].to;

		if (bigInt.isInstance(from) && from.isNegative()) {
			from = bigInt.zero;
			data[memoryMarker].from = from; // used for checkIfSame
		}

		if (bigInt.isInstance(to) && to.isNegative()) {
			to = bigInt.zero;
			data[memoryMarker].to = to; // used for checkIfSame
		}

		if (this.checkIfSame(data, type)) {
			return sameMarker;
		}

		return {
			[changedMarker]: {
				from: this.toHexJs(from),
				to: this.toHexJs(to),
			}
		}
	},

	checkIfSame: function(data, type) {
		type = type || this.diffMarkers.Changed;

		var memoryMarker = this.diffMarkers.Memory;

		var from = data[memoryMarker].from;
		var to = data[memoryMarker].to;

		return (to === undefined ||
			(bigInt.isInstance(from) && from.compare(to) === 0) ||
			from === to);
	},

	hasAccountChanges: function(data) {
		var sameMarker = this.diffMarkers.Same;
		// var bornMarker = this.diffMarkers.Born;

		if (data.balance === sameMarker
				&& data.nonce === sameMarker
				&& data.code === sameMarker
				&& this.isObjectEmpty(data.storage)
		) {
			return false;
		// } else if (data.balance[bornMarker] === "0x0" &&
		// 		data.nonce[bornMarker] === "0x0" &&
		// 		data.code[bornMarker] === "0x" &&
		// 		this.isObjectEmpty(data.storage)
		// ) {
		// 	return false;
		}
		return true;
	},

	format: function(db) {
		for (var acc in this.stateDiff) {
			var accountAddress = toAddress(acc);
			// fetch latest balance
			// TODO: optimise
			this.lookupAccount(accountAddress, db);

			// has been cleared in lookupAccount
			if (this.stateDiff[acc] === undefined) {
				continue;
			}

			var memoryMarker = this.diffMarkers.Memory;
			var changedMarker = this.diffMarkers.Changed;

			var accountData = this.stateDiff[acc];

			// remove accounts with errors. Do this check last (not in another method),
			// otherwise account might be re-added if one of ctx.[from|to|coinbase] from within this.result()
			if (accountData._error || accountData._removed) {
				delete this.stateDiff[acc];
				continue;
			}

			// check if it is a new borned account
			if (accountData.balance[memoryMarker].from == 0
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

			for (var idx in accountData.storage) {
				// fetch latest value in storage
				this.lookupStorage(accountAddress, toWord(idx), null, db);

				var sti = accountData.storage[idx];
				if (sti[memoryMarker] === undefined
						|| sti[memoryMarker].to === undefined
						|| !sti[memoryMarker]._changed) {
					delete this.stateDiff[acc].storage[idx];
				// } else if (accountData.storage[idx][type] === undefined ||
				// 		/^(0x)?0*$/.test(accountData.storage[idx][type])) {
				// 	delete this.stateDiff[acc].storage[idx];
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

	includeOpError: function(err) {
		return err
			&& err.indexOf('contract address collision') > -1;
	},

	// init is invoked on the first call VM executes.
	// IMPORTANT: it is being called only on contract calls and not on transfers,
	//            this is being handled in result()
	init: function(ctx, log, db) {
		this.hasInitCalled = true;

		// if (this.stateDiff === null) {
		// 	this.stateDiff = {};
		// }

		// Balances will potentially be wrong here, since they will include the value
		// sent along with the message or the paid full gas cost. We fix that in "result()".

		// Get actual "to" values for from, to, coinbase accounts
		this.lookupAccount(ctx.from, db);
		this.lookupAccount(ctx.coinbase, db);

		// TODO: do we need to check for contractAddress? is ctx.to === contractAddress?
		var toAccHex = toHex(ctx.to);
		if(!/^(0x)?0*$/.test(toAccHex)) {
			this.lookupAccount(ctx.to, db);
		}

		var contractAddress = log.contract.getAddress();
		if (toHex(contractAddress) !== toAccHex) {
			this.lookupAccount(contractAddress, db);
		}
	},

	// step is invoked for every opcode that the VM executes.
	step: function(log, db) {
		// Capture any errors immediately
		var error = log.getError();
		var opError = log.getCallError();
		if (!this.hasError && (error !== undefined || opError !== undefined)) {
			this.hasError = true;
		}

		// if (this.hasError) {
		// 	console.log('ERROR', log.op.toString())
		// }

		if ((error !== undefined || this.includeOpError(opError))
				&& this.lastAccessedAccount !== null
				&& this.stateDiff[this.lastAccessedAccount] !== undefined) {
        // console.log('damn line 318 \t this.lastAccessedAccount',log.op.toString(), this.lastAccessedAccount,error)

			this.stateDiff[this.lastAccessedAccount]._error = true;  // mark account that had an error
			this.lastAccessedAccount = null;
			return;
		}

		// Add the current account if we just started tracing
		// if (this.stateDiff === null){
		// 	this.stateDiff = {};

		// 	// var contractAddr = log.contract.getAddress();
		// 	// console.log('ex',db.exists(contractAddr));
		// 	// var type = db.exists(contractAddr) ? this.diffMarkers.Changed : this.diffMarkers.Born;
		// 	// Balance will potentially be wrong here, since this will include the value
		// 	// sent along with the message. We fix that in "result()".
		// 	this.lookupAccount(log.contract.getAddress(), db);
		// 	// console.log(log.contract.getAddress())
		// 	// this.lookupAccount(toAddress('0x893defcfe8dc3b7fe2b95c2ddd6415ac2f1eb582'), db);
		// }


		this.lastGasIn = log.getGas();
		this.lastRefund = log.getRefund();
    // console.log('347 \t this.lastRefund', this.lastRefund)
    // console.log('this.lastGasIn', 106376 - this.lastGasIn, 106376-log.getAvailableGas())
		// var loga = {
		// 	gasAvailableGas: log.getAvailableGas(),
		// 	gasGas:   log.getGas(),
		// 	gasCost: log.getCost(),
		// 	getValue: log.contract.getValue(),
		// }
		// console.log('loga', loga)


		// if (loga.gasAvailable > 0) {
		// 	console.log('AMAZINIIG', loga)
		// }
		// this.logger.push(loga)
		// console.log(log.op.toString(), log.stack.peek(0).toString(16), log.stack.peek(1).toString(16), log.stack.peek(2).toString(16), log.stack.peek(3).toString(16), log.stack.peek(4).toString(16));
		// console.log(log.op.toString());

		// Whenever new state is accessed, add it to the stateDiff
		switch (log.op.toString()) {
			// case "ORIGIN": case "CALLER": case "ADDRESS":
			// 	console.log(toHex(toAddress(log.stack.peek(0).toString(16))));
			// 	break;
			case "EXTCODECOPY": case "EXTCODESIZE": case "BALANCE":
				this.lookupAccount(toAddress(log.stack.peek(0).toString(16)), db);
				break;
			case "CREATE":
				var from = log.contract.getAddress();
				this.lookupAccount(toContract(from, db.getNonce(from)), db, this.diffMarkers.Born);
				break;
			case "CREATE2":
				var from = log.contract.getAddress();
				// stack: salt, size, offset, endowment
				var offset = log.stack.peek(1).valueOf()
				var size = log.stack.peek(2).valueOf()
				var end = offset + size
				this.lookupAccount(toContract2(from, log.stack.peek(3).toString(16), log.memory.slice(offset, end)), db, this.diffMarkers.Born);
				break;
			case "CALL": case "CALLCODE": case "DELEGATECALL": case "STATICCALL":
				var address = toAddress(log.stack.peek(1).toString(16));

				// No need to handle anything in pre-compiles
				// While also helping, to maintain the lastAccessedAccount logic
				if (isPrecompiled(address)) {
					break;
				}

				// if (log.op.toString() == "CALLCODE") {
				// 	var value = log.stack.peek(2);
        //   console.log('TCL: \t file: state_diff_tracer.js \t line 394 \t value', value)

				// 	if (value > 0) {
				// 		this.lastGasIn = bigInt(this.lastGasIn + 2300);
				// 	}
				// }
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

	// fault is invoked when the actual execution of an opcode fails.
	fault: function(log, db) {
		// this.lastGasIn = log.getGas();
		// console.log('fault this.lastGasIn', 106376 - this.lastGasIn, 106376-log.getAvailableGas())
		// console.log('fault', log.op.toString(), toHex(log.contract.getAddress()), log.stack.peek(0).toString(16), log.stack.peek(1).toString(16), log.stack.peek(2).toString(16), log.stack.peek(3).toString(16), log.stack.peek(4).toString(16));

		var error = log.getError();
    // console.log('TCL: \t file: state_diff_tracer.js \t line 325 \t error', error)
		var opError = log.getCallError();
		// console.log('TCL: \t file: state_diff_tracer.js \t line 327 \t opError', opError)
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

		// Reset lastAccessedAccount cleanup logic as it is being used only for `step` method
		// NOTE: it's safe to be removed, as it's not being utilised from this point
		// this.lastAccessedAccount = null;

		var memoryMarker = this.diffMarkers.Memory;

		if (this.lastGasIn === null) {
			this.lastGasIn = ctx.gas;
		}

		// Get actual "to" values for from, to, coinbase accounts
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
		var gasUsed = ctx.gasLimit.subtract(ctx.gas).add(ctx.gasUsed).add(refund); // bigInt(this.lastGasIn) == ctx.gasLimit.subtract(ctx.gas).add(ctx.gasUsed);

		// NOTE: both are equal for getting gasUsed
		// 1. gasLimit - gas (=intrinsic gas) + gasUsed
		// 2. gasLimit - this.lastGasIn
		var gasLeft = ctx.gasLimit.subtract(gasUsed);
		var refundValue = gasLeft.multiply(ctx.gasPrice);
		var feesValue = gasUsed.multiply(ctx.gasPrice);

		var fullGasCost = ctx.gasLimit.multiply(ctx.gasPrice);

		var gasCost = ctx.gasLimit					// full gas for tx
										.subtract(ctx.gas)	// gas given for this internal tx, calculates mostly the IntrinsicGas
										.add(ctx.gasUsed)		// gas used for this internal tx
										// .subtract(refund) 					// refund gas given
										.multiply(ctx.gasPrice);

		var coinbaseFees = ctx.gasLimit					// full gas for tx
												.subtract(ctx.gas)	// gas given for this internal tx, calculates mostly the IntrinsicGas
												.add(ctx.gasUsed)		// gas used for this internal tx
												// .subtract(refund) 					// refund gas given
												.multiply(ctx.gasPrice);


		console.log('coinbaseFees\t', coinbaseFees)
		console.log('fullGasCost\t', fullGasCost)
		var hasFromSufficientBalanceForValueAndGasCost = ctx.hasFromSufficientBalanceForValueAndGasCost || false;
		var hasFromSufficientBalanceForGasCost = ctx.hasFromSufficientBalanceForGasCost || false;

		var isCreateType = ctx.type == "CREATE" || ctx.type == "CREATE2";
		var isCallTypeWithZeroCodeForContract = !isCreateType && toHex(db.getCode(ctx.to)) == "0x"
		var isCallTypeOnNonExistingAccount = ctx.type == "CALL" && ctx.value.isZero() && !db.exists(ctx.to)
		// At this point, we need to deduct the "value" from the
		// outer transaction, and move it back to the origin

		// DEBUGGING GAS
		console.log('\nCtx -----');
		ctx.from = toHex(ctx.from);
		ctx.to = toHex(ctx.to);
		ctx.coinbase = toHex(ctx.coinbase);
		ctx.parentBlockCoinbase = toHex(ctx.parentBlockCoinbase);
		ctx.input = toHex(ctx.input).slice(0, 10) + '...<trimmed>';
		ctx.output = toHex(ctx.output).slice(0, 10) + '...<trimmed>';
		console.log(JSON.stringify(ctx, null, 2));

		console.log('\nto exists\t\t', db.exists(ctx.to));
		console.log('hasInitCalled\t\t', this.hasInitCalled);

		console.log('hasError\t\t', this.hasError)

		console.log('\nhasFromSufficientBalanceForValueAndGasCost\t', hasFromSufficientBalanceForValueAndGasCost)
		console.log('hasFromSufficientBalanceForGasCost\t\t', hasFromSufficientBalanceForGasCost)

		console.log('\nisCreateType\t\t\t\t', isCreateType)
		console.log('isCallTypeWithZeroCodeForContract\t', isCallTypeWithZeroCodeForContract)
		console.log('isCallTypeOnNonExistingAccount\t\t', isCallTypeOnNonExistingAccount)


		console.log('\nCalcs ----');
		console.log('\nvalue\t\t', ctx.value)
		console.log('refund\t\t', refund)

		console.log('\ngasUsed\t\t', gasUsed)
		console.log('gasLeft\t\t', gasLeft)
    console.log('refundValue\t', refundValue)
    console.log('feesValue\t', feesValue)
		console.log('fullGasCost\t\t', fullGasCost)

		console.log('gasCost\t\t', gasCost)
		console.log('lastGasIn\t', this.lastGasIn, '\t\t-gasUsed\t', this.lastGasIn - ctx.gasUsed)

		console.log('\ngas\t\t', ctx.gas)
		console.log('gasLimit\t', ctx.gasLimit, '\t-gasUsed\t', ctx.gasLimit - ctx.gasUsed)
		console.log('gasLimit\t', ctx.gasLimit, '\t-gas\t\t', ctx.gasLimit - ctx.gas)
		console.log('gasUsed\t\t', ctx.gasUsed)
		console.log('gasPrice\t', ctx.gasPrice)

		if (this.stateDiff[fromAccountHex] !== undefined) {
			console.log('\nbf.b fromA from\t\t', this.stateDiff[fromAccountHex].balance[memoryMarker].from)
			console.log('bf.b fromA to\t\t', this.stateDiff[fromAccountHex].balance[memoryMarker].to)
		}
		if (this.stateDiff[toAccountHex] !== undefined) {
			console.log('bf.b toA from\t\t', this.stateDiff[toAccountHex].balance[memoryMarker].from)
			console.log('bf.b toA to\t\t', this.stateDiff[toAccountHex].balance[memoryMarker].to)
		}
		if (this.stateDiff[coinbaseHex] !== undefined) {
			console.log('bf.b coinbase from\t', this.stateDiff[coinbaseHex].balance[memoryMarker].from)
			console.log('bf.b coinbase to\t', this.stateDiff[coinbaseHex].balance[memoryMarker].to)
		}
		// END DEBUGGING GAS




		// A transaction with value set, while from account has not enough balance to pay even for the gas cost,
		// will not be run at all, though will not change the state
		if (!hasFromSufficientBalanceForGasCost && !hasFromSufficientBalanceForValueAndGasCost && ctx.value.isPositive()) {
			return {};
		}

		if (this.stateDiff[fromAccountHex] !== undefined) {
			var fromAcc = this.stateDiff[fromAccountHex];
			var fromAccB = fromAcc.balance[memoryMarker];

			// Add back call value, gasCost and refunds to the start balance,
			// as it has been transfered before the CaptureStart and the interpreter execution
			var fromBal = fromAccB.from;
			if (hasFromSufficientBalanceForGasCost) {
				// TODO: check for GetEIP161abcTransition too?
				if (!this.hasInitCalled
						|| isCallTypeOnNonExistingAccount) {
					fromBal = fromBal.add(feesValue);
				} else {
					console.log('from 2')
					fromBal = fromBal.add(fullGasCost);
				}
				}
			}
			if (hasFromSufficientBalanceForValueAndGasCost && ctx.value.isPositive()) {
				fromBal = fromBal.add(ctx.value);
			}
			fromAccB.from = fromBal;
			var toBal = fromAccB.to;

			// In case account doesn't have enough balance, the transfer won't happen,
			// though the gasCost will still have to be paid
			if (!hasFromSufficientBalanceForGasCost) {
				// Remove gasCost paid for the transaction as it will happen after the interpreter execution
				// var toBal = bigInt(fromAccB.to.slice(2), 16);
				fromAccB.to = fromBal.subtract(feesValue);
			}

			// Decrement the caller's nonce,
			// as it has been increased before the CaptureStart and the interpreter execution
			var fromNonce = fromAcc.nonce[memoryMarker].from;
			fromAcc.nonce[memoryMarker].from = fromNonce - 1;

			// var toNonce = fromAcc.nonce[memoryMarker].to;
			// fromAcc.nonce[memoryMarker].to = "0x" + (toNonce + 1).toString(16);

			// remove any errors marked on the from account, as it has to be included on output
			// happens on mordor tx: 0x8f26c1acfce0178a2b037d85feeea99bb961bb46f541ad8c01c6668455952221
			fromAcc._error = false;

			fromAcc._final = true;
		}

		if (this.stateDiff[toAccountHex] !== undefined) {
			var toAcc = this.stateDiff[toAccountHex];

			var fromBal = toAcc.balance[memoryMarker].from;

			if (hasFromSufficientBalanceForValueAndGasCost) {
				// Remove transferred value, as it has been transfered before the CaptureStart and the interpreter execution
				toAcc.balance[memoryMarker].from = fromBal.subtract(ctx.value);
			} else {
				// In case from account doesn't have enough balance, the transfer won't happen
				toAcc.balance[memoryMarker].to = toAcc.balance[memoryMarker].from;
			}

			// var toBal = bigInt(toAcc.balance[memoryMarker].to.slice(2), 16);
			// // Remove transferred value, as it has been transfered before the CaptureStart and the interpreter execution
			// toAcc.balance[memoryMarker].to   = "0x" + toBal.subtract(ctx.value).toString(16);

			toAcc._final = true;

			// Mark new created contracts address
			if (isCreateType) {
				if (toAcc._type !== this.diffMarkers.Died) {
					toAcc._type = this.diffMarkers.Born;
				} else {
					// if the new contract has not be persisted to state then remove it from state_diff output
					toAcc._removed = true;
				}
			} else {
				// remove any errors marked on the to account (if not type = CREATE*), as it has to be included on output
				// happens on mordor tx: 0x89fd95d97374ccb9cdac249c74efdc57907c53beecb3e6ebce03b4ca31b0df2f
				toAcc._error = false;
			}
		}

		if (this.stateDiff[coinbaseHex] !== undefined) {
			var coinbaseAcc = this.stateDiff[coinbaseHex];
			// var blockCoinbaseReward = ctx.blockCoinbaseReward ? bigInt(ctx.blockCoinbaseReward) : 0;

			// FIXME: check if to or from is used
			// var fromBal = coinbaseAcc.balance[memoryMarker].to;
			var fromBal = coinbaseAcc.balance[memoryMarker].from;

			// Remove gas cost as it will happen after the interpreter execution
			if (hasFromSufficientBalanceForGasCost) {
				// TODO: check for GetEIP161abcTransition too?
				if (!this.hasInitCalled
						|| isCallTypeOnNonExistingAccount) {
					fromBal = fromBal.subtract(feesValue);
				} else if (fromAccountHex === coinbaseHex) {
          console.log('coinbase 2')
					// fromBal = fromBal.add(refundValue);
				}
				}
			}

			coinbaseAcc.balance[memoryMarker].from = fromBal;

			// Remove block reward, as state diff tracer applies to txs and not to blocks

			// Add block reward from previous block
			// var fromBal = bigInt(coinbaseAcc.balance[memoryMarker].from.slice(2), 16);
			// coinbaseAcc.balance[memoryMarker].from = fromBal.add(blockCoinbaseReward);


			// Remove block reward, as state diff tracer applies to txs and not to blocks

			// Add block reward from previous block
			// coinbaseToBal = coinbaseToBal.add(blockCoinbaseReward);

			// // Add gas cost from current block as it will happen after the interpreter execution
			// if (fromAccountHex !== coinbaseHex) {
			// 	coinbaseToBal = coinbaseToBal.add(gasCost);
			// }
			// remove any errors marked on the coinbase account, as it has to be included on output
			// happens on mordor tx: 0xbfca41d82781ba1888c10d96de84ff68799e328c658b34964d382eba019b3752
			coinbaseAcc._error = false;

			coinbaseAcc._final = true;
		}

		if (this.stateDiff[fromAccountHex] !== undefined) {
			console.log('\naf.b fromA from\t\t', this.stateDiff[fromAccountHex].balance[memoryMarker].from,
				this.stateDiff[fromAccountHex].balance[memoryMarker].from.subtract(targetFrom))
			console.log('af.b fromA to\t\t', this.stateDiff[fromAccountHex].balance[memoryMarker].to)
		}
		if (this.stateDiff[toAccountHex] !== undefined) {
			console.log('af.b toA from\t\t', this.stateDiff[toAccountHex].balance[memoryMarker].from)
			console.log('af.b toA to\t\t', this.stateDiff[toAccountHex].balance[memoryMarker].to)
		}
		if (this.stateDiff[coinbaseHex] !== undefined) {
			console.log('af.b coinbase from\t', this.stateDiff[coinbaseHex].balance[memoryMarker].from, this.stateDiff
				[coinbaseHex].balance[memoryMarker].from.subtract(targetCoinbase))
			console.log('af.b coinbase to\t', this.stateDiff[coinbaseHex].balance[memoryMarker].to)
		}

		this.format(db);

		// Return the assembled allocations (stateDiff)
		return this.stateDiff;
	}
}
