






| Entity | Version |
| --- | --- |
| Source | <code>1.12.14-unstable/generated-at:2023-09-04T08:02:34-06:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### personal_deriveAccount

DeriveAccount requests an HD wallet to derive a new account, optionally pinning
it for later reuse.


#### Params (3)

Parameters must be given _by position_.


__1:__ 
url <code>string</code> 

  + Required: ✓ Yes





__2:__ 
path <code>string</code> 

  + Required: ✓ Yes





__3:__ 
pin <code>*bool</code> 

  + Required: ✓ Yes






#### Result




<code>accounts.Account</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- address: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- url: 
			- additionalProperties: `false`
			- properties: 
				- Path: 
					- type: `string`

				- Scheme: 
					- type: `string`


			- type: `object`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "address": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "url": {
                "additionalProperties": false,
                "properties": {
                    "Path": {
                        "type": "string"
                    },
                    "Scheme": {
                        "type": "string"
                    }
                },
                "type": "object"
            }
        },
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_deriveAccount", "params": [<url>, <path>, <pin>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "personal_deriveAccount", "params": [<url>, <path>, <pin>]}'
	```


=== "Javascript Console"

	``` js
	personal.deriveAccount(url,path,pin);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *PersonalAccountAPI) DeriveAccount(url string, path string, pin *bool) (accounts.Account, error) {
	wallet, err := s.am.Wallet(url)
	if err != nil {
		return accounts.Account{}, err
	}
	derivPath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return accounts.Account{}, err
	}
	if pin == nil {
		pin = new(bool)
	}
	return wallet.Derive(derivPath, *pin)
}// DeriveAccount requests an HD wallet to derive a new account, optionally pinning
// it for later reuse.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L344" target="_">View on GitHub →</a>
</p>
</details>

---



### personal_ecRecover

EcRecover returns the address for the account that was used to create the signature.
Note, this function is compatible with eth_sign and personal_sign. As such it recovers
the address of:
hash = keccak256("\x19Ethereum Signed Message:\n"${message length}${message})
addr = ecrecover(hash, signature)

Note, the signature must conform to the secp256k1 curve R, S and V values, where
the V value must be 27 or 28 for legacy reasons.

https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_ecRecover


#### Params (2)

Parameters must be given _by position_.


__1:__ 
data <code>hexutil.Bytes</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `dataWord`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of some bytes",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "dataWord",
        "type": [
            "string"
        ]
    }
	```




__2:__ 
sig <code>hexutil.Bytes</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `dataWord`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of some bytes",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "dataWord",
        "type": [
            "string"
        ]
    }
	```





#### Result




<code>common.Address</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a Keccak 256 hash POINTER`
	- pattern: `^0x[a-fA-F\d]{64}$`
	- title: `keccak`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a Keccak 256 hash POINTER",
        "pattern": "^0x[a-fA-F\\d]{64}$",
        "title": "keccak",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_ecRecover", "params": [<data>, <sig>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "personal_ecRecover", "params": [<data>, <sig>]}'
	```


=== "Javascript Console"

	``` js
	personal.ecRecover(data,sig);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *PersonalAccountAPI) EcRecover(ctx context.Context, data, sig hexutil.Bytes) (common.Address, error) {
	if len(sig) != crypto.SignatureLength {
		return common.Address{}, fmt.Errorf("signature must be %d bytes long", crypto.SignatureLength)
	}
	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return common.Address{}, errors.New("invalid Ethereum signature (V is not 27 or 28)")
	}
	sig[crypto.RecoveryIDOffset] -= 27
	rpk, err := crypto.SigToPub(accounts.TextHash(data), sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*rpk), nil
}// EcRecover returns the address for the account that was used to create the signature.
// Note, this function is compatible with eth_sign and personal_sign. As such it recovers
// the address of:
// hash = keccak256("\x19Ethereum Signed Message:\n"${message length}${message})
// addr = ecrecover(hash, signature)
//
// Note, the signature must conform to the secp256k1 curve R, S and V values, where
// the V value must be 27 or 28 for legacy reasons.
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_ecRecover

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L549" target="_">View on GitHub →</a>
</p>
</details>

---



### personal_importRawKey

ImportRawKey stores the given hex encoded ECDSA key into the key directory,
encrypting it with the passphrase.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
privkey <code>string</code> 

  + Required: ✓ Yes





__2:__ 
password <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>common.Address</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a Keccak 256 hash POINTER`
	- pattern: `^0x[a-fA-F\d]{64}$`
	- title: `keccak`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a Keccak 256 hash POINTER",
        "pattern": "^0x[a-fA-F\\d]{64}$",
        "title": "keccak",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_importRawKey", "params": [<privkey>, <password>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "personal_importRawKey", "params": [<privkey>, <password>]}'
	```


=== "Javascript Console"

	``` js
	personal.importRawKey(privkey,password);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *PersonalAccountAPI) ImportRawKey(privkey string, password string) (common.Address, error) {
	key, err := crypto.HexToECDSA(privkey)
	if err != nil {
		return common.Address{}, err
	}
	ks, err := fetchKeystore(s.am)
	if err != nil {
		return common.Address{}, err
	}
	acc, err := ks.ImportECDSA(key, password)
	return acc.Address, err
}// ImportRawKey stores the given hex encoded ECDSA key into the key directory,
// encrypting it with the passphrase.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L386" target="_">View on GitHub →</a>
</p>
</details>

---



### personal_initializeWallet

InitializeWallet initializes a new wallet at the provided URL, by generating and returning a new private key.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
url <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>string</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_initializeWallet", "params": [<url>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "personal_initializeWallet", "params": [<url>]}'
	```


=== "Javascript Console"

	``` js
	personal.initializeWallet(url);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *PersonalAccountAPI) InitializeWallet(ctx context.Context, url string) (string, error) {
	wallet, err := s.am.Wallet(url)
	if err != nil {
		return "", err
	}
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	seed := bip39.NewSeed(mnemonic, "")
	switch wallet := wallet.( // InitializeWallet initializes a new wallet at the provided URL, by generating and returning a new private key.
	type) {
	case *scwallet.Wallet:
		return mnemonic, wallet.Initialize(seed)
	default:
		return "", errors.New("specified wallet does not support initialization")
	}
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L566" target="_">View on GitHub →</a>
</p>
</details>

---



### personal_listAccounts

ListAccounts will return a list of addresses for accounts this node manages.


#### Params (0)

_None_

#### Result



commonAddress <code>[]common.Address</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- items: 

			- description: `Hex representation of a Keccak 256 hash POINTER`
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: string


	- type: array


	```

=== "Raw"

	``` Raw
	{
        "items": [
            {
                "description": "Hex representation of a Keccak 256 hash POINTER",
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": [
                    "string"
                ]
            }
        ],
        "type": [
            "array"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_listAccounts", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "personal_listAccounts", "params": []}'
	```


=== "Javascript Console"

	``` js
	personal.listAccounts();
	```



<details><summary>Source code</summary>
<p>
```go
func (s *PersonalAccountAPI) ListAccounts() [ // ListAccounts will return a list of addresses for accounts this node manages.
]common.Address {
	return s.am.Accounts()
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L294" target="_">View on GitHub →</a>
</p>
</details>

---



### personal_listWallets

ListWallets will return a list of wallets this node manages.


#### Params (0)

_None_

#### Result



rawWallet <code>[]rawWallet</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- items: 

			- additionalProperties: `false`
			- properties: 
				- accounts: 
					- items: 
						- additionalProperties: `false`
						- properties: 
							- address: 
								- pattern: `^0x[a-fA-F\d]{64}$`
								- title: `keccak`
								- type: `string`

							- url: 
								- additionalProperties: `false`
								- properties: 
									- Path: 
										- type: `string`

									- Scheme: 
										- type: `string`


								- type: `object`


						- type: `object`

					- type: `array`

				- failure: 
					- type: `string`

				- status: 
					- type: `string`

				- url: 
					- type: `string`


			- type: object


	- type: array


	```

=== "Raw"

	``` Raw
	{
        "items": [
            {
                "additionalProperties": false,
                "properties": {
                    "accounts": {
                        "items": {
                            "additionalProperties": false,
                            "properties": {
                                "address": {
                                    "pattern": "^0x[a-fA-F\\d]{64}$",
                                    "title": "keccak",
                                    "type": "string"
                                },
                                "url": {
                                    "additionalProperties": false,
                                    "properties": {
                                        "Path": {
                                            "type": "string"
                                        },
                                        "Scheme": {
                                            "type": "string"
                                        }
                                    },
                                    "type": "object"
                                }
                            },
                            "type": "object"
                        },
                        "type": "array"
                    },
                    "failure": {
                        "type": "string"
                    },
                    "status": {
                        "type": "string"
                    },
                    "url": {
                        "type": "string"
                    }
                },
                "type": [
                    "object"
                ]
            }
        ],
        "type": [
            "array"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_listWallets", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "personal_listWallets", "params": []}'
	```


=== "Javascript Console"

	``` js
	personal.listWallets();
	```



<details><summary>Source code</summary>
<p>
```go
func (s *PersonalAccountAPI) ListWallets() [ // ListWallets will return a list of wallets this node manages.
]rawWallet {
	wallets := make([]rawWallet, 0)
	for _, wallet := range s.am.Wallets() {
		status, failure := wallet.Status()
		raw := rawWallet{URL: wallet.URL().String(), Status: status, Accounts: wallet.Accounts()}
		if failure != nil {
			raw.Failure = failure.Error()
		}
		wallets = append(wallets, raw)
	}
	return wallets
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L308" target="_">View on GitHub →</a>
</p>
</details>

---



### personal_lockAccount

LockAccount will lock the account associated with the given address when it's unlocked.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
addr <code>common.Address</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a Keccak 256 hash POINTER`
	- pattern: `^0x[a-fA-F\d]{64}$`
	- title: `keccak`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a Keccak 256 hash POINTER",
        "pattern": "^0x[a-fA-F\\d]{64}$",
        "title": "keccak",
        "type": [
            "string"
        ]
    }
	```





#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_lockAccount", "params": [<addr>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "personal_lockAccount", "params": [<addr>]}'
	```


=== "Javascript Console"

	``` js
	personal.lockAccount(addr);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *PersonalAccountAPI) LockAccount(addr common.Address) bool {
	if ks, err := fetchKeystore(s.am); err == nil {
		return ks.Lock(addr) == nil
	}
	return false
}// LockAccount will lock the account associated with the given address when it's unlocked.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L431" target="_">View on GitHub →</a>
</p>
</details>

---



### personal_newAccount

NewAccount will create a new account and returns the address for the new account.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
password <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>common.AddressEIP55</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- items: 

			- description: `Hex representation of the integer`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: string


	- maxItems: `20`
	- minItems: `20`
	- type: array


	```

=== "Raw"

	``` Raw
	{
        "items": [
            {
                "description": "Hex representation of the integer",
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": [
                    "string"
                ]
            }
        ],
        "maxItems": 20,
        "minItems": 20,
        "type": [
            "array"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_newAccount", "params": [<password>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "personal_newAccount", "params": [<password>]}'
	```


=== "Javascript Console"

	``` js
	personal.newAccount(password);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *PersonalAccountAPI) NewAccount(password string) (common.AddressEIP55, error) {
	ks, err := fetchKeystore(s.am)
	if err != nil {
		return common.AddressEIP55{}, err
	}
	acc, err := ks.NewAccount(password)
	if err == nil {
		addrEIP55 := common.AddressEIP55(acc.Address)
		log.Info("Your new key was generated", "address", addrEIP55.String())
		log.Warn("Please backup your key file!", "path", acc.URL.Path)
		log.Warn("Please remember your password!")
		return addrEIP55, nil
	}
	return common.AddressEIP55{}, err
}// NewAccount will create a new account and returns the address for the new account.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L360" target="_">View on GitHub →</a>
</p>
</details>

---



### personal_openWallet

OpenWallet initiates a hardware wallet opening procedure, establishing a USB
connection and attempting to authenticate via the provided passphrase. Note,
the method may return an extra challenge requiring a second open (e.g. the
Trezor PIN matrix challenge).


#### Params (2)

Parameters must be given _by position_.


__1:__ 
url <code>string</code> 

  + Required: ✓ Yes





__2:__ 
passphrase <code>*string</code> 

  + Required: ✓ Yes






#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_openWallet", "params": [<url>, <passphrase>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "personal_openWallet", "params": [<url>, <passphrase>]}'
	```


=== "Javascript Console"

	``` js
	personal.openWallet(url,passphrase);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *PersonalAccountAPI) OpenWallet(url string, passphrase *string) error {
	wallet, err := s.am.Wallet(url)
	if err != nil {
		return err
	}
	pass := ""
	if passphrase != nil {
		pass = *passphrase
	}
	return wallet.Open(pass)
}// OpenWallet initiates a hardware wallet opening procedure, establishing a USB
// connection and attempting to authenticate via the provided passphrase. Note,
// the method may return an extra challenge requiring a second open (e.g. the
// Trezor PIN matrix challenge).

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L330" target="_">View on GitHub →</a>
</p>
</details>

---



### personal_sendTransaction

SendTransaction will create a transaction from the given arguments and
tries to sign it with the key associated with args.From. If the given
passwd isn't able to decrypt the key it fails.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
args <code>TransactionArgs</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- accessList: 
			- items: 
				- additionalProperties: `false`
				- properties: 
					- address: 
						- pattern: `^0x[a-fA-F\d]{64}$`
						- title: `keccak`
						- type: `string`

					- storageKeys: 
						- items: 
							- description: `Hex representation of a Keccak 256 hash`
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`

						- type: `array`


				- type: `object`

			- type: `array`

		- chainId: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- data: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- from: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- gas: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- gasPrice: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- input: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- maxFeePerGas: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- maxPriorityFeePerGas: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- nonce: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- to: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- value: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "accessList": {
                "items": {
                    "additionalProperties": false,
                    "properties": {
                        "address": {
                            "pattern": "^0x[a-fA-F\\d]{64}$",
                            "title": "keccak",
                            "type": "string"
                        },
                        "storageKeys": {
                            "items": {
                                "description": "Hex representation of a Keccak 256 hash",
                                "pattern": "^0x[a-fA-F\\d]{64}$",
                                "title": "keccak",
                                "type": "string"
                            },
                            "type": "array"
                        }
                    },
                    "type": "object"
                },
                "type": "array"
            },
            "chainId": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "data": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "dataWord",
                "type": "string"
            },
            "from": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "gas": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "gasPrice": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "input": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "dataWord",
                "type": "string"
            },
            "maxFeePerGas": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "maxPriorityFeePerGas": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "nonce": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "to": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "value": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            }
        },
        "type": [
            "object"
        ]
    }
	```




__2:__ 
passwd <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>common.Hash</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
	- title: `keccak`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a Keccak 256 hash",
        "pattern": "^0x[a-fA-F\\d]{64}$",
        "title": "keccak",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_sendTransaction", "params": [<args>, <passwd>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "personal_sendTransaction", "params": [<args>, <passwd>]}'
	```


=== "Javascript Console"

	``` js
	personal.sendTransaction(args,passwd);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *PersonalAccountAPI) SendTransaction(ctx context.Context, args TransactionArgs, passwd string) (common.Hash, error) {
	if args.Nonce == nil {
		s.nonceLock.LockAddr(args.from())
		defer s.nonceLock.UnlockAddr(args.from())
	}
	signed, err := s.signTransaction(ctx, &args, passwd)
	if err != nil {
		log.Warn("Failed transaction send attempt", "from", args.from(), "to", args.To, "value", args.Value.ToInt(), "err", err)
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.b, signed)
}// SendTransaction will create a transaction from the given arguments and
// tries to sign it with the key associated with args.From. If the given
// passwd isn't able to decrypt the key it fails.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L461" target="_">View on GitHub →</a>
</p>
</details>

---



### personal_sign

Sign calculates an Ethereum ECDSA signature for:
keccak256("\x19Ethereum Signed Message:\n" + len(message) + message))

Note, the produced signature conforms to the secp256k1 curve R, S and V values,
where the V value will be 27 or 28 for legacy reasons.

The key used to calculate the signature is decrypted with the given password.

https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_sign


#### Params (3)

Parameters must be given _by position_.


__1:__ 
data <code>hexutil.Bytes</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `dataWord`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of some bytes",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "dataWord",
        "type": [
            "string"
        ]
    }
	```




__2:__ 
addr <code>common.Address</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a Keccak 256 hash POINTER`
	- pattern: `^0x[a-fA-F\d]{64}$`
	- title: `keccak`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a Keccak 256 hash POINTER",
        "pattern": "^0x[a-fA-F\\d]{64}$",
        "title": "keccak",
        "type": [
            "string"
        ]
    }
	```




__3:__ 
passwd <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>hexutil.Bytes</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `dataWord`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of some bytes",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "dataWord",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_sign", "params": [<data>, <addr>, <passwd>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "personal_sign", "params": [<data>, <addr>, <passwd>]}'
	```


=== "Javascript Console"

	``` js
	personal.sign(data,addr,passwd);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *PersonalAccountAPI) Sign(ctx context.Context, data hexutil.Bytes, addr common.Address, passwd string) (hexutil.Bytes, error) {
	account := accounts.Account{Address: addr}
	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	signature, err := wallet.SignTextWithPassphrase(account, passwd, data)
	if err != nil {
		log.Warn("Failed data sign attempt", "address", addr, "err", err)
		return nil, err
	}
	signature[crypto.RecoveryIDOffset] += 27
	return signature, nil
}// Sign calculates an Ethereum ECDSA signature for:
// keccak256("\x19Ethereum Signed Message:\n" + len(message) + message))
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The key used to calculate the signature is decrypted with the given password.
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_sign

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L521" target="_">View on GitHub →</a>
</p>
</details>

---



### personal_signTransaction

SignTransaction will create a transaction from the given arguments and
tries to sign it with the key associated with args.From. If the given passwd isn't
able to decrypt the key it fails. The transaction is returned in RLP-form, not broadcast
to other nodes


#### Params (2)

Parameters must be given _by position_.


__1:__ 
args <code>TransactionArgs</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- accessList: 
			- items: 
				- additionalProperties: `false`
				- properties: 
					- address: 
						- pattern: `^0x[a-fA-F\d]{64}$`
						- title: `keccak`
						- type: `string`

					- storageKeys: 
						- items: 
							- description: `Hex representation of a Keccak 256 hash`
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`

						- type: `array`


				- type: `object`

			- type: `array`

		- chainId: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- data: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- from: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- gas: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- gasPrice: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- input: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- maxFeePerGas: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- maxPriorityFeePerGas: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- nonce: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- to: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- value: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "accessList": {
                "items": {
                    "additionalProperties": false,
                    "properties": {
                        "address": {
                            "pattern": "^0x[a-fA-F\\d]{64}$",
                            "title": "keccak",
                            "type": "string"
                        },
                        "storageKeys": {
                            "items": {
                                "description": "Hex representation of a Keccak 256 hash",
                                "pattern": "^0x[a-fA-F\\d]{64}$",
                                "title": "keccak",
                                "type": "string"
                            },
                            "type": "array"
                        }
                    },
                    "type": "object"
                },
                "type": "array"
            },
            "chainId": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "data": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "dataWord",
                "type": "string"
            },
            "from": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "gas": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "gasPrice": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "input": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "dataWord",
                "type": "string"
            },
            "maxFeePerGas": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "maxPriorityFeePerGas": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "nonce": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "to": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "value": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            }
        },
        "type": [
            "object"
        ]
    }
	```




__2:__ 
passwd <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>*SignTransactionResult</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- raw: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- tx: 
			- additionalProperties: `false`
			- type: `object`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "raw": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "dataWord",
                "type": "string"
            },
            "tx": {
                "additionalProperties": false,
                "type": "object"
            }
        },
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_signTransaction", "params": [<args>, <passwd>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "personal_signTransaction", "params": [<args>, <passwd>]}'
	```


=== "Javascript Console"

	``` js
	personal.signTransaction(args,passwd);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *PersonalAccountAPI) SignTransaction(ctx context.Context, args TransactionArgs, passwd string) (*SignTransactionResult, error) {
	if args.From == nil {
		return nil, errors.New("sender not specified")
	}
	if args.Gas == nil {
		return nil, errors.New("gas not specified")
	}
	if args.GasPrice == nil && (args.MaxFeePerGas == nil || args.MaxPriorityFeePerGas == nil) {
		return nil, errors.New("missing gasPrice or maxFeePerGas/maxPriorityFeePerGas")
	}
	if args.Nonce == nil {
		return nil, errors.New("nonce not specified")
	}
	tx := args.toTransaction()
	if err := checkTxFee(tx.GasPrice(), tx.Gas(), s.b.RPCTxFeeCap()); err != nil {
		return nil, err
	}
	signed, err := s.signTransaction(ctx, &args, passwd)
	if err != nil {
		log.Warn("Failed transaction sign attempt", "from", args.from(), "to", args.To, "value", args.Value.ToInt(), "err", err)
		return nil, err
	}
	data, err := signed.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, signed}, nil
}// SignTransaction will create a transaction from the given arguments and
// tries to sign it with the key associated with args.From. If the given passwd isn't
// able to decrypt the key it fails. The transaction is returned in RLP-form, not broadcast
// to other nodes

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L480" target="_">View on GitHub →</a>
</p>
</details>

---



### personal_unlockAccount

UnlockAccount will unlock the account associated with the given address with
the given password for duration seconds. If duration is nil it will use a
default of 300 seconds. It returns an indication if the account was unlocked.


#### Params (3)

Parameters must be given _by position_.


__1:__ 
addr <code>common.Address</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a Keccak 256 hash POINTER`
	- pattern: `^0x[a-fA-F\d]{64}$`
	- title: `keccak`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a Keccak 256 hash POINTER",
        "pattern": "^0x[a-fA-F\\d]{64}$",
        "title": "keccak",
        "type": [
            "string"
        ]
    }
	```




__2:__ 
password <code>string</code> 

  + Required: ✓ Yes





__3:__ 
duration <code>*uint64</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```





#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_unlockAccount", "params": [<addr>, <password>, <duration>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "personal_unlockAccount", "params": [<addr>, <password>, <duration>]}'
	```


=== "Javascript Console"

	``` js
	personal.unlockAccount(addr,password,duration);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *PersonalAccountAPI) UnlockAccount(ctx context.Context, addr common.Address, password string, duration *uint64) (bool, error) {
	if s.b.ExtRPCEnabled() && !s.b.AccountManager().Config().InsecureUnlockAllowed {
		return false, errors.New("account unlock with HTTP access is forbidden")
	}
	const max = uint64(time.Duration(math.MaxInt64) / time.Second)
	var d time.Duration
	if duration == nil {
		d = 300 * time.Second
	} else if *duration > max {
		return false, errors.New("unlock duration too large")
	} else {
		d = time.Duration(*duration) * time.Second
	}
	ks, err := fetchKeystore(s.am)
	if err != nil {
		return false, err
	}
	err = ks.TimedUnlock(accounts.Account{Address: addr}, password, d)
	if err != nil {
		log.Warn("Failed account unlock attempt", "address", addr, "err", err)
	}
	return err == nil, err
}// UnlockAccount will unlock the account associated with the given address with
// the given password for duration seconds. If duration is nil it will use a
// default of 300 seconds. It returns an indication if the account was unlocked.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L402" target="_">View on GitHub →</a>
</p>
</details>

---



### personal_unpair

Unpair deletes a pairing between wallet and geth.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
url <code>string</code> 

  + Required: ✓ Yes





__2:__ 
pin <code>string</code> 

  + Required: ✓ Yes






#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "personal_unpair", "params": [<url>, <pin>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "personal_unpair", "params": [<url>, <pin>]}'
	```


=== "Javascript Console"

	``` js
	personal.unpair(url,pin);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *PersonalAccountAPI) Unpair(ctx context.Context, url string, pin string) error {
	wallet, err := s.am.Wallet(url)
	if err != nil {
		return err
	}
	switch wallet := wallet.( // Unpair deletes a pairing between wallet and geth.
	type) {
	case *scwallet.Wallet:
		return wallet.Unpair([]byte(pin))
	default:
		return errors.New("specified wallet does not support pairing")
	}
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L593" target="_">View on GitHub →</a>
</p>
</details>

---

