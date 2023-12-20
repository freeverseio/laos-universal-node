# LAOS Universal Node

For a comprehensive understanding of the LAOS system, please refer to the [LAOS Technical Whitepaper](https://github.com/freeverseio/laos-whitepaper/blob/main/laos.pdf), which extensively covers all components.

The LAOS Universal Node (uNode) streamlines the integration process for DApps aiming to incorporate bridgeless minting and evolution across various chains, including Ethereum, by merely adjusting the RPC endpoint to connect to the relevant Universal Nodes.

## Starting and Syncing the Universal Node Locally

Launch and synchronize the uNode locally using the command:
```
$ docker run -p 5001:5001 freeverseio/laos-universal-node:<release> -rpc=<ownership-node-rpc> -evo_rpc=<evochain-node-rpc>
```
The specified port (`5001`) is used for the JSON-RPC interface.

Use the following command to display all available command-line options:
```
$ docker run freeverseio/laos-universal-node --help 
```

Please be aware that the current version of the uNode does not handle blockchain reorganizations (reorgs). As a precaution, we strongly encourage operating with a heightened safety margin in your ownership chain management.
We are actively working to address this in future updates. Your understanding and cooperation are greatly appreciated as we strive to enhance the capabilities and security of the uNode.

## Contributing

Contributions to the LAOS Universal Node project are welcome. Please adhere to [GitHub's contribution guidelines](https://docs.github.com/en/get-started/quickstart/contributing-to-projects) to ensure a smooth collaboration process.

#### Pre-commit hook

To ensure code consistency, please set up the pre-commit hook. Execute `cp git/pre-commit ./.git/hooks/pre-commit` to copy the pre-defined hook into your `.git` folder. This will automatically perform essential checks before each commit.

For more information, please refer to the [pre-commit details](./git/pre-commit).