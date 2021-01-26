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
	stateDiff: null,

	lastGasIn: null,

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
	lookupAccount: function(addr, db, type){
		type = type || this.diffMarkers.Changed;

		var memoryMarker = this.diffMarkers.Memory;

		var acc = toHex(addr);
		// if (acc == '0xd6758d1907ed647605429d40cd19c58a6d05eb8b') {
		// 	var ba = db.getBalance(toAddress('0x893defcfe8dc3b7fe2b95c2ddd6415ac2f1eb582'));
		// 	console.log('miner', ba, "0x" + ba.toString(16))
		// }
		var balance = "0x" + db.getBalance(addr).toString(16);
    console.log('acc', acc, balance)
		var code = toHex(db.getCode(addr));
		var nonce = db.getNonce(addr);

		if (this.stateDiff[acc] === undefined) {
			// if (nonce == 1 && code != "0x") {
			// 	type = this.diffMarkers.Born;
			// }


			this.stateDiff[acc] = {
				type: type,
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

		// console.log('nonce', type, acc, accountData.nonce, '=>',nonce)


		// re-read type from stateDiff
		// accountData.type = type;
		if (balance) {
			accountData.balance[memoryMarker].to = balance;
		}

		// var latestKnownNonce = accountData.nonce[memoryMarker].to || accountData.nonce[memoryMarker].from
		// if (nonce && type === this.diffMarkers.Born) {
		// 	accountData.nonce[memoryMarker].to = latestKnownNonce + 1;
		if (nonce) {
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

		var memoryMarker = this.diffMarkers.Memory;

		// if (this.stateDiff[acc] !== undefined) {
		// 	return;
		// }

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

		if (data.balance === sameMarker &&
				data.nonce === sameMarker &&
				data.code === sameMarker &&
				this.isObjectEmpty(data.storage)
			) {
				return false;
			}
			return true;
	},

	format: function(db) {
		for (var acc in this.stateDiff) {
			// Fetch latest balance
			// TODO: optimise
			this.lookupAccount(toAddress(acc), db);

			var accountData = this.stateDiff[acc];
			var type = accountData.type;
      console.log('149 \t type', type)
			delete accountData.type;

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

	// result is invoked when all the opcodes have been iterated over and returns
	// the final result of the tracing.
	result: function(ctx, db) {

		console.log('----- RESULT -------');
		// 0x893defcfe8dc3b7fe2b95c2ddd6415ac2f1eb582

		var memoryMarker = this.diffMarkers.Memory;


		this.lookupAccount(toAddress(ctx.coinbase), db);

		var gasCost = (ctx.gasLimit - this.lastGasIn) * ctx.gasPrice;
    console.log('255 \t gasCost', gasCost)


		var coinbaseHex = toHex(ctx.coinbase);
		var coinbaseFromBal = bigInt(this.stateDiff[coinbaseHex].balance[memoryMarker].from.slice(2), 16);
		this.stateDiff[coinbaseHex].balance[memoryMarker].from = "0x" + coinbaseFromBal.subtract(gasCost).toString(16);



		// At this point, we need to deduct the "value" from the
		// outer transaction, and move it back to the origin
		this.lookupAccount(toAddress(ctx.from), db);

		var fromAccountHex = toHex(ctx.from);
		var toAccountHex = toHex(ctx.to);

		// var fromAccountData = this.stateDiff[fromAccountHex] || {};
		// var toAccountData = this.stateDiff[toAccountHex] || {};


		// var gas = ctx.gas;
    // console.log('262 \t gas', gas)
		// var gasUsed = ctx.gasUsed;
    // console.log('264 \t gasUsed', gasUsed)

		// var gasSet = 250000;
		// var gasPrice = 20;

		console.log('ctx', ctx);

		console.log('0xb62da4cfbd8a3206e, 0xb62d82756f99d186e', bigInt(0xb62da4cfbd8a3206e).subtract(bigInt(0xb62d82756f99d186e)))

		// console.log('cost', gas * gasUsed)
		// console.log('cost g', gas * gasPrice)
		// console.log('cost gu', gasUsed * gasPrice)
		// console.log('cost gs', gasSet * gasPrice)
		// console.log('cost gs 137,785', 137785 * gasPrice)
		// console.log('cost tgu', (gas * gasPrice) + (gasUsed * gasPrice))

		// console.log('refund', this.logger.refund);

		// console.log('this.logger', JSON.stringify(this.logger))


		var fromBal = bigInt(this.stateDiff[fromAccountHex].balance[memoryMarker].from.slice(2), 16);
		this.stateDiff[fromAccountHex].balance[memoryMarker].from = "0x" + fromBal.add(ctx.value).add(gasCost).toString(16);
    console.log('305 \t fromBal', fromBal, ctx.value, gasCost, fromBal.add(ctx.value), fromBal.add(ctx.value).add(gasCost))

		console.log('261 \t this.stateDiff[fromAccountHex].balance[memoryMarker].from', fromAccountHex, this.stateDiff[fromAccountHex].balance[memoryMarker].from)

		var toBal   = bigInt(this.stateDiff[toAccountHex].balance[memoryMarker].from.slice(2), 16);
		this.stateDiff[toAccountHex].balance[memoryMarker].from   = "0x" + toBal.subtract(ctx.value).toString(16);
    console.log('265 \t this.stateDiff[toAccountHex].balance[memoryMarker].from', toAccountHex, this.stateDiff[toAccountHex].balance[memoryMarker].from)

		// Decrement the caller's nonce, and remove empty create targets
		var toNonce = this.stateDiff[fromAccountHex].nonce[memoryMarker].from;
		this.stateDiff[fromAccountHex].nonce[memoryMarker].from = "0x" + (toNonce - 1).toString(16);
		this.stateDiff[fromAccountHex].nonce[memoryMarker].to = "0x" + toNonce.toString(16);

		// Mark new contracts address as new/born ones
		if (ctx.type == "CREATE" || ctx.type == "CREATE2") {
			this.stateDiff[toAccountHex].type = this.diffMarkers.Born;
		}

		this.format(db);

		// Return the assembled allocations (stateDiff)
		return this.stateDiff;
	},

	// step is invoked for every opcode that the VM executes.
	step: function(log, db) {
		// Add the current account if we just started tracing
		if (this.stateDiff === null){
			this.stateDiff = {};

			// var contractAddr = log.contract.getAddress();
			// console.log('ex',db.exists(contractAddr));
			// var type = db.exists(contractAddr) ? this.diffMarkers.Changed : this.diffMarkers.Born;
			// Balance will potentially be wrong here, since this will include the value
			// sent along with the message. We fix that in "result()".
			this.lookupAccount(log.contract.getAddress(), db);

			// this.lookupAccount(toAddress('0x893defcfe8dc3b7fe2b95c2ddd6415ac2f1eb582'), db);
		}

		// var refund = log.getRefund();
    // console.log('refund', refund)
		// this.logger.refund = refund

		this.lastGasIn = log.getGas();
		// var loga = {
		// 	gasAvailable: log.getAvailableGas(),
		// 	gasIn:   log.getGas(),
		// 	gasCost: log.getCost(),
		// }
		// console.log('gasIn', loga.gasIn)
		// if (loga.gasAvailable > 0) {
		// 	console.log('AMAZINIIG', loga)
		// }
		// this.logger.push(loga)
		// console.log(log.op.toString(), log.stack.peek(0).toString(16), log.stack.peek(1).toString(16), log.stack.peek(2).toString(16), log.stack.peek(3).toString(16), log.stack.peek(4).toString(16));
		console.log(log.op.toString());

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
				// if (isPrecompiled(address)) {
				// 	break;
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
		console.log('fault', log.op.toString(), toHex(log.contract.getAddress()), log.stack.peek(0).toString(16), log.stack.peek(1).toString(16), log.stack.peek(2).toString(16), log.stack.peek(3).toString(16), log.stack.peek(4).toString(16));

		var error = log.getError();
    console.log('TCL: \t file: state_diff_tracer.js \t line 325 \t error', error)
		var opError = log.getCallError();
    console.log('TCL: \t file: state_diff_tracer.js \t line 327 \t opError', opError)
	}
}
