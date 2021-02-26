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

	debugState: {},

	lastUsedContractAddress: null,
	lastAccessedAccount: null,
	lastGasIn: null,
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
		type = type || this.diffMarkers.Changed;

		var memoryMarker = this.diffMarkers.Memory;

		var acc = toHex(addr);
    // console.log('TCL: \t file: state_diff_tracer.js \t line 50 \t acc', acc)

		this.lastAccessedAccount = acc;

		// if (acc == '0xd6758d1907ed647605429d40cd19c58a6d05eb8b') {
		// 	var ba = db.getBalance(toAddress('0x893defcfe8dc3b7fe2b95c2ddd6415ac2f1eb582'));
		// 	console.log('miner', ba, "0x" + ba.toString(16))
		// }
		var balance = "0x" + db.getBalance(addr).toString(16);
    // console.log('acc', acc, balance)
		var code = toHex(db.getCode(addr));
		var nonce = db.getNonce(addr);

		if (this.stateDiff[acc] === undefined) {
			// if (nonce == 1 && code != "0x") {
			// 	type = this.diffMarkers.Born;
			// }

			this.stateDiff[acc] = {
				_type: type,  // temp storage of account's initial type
				_error: false, // evm returned an error
				_final: false, // stop updating state if account's state marked as final
				balance: {
					[memoryMarker]: {
						"from": balance,
					}
				},
				nonce: {
					[memoryMarker]: {
						"from": nonce,
					}
				},
				code: {
					[memoryMarker]: {
						"from": code,
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


		// re-read type from stateDiff
		// accountData.type = type;
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
			accountData.storage[idx] = {
				[memoryMarker]: {
					"from": toHex(db.getState(addr, key))
				}
			};
		}

		if (val) {
			accountData.storage[idx][memoryMarker].to = toHex(val);
		}
	},

	formatSingle: function(data, type) {
		type = type || this.diffMarkers.Changed;

		var memoryMarker = this.diffMarkers.Memory;
		var val = data[memoryMarker].to || data[memoryMarker].from;

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

		if (to === undefined ||
			from === to) {
			return sameMarker;
		}

		return {
			[changedMarker]: {
				from: this.toHexJs(from),
				to: this.toHexJs(to),
			}
		}
	},

	hasAccountChanges: function(data) {
		var sameMarker = this.diffMarkers.Same;
		var bornMarker = this.diffMarkers.Born;

		if (data.balance === sameMarker &&
				data.nonce === sameMarker &&
				data.code === sameMarker &&
				this.isObjectEmpty(data.storage)
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
			// fetch latest balance
			// TODO: optimise
			this.lookupAccount(toAddress(acc), db);

			// has been cleared in lookupAccount
			if (this.stateDiff[acc] === undefined) {
				continue;
			}

			var accountData = this.stateDiff[acc];

			// remove accounts with errors from output. Do this check last (not in another method),
			// otherwise account might be re-added if one of ctx.[from|to|coinbase] from within this.result()
			if (accountData._error) {
				delete this.stateDiff[acc];
				continue;
			}

			var type = accountData._type;
			delete accountData._type;
			delete accountData._error;
			delete accountData._final;

			var memoryMarker = this.diffMarkers.Memory;

			var changedMarker = this.diffMarkers.Changed;
			var sameMarker = this.diffMarkers.Same;

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
				var sti = accountData.storage[idx];
				if (sti[memoryMarker] === undefined ||
						sti[memoryMarker].to === undefined ||
						/^(0x)?0*$/.test(sti[memoryMarker].to)) {
					delete this.stateDiff[acc].storage[idx];
				// } else if (accountData.storage[idx][type] === undefined ||
				// 		/^(0x)?0*$/.test(accountData.storage[idx][type])) {
				// 	delete this.stateDiff[acc].storage[idx];
					continue;
				}

				if (type === changedMarker) {
					var res = this.formatChanged(sti, type);
					if (res === sameMarker) {
						delete this.stateDiff[acc].storage[idx];
					} else {
						accountData.storage[idx] = res;
					}
				} else {
					accountData.storage[idx] = this.formatSingle(sti, type);
				}
			}

			// remove unchanged accounts
			if (!this.hasAccountChanges(accountData)) {
				delete this.stateDiff[acc];
				continue;
			}
		}
	},

	// init is invoked on the first call VM executes.
	init: function(ctx, log, db) {
		// Balances will potentially be wrong here, since they will include the value
		// sent along with the message or the paid full gas cost. We fix that in "result()".

		var toAccAddress = toAddress(ctx.to);

		// Get actual "to" values for from, to, coinbase accounts
		this.lookupAccount(toAddress(ctx.from), db);
		this.lookupAccount(toAccAddress, db);
		this.lookupAccount(toAddress(ctx.coinbase), db);

		var contractAddress = log.contract.getAddress();
		if (toHex(contractAddress) !== toHex(toAccAddress)) {
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
		if (error !== undefined &&
				this.lastAccessedAccount !== null) {
        // console.log('damn line 318 \t this.lastAccessedAccount',log.op.toString(), this.lastAccessedAccount,error)

			this.stateDiff[this.lastAccessedAccount]._error = true;  // mark account that had an error
			this.lastAccessedAccount = null;
			return;
		}

		// Add the current account if we just started tracing

		// var refund = log.getRefund();
    // console.log('refund', refund)
		// this.logger.refund = refund

		this.lastGasIn = log.getGas();
		this.lastRefund = log.getRefund();
    // console.log('this.lastGasIn', 106376 - this.lastGasIn, 106376-log.getAvailableGas())
		// var loga = {
		// 	gasAvailableGas: log.getAvailableGas(),
		// 	gasGas:   log.getGas(),
		// 	gasCost: log.getCost(),
		// 	getValue: log.contract.getValue(),
		// }
		// console.log('loga', loga)a


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

		if (this.hasError && this.lastAccessedAccount !== null) {
			this.stateDiff[this.lastAccessedAccount]._error = true;  // mark account that had an error
			this.lastAccessedAccount = null;
			return;
		}
	},

	// result is invoked when all the opcodes have been iterated over and returns
	// the final result of the tracing
	result: function(ctx, db) {

		// Reset lastAccessedAccount cleanup logic as it is being used only for `step` method
		// NOTE: it's safe to be removed, as it's not being utilised from this point
		// this.lastAccessedAccount = null;

		// console.log('----- RESULT -------');
		// 0x893defcfe8dc3b7fe2b95c2ddd6415ac2f1eb582

		var memoryMarker = this.diffMarkers.Memory;



		// EIP161, when calling a non existing account, passing no value,
		// then nothing happens and `CaptureState` (and  inline `step`) method are not being called,
		// so no calculations or logic is being applied in the tracer.
		// For this reason we initiate the state here.
		if (this.lastGasIn === null) {
			this.lastGasIn = ctx.gas;
		}
		this.lookupAccount(toAddress(ctx.from), db);
		this.lookupAccount(toAddress(ctx.to), db);
		this.lookupAccount(toAddress(ctx.coinbase), db);

		var fromAccountHex = toHex(ctx.from);
		var toAccountHex = toHex(ctx.to);
		var coinbaseHex = toHex(ctx.coinbase);
		var gasCost = bigInt(ctx.gasLimit)
										.subtract(bigInt(ctx.gas))
										.add(bigInt(ctx.gasUsed))
										.subtract(refund)
										.multiply(bigInt(ctx.gasPrice));




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

		console.log('\nCalcs ----');
		console.log('\nvalue\t\t', ctx.value)
		console.log('\ncanTransferBalanceFrom\t\t', ctx.canTransferBalanceFrom)

		console.log('\ngasCost\t\t', gasCost)
		console.log('lastGasIn\t', this.lastGasIn, '\t\t-gasUsed\t', this.lastGasIn - ctx.gasUsed)

		console.log('\ngas\t\t', ctx.gas)
		console.log('gasLimit\t', ctx.gasLimit, '\t-gasUsed\t', ctx.gasLimit - ctx.gasUsed)
		console.log('gasLimit\t', ctx.gasLimit, '\t-gas\t\t', ctx.gasLimit - ctx.gas)
		console.log('gasUsed\t\t', ctx.gasUsed)
		console.log('gasPrice\t', ctx.gasPrice)

		console.log('\nbf.b fromA from\t\t', bigInt(this.stateDiff[fromAccountHex].balance[memoryMarker].from.slice(2), 16))
		console.log('bf.b fromA to\t\t', bigInt(this.stateDiff[fromAccountHex].balance[memoryMarker].to.slice(2), 16))
		console.log('bf.b toA from\t\t', bigInt(this.stateDiff[toAccountHex].balance[memoryMarker].from.slice(2), 16))
		console.log('bf.b toA to\t\t', bigInt(this.stateDiff[toAccountHex].balance[memoryMarker].to.slice(2), 16))
		console.log('bf.b coinbase from\t', bigInt(this.stateDiff[coinbaseHex].balance[memoryMarker].from.slice(2), 16))
		console.log('bf.b coinbase to\t', bigInt(this.stateDiff[coinbaseHex].balance[memoryMarker].to.slice(2), 16))
		// END DEBUGGING GAS

		// var fromAccountData = this.stateDiff[fromAccountHex] || {};
		// var toAccountData = this.stateDiff[toAccountHex] || {};
		// In case from balance is negative because the tracer has disabled the CanTransfer check,
		// and the Transfer happened before the CaptureStart and the interpreter execution
		var canTransferBalance = ctx.canTransferBalanceFrom || false;

		if (this.stateDiff[fromAccountHex] !== undefined) {
			var fromAcc = this.stateDiff[fromAccountHex];
			var fromAccB = fromAcc.balance[memoryMarker];

			// Add back call value, gasCost and refunds to the start balance,
			// as it has been transfered before the CaptureStart and the interpreter execution
			var fromBal = bigInt(fromAccB.from.slice(2), 16);
			if (fromAccountHex !== coinbaseHex) {
				fromBal = fromBal.add(gasCost);
			}
			if (!this.hasError) {
				fromBal = fromBal.add(ctx.value);
			}
			fromAccB.from = "0x" + fromBal.toString(16);

			// In case account doesn't have enough balance, the transfer won't happen,
			// though the gasCost will still have to be paid
			if (!canTransferBalance) {
				// Remove gasCost paid for the transaction as it will happen after the interpreter execution
				// var toBal = bigInt(fromAccB.to.slice(2), 16);
				fromAccB.to = "0x" + fromBal.subtract(gasCost).toString(16);
			}

			// Decrement the caller's nonce,
			// as it has been increased before the CaptureStart and the interpreter execution
			var fromNonce = fromAcc.nonce[memoryMarker].from;
			fromAcc.nonce[memoryMarker].from = "0x" + (fromNonce - 1).toString(16);

			// var toNonce = fromAcc.nonce[memoryMarker].to;
			// fromAcc.nonce[memoryMarker].to = "0x" + (toNonce + 1).toString(16);

			fromAcc._final = true;
		}

		if (this.stateDiff[toAccountHex] !== undefined) {
			var toAcc = this.stateDiff[toAccountHex];

			var fromBal = bigInt(toAcc.balance[memoryMarker].from.slice(2), 16);

			// Remove transferred value, as it has been transfered before the CaptureStart and the interpreter execution
			toAcc.balance[memoryMarker].from   = "0x" + fromBal.subtract(ctx.value).toString(16);

			// In case from account doesn't have enough balance, the transfer won't happen
			if (!canTransferBalance) {
				toAcc.balance[memoryMarker].to = toAcc.balance[memoryMarker].from;
			}

			// var toBal = bigInt(toAcc.balance[memoryMarker].to.slice(2), 16);
			// // Remove transferred value, as it has been transfered before the CaptureStart and the interpreter execution
			// toAcc.balance[memoryMarker].to   = "0x" + toBal.subtract(ctx.value).toString(16);

			toAcc._final = true;

			// Mark new contracts address as new/born ones
			if (ctx.type == "CREATE" || ctx.type == "CREATE2") {
				toAcc._type = this.diffMarkers.Born;
			}
		}

		if (this.stateDiff[coinbaseHex] !== undefined) {
			var coinbaseAcc = this.stateDiff[coinbaseHex];
			var blockCoinbaseReward = ctx.blockCoinbaseReward ? bigInt(ctx.blockCoinbaseReward) : 0;

			var fromBal = bigInt(coinbaseAcc.balance[memoryMarker].from.slice(2), 16);

			// Remove gas cost as it will happen after the interpreter execution
			if (fromAccountHex !== coinbaseHex) {
				fromBal = fromBal.subtract(gasCost);
			}

			// Remove block reward, as state diff tracer applies to txs and not to blocks
			coinbaseAcc.balance[memoryMarker].from = "0x" + fromBal.subtract(blockCoinbaseReward).toString(16);

			// Add block reward from previous block
			// var fromBal = bigInt(coinbaseAcc.balance[memoryMarker].from.slice(2), 16);
			// coinbaseAcc.balance[memoryMarker].from = "0x" + fromBal.add(blockCoinbaseReward).toString(16);


			var coinbaseToBal = bigInt(coinbaseAcc.balance[memoryMarker].to.slice(2), 16);

			// Remove block reward, as state diff tracer applies to txs and not to blocks
			coinbaseToBal = coinbaseToBal.subtract(blockCoinbaseReward);

			// Add block reward from previous block
			// coinbaseToBal = coinbaseToBal.add(blockCoinbaseReward);

			// // Add gas cost from current block as it will happen after the interpreter execution
			// if (fromAccountHex !== coinbaseHex) {
			// 	coinbaseToBal = coinbaseToBal.add(gasCost);
			// }
			coinbaseAcc.balance[memoryMarker].to = "0x" + coinbaseToBal.toString(16);

			coinbaseAcc._final = true;
		}

		console.log('\naf.b fromA from\t\t', bigInt(this.stateDiff[fromAccountHex].balance[memoryMarker].from.slice(2), 16))
		console.log('af.b fromA to\t\t', bigInt(this.stateDiff[fromAccountHex].balance[memoryMarker].to.slice(2), 16))
		console.log('af.b toA from\t\t', bigInt(this.stateDiff[toAccountHex].balance[memoryMarker].from.slice(2), 16))
		console.log('af.b toA to\t\t', bigInt(this.stateDiff[toAccountHex].balance[memoryMarker].to.slice(2), 16))
		console.log('af.b coinbase from\t', bigInt(this.stateDiff[coinbaseHex].balance[memoryMarker].from.slice(2), 16))
		console.log('af.b coinbase to\t', bigInt(this.stateDiff[coinbaseHex].balance[memoryMarker].to.slice(2), 16))

		this.format(db);

		// Return the assembled allocations (stateDiff)
		return this.stateDiff;
	}
}
