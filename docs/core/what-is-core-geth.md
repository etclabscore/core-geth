# What makes Core-Geth what it is?

### Ancestry

It's a fork of ethereum/go-ethereum. Another fork?? Why??

ethereum/go-ethereum assumes that you want to run their ~~network(s)~~ forks only.

5 networks are supported:
- ETH (main net)
- 4 other testnets

Plus, "custom" networks, where the fork activation parameters are editable in a config file (genesis.json)
and you can combine any permutation you want of those with either Ethash PoW or Clique PoA consensus engines,
and have your own network. 

Although these custom/development networks are often touted, for example their support under Puppeth,
their relevance is significantly and, in my opinion, critically handicapped by the opacity and
ETH-fork-level configuration options.

For example, you can run your own network with the opcodes from the Byzantium and Constantinople
ETH-network forks, but you'll also need accept the block reward (reductions) and difficulty bomb (delays)
that came on the ETH network hand-in-hand with the opcodes. You'll need to activate the opcodes with the
same chronology they used on Ethereum mainnet.

These bundles of ETH-fork-level protocol configuration options, when applied to any network that is not ETH,
demonstrate very clearly that networks beside ETH mainnet and a few officially-endorsed testnets are
-- at best -- afterthoughts. (Even the testnets -- they used squashed configs which are NOT analagous to mainnet...)

### Current

This is why core-geth exists.

We believe that the Ethereum _protocol_ -- the yellow paper, the EVM, the peer-to-peer protocols -- these are
technologies that can and should be made available and useable, for the common good.

Ethereum Classic benefits from this kind of thinking, of course, but is not alone in having alternative ideas and opinions
than those officially endorsed by the Ethereum Foundation.

- Core-geth explodes the fork-level configuration options into smaller atomic configuration options at an individual EIP/ECIP-level. This lets you truly build any blockchain theoretically supported by any combination of the features introduced in Ethereum
and Ethereum Classic _at will_, without finding yourself locked in to accept the difficulty bomb or use hardcoded arbitrarily-diminishing block rewards. You an build your own technology with your own economy.
- Beyond Improvement Proposal-level configuration options, Core-geth has completely rethought the way configuration is 
implemented in the client, shifing from a 1-of-a-few default singleton pattern to a throughly interface-driven pattern, enabling
out-of-the-box support for ethereum/go-ethereum, openethereum/openethereum, and multi-geth/multi-geth configuration files and
data types. Core-geth includes a command tool `echainspec` which can freely convert between these data types (although to use your Parity config file with core-geth, it doesn't even need to be converted).

### Future



The principle is: ethereum/go-ethereum is built in support of one network -- ETH. Where a network is a single implementation of
what is really a very generic set of protocols. ETH is one case, it is one set of data, and is governed and maintained by a handful
of familiar faces. Core-geth hopes to bring back to Ethereum an emphasis on technology-first; pushing protocol and configuration and interoperability instead of a priority on a single network. Maybe ETH will continue to dominate and maybe that one network will be enough. But maybe there will be more, others. Maybe the death-prophecy of ethereum/go-ethereum and ETH (about itself!) from Day 0 will actually come true, and maybe in the meantime we can start to see a forest instead of one tree.








