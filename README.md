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

Right now there is only one method supported:

```
{"jsonrpc":"2.0", "method":"System.Up", "params":[]}
```


## Contributing

Contributions to the LAOS Universal Node project are welcome. When you work on this repo, you adhere to the following rules:

### Git

#### Rebase

Whenever `main` is updated and you have to pull changes into your branch, you will rebase your branch with `main` and force-push the changes to the repo. This keeps a clean history in feature branches by avoiding additional `merge` commits

Instructions to rebase:
```shell
# first, pull the latest changes from main
git checkout main && git pull
git checkout $MY_BRANCH
# now rebase
git rebase main
git push -f
```

#### Pre-commit hook

Run `cp git/pre-commit ./.git/hooks/pre-commit` to copy the pre-commit hook to your `.git` folder. This runs a few important checks whenever you have something to commit.

For further details, check [pre-commit](./git/pre-commit).
