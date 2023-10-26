// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ERC721BridgelessMinting

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

// ERC721BridgelessMintingMetaData contains all meta data concerning the ERC721BridgelessMinting contract.
var ERC721BridgelessMintingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"baseURI_\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721IncorrectOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721InsufficientApproval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC721InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721NonexistentToken\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"baseURI\",\"type\":\"string\"}],\"name\":\"NewERC721BridgelessMinting\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"initOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"isBurnedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801562000010575f80fd5b506040516200250638038062002506833981810160405281019062000036919062000238565b8282815f908162000048919062000525565b5080600190816200005a919062000525565b50505080600790816200006e919062000525565b507f821a490a0b4f9fa6744efb226f24ce4c3917ff2fca72c1750947d75a992546103082604051620000a29291906200069c565b60405180910390a1505050620006ce565b5f604051905090565b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6200011482620000cc565b810181811067ffffffffffffffff82111715620001365762000135620000dc565b5b80604052505050565b5f6200014a620000b3565b905062000158828262000109565b919050565b5f67ffffffffffffffff8211156200017a5762000179620000dc565b5b6200018582620000cc565b9050602081019050919050565b5f5b83811015620001b157808201518184015260208101905062000194565b5f8484015250505050565b5f620001d2620001cc846200015d565b6200013f565b905082815260208101848484011115620001f157620001f0620000c8565b5b620001fe84828562000192565b509392505050565b5f82601f8301126200021d576200021c620000c4565b5b81516200022f848260208601620001bc565b91505092915050565b5f805f60608486031215620002525762000251620000bc565b5b5f84015167ffffffffffffffff811115620002725762000271620000c0565b5b620002808682870162000206565b935050602084015167ffffffffffffffff811115620002a457620002a3620000c0565b5b620002b28682870162000206565b925050604084015167ffffffffffffffff811115620002d657620002d5620000c0565b5b620002e48682870162000206565b9150509250925092565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806200033d57607f821691505b602082108103620003535762000352620002f8565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f60088302620003b77fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826200037a565b620003c386836200037a565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f6200040d620004076200040184620003db565b620003e4565b620003db565b9050919050565b5f819050919050565b6200042883620003ed565b62000440620004378262000414565b84845462000386565b825550505050565b5f90565b6200045662000448565b620004638184846200041d565b505050565b5b818110156200048a576200047e5f826200044c565b60018101905062000469565b5050565b601f821115620004d957620004a38162000359565b620004ae846200036b565b81016020851015620004be578190505b620004d6620004cd856200036b565b83018262000468565b50505b505050565b5f82821c905092915050565b5f620004fb5f1984600802620004de565b1980831691505092915050565b5f620005158383620004ea565b9150826002028217905092915050565b6200053082620002ee565b67ffffffffffffffff8111156200054c576200054b620000dc565b5b62000558825462000325565b620005658282856200048e565b5f60209050601f8311600181146200059b575f841562000586578287015190505b62000592858262000508565b86555062000601565b601f198416620005ab8662000359565b5f5b82811015620005d457848901518255600182019150602085019450602081019050620005ad565b86831015620005f45784890151620005f0601f891682620004ea565b8355505b6001600288020188555050505b505050505050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f620006348262000609565b9050919050565b620006468162000628565b82525050565b5f82825260208201905092915050565b5f6200066882620002ee565b6200067481856200064c565b93506200068681856020860162000192565b6200069181620000cc565b840191505092915050565b5f604082019050620006b15f8301856200063b565b8181036020830152620006c581846200065c565b90509392505050565b611e2a80620006dc5f395ff3fe608060405234801561000f575f80fd5b5060043610610109575f3560e01c80636352211e116100a057806395d89b411161006f57806395d89b41146102d9578063a22cb465146102f7578063b88d4fde14610313578063c87b56dd1461032f578063e985e9c51461035f57610109565b80636352211e1461022b5780636c0360eb1461025b57806370a08231146102795780638a33a14e146102a957610109565b806323b872dd116100dc57806323b872dd146101a757806342842e0e146101c357806342966c68146101df57806357854508146101fb57610109565b806301ffc9a71461010d57806306fdde031461013d578063081812fc1461015b578063095ea7b31461018b575b5f80fd5b61012760048036038101906101229190611694565b61038f565b60405161013491906116d9565b60405180910390f35b610145610470565b604051610152919061177c565b60405180910390f35b610175600480360381019061017091906117cf565b6104ff565b6040516101829190611839565b60405180910390f35b6101a560048036038101906101a0919061187c565b61051a565b005b6101c160048036038101906101bc91906118ba565b610530565b005b6101dd60048036038101906101d891906118ba565b61062f565b005b6101f960048036038101906101f491906117cf565b61064e565b005b610215600480360381019061021091906117cf565b61068d565b6040516102229190611839565b60405180910390f35b610245600480360381019061024091906117cf565b610696565b6040516102529190611839565b60405180910390f35b6102636106a7565b604051610270919061177c565b60405180910390f35b610293600480360381019061028e919061190a565b610733565b6040516102a09190611944565b60405180910390f35b6102c360048036038101906102be91906117cf565b610749565b6040516102d091906116d9565b60405180910390f35b6102e1610766565b6040516102ee919061177c565b60405180910390f35b610311600480360381019061030c9190611987565b6107f6565b005b61032d60048036038101906103289190611af1565b61080c565b005b610349600480360381019061034491906117cf565b610829565b604051610356919061177c565b60405180910390f35b61037960048036038101906103749190611b71565b61088f565b60405161038691906116d9565b60405180910390f35b5f7f80ac58cd000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916148061045957507f5b5e139f000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916145b8061046957506104688261091d565b5b9050919050565b60605f805461047e90611bdc565b80601f01602080910402602001604051908101604052809291908181526020018280546104aa90611bdc565b80156104f55780601f106104cc576101008083540402835291602001916104f5565b820191905f5260205f20905b8154815290600101906020018083116104d857829003601f168201915b5050505050905090565b5f61050982610986565b5061051382610a0c565b9050919050565b61052c8282610527610a45565b610a4c565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036105a0575f6040517f64a0ae920000000000000000000000000000000000000000000000000000000081526004016105979190611839565b60405180910390fd5b5f6105b383836105ae610a45565b610a5e565b90508373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610629578382826040517f64283d7b00000000000000000000000000000000000000000000000000000000815260040161062093929190611c0c565b60405180910390fd5b50505050565b61064983838360405180602001604052805f81525061080c565b505050565b6106605f8261065b610a45565b610a5e565b50600160065f8381526020019081526020015f205f6101000a81548160ff02191690831515021790555050565b5f819050919050565b5f6106a082610986565b9050919050565b600780546106b490611bdc565b80601f01602080910402602001604051908101604052809291908181526020018280546106e090611bdc565b801561072b5780601f106107025761010080835404028352916020019161072b565b820191905f5260205f20905b81548152906001019060200180831161070e57829003601f168201915b505050505081565b5f6c010000000000000000000000009050919050565b6006602052805f5260405f205f915054906101000a900460ff1681565b60606001805461077590611bdc565b80601f01602080910402602001604051908101604052809291908181526020018280546107a190611bdc565b80156107ec5780601f106107c3576101008083540402835291602001916107ec565b820191905f5260205f20905b8154815290600101906020018083116107cf57829003601f168201915b5050505050905090565b610808610801610a45565b8383610c69565b5050565b610817848484610530565b61082384848484610dd2565b50505050565b606061083482610986565b505f61083e610f84565b90505f81511161085c5760405180602001604052805f815250610887565b8061086684611014565b604051602001610877929190611c7b565b6040516020818303038152906040525b915050919050565b5f60055f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905092915050565b5f7f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149050919050565b5f80610991836110de565b90505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610a0357826040517f7e2732890000000000000000000000000000000000000000000000000000000081526004016109fa9190611944565b60405180910390fd5b80915050919050565b5f60045f8381526020019081526020015f205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b5f33905090565b610a598383836001611162565b505050565b5f80610a69846110de565b90505f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614610aaa57610aa9818486611321565b5b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610b3557610ae95f855f80611162565b600160035f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825403925050819055505b5f73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1614610bb457600160035f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8460025f8681526020019081526020015f205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550838573ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a4809150509392505050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610cd957816040517f5b08ba18000000000000000000000000000000000000000000000000000000008152600401610cd09190611839565b60405180910390fd5b8060055f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c3183604051610dc591906116d9565b60405180910390a3505050565b5f8373ffffffffffffffffffffffffffffffffffffffff163b1115610f7e578273ffffffffffffffffffffffffffffffffffffffff1663150b7a02610e15610a45565b8685856040518563ffffffff1660e01b8152600401610e379493929190611cf0565b6020604051808303815f875af1925050508015610e7257506040513d601f19601f82011682018060405250810190610e6f9190611d4e565b60015b610ef3573d805f8114610ea0576040519150601f19603f3d011682016040523d82523d5f602084013e610ea5565b606091505b505f815103610eeb57836040517f64a0ae92000000000000000000000000000000000000000000000000000000008152600401610ee29190611839565b60405180910390fd5b805181602001fd5b63150b7a0260e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614610f7c57836040517f64a0ae92000000000000000000000000000000000000000000000000000000008152600401610f739190611839565b60405180910390fd5b505b50505050565b606060078054610f9390611bdc565b80601f0160208091040260200160405190810160405280929190818152602001828054610fbf90611bdc565b801561100a5780601f10610fe15761010080835404028352916020019161100a565b820191905f5260205f20905b815481529060010190602001808311610fed57829003601f168201915b5050505050905090565b60605f6001611022846113e4565b0190505f8167ffffffffffffffff8111156110405761103f6119cd565b5b6040519080825280601f01601f1916602001820160405280156110725781602001600182028036833780820191505090505b5090505f82602001820190505b6001156110d3578080600190039150507f3031323334353637383961626364656600000000000000000000000000000000600a86061a8153600a85816110c8576110c7611d79565b5b0494505f850361107f575b819350505050919050565b5f60065f8381526020019081526020015f205f9054906101000a900460ff161561110a575f905061115d565b5f61111483611535565b90505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161461114f5780611159565b6111588361068d565b5b9150505b919050565b808061119a57505f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b156112cc575f6111a984610986565b90505f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415801561121357508273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614155b80156112265750611224818461088f565b155b1561126857826040517fa9fbf51f00000000000000000000000000000000000000000000000000000000815260040161125f9190611839565b60405180910390fd5b81156112ca57838573ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b8360045f8581526020019081526020015f205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050505050565b61132c83838361156e565b6113df575f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036113a057806040517f7e2732890000000000000000000000000000000000000000000000000000000081526004016113979190611944565b60405180910390fd5b81816040517f177e802f0000000000000000000000000000000000000000000000000000000081526004016113d6929190611da6565b60405180910390fd5b505050565b5f805f90507a184f03e93ff9f4daa797ed6e38ed64bf6a1f0100000000000000008310611440577a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000838161143657611435611d79565b5b0492506040810190505b6d04ee2d6d415b85acef8100000000831061147d576d04ee2d6d415b85acef8100000000838161147357611472611d79565b5b0492506020810190505b662386f26fc1000083106114ac57662386f26fc1000083816114a2576114a1611d79565b5b0492506010810190505b6305f5e10083106114d5576305f5e10083816114cb576114ca611d79565b5b0492506008810190505b61271083106114fa5761271083816114f0576114ef611d79565b5b0492506004810190505b6064831061151d576064838161151357611512611d79565b5b0492506002810190505b600a831061152c576001810190505b80915050919050565b5f60025f8381526020019081526020015f205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b5f8073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415801561162557508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1614806115e657506115e5848461088f565b5b8061162457508273ffffffffffffffffffffffffffffffffffffffff1661160c83610a0c565b73ffffffffffffffffffffffffffffffffffffffff16145b5b90509392505050565b5f604051905090565b5f80fd5b5f80fd5b5f7fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b6116738161163f565b811461167d575f80fd5b50565b5f8135905061168e8161166a565b92915050565b5f602082840312156116a9576116a8611637565b5b5f6116b684828501611680565b91505092915050565b5f8115159050919050565b6116d3816116bf565b82525050565b5f6020820190506116ec5f8301846116ca565b92915050565b5f81519050919050565b5f82825260208201905092915050565b5f5b8381101561172957808201518184015260208101905061170e565b5f8484015250505050565b5f601f19601f8301169050919050565b5f61174e826116f2565b61175881856116fc565b935061176881856020860161170c565b61177181611734565b840191505092915050565b5f6020820190508181035f8301526117948184611744565b905092915050565b5f819050919050565b6117ae8161179c565b81146117b8575f80fd5b50565b5f813590506117c9816117a5565b92915050565b5f602082840312156117e4576117e3611637565b5b5f6117f1848285016117bb565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f611823826117fa565b9050919050565b61183381611819565b82525050565b5f60208201905061184c5f83018461182a565b92915050565b61185b81611819565b8114611865575f80fd5b50565b5f8135905061187681611852565b92915050565b5f806040838503121561189257611891611637565b5b5f61189f85828601611868565b92505060206118b0858286016117bb565b9150509250929050565b5f805f606084860312156118d1576118d0611637565b5b5f6118de86828701611868565b93505060206118ef86828701611868565b9250506040611900868287016117bb565b9150509250925092565b5f6020828403121561191f5761191e611637565b5b5f61192c84828501611868565b91505092915050565b61193e8161179c565b82525050565b5f6020820190506119575f830184611935565b92915050565b611966816116bf565b8114611970575f80fd5b50565b5f813590506119818161195d565b92915050565b5f806040838503121561199d5761199c611637565b5b5f6119aa85828601611868565b92505060206119bb85828601611973565b9150509250929050565b5f80fd5b5f80fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b611a0382611734565b810181811067ffffffffffffffff82111715611a2257611a216119cd565b5b80604052505050565b5f611a3461162e565b9050611a4082826119fa565b919050565b5f67ffffffffffffffff821115611a5f57611a5e6119cd565b5b611a6882611734565b9050602081019050919050565b828183375f83830152505050565b5f611a95611a9084611a45565b611a2b565b905082815260208101848484011115611ab157611ab06119c9565b5b611abc848285611a75565b509392505050565b5f82601f830112611ad857611ad76119c5565b5b8135611ae8848260208601611a83565b91505092915050565b5f805f8060808587031215611b0957611b08611637565b5b5f611b1687828801611868565b9450506020611b2787828801611868565b9350506040611b38878288016117bb565b925050606085013567ffffffffffffffff811115611b5957611b5861163b565b5b611b6587828801611ac4565b91505092959194509250565b5f8060408385031215611b8757611b86611637565b5b5f611b9485828601611868565b9250506020611ba585828601611868565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680611bf357607f821691505b602082108103611c0657611c05611baf565b5b50919050565b5f606082019050611c1f5f83018661182a565b611c2c6020830185611935565b611c39604083018461182a565b949350505050565b5f81905092915050565b5f611c55826116f2565b611c5f8185611c41565b9350611c6f81856020860161170c565b80840191505092915050565b5f611c868285611c4b565b9150611c928284611c4b565b91508190509392505050565b5f81519050919050565b5f82825260208201905092915050565b5f611cc282611c9e565b611ccc8185611ca8565b9350611cdc81856020860161170c565b611ce581611734565b840191505092915050565b5f608082019050611d035f83018761182a565b611d10602083018661182a565b611d1d6040830185611935565b8181036060830152611d2f8184611cb8565b905095945050505050565b5f81519050611d488161166a565b92915050565b5f60208284031215611d6357611d62611637565b5b5f611d7084828501611d3a565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f604082019050611db95f83018561182a565b611dc66020830184611935565b939250505056fea26469706673582212200320f0eaf9b27243e53bb6091d338330536ee9db430a94506f883c976e04d41164736f6c637829302e382e32322d646576656c6f702e323032332e31302e31372b636f6d6d69742e6539386631373464005a",
}

// ERC721BridgelessMintingABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC721BridgelessMintingMetaData.ABI instead.
var ERC721BridgelessMintingABI = ERC721BridgelessMintingMetaData.ABI

// ERC721BridgelessMintingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC721BridgelessMintingMetaData.Bin instead.
var ERC721BridgelessMintingBin = ERC721BridgelessMintingMetaData.Bin

// DeployERC721BridgelessMinting deploys a new Ethereum contract, binding an instance of ERC721BridgelessMinting to it.
func DeployERC721BridgelessMinting(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, baseURI_ string) (common.Address, *types.Transaction, *ERC721BridgelessMinting, error) {
	parsed, err := ERC721BridgelessMintingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC721BridgelessMintingBin), backend, name, symbol, baseURI_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC721BridgelessMinting{ERC721BridgelessMintingCaller: ERC721BridgelessMintingCaller{contract: contract}, ERC721BridgelessMintingTransactor: ERC721BridgelessMintingTransactor{contract: contract}, ERC721BridgelessMintingFilterer: ERC721BridgelessMintingFilterer{contract: contract}}, nil
}

// ERC721BridgelessMinting is an auto generated Go binding around an Ethereum contract.
type ERC721BridgelessMinting struct {
	ERC721BridgelessMintingCaller     // Read-only binding to the contract
	ERC721BridgelessMintingTransactor // Write-only binding to the contract
	ERC721BridgelessMintingFilterer   // Log filterer for contract events
}

// ERC721BridgelessMintingCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC721BridgelessMintingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721BridgelessMintingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC721BridgelessMintingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721BridgelessMintingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC721BridgelessMintingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721BridgelessMintingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC721BridgelessMintingSession struct {
	Contract     *ERC721BridgelessMinting // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ERC721BridgelessMintingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC721BridgelessMintingCallerSession struct {
	Contract *ERC721BridgelessMintingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// ERC721BridgelessMintingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC721BridgelessMintingTransactorSession struct {
	Contract     *ERC721BridgelessMintingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// ERC721BridgelessMintingRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC721BridgelessMintingRaw struct {
	Contract *ERC721BridgelessMinting // Generic contract binding to access the raw methods on
}

// ERC721BridgelessMintingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC721BridgelessMintingCallerRaw struct {
	Contract *ERC721BridgelessMintingCaller // Generic read-only contract binding to access the raw methods on
}

// ERC721BridgelessMintingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC721BridgelessMintingTransactorRaw struct {
	Contract *ERC721BridgelessMintingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC721BridgelessMinting creates a new instance of ERC721BridgelessMinting, bound to a specific deployed contract.
func NewERC721BridgelessMinting(address common.Address, backend bind.ContractBackend) (*ERC721BridgelessMinting, error) {
	contract, err := bindERC721BridgelessMinting(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC721BridgelessMinting{ERC721BridgelessMintingCaller: ERC721BridgelessMintingCaller{contract: contract}, ERC721BridgelessMintingTransactor: ERC721BridgelessMintingTransactor{contract: contract}, ERC721BridgelessMintingFilterer: ERC721BridgelessMintingFilterer{contract: contract}}, nil
}

// NewERC721BridgelessMintingCaller creates a new read-only instance of ERC721BridgelessMinting, bound to a specific deployed contract.
func NewERC721BridgelessMintingCaller(address common.Address, caller bind.ContractCaller) (*ERC721BridgelessMintingCaller, error) {
	contract, err := bindERC721BridgelessMinting(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721BridgelessMintingCaller{contract: contract}, nil
}

// NewERC721BridgelessMintingTransactor creates a new write-only instance of ERC721BridgelessMinting, bound to a specific deployed contract.
func NewERC721BridgelessMintingTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC721BridgelessMintingTransactor, error) {
	contract, err := bindERC721BridgelessMinting(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721BridgelessMintingTransactor{contract: contract}, nil
}

// NewERC721BridgelessMintingFilterer creates a new log filterer instance of ERC721BridgelessMinting, bound to a specific deployed contract.
func NewERC721BridgelessMintingFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC721BridgelessMintingFilterer, error) {
	contract, err := bindERC721BridgelessMinting(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC721BridgelessMintingFilterer{contract: contract}, nil
}

// bindERC721BridgelessMinting binds a generic wrapper to an already deployed contract.
func bindERC721BridgelessMinting(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC721BridgelessMintingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721BridgelessMinting *ERC721BridgelessMintingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC721BridgelessMinting.Contract.ERC721BridgelessMintingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721BridgelessMinting *ERC721BridgelessMintingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.ERC721BridgelessMintingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721BridgelessMinting *ERC721BridgelessMintingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.ERC721BridgelessMintingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC721BridgelessMinting.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721BridgelessMinting *ERC721BridgelessMintingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721BridgelessMinting *ERC721BridgelessMintingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) pure returns(uint256)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC721BridgelessMinting.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) pure returns(uint256)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC721BridgelessMinting.Contract.BalanceOf(&_ERC721BridgelessMinting.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) pure returns(uint256)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC721BridgelessMinting.Contract.BalanceOf(&_ERC721BridgelessMinting.CallOpts, owner)
}

// BaseURI is a free data retrieval call binding the contract method 0x6c0360eb.
//
// Solidity: function baseURI() view returns(string)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCaller) BaseURI(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC721BridgelessMinting.contract.Call(opts, &out, "baseURI")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// BaseURI is a free data retrieval call binding the contract method 0x6c0360eb.
//
// Solidity: function baseURI() view returns(string)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) BaseURI() (string, error) {
	return _ERC721BridgelessMinting.Contract.BaseURI(&_ERC721BridgelessMinting.CallOpts)
}

// BaseURI is a free data retrieval call binding the contract method 0x6c0360eb.
//
// Solidity: function baseURI() view returns(string)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCallerSession) BaseURI() (string, error) {
	return _ERC721BridgelessMinting.Contract.BaseURI(&_ERC721BridgelessMinting.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ERC721BridgelessMinting.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ERC721BridgelessMinting.Contract.GetApproved(&_ERC721BridgelessMinting.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ERC721BridgelessMinting.Contract.GetApproved(&_ERC721BridgelessMinting.CallOpts, tokenId)
}

// InitOwner is a free data retrieval call binding the contract method 0x57854508.
//
// Solidity: function initOwner(uint256 tokenId) pure returns(address)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCaller) InitOwner(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ERC721BridgelessMinting.contract.Call(opts, &out, "initOwner", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// InitOwner is a free data retrieval call binding the contract method 0x57854508.
//
// Solidity: function initOwner(uint256 tokenId) pure returns(address)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) InitOwner(tokenId *big.Int) (common.Address, error) {
	return _ERC721BridgelessMinting.Contract.InitOwner(&_ERC721BridgelessMinting.CallOpts, tokenId)
}

// InitOwner is a free data retrieval call binding the contract method 0x57854508.
//
// Solidity: function initOwner(uint256 tokenId) pure returns(address)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCallerSession) InitOwner(tokenId *big.Int) (common.Address, error) {
	return _ERC721BridgelessMinting.Contract.InitOwner(&_ERC721BridgelessMinting.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _ERC721BridgelessMinting.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC721BridgelessMinting.Contract.IsApprovedForAll(&_ERC721BridgelessMinting.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC721BridgelessMinting.Contract.IsApprovedForAll(&_ERC721BridgelessMinting.CallOpts, owner, operator)
}

// IsBurnedToken is a free data retrieval call binding the contract method 0x8a33a14e.
//
// Solidity: function isBurnedToken(uint256 tokenId) view returns(bool)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCaller) IsBurnedToken(opts *bind.CallOpts, tokenId *big.Int) (bool, error) {
	var out []interface{}
	err := _ERC721BridgelessMinting.contract.Call(opts, &out, "isBurnedToken", tokenId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBurnedToken is a free data retrieval call binding the contract method 0x8a33a14e.
//
// Solidity: function isBurnedToken(uint256 tokenId) view returns(bool)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) IsBurnedToken(tokenId *big.Int) (bool, error) {
	return _ERC721BridgelessMinting.Contract.IsBurnedToken(&_ERC721BridgelessMinting.CallOpts, tokenId)
}

// IsBurnedToken is a free data retrieval call binding the contract method 0x8a33a14e.
//
// Solidity: function isBurnedToken(uint256 tokenId) view returns(bool)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCallerSession) IsBurnedToken(tokenId *big.Int) (bool, error) {
	return _ERC721BridgelessMinting.Contract.IsBurnedToken(&_ERC721BridgelessMinting.CallOpts, tokenId)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC721BridgelessMinting.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) Name() (string, error) {
	return _ERC721BridgelessMinting.Contract.Name(&_ERC721BridgelessMinting.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCallerSession) Name() (string, error) {
	return _ERC721BridgelessMinting.Contract.Name(&_ERC721BridgelessMinting.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ERC721BridgelessMinting.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ERC721BridgelessMinting.Contract.OwnerOf(&_ERC721BridgelessMinting.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ERC721BridgelessMinting.Contract.OwnerOf(&_ERC721BridgelessMinting.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ERC721BridgelessMinting.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC721BridgelessMinting.Contract.SupportsInterface(&_ERC721BridgelessMinting.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC721BridgelessMinting.Contract.SupportsInterface(&_ERC721BridgelessMinting.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC721BridgelessMinting.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) Symbol() (string, error) {
	return _ERC721BridgelessMinting.Contract.Symbol(&_ERC721BridgelessMinting.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCallerSession) Symbol() (string, error) {
	return _ERC721BridgelessMinting.Contract.Symbol(&_ERC721BridgelessMinting.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _ERC721BridgelessMinting.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) TokenURI(tokenId *big.Int) (string, error) {
	return _ERC721BridgelessMinting.Contract.TokenURI(&_ERC721BridgelessMinting.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _ERC721BridgelessMinting.Contract.TokenURI(&_ERC721BridgelessMinting.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.Approve(&_ERC721BridgelessMinting.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.Approve(&_ERC721BridgelessMinting.TransactOpts, to, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.Burn(&_ERC721BridgelessMinting.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.Burn(&_ERC721BridgelessMinting.TransactOpts, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.SafeTransferFrom(&_ERC721BridgelessMinting.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.SafeTransferFrom(&_ERC721BridgelessMinting.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.SafeTransferFrom0(&_ERC721BridgelessMinting.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.SafeTransferFrom0(&_ERC721BridgelessMinting.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.SetApprovalForAll(&_ERC721BridgelessMinting.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.SetApprovalForAll(&_ERC721BridgelessMinting.TransactOpts, operator, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.TransferFrom(&_ERC721BridgelessMinting.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721BridgelessMinting *ERC721BridgelessMintingTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721BridgelessMinting.Contract.TransferFrom(&_ERC721BridgelessMinting.TransactOpts, from, to, tokenId)
}

// ERC721BridgelessMintingApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC721BridgelessMinting contract.
type ERC721BridgelessMintingApprovalIterator struct {
	Event *ERC721BridgelessMintingApproval // Event containing the contract specifics and raw log

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
func (it *ERC721BridgelessMintingApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721BridgelessMintingApproval)
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
		it.Event = new(ERC721BridgelessMintingApproval)
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
func (it *ERC721BridgelessMintingApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721BridgelessMintingApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721BridgelessMintingApproval represents a Approval event raised by the ERC721BridgelessMinting contract.
type ERC721BridgelessMintingApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*ERC721BridgelessMintingApprovalIterator, error) {

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

	logs, sub, err := _ERC721BridgelessMinting.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ERC721BridgelessMintingApprovalIterator{contract: _ERC721BridgelessMinting.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC721BridgelessMintingApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _ERC721BridgelessMinting.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721BridgelessMintingApproval)
				if err := _ERC721BridgelessMinting.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ERC721BridgelessMinting *ERC721BridgelessMintingFilterer) ParseApproval(log types.Log) (*ERC721BridgelessMintingApproval, error) {
	event := new(ERC721BridgelessMintingApproval)
	if err := _ERC721BridgelessMinting.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC721BridgelessMintingApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the ERC721BridgelessMinting contract.
type ERC721BridgelessMintingApprovalForAllIterator struct {
	Event *ERC721BridgelessMintingApprovalForAll // Event containing the contract specifics and raw log

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
func (it *ERC721BridgelessMintingApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721BridgelessMintingApprovalForAll)
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
		it.Event = new(ERC721BridgelessMintingApprovalForAll)
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
func (it *ERC721BridgelessMintingApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721BridgelessMintingApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721BridgelessMintingApprovalForAll represents a ApprovalForAll event raised by the ERC721BridgelessMinting contract.
type ERC721BridgelessMintingApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*ERC721BridgelessMintingApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC721BridgelessMinting.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ERC721BridgelessMintingApprovalForAllIterator{contract: _ERC721BridgelessMinting.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *ERC721BridgelessMintingApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC721BridgelessMinting.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721BridgelessMintingApprovalForAll)
				if err := _ERC721BridgelessMinting.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_ERC721BridgelessMinting *ERC721BridgelessMintingFilterer) ParseApprovalForAll(log types.Log) (*ERC721BridgelessMintingApprovalForAll, error) {
	event := new(ERC721BridgelessMintingApprovalForAll)
	if err := _ERC721BridgelessMinting.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC721BridgelessMintingNewERC721BridgelessMintingIterator is returned from FilterNewERC721BridgelessMinting and is used to iterate over the raw logs and unpacked data for NewERC721BridgelessMinting events raised by the ERC721BridgelessMinting contract.
type ERC721BridgelessMintingNewERC721BridgelessMintingIterator struct {
	Event *ERC721BridgelessMintingNewERC721BridgelessMinting // Event containing the contract specifics and raw log

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
func (it *ERC721BridgelessMintingNewERC721BridgelessMintingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721BridgelessMintingNewERC721BridgelessMinting)
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
		it.Event = new(ERC721BridgelessMintingNewERC721BridgelessMinting)
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
func (it *ERC721BridgelessMintingNewERC721BridgelessMintingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721BridgelessMintingNewERC721BridgelessMintingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721BridgelessMintingNewERC721BridgelessMinting represents a NewERC721BridgelessMinting event raised by the ERC721BridgelessMinting contract.
type ERC721BridgelessMintingNewERC721BridgelessMinting struct {
	NewContractAddress common.Address
	BaseURI            string
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterNewERC721BridgelessMinting is a free log retrieval operation binding the contract event 0x821a490a0b4f9fa6744efb226f24ce4c3917ff2fca72c1750947d75a99254610.
//
// Solidity: event NewERC721BridgelessMinting(address newContractAddress, string baseURI)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingFilterer) FilterNewERC721BridgelessMinting(opts *bind.FilterOpts) (*ERC721BridgelessMintingNewERC721BridgelessMintingIterator, error) {

	logs, sub, err := _ERC721BridgelessMinting.contract.FilterLogs(opts, "NewERC721BridgelessMinting")
	if err != nil {
		return nil, err
	}
	return &ERC721BridgelessMintingNewERC721BridgelessMintingIterator{contract: _ERC721BridgelessMinting.contract, event: "NewERC721BridgelessMinting", logs: logs, sub: sub}, nil
}

// WatchNewERC721BridgelessMinting is a free log subscription operation binding the contract event 0x821a490a0b4f9fa6744efb226f24ce4c3917ff2fca72c1750947d75a99254610.
//
// Solidity: event NewERC721BridgelessMinting(address newContractAddress, string baseURI)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingFilterer) WatchNewERC721BridgelessMinting(opts *bind.WatchOpts, sink chan<- *ERC721BridgelessMintingNewERC721BridgelessMinting) (event.Subscription, error) {

	logs, sub, err := _ERC721BridgelessMinting.contract.WatchLogs(opts, "NewERC721BridgelessMinting")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721BridgelessMintingNewERC721BridgelessMinting)
				if err := _ERC721BridgelessMinting.contract.UnpackLog(event, "NewERC721BridgelessMinting", log); err != nil {
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

// ParseNewERC721BridgelessMinting is a log parse operation binding the contract event 0x821a490a0b4f9fa6744efb226f24ce4c3917ff2fca72c1750947d75a99254610.
//
// Solidity: event NewERC721BridgelessMinting(address newContractAddress, string baseURI)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingFilterer) ParseNewERC721BridgelessMinting(log types.Log) (*ERC721BridgelessMintingNewERC721BridgelessMinting, error) {
	event := new(ERC721BridgelessMintingNewERC721BridgelessMinting)
	if err := _ERC721BridgelessMinting.contract.UnpackLog(event, "NewERC721BridgelessMinting", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC721BridgelessMintingTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC721BridgelessMinting contract.
type ERC721BridgelessMintingTransferIterator struct {
	Event *ERC721BridgelessMintingTransfer // Event containing the contract specifics and raw log

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
func (it *ERC721BridgelessMintingTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721BridgelessMintingTransfer)
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
		it.Event = new(ERC721BridgelessMintingTransfer)
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
func (it *ERC721BridgelessMintingTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721BridgelessMintingTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721BridgelessMintingTransfer represents a Transfer event raised by the ERC721BridgelessMinting contract.
type ERC721BridgelessMintingTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*ERC721BridgelessMintingTransferIterator, error) {

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

	logs, sub, err := _ERC721BridgelessMinting.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ERC721BridgelessMintingTransferIterator{contract: _ERC721BridgelessMinting.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ERC721BridgelessMinting *ERC721BridgelessMintingFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC721BridgelessMintingTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _ERC721BridgelessMinting.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721BridgelessMintingTransfer)
				if err := _ERC721BridgelessMinting.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ERC721BridgelessMinting *ERC721BridgelessMintingFilterer) ParseTransfer(log types.Log) (*ERC721BridgelessMintingTransfer, error) {
	event := new(ERC721BridgelessMintingTransfer)
	if err := _ERC721BridgelessMinting.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
