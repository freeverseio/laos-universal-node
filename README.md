# LAOS Universal Node

For a comprehensive understanding of the LAOS system, please refer to the [LAOS Technical Whitepaper](https://github.com/freeverseio/laos-whitepaper/blob/main/laos.pdf), which extensively covers all components.

The LAOS Universal Node is a decentralized node built using Go, focusing primarily on the remote minting and evolution of NFTs (Non-Fungible Tokens).

It allows users to mint assets on any EVM-Compatible Chain without paying gas fees on those chains.

## Run your own Node

You can start and sync the universal node locally with the following command:
```
$ docker run -p 5001:5001 freeverseio/laos-universal-node:<release> 
```

The port is for the json-rpc interface.

Right now there is only on method supported:

```
{"jsonrpc":"2.0", "method":"System.Up", "params":[]}
```


## Contributing

Contributions to the LAOS Universal Node project are welcome.