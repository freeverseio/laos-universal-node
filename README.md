# LAOS Universal Node

For a comprehensive understanding of the LAOS system, please refer to the [LAOS Technical Whitepaper](https://github.com/freeverseio/laos-whitepaper/blob/main/laos.pdf), which extensively covers all components.

The LAOS Universal Node is a decentralized node built using Go, focusing primarily on the remote minting and evolution of NFTs (Non-Fungible Tokens).

It allows users to mint assets on any EVM-Compatible Chain without paying gas fees on those chains.

## Running Your Own Node

You can start and sync the universal node locally with the following command:
```
$ docker run -p 5001:5001 freeverseio/laos-universal-node:<release> -rpc=<ownership-node-rpc> -evo_rpc=<evochain-node-rpc>
```
The port is for the json-rpc interface.

Please be aware that this version currently does not handle blockchain reorganizations (reorgs). As a precaution, we strongly encourage operating with a heightened safety margin in your ownership chain management.
We are actively working to address this in future updates. Your understanding and cooperation are greatly appreciated as we strive to enhance the capabilities and security of the Universal Node.

## Contributing

We welcome your contributions to the LAOS Universal Node project. By participating, you agree to adhere to our guidelines:

### Git Practices

#### Pre-commit Hook

Ensure code quality by running `cp git/pre-commit ./.git/hooks/pre-commit`, integrating our pre-commit hook into your workflow. This executes several crucial checks prior to committing changes.

For more information, see [pre-commit](./git/pre-commit).

#### Versioning and Constraints

- **Evo Range:** For now it must not be changed. Default value is 1.
- **Tag Version:** The codebase should be tagged as “0.1”, reflecting its current stage.

### Project Status

Please note, the LAOS Universal Node is in its Beta phase and is not yet ready for production use. We're dedicated to refining its functionalities for a seamless experience in the future.

### Minimum System Requirements

To ensure optimal performance of the LAOS Universal Node, your system should meet the following specifications:

- **CPU:** minimum: 4 vCPU / recommended: 6 vCPU
- **Memory:** minimum: 10 GB RAM / recommended: 12 GB RA
- **Storage:** minimum: 512 GB / recommended: 1 TB

We're excited to see how you'll leverage the LAOS Universal Node. Your feedback and contributions are invaluable as we strive to revolutionize the NFT landscape.
