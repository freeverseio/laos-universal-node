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

Contributions to the LAOS Universal Node project are welcome. Please adhere to [Github's contribution guidelines](https://docs.github.com/en/get-started/quickstart/contributing-to-projects) to ensure a smooth collaboration process.

#### Pre-commit hook

To ensure code consistency, please set up the pre-commit hook. Execute `cp git/pre-commit ./.git/hooks/pre-commit` to copy the pre-defined hook into your `.git` folder. This will automatically perform essential checks before each commit.

For more information, please refer to the [pre-commit details](./git/pre-commit).