# LAOS Universal Node

For a comprehensive understanding of the LAOS system, please refer to the [LAOS Technical Whitepaper](https://github.com/freeverseio/laos-whitepaper/blob/main/laos.pdf), which extensively covers all components.

The LAOS Universal Node is a decentralized node built using Go, focusing primarily on the remote minting and evolution of NFTs (Non-Fungible Tokens).

It allows users to mint assets on any EVM-Compatible Chain without paying gas fees on those chains.

## Run your own Node

You can start and sync the universal node locally with the following command:
```
$ docker run -p 5001:5001 freeverseio/laos-universal-node:<release> -rpc=<ownership-node-rpc> -evo_rpc=<evochain-node-rpc>
```
The port is for the json-rpc interface.

Please be aware that this version currently does not handle blockchain reorganizations (reorgs). As a precaution, we strongly encourage operating with a heightened safety margin in your ownership chain management.
We are actively working to address this in future updates. Your understanding and cooperation are greatly appreciated as we strive to enhance the capabilities and security of the Universal Node.

## Contributing

Contributions to the LAOS Universal Node project are welcome. When you work on this repo, you adhere to the following rules:

### Git

#### Pre-commit hook

Run `cp git/pre-commit ./.git/hooks/pre-commit` to copy the pre-commit hook to your `.git` folder. This runs a few important checks whenever you have something to commit.

For further details, check [pre-commit](./git/pre-commit).
