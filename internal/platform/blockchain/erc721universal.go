// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc721universal

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// Erc721universalMetaData contains all meta data concerning the Erc721universal contract.
var Erc721universalMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"baseURI_\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721IncorrectOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721InsufficientApproval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC721InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721NonexistentToken\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"baseURI\",\"type\":\"string\"}],\"name\":\"NewERC721Universal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"initOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"isBurnedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801562000010575f80fd5b506040516200257f3803806200257f833981810160405281019062000036919062000238565b8282815f908162000048919062000525565b5080600190816200005a919062000525565b50505080600790816200006e919062000525565b507f74b81bc88402765a52dad72d3d893684f472a679558f3641500e0ee14924a10a3082604051620000a29291906200069c565b60405180910390a1505050620006ce565b5f604051905090565b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6200011482620000cc565b810181811067ffffffffffffffff82111715620001365762000135620000dc565b5b80604052505050565b5f6200014a620000b3565b905062000158828262000109565b919050565b5f67ffffffffffffffff8211156200017a5762000179620000dc565b5b6200018582620000cc565b9050602081019050919050565b5f5b83811015620001b157808201518184015260208101905062000194565b5f8484015250505050565b5f620001d2620001cc846200015d565b6200013f565b905082815260208101848484011115620001f157620001f0620000c8565b5b620001fe84828562000192565b509392505050565b5f82601f8301126200021d576200021c620000c4565b5b81516200022f848260208601620001bc565b91505092915050565b5f805f60608486031215620002525762000251620000bc565b5b5f84015167ffffffffffffffff811115620002725762000271620000c0565b5b620002808682870162000206565b935050602084015167ffffffffffffffff811115620002a457620002a3620000c0565b5b620002b28682870162000206565b925050604084015167ffffffffffffffff811115620002d657620002d5620000c0565b5b620002e48682870162000206565b9150509250925092565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806200033d57607f821691505b602082108103620003535762000352620002f8565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f60088302620003b77fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826200037a565b620003c386836200037a565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f6200040d620004076200040184620003db565b620003e4565b620003db565b9050919050565b5f819050919050565b6200042883620003ed565b62000440620004378262000414565b84845462000386565b825550505050565b5f90565b6200045662000448565b620004638184846200041d565b505050565b5b818110156200048a576200047e5f826200044c565b60018101905062000469565b5050565b601f821115620004d957620004a38162000359565b620004ae846200036b565b81016020851015620004be578190505b620004d6620004cd856200036b565b83018262000468565b50505b505050565b5f82821c905092915050565b5f620004fb5f1984600802620004de565b1980831691505092915050565b5f620005158383620004ea565b9150826002028217905092915050565b6200053082620002ee565b67ffffffffffffffff8111156200054c576200054b620000dc565b5b62000558825462000325565b620005658282856200048e565b5f60209050601f8311600181146200059b575f841562000586578287015190505b62000592858262000508565b86555062000601565b601f198416620005ab8662000359565b5f5b82811015620005d457848901518255600182019150602085019450602081019050620005ad565b86831015620005f45784890151620005f0601f891682620004ea565b8355505b6001600288020188555050505b505050505050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f620006348262000609565b9050919050565b620006468162000628565b82525050565b5f82825260208201905092915050565b5f6200066882620002ee565b6200067481856200064c565b93506200068681856020860162000192565b6200069181620000cc565b840191505092915050565b5f604082019050620006b15f8301856200063b565b8181036020830152620006c581846200065c565b90509392505050565b611ea380620006dc5f395ff3fe608060405234801561000f575f80fd5b5060043610610109575f3560e01c80636352211e116100a057806395d89b411161006f57806395d89b41146102d9578063a22cb465146102f7578063b88d4fde14610313578063c87b56dd1461032f578063e985e9c51461035f57610109565b80636352211e1461022b5780636c0360eb1461025b57806370a08231146102795780638a33a14e146102a957610109565b806323b872dd116100dc57806323b872dd146101a757806342842e0e146101c357806342966c68146101df57806357854508146101fb57610109565b806301ffc9a71461010d57806306fdde031461013d578063081812fc1461015b578063095ea7b31461018b575b5f80fd5b6101276004803603810190610122919061170d565b61038f565b6040516101349190611752565b60405180910390f35b610145610408565b60405161015291906117f5565b60405180910390f35b61017560048036038101906101709190611848565b610497565b60405161018291906118b2565b60405180910390f35b6101a560048036038101906101a091906118f5565b6104b2565b005b6101c160048036038101906101bc9190611933565b6104c8565b005b6101dd60048036038101906101d89190611933565b6105c7565b005b6101f960048036038101906101f49190611848565b6105e6565b005b61021560048036038101906102109190611848565b610625565b60405161022291906118b2565b60405180910390f35b61024560048036038101906102409190611848565b61062e565b60405161025291906118b2565b60405180910390f35b61026361063f565b60405161027091906117f5565b60405180910390f35b610293600480360381019061028e9190611983565b6106cb565b6040516102a091906119bd565b60405180910390f35b6102c360048036038101906102be9190611848565b6106e1565b6040516102d09190611752565b60405180910390f35b6102e16106fe565b6040516102ee91906117f5565b60405180910390f35b610311600480360381019061030c9190611a00565b61078e565b005b61032d60048036038101906103289190611b6a565b6107a4565b005b61034960048036038101906103449190611848565b6107c1565b60405161035691906117f5565b60405180910390f35b61037960048036038101906103749190611bea565b610827565b6040516103869190611752565b60405180910390f35b5f7f57854508000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614806104015750610400826108b5565b5b9050919050565b60605f805461041690611c55565b80601f016020809104026020016040519081016040528092919081815260200182805461044290611c55565b801561048d5780601f106104645761010080835404028352916020019161048d565b820191905f5260205f20905b81548152906001019060200180831161047057829003601f168201915b5050505050905090565b5f6104a182610996565b506104ab82610a1c565b9050919050565b6104c482826104bf610a55565b610a5c565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610538575f6040517f64a0ae9200000000000000000000000000000000000000000000000000000000815260040161052f91906118b2565b60405180910390fd5b5f61054b8383610546610a55565b610a6e565b90508373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146105c1578382826040517f64283d7b0000000000000000000000000000000000000000000000000000000081526004016105b893929190611c85565b60405180910390fd5b50505050565b6105e183838360405180602001604052805f8152506107a4565b505050565b6105f85f826105f3610a55565b610a6e565b50600160065f8381526020019081526020015f205f6101000a81548160ff02191690831515021790555050565b5f819050919050565b5f61063882610996565b9050919050565b6007805461064c90611c55565b80601f016020809104026020016040519081016040528092919081815260200182805461067890611c55565b80156106c35780601f1061069a576101008083540402835291602001916106c3565b820191905f5260205f20905b8154815290600101906020018083116106a657829003601f168201915b505050505081565b5f6c010000000000000000000000009050919050565b6006602052805f5260405f205f915054906101000a900460ff1681565b60606001805461070d90611c55565b80601f016020809104026020016040519081016040528092919081815260200182805461073990611c55565b80156107845780601f1061075b57610100808354040283529160200191610784565b820191905f5260205f20905b81548152906001019060200180831161076757829003601f168201915b5050505050905090565b6107a0610799610a55565b8383610c79565b5050565b6107af8484846104c8565b6107bb84848484610de2565b50505050565b60606107cc82610996565b505f6107d6610f94565b90505f8151116107f45760405180602001604052805f81525061081f565b806107fe84611024565b60405160200161080f929190611cf4565b6040516020818303038152906040525b915050919050565b5f60055f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905092915050565b5f7f80ac58cd000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916148061097f57507f5b5e139f000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916145b8061098f575061098e826110ee565b5b9050919050565b5f806109a183611157565b90505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610a1357826040517f7e273289000000000000000000000000000000000000000000000000000000008152600401610a0a91906119bd565b60405180910390fd5b80915050919050565b5f60045f8381526020019081526020015f205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b5f33905090565b610a6983838360016111db565b505050565b5f80610a7984611157565b90505f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614610aba57610ab981848661139a565b5b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610b4557610af95f855f806111db565b600160035f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825403925050819055505b5f73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1614610bc457600160035f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8460025f8681526020019081526020015f205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550838573ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a4809150509392505050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610ce957816040517f5b08ba18000000000000000000000000000000000000000000000000000000008152600401610ce091906118b2565b60405180910390fd5b8060055f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c3183604051610dd59190611752565b60405180910390a3505050565b5f8373ffffffffffffffffffffffffffffffffffffffff163b1115610f8e578273ffffffffffffffffffffffffffffffffffffffff1663150b7a02610e25610a55565b8685856040518563ffffffff1660e01b8152600401610e479493929190611d69565b6020604051808303815f875af1925050508015610e8257506040513d601f19601f82011682018060405250810190610e7f9190611dc7565b60015b610f03573d805f8114610eb0576040519150601f19603f3d011682016040523d82523d5f602084013e610eb5565b606091505b505f815103610efb57836040517f64a0ae92000000000000000000000000000000000000000000000000000000008152600401610ef291906118b2565b60405180910390fd5b805181602001fd5b63150b7a0260e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614610f8c57836040517f64a0ae92000000000000000000000000000000000000000000000000000000008152600401610f8391906118b2565b60405180910390fd5b505b50505050565b606060078054610fa390611c55565b80601f0160208091040260200160405190810160405280929190818152602001828054610fcf90611c55565b801561101a5780601f10610ff15761010080835404028352916020019161101a565b820191905f5260205f20905b815481529060010190602001808311610ffd57829003601f168201915b5050505050905090565b60605f60016110328461145d565b0190505f8167ffffffffffffffff8111156110505761104f611a46565b5b6040519080825280601f01601f1916602001820160405280156110825781602001600182028036833780820191505090505b5090505f82602001820190505b6001156110e3578080600190039150507f3031323334353637383961626364656600000000000000000000000000000000600a86061a8153600a85816110d8576110d7611df2565b5b0494505f850361108f575b819350505050919050565b5f7f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149050919050565b5f60065f8381526020019081526020015f205f9054906101000a900460ff1615611183575f90506111d6565b5f61118d836115ae565b90505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146111c857806111d2565b6111d183610625565b5b9150505b919050565b808061121357505f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b15611345575f61122284610996565b90505f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415801561128c57508273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614155b801561129f575061129d8184610827565b155b156112e157826040517fa9fbf51f0000000000000000000000000000000000000000000000000000000081526004016112d891906118b2565b60405180910390fd5b811561134357838573ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b8360045f8581526020019081526020015f205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050505050565b6113a58383836115e7565b611458575f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361141957806040517f7e27328900000000000000000000000000000000000000000000000000000000815260040161141091906119bd565b60405180910390fd5b81816040517f177e802f00000000000000000000000000000000000000000000000000000000815260040161144f929190611e1f565b60405180910390fd5b505050565b5f805f90507a184f03e93ff9f4daa797ed6e38ed64bf6a1f01000000000000000083106114b9577a184f03e93ff9f4daa797ed6e38ed64bf6a1f01000000000000000083816114af576114ae611df2565b5b0492506040810190505b6d04ee2d6d415b85acef810000000083106114f6576d04ee2d6d415b85acef810000000083816114ec576114eb611df2565b5b0492506020810190505b662386f26fc10000831061152557662386f26fc10000838161151b5761151a611df2565b5b0492506010810190505b6305f5e100831061154e576305f5e100838161154457611543611df2565b5b0492506008810190505b612710831061157357612710838161156957611568611df2565b5b0492506004810190505b60648310611596576064838161158c5761158b611df2565b5b0492506002810190505b600a83106115a5576001810190505b80915050919050565b5f60025f8381526020019081526020015f205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b5f8073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415801561169e57508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16148061165f575061165e8484610827565b5b8061169d57508273ffffffffffffffffffffffffffffffffffffffff1661168583610a1c565b73ffffffffffffffffffffffffffffffffffffffff16145b5b90509392505050565b5f604051905090565b5f80fd5b5f80fd5b5f7fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b6116ec816116b8565b81146116f6575f80fd5b50565b5f81359050611707816116e3565b92915050565b5f60208284031215611722576117216116b0565b5b5f61172f848285016116f9565b91505092915050565b5f8115159050919050565b61174c81611738565b82525050565b5f6020820190506117655f830184611743565b92915050565b5f81519050919050565b5f82825260208201905092915050565b5f5b838110156117a2578082015181840152602081019050611787565b5f8484015250505050565b5f601f19601f8301169050919050565b5f6117c78261176b565b6117d18185611775565b93506117e1818560208601611785565b6117ea816117ad565b840191505092915050565b5f6020820190508181035f83015261180d81846117bd565b905092915050565b5f819050919050565b61182781611815565b8114611831575f80fd5b50565b5f813590506118428161181e565b92915050565b5f6020828403121561185d5761185c6116b0565b5b5f61186a84828501611834565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61189c82611873565b9050919050565b6118ac81611892565b82525050565b5f6020820190506118c55f8301846118a3565b92915050565b6118d481611892565b81146118de575f80fd5b50565b5f813590506118ef816118cb565b92915050565b5f806040838503121561190b5761190a6116b0565b5b5f611918858286016118e1565b925050602061192985828601611834565b9150509250929050565b5f805f6060848603121561194a576119496116b0565b5b5f611957868287016118e1565b9350506020611968868287016118e1565b925050604061197986828701611834565b9150509250925092565b5f60208284031215611998576119976116b0565b5b5f6119a5848285016118e1565b91505092915050565b6119b781611815565b82525050565b5f6020820190506119d05f8301846119ae565b92915050565b6119df81611738565b81146119e9575f80fd5b50565b5f813590506119fa816119d6565b92915050565b5f8060408385031215611a1657611a156116b0565b5b5f611a23858286016118e1565b9250506020611a34858286016119ec565b9150509250929050565b5f80fd5b5f80fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b611a7c826117ad565b810181811067ffffffffffffffff82111715611a9b57611a9a611a46565b5b80604052505050565b5f611aad6116a7565b9050611ab98282611a73565b919050565b5f67ffffffffffffffff821115611ad857611ad7611a46565b5b611ae1826117ad565b9050602081019050919050565b828183375f83830152505050565b5f611b0e611b0984611abe565b611aa4565b905082815260208101848484011115611b2a57611b29611a42565b5b611b35848285611aee565b509392505050565b5f82601f830112611b5157611b50611a3e565b5b8135611b61848260208601611afc565b91505092915050565b5f805f8060808587031215611b8257611b816116b0565b5b5f611b8f878288016118e1565b9450506020611ba0878288016118e1565b9350506040611bb187828801611834565b925050606085013567ffffffffffffffff811115611bd257611bd16116b4565b5b611bde87828801611b3d565b91505092959194509250565b5f8060408385031215611c0057611bff6116b0565b5b5f611c0d858286016118e1565b9250506020611c1e858286016118e1565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680611c6c57607f821691505b602082108103611c7f57611c7e611c28565b5b50919050565b5f606082019050611c985f8301866118a3565b611ca560208301856119ae565b611cb260408301846118a3565b949350505050565b5f81905092915050565b5f611cce8261176b565b611cd88185611cba565b9350611ce8818560208601611785565b80840191505092915050565b5f611cff8285611cc4565b9150611d0b8284611cc4565b91508190509392505050565b5f81519050919050565b5f82825260208201905092915050565b5f611d3b82611d17565b611d458185611d21565b9350611d55818560208601611785565b611d5e816117ad565b840191505092915050565b5f608082019050611d7c5f8301876118a3565b611d8960208301866118a3565b611d9660408301856119ae565b8181036060830152611da88184611d31565b905095945050505050565b5f81519050611dc1816116e3565b92915050565b5f60208284031215611ddc57611ddb6116b0565b5b5f611de984828501611db3565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f604082019050611e325f8301856118a3565b611e3f60208301846119ae565b939250505056fea264697066735822122067fe11946428bed1091fcd5d480f79f62eb372cf048ec6e97578860bf13ab35a64736f6c637829302e382e32332d646576656c6f702e323032332e31302e32372b636f6d6d69742e6438646539376430005a",
}

// Erc721universalABI is the input ABI used to generate the binding from.
// Deprecated: Use Erc721universalMetaData.ABI instead.
var Erc721universalABI = Erc721universalMetaData.ABI

// Erc721universalBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Erc721universalMetaData.Bin instead.
var Erc721universalBin = Erc721universalMetaData.Bin

// DeployErc721universal deploys a new Ethereum contract, binding an instance of Erc721universal to it.
func DeployErc721universal(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, baseURI_ string) (common.Address, *types.Transaction, *Erc721universal, error) {
	parsed, err := Erc721universalMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Erc721universalBin), backend, name, symbol, baseURI_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Erc721universal{Erc721universalCaller: Erc721universalCaller{contract: contract}, Erc721universalTransactor: Erc721universalTransactor{contract: contract}, Erc721universalFilterer: Erc721universalFilterer{contract: contract}}, nil
}

// Erc721universal is an auto generated Go binding around an Ethereum contract.
type Erc721universal struct {
	Erc721universalCaller     // Read-only binding to the contract
	Erc721universalTransactor // Write-only binding to the contract
	Erc721universalFilterer   // Log filterer for contract events
}

// Erc721universalCaller is an auto generated read-only Go binding around an Ethereum contract.
type Erc721universalCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc721universalTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Erc721universalTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc721universalFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Erc721universalFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc721universalSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Erc721universalSession struct {
	Contract     *Erc721universal  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Erc721universalCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Erc721universalCallerSession struct {
	Contract *Erc721universalCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// Erc721universalTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Erc721universalTransactorSession struct {
	Contract     *Erc721universalTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// Erc721universalRaw is an auto generated low-level Go binding around an Ethereum contract.
type Erc721universalRaw struct {
	Contract *Erc721universal // Generic contract binding to access the raw methods on
}

// Erc721universalCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Erc721universalCallerRaw struct {
	Contract *Erc721universalCaller // Generic read-only contract binding to access the raw methods on
}

// Erc721universalTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Erc721universalTransactorRaw struct {
	Contract *Erc721universalTransactor // Generic write-only contract binding to access the raw methods on
}

// NewErc721universal creates a new instance of Erc721universal, bound to a specific deployed contract.
func NewErc721universal(address common.Address, backend bind.ContractBackend) (*Erc721universal, error) {
	contract, err := bindErc721universal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Erc721universal{Erc721universalCaller: Erc721universalCaller{contract: contract}, Erc721universalTransactor: Erc721universalTransactor{contract: contract}, Erc721universalFilterer: Erc721universalFilterer{contract: contract}}, nil
}

// NewErc721universalCaller creates a new read-only instance of Erc721universal, bound to a specific deployed contract.
func NewErc721universalCaller(address common.Address, caller bind.ContractCaller) (*Erc721universalCaller, error) {
	contract, err := bindErc721universal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Erc721universalCaller{contract: contract}, nil
}

// NewErc721universalTransactor creates a new write-only instance of Erc721universal, bound to a specific deployed contract.
func NewErc721universalTransactor(address common.Address, transactor bind.ContractTransactor) (*Erc721universalTransactor, error) {
	contract, err := bindErc721universal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Erc721universalTransactor{contract: contract}, nil
}

// NewErc721universalFilterer creates a new log filterer instance of Erc721universal, bound to a specific deployed contract.
func NewErc721universalFilterer(address common.Address, filterer bind.ContractFilterer) (*Erc721universalFilterer, error) {
	contract, err := bindErc721universal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Erc721universalFilterer{contract: contract}, nil
}

// bindErc721universal binds a generic wrapper to an already deployed contract.
func bindErc721universal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Erc721universalABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc721universal *Erc721universalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc721universal.Contract.Erc721universalCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc721universal *Erc721universalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc721universal.Contract.Erc721universalTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc721universal *Erc721universalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc721universal.Contract.Erc721universalTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc721universal *Erc721universalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc721universal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc721universal *Erc721universalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc721universal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc721universal *Erc721universalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc721universal.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) pure returns(uint256)
func (_Erc721universal *Erc721universalCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Erc721universal.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) pure returns(uint256)
func (_Erc721universal *Erc721universalSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _Erc721universal.Contract.BalanceOf(&_Erc721universal.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) pure returns(uint256)
func (_Erc721universal *Erc721universalCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _Erc721universal.Contract.BalanceOf(&_Erc721universal.CallOpts, owner)
}

// BaseURI is a free data retrieval call binding the contract method 0x6c0360eb.
//
// Solidity: function baseURI() view returns(string)
func (_Erc721universal *Erc721universalCaller) BaseURI(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Erc721universal.contract.Call(opts, &out, "baseURI")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// BaseURI is a free data retrieval call binding the contract method 0x6c0360eb.
//
// Solidity: function baseURI() view returns(string)
func (_Erc721universal *Erc721universalSession) BaseURI() (string, error) {
	return _Erc721universal.Contract.BaseURI(&_Erc721universal.CallOpts)
}

// BaseURI is a free data retrieval call binding the contract method 0x6c0360eb.
//
// Solidity: function baseURI() view returns(string)
func (_Erc721universal *Erc721universalCallerSession) BaseURI() (string, error) {
	return _Erc721universal.Contract.BaseURI(&_Erc721universal.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Erc721universal *Erc721universalCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Erc721universal.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Erc721universal *Erc721universalSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _Erc721universal.Contract.GetApproved(&_Erc721universal.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Erc721universal *Erc721universalCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _Erc721universal.Contract.GetApproved(&_Erc721universal.CallOpts, tokenId)
}

// InitOwner is a free data retrieval call binding the contract method 0x57854508.
//
// Solidity: function initOwner(uint256 tokenId) pure returns(address)
func (_Erc721universal *Erc721universalCaller) InitOwner(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Erc721universal.contract.Call(opts, &out, "initOwner", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// InitOwner is a free data retrieval call binding the contract method 0x57854508.
//
// Solidity: function initOwner(uint256 tokenId) pure returns(address)
func (_Erc721universal *Erc721universalSession) InitOwner(tokenId *big.Int) (common.Address, error) {
	return _Erc721universal.Contract.InitOwner(&_Erc721universal.CallOpts, tokenId)
}

// InitOwner is a free data retrieval call binding the contract method 0x57854508.
//
// Solidity: function initOwner(uint256 tokenId) pure returns(address)
func (_Erc721universal *Erc721universalCallerSession) InitOwner(tokenId *big.Int) (common.Address, error) {
	return _Erc721universal.Contract.InitOwner(&_Erc721universal.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Erc721universal *Erc721universalCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _Erc721universal.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Erc721universal *Erc721universalSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _Erc721universal.Contract.IsApprovedForAll(&_Erc721universal.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Erc721universal *Erc721universalCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _Erc721universal.Contract.IsApprovedForAll(&_Erc721universal.CallOpts, owner, operator)
}

// IsBurnedToken is a free data retrieval call binding the contract method 0x8a33a14e.
//
// Solidity: function isBurnedToken(uint256 tokenId) view returns(bool)
func (_Erc721universal *Erc721universalCaller) IsBurnedToken(opts *bind.CallOpts, tokenId *big.Int) (bool, error) {
	var out []interface{}
	err := _Erc721universal.contract.Call(opts, &out, "isBurnedToken", tokenId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBurnedToken is a free data retrieval call binding the contract method 0x8a33a14e.
//
// Solidity: function isBurnedToken(uint256 tokenId) view returns(bool)
func (_Erc721universal *Erc721universalSession) IsBurnedToken(tokenId *big.Int) (bool, error) {
	return _Erc721universal.Contract.IsBurnedToken(&_Erc721universal.CallOpts, tokenId)
}

// IsBurnedToken is a free data retrieval call binding the contract method 0x8a33a14e.
//
// Solidity: function isBurnedToken(uint256 tokenId) view returns(bool)
func (_Erc721universal *Erc721universalCallerSession) IsBurnedToken(tokenId *big.Int) (bool, error) {
	return _Erc721universal.Contract.IsBurnedToken(&_Erc721universal.CallOpts, tokenId)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Erc721universal *Erc721universalCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Erc721universal.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Erc721universal *Erc721universalSession) Name() (string, error) {
	return _Erc721universal.Contract.Name(&_Erc721universal.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Erc721universal *Erc721universalCallerSession) Name() (string, error) {
	return _Erc721universal.Contract.Name(&_Erc721universal.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Erc721universal *Erc721universalCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Erc721universal.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Erc721universal *Erc721universalSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _Erc721universal.Contract.OwnerOf(&_Erc721universal.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Erc721universal *Erc721universalCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _Erc721universal.Contract.OwnerOf(&_Erc721universal.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Erc721universal *Erc721universalCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Erc721universal.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Erc721universal *Erc721universalSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Erc721universal.Contract.SupportsInterface(&_Erc721universal.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Erc721universal *Erc721universalCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Erc721universal.Contract.SupportsInterface(&_Erc721universal.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Erc721universal *Erc721universalCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Erc721universal.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Erc721universal *Erc721universalSession) Symbol() (string, error) {
	return _Erc721universal.Contract.Symbol(&_Erc721universal.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Erc721universal *Erc721universalCallerSession) Symbol() (string, error) {
	return _Erc721universal.Contract.Symbol(&_Erc721universal.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Erc721universal *Erc721universalCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _Erc721universal.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Erc721universal *Erc721universalSession) TokenURI(tokenId *big.Int) (string, error) {
	return _Erc721universal.Contract.TokenURI(&_Erc721universal.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Erc721universal *Erc721universalCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _Erc721universal.Contract.TokenURI(&_Erc721universal.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Erc721universal *Erc721universalTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721universal.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Erc721universal *Erc721universalSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721universal.Contract.Approve(&_Erc721universal.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Erc721universal *Erc721universalTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721universal.Contract.Approve(&_Erc721universal.TransactOpts, to, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_Erc721universal *Erc721universalTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721universal.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_Erc721universal *Erc721universalSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721universal.Contract.Burn(&_Erc721universal.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_Erc721universal *Erc721universalTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721universal.Contract.Burn(&_Erc721universal.TransactOpts, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721universal *Erc721universalTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721universal.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721universal *Erc721universalSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721universal.Contract.SafeTransferFrom(&_Erc721universal.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721universal *Erc721universalTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721universal.Contract.SafeTransferFrom(&_Erc721universal.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Erc721universal *Erc721universalTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Erc721universal.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Erc721universal *Erc721universalSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Erc721universal.Contract.SafeTransferFrom0(&_Erc721universal.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Erc721universal *Erc721universalTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Erc721universal.Contract.SafeTransferFrom0(&_Erc721universal.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Erc721universal *Erc721universalTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _Erc721universal.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Erc721universal *Erc721universalSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Erc721universal.Contract.SetApprovalForAll(&_Erc721universal.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Erc721universal *Erc721universalTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Erc721universal.Contract.SetApprovalForAll(&_Erc721universal.TransactOpts, operator, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721universal *Erc721universalTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721universal.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721universal *Erc721universalSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721universal.Contract.TransferFrom(&_Erc721universal.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721universal *Erc721universalTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721universal.Contract.TransferFrom(&_Erc721universal.TransactOpts, from, to, tokenId)
}

// Erc721universalApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Erc721universal contract.
type Erc721universalApprovalIterator struct {
	Event *Erc721universalApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Erc721universalApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721universalApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Erc721universalApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Erc721universalApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721universalApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721universalApproval represents a Approval event raised by the Erc721universal contract.
type Erc721universalApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_Erc721universal *Erc721universalFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*Erc721universalApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Erc721universal.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &Erc721universalApprovalIterator{contract: _Erc721universal.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_Erc721universal *Erc721universalFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *Erc721universalApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Erc721universal.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721universalApproval)
				if err := _Erc721universal.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_Erc721universal *Erc721universalFilterer) ParseApproval(log types.Log) (*Erc721universalApproval, error) {
	event := new(Erc721universalApproval)
	if err := _Erc721universal.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc721universalApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the Erc721universal contract.
type Erc721universalApprovalForAllIterator struct {
	Event *Erc721universalApprovalForAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Erc721universalApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721universalApprovalForAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Erc721universalApprovalForAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Erc721universalApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721universalApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721universalApprovalForAll represents a ApprovalForAll event raised by the Erc721universal contract.
type Erc721universalApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_Erc721universal *Erc721universalFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*Erc721universalApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Erc721universal.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &Erc721universalApprovalForAllIterator{contract: _Erc721universal.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_Erc721universal *Erc721universalFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *Erc721universalApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Erc721universal.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721universalApprovalForAll)
				if err := _Erc721universal.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_Erc721universal *Erc721universalFilterer) ParseApprovalForAll(log types.Log) (*Erc721universalApprovalForAll, error) {
	event := new(Erc721universalApprovalForAll)
	if err := _Erc721universal.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc721universalNewERC721UniversalIterator is returned from FilterNewERC721Universal and is used to iterate over the raw logs and unpacked data for NewERC721Universal events raised by the Erc721universal contract.
type Erc721universalNewERC721UniversalIterator struct {
	Event *Erc721universalNewERC721Universal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Erc721universalNewERC721UniversalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721universalNewERC721Universal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Erc721universalNewERC721Universal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Erc721universalNewERC721UniversalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721universalNewERC721UniversalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721universalNewERC721Universal represents a NewERC721Universal event raised by the Erc721universal contract.
type Erc721universalNewERC721Universal struct {
	NewContractAddress common.Address
	BaseURI            string
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterNewERC721Universal is a free log retrieval operation binding the contract event 0x74b81bc88402765a52dad72d3d893684f472a679558f3641500e0ee14924a10a.
//
// Solidity: event NewERC721Universal(address newContractAddress, string baseURI)
func (_Erc721universal *Erc721universalFilterer) FilterNewERC721Universal(opts *bind.FilterOpts) (*Erc721universalNewERC721UniversalIterator, error) {

	logs, sub, err := _Erc721universal.contract.FilterLogs(opts, "NewERC721Universal")
	if err != nil {
		return nil, err
	}
	return &Erc721universalNewERC721UniversalIterator{contract: _Erc721universal.contract, event: "NewERC721Universal", logs: logs, sub: sub}, nil
}

// WatchNewERC721Universal is a free log subscription operation binding the contract event 0x74b81bc88402765a52dad72d3d893684f472a679558f3641500e0ee14924a10a.
//
// Solidity: event NewERC721Universal(address newContractAddress, string baseURI)
func (_Erc721universal *Erc721universalFilterer) WatchNewERC721Universal(opts *bind.WatchOpts, sink chan<- *Erc721universalNewERC721Universal) (event.Subscription, error) {

	logs, sub, err := _Erc721universal.contract.WatchLogs(opts, "NewERC721Universal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721universalNewERC721Universal)
				if err := _Erc721universal.contract.UnpackLog(event, "NewERC721Universal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewERC721Universal is a log parse operation binding the contract event 0x74b81bc88402765a52dad72d3d893684f472a679558f3641500e0ee14924a10a.
//
// Solidity: event NewERC721Universal(address newContractAddress, string baseURI)
func (_Erc721universal *Erc721universalFilterer) ParseNewERC721Universal(log types.Log) (*Erc721universalNewERC721Universal, error) {
	event := new(Erc721universalNewERC721Universal)
	if err := _Erc721universal.contract.UnpackLog(event, "NewERC721Universal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc721universalTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Erc721universal contract.
type Erc721universalTransferIterator struct {
	Event *Erc721universalTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Erc721universalTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721universalTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Erc721universalTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Erc721universalTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721universalTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721universalTransfer represents a Transfer event raised by the Erc721universal contract.
type Erc721universalTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_Erc721universal *Erc721universalFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*Erc721universalTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Erc721universal.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &Erc721universalTransferIterator{contract: _Erc721universal.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_Erc721universal *Erc721universalFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *Erc721universalTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Erc721universal.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721universalTransfer)
				if err := _Erc721universal.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_Erc721universal *Erc721universalFilterer) ParseTransfer(log types.Log) (*Erc721universalTransfer, error) {
	event := new(Erc721universalTransfer)
	if err := _Erc721universal.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
