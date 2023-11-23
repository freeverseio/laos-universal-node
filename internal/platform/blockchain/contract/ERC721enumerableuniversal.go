// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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
	_ = abi.ConvertType
)

// EnumerableMetaData contains all meta data concerning the Enumerable contract.
var EnumerableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ERC721EnumerableForbiddenBatchMint\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721IncorrectOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721InsufficientApproval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC721InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721NonexistentToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"ERC721OutOfBoundsIndex\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"tokenByIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"tokenOfOwnerByIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801562000010575f80fd5b506040518060400160405280601981526020017f455243373231456e756d657261626c65556e6976657273616c000000000000008152506040518060400160405280600381526020017f4d544b0000000000000000000000000000000000000000000000000000000000815250815f90816200008d91906200030c565b5080600190816200009f91906200030c565b505050620003f0565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806200012457607f821691505b6020821081036200013a5762000139620000df565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026200019e7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8262000161565b620001aa868362000161565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f620001f4620001ee620001e884620001c2565b620001cb565b620001c2565b9050919050565b5f819050919050565b6200020f83620001d4565b620002276200021e82620001fb565b8484546200016d565b825550505050565b5f90565b6200023d6200022f565b6200024a81848462000204565b505050565b5b818110156200027157620002655f8262000233565b60018101905062000250565b5050565b601f821115620002c0576200028a8162000140565b620002958462000152565b81016020851015620002a5578190505b620002bd620002b48562000152565b8301826200024f565b50505b505050565b5f82821c905092915050565b5f620002e25f1984600802620002c5565b1980831691505092915050565b5f620002fc8383620002d1565b9150826002028217905092915050565b6200031782620000a8565b67ffffffffffffffff811115620003335762000332620000b2565b5b6200033f82546200010c565b6200034c82828562000275565b5f60209050601f83116001811462000382575f84156200036d578287015190505b620003798582620002ef565b865550620003e8565b601f198416620003928662000140565b5f5b82811015620003bb5784890151825560018201915060208501945060208101905062000394565b86831015620003db5784890151620003d7601f891682620002d1565b8355505b6001600288020188555050505b505050505050565b6122f480620003fe5f395ff3fe608060405234801561000f575f80fd5b50600436106100fe575f3560e01c80634f6ccce711610095578063a22cb46511610064578063a22cb465146102d0578063b88d4fde146102ec578063c87b56dd14610308578063e985e9c514610338576100fe565b80634f6ccce7146102225780636352211e1461025257806370a082311461028257806395d89b41146102b2576100fe565b806318160ddd116100d157806318160ddd1461019c57806323b872dd146101ba5780632f745c59146101d657806342842e0e14610206576100fe565b806301ffc9a71461010257806306fdde0314610132578063081812fc14610150578063095ea7b314610180575b5f80fd5b61011c60048036038101906101179190611acb565b610368565b6040516101299190611b10565b60405180910390f35b61013a610379565b6040516101479190611bb3565b60405180910390f35b61016a60048036038101906101659190611c06565b610408565b6040516101779190611c70565b60405180910390f35b61019a60048036038101906101959190611cb3565b610423565b005b6101a4610439565b6040516101b19190611d00565b60405180910390f35b6101d460048036038101906101cf9190611d19565b610445565b005b6101f060048036038101906101eb9190611cb3565b610544565b6040516101fd9190611d00565b60405180910390f35b610220600480360381019061021b9190611d19565b6105e8565b005b61023c60048036038101906102379190611c06565b610607565b6040516102499190611d00565b60405180910390f35b61026c60048036038101906102679190611c06565b610679565b6040516102799190611c70565b60405180910390f35b61029c60048036038101906102979190611d69565b61068a565b6040516102a99190611d00565b60405180910390f35b6102ba610740565b6040516102c79190611bb3565b60405180910390f35b6102ea60048036038101906102e59190611dbe565b6107d0565b005b61030660048036038101906103019190611f28565b6107e6565b005b610322600480360381019061031d9190611c06565b610803565b60405161032f9190611bb3565b60405180910390f35b610352600480360381019061034d9190611fa8565b610869565b60405161035f9190611b10565b60405180910390f35b5f610372826108f7565b9050919050565b60605f805461038790612013565b80601f01602080910402602001604051908101604052809291908181526020018280546103b390612013565b80156103fe5780601f106103d5576101008083540402835291602001916103fe565b820191905f5260205f20905b8154815290600101906020018083116103e157829003601f168201915b5050505050905090565b5f61041282610970565b5061041c826109f6565b9050919050565b6104358282610430610a2f565b610a36565b5050565b5f600880549050905090565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036104b5575f6040517f64a0ae920000000000000000000000000000000000000000000000000000000081526004016104ac9190611c70565b60405180910390fd5b5f6104c883836104c3610a2f565b610a48565b90508373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161461053e578382826040517f64283d7b00000000000000000000000000000000000000000000000000000000815260040161053593929190612043565b60405180910390fd5b50505050565b5f61054e8361068a565b82106105935782826040517fa57d13dc00000000000000000000000000000000000000000000000000000000815260040161058a929190612078565b60405180910390fd5b60065f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8381526020019081526020015f2054905092915050565b61060283838360405180602001604052805f8152506107e6565b505050565b5f610610610439565b8210610655575f826040517fa57d13dc00000000000000000000000000000000000000000000000000000000815260040161064c929190612078565b60405180910390fd5b600882815481106106695761066861209f565b5b905f5260205f2001549050919050565b5f61068382610970565b9050919050565b5f8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036106fb575f6040517f89c62b640000000000000000000000000000000000000000000000000000000081526004016106f29190611c70565b60405180910390fd5b60035f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b60606001805461074f90612013565b80601f016020809104026020016040519081016040528092919081815260200182805461077b90612013565b80156107c65780601f1061079d576101008083540402835291602001916107c6565b820191905f5260205f20905b8154815290600101906020018083116107a957829003601f168201915b5050505050905090565b6107e26107db610a2f565b8383610a5d565b5050565b6107f1848484610445565b6107fd84848484610bc6565b50505050565b606061080e82610970565b505f610818610d78565b90505f8151116108365760405180602001604052805f815250610861565b8061084084610d8e565b604051602001610851929190612106565b6040516020818303038152906040525b915050919050565b5f60055f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905092915050565b5f7f780e9d63000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161480610969575061096882610e58565b5b9050919050565b5f8061097b83610f39565b90505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036109ed57826040517f7e2732890000000000000000000000000000000000000000000000000000000081526004016109e49190611d00565b60405180910390fd5b80915050919050565b5f60045f8381526020019081526020015f205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b5f33905090565b610a438383836001610f72565b505050565b5f610a54848484611131565b90509392505050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610acd57816040517f5b08ba18000000000000000000000000000000000000000000000000000000008152600401610ac49190611c70565b60405180910390fd5b8060055f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c3183604051610bb99190611b10565b60405180910390a3505050565b5f8373ffffffffffffffffffffffffffffffffffffffff163b1115610d72578273ffffffffffffffffffffffffffffffffffffffff1663150b7a02610c09610a2f565b8685856040518563ffffffff1660e01b8152600401610c2b949392919061217b565b6020604051808303815f875af1925050508015610c6657506040513d601f19601f82011682018060405250810190610c6391906121d9565b60015b610ce7573d805f8114610c94576040519150601f19603f3d011682016040523d82523d5f602084013e610c99565b606091505b505f815103610cdf57836040517f64a0ae92000000000000000000000000000000000000000000000000000000008152600401610cd69190611c70565b60405180910390fd5b805181602001fd5b63150b7a0260e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614610d7057836040517f64a0ae92000000000000000000000000000000000000000000000000000000008152600401610d679190611c70565b60405180910390fd5b505b50505050565b606060405180602001604052805f815250905090565b60605f6001610d9c8461124b565b0190505f8167ffffffffffffffff811115610dba57610db9611e04565b5b6040519080825280601f01601f191660200182016040528015610dec5781602001600182028036833780820191505090505b5090505f82602001820190505b600115610e4d578080600190039150507f3031323334353637383961626364656600000000000000000000000000000000600a86061a8153600a8581610e4257610e41612204565b5b0494505f8503610df9575b819350505050919050565b5f7f80ac58cd000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161480610f2257507f5b5e139f000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916145b80610f325750610f318261139c565b5b9050919050565b5f60025f8381526020019081526020015f205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b8080610faa57505f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b156110dc575f610fb984610970565b90505f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415801561102357508273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614155b801561103657506110348184610869565b155b1561107857826040517fa9fbf51f00000000000000000000000000000000000000000000000000000000815260040161106f9190611c70565b60405180910390fd5b81156110da57838573ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b8360045f8581526020019081526020015f205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050505050565b5f8061113e858585611405565b90505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036111815761117c84611610565b6111c0565b8473ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146111bf576111be8185611654565b5b5b5f73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1603611201576111fc8461179e565b611240565b8473ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161461123f5761123e858561185e565b5b5b809150509392505050565b5f805f90507a184f03e93ff9f4daa797ed6e38ed64bf6a1f01000000000000000083106112a7577a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000838161129d5761129c612204565b5b0492506040810190505b6d04ee2d6d415b85acef810000000083106112e4576d04ee2d6d415b85acef810000000083816112da576112d9612204565b5b0492506020810190505b662386f26fc10000831061131357662386f26fc10000838161130957611308612204565b5b0492506010810190505b6305f5e100831061133c576305f5e100838161133257611331612204565b5b0492506008810190505b612710831061136157612710838161135757611356612204565b5b0492506004810190505b60648310611384576064838161137a57611379612204565b5b0492506002810190505b600a8310611393576001810190505b80915050919050565b5f7f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149050919050565b5f8061141084610f39565b90505f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614611451576114508184866118e2565b5b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146114dc576114905f855f80610f72565b600160035f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825403925050819055505b5f73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff161461155b57600160035f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8460025f8681526020019081526020015f205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550838573ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a4809150509392505050565b60088054905060095f8381526020019081526020015f2081905550600881908060018154018082558091505060019003905f5260205f20015f909190919091505550565b5f61165e8361068a565b90505f60075f8481526020019081526020015f20549050818114611735575f60065f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8481526020019081526020015f205490508060065f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8481526020019081526020015f20819055508160075f8381526020019081526020015f2081905550505b60075f8481526020019081526020015f205f905560065f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8381526020019081526020015f205f905550505050565b5f60016008805490506117b1919061225e565b90505f60095f8481526020019081526020015f205490505f600883815481106117dd576117dc61209f565b5b905f5260205f200154905080600883815481106117fd576117fc61209f565b5b905f5260205f2001819055508160095f8381526020019081526020015f208190555060095f8581526020019081526020015f205f9055600880548061184557611844612291565b5b600190038181905f5260205f20015f9055905550505050565b5f600161186a8461068a565b611874919061225e565b90508160065f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8381526020019081526020015f20819055508060075f8481526020019081526020015f2081905550505050565b6118ed8383836119a5565b6119a0575f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361196157806040517f7e2732890000000000000000000000000000000000000000000000000000000081526004016119589190611d00565b60405180910390fd5b81816040517f177e802f000000000000000000000000000000000000000000000000000000008152600401611997929190612078565b60405180910390fd5b505050565b5f8073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614158015611a5c57508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161480611a1d5750611a1c8484610869565b5b80611a5b57508273ffffffffffffffffffffffffffffffffffffffff16611a43836109f6565b73ffffffffffffffffffffffffffffffffffffffff16145b5b90509392505050565b5f604051905090565b5f80fd5b5f80fd5b5f7fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b611aaa81611a76565b8114611ab4575f80fd5b50565b5f81359050611ac581611aa1565b92915050565b5f60208284031215611ae057611adf611a6e565b5b5f611aed84828501611ab7565b91505092915050565b5f8115159050919050565b611b0a81611af6565b82525050565b5f602082019050611b235f830184611b01565b92915050565b5f81519050919050565b5f82825260208201905092915050565b5f5b83811015611b60578082015181840152602081019050611b45565b5f8484015250505050565b5f601f19601f8301169050919050565b5f611b8582611b29565b611b8f8185611b33565b9350611b9f818560208601611b43565b611ba881611b6b565b840191505092915050565b5f6020820190508181035f830152611bcb8184611b7b565b905092915050565b5f819050919050565b611be581611bd3565b8114611bef575f80fd5b50565b5f81359050611c0081611bdc565b92915050565b5f60208284031215611c1b57611c1a611a6e565b5b5f611c2884828501611bf2565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f611c5a82611c31565b9050919050565b611c6a81611c50565b82525050565b5f602082019050611c835f830184611c61565b92915050565b611c9281611c50565b8114611c9c575f80fd5b50565b5f81359050611cad81611c89565b92915050565b5f8060408385031215611cc957611cc8611a6e565b5b5f611cd685828601611c9f565b9250506020611ce785828601611bf2565b9150509250929050565b611cfa81611bd3565b82525050565b5f602082019050611d135f830184611cf1565b92915050565b5f805f60608486031215611d3057611d2f611a6e565b5b5f611d3d86828701611c9f565b9350506020611d4e86828701611c9f565b9250506040611d5f86828701611bf2565b9150509250925092565b5f60208284031215611d7e57611d7d611a6e565b5b5f611d8b84828501611c9f565b91505092915050565b611d9d81611af6565b8114611da7575f80fd5b50565b5f81359050611db881611d94565b92915050565b5f8060408385031215611dd457611dd3611a6e565b5b5f611de185828601611c9f565b9250506020611df285828601611daa565b9150509250929050565b5f80fd5b5f80fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b611e3a82611b6b565b810181811067ffffffffffffffff82111715611e5957611e58611e04565b5b80604052505050565b5f611e6b611a65565b9050611e778282611e31565b919050565b5f67ffffffffffffffff821115611e9657611e95611e04565b5b611e9f82611b6b565b9050602081019050919050565b828183375f83830152505050565b5f611ecc611ec784611e7c565b611e62565b905082815260208101848484011115611ee857611ee7611e00565b5b611ef3848285611eac565b509392505050565b5f82601f830112611f0f57611f0e611dfc565b5b8135611f1f848260208601611eba565b91505092915050565b5f805f8060808587031215611f4057611f3f611a6e565b5b5f611f4d87828801611c9f565b9450506020611f5e87828801611c9f565b9350506040611f6f87828801611bf2565b925050606085013567ffffffffffffffff811115611f9057611f8f611a72565b5b611f9c87828801611efb565b91505092959194509250565b5f8060408385031215611fbe57611fbd611a6e565b5b5f611fcb85828601611c9f565b9250506020611fdc85828601611c9f565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061202a57607f821691505b60208210810361203d5761203c611fe6565b5b50919050565b5f6060820190506120565f830186611c61565b6120636020830185611cf1565b6120706040830184611c61565b949350505050565b5f60408201905061208b5f830185611c61565b6120986020830184611cf1565b9392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f81905092915050565b5f6120e082611b29565b6120ea81856120cc565b93506120fa818560208601611b43565b80840191505092915050565b5f61211182856120d6565b915061211d82846120d6565b91508190509392505050565b5f81519050919050565b5f82825260208201905092915050565b5f61214d82612129565b6121578185612133565b9350612167818560208601611b43565b61217081611b6b565b840191505092915050565b5f60808201905061218e5f830187611c61565b61219b6020830186611c61565b6121a86040830185611cf1565b81810360608301526121ba8184612143565b905095945050505050565b5f815190506121d381611aa1565b92915050565b5f602082840312156121ee576121ed611a6e565b5b5f6121fb848285016121c5565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61226882611bd3565b915061227383611bd3565b925082820390508181111561228b5761228a612231565b5b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603160045260245ffdfea2646970667358221220e96a3b3bfbc93f98c67d74ea8de94ef895d3ec16d5110497af355b2b8f2c90bf64736f6c63430008150033",
}

// EnumerableABI is the input ABI used to generate the binding from.
// Deprecated: Use EnumerableMetaData.ABI instead.
var EnumerableABI = EnumerableMetaData.ABI

// EnumerableBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EnumerableMetaData.Bin instead.
var EnumerableBin = EnumerableMetaData.Bin

// DeployEnumerable deploys a new Ethereum contract, binding an instance of Enumerable to it.
func DeployEnumerable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Enumerable, error) {
	parsed, err := EnumerableMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EnumerableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Enumerable{EnumerableCaller: EnumerableCaller{contract: contract}, EnumerableTransactor: EnumerableTransactor{contract: contract}, EnumerableFilterer: EnumerableFilterer{contract: contract}}, nil
}

// Enumerable is an auto generated Go binding around an Ethereum contract.
type Enumerable struct {
	EnumerableCaller     // Read-only binding to the contract
	EnumerableTransactor // Write-only binding to the contract
	EnumerableFilterer   // Log filterer for contract events
}

// EnumerableCaller is an auto generated read-only Go binding around an Ethereum contract.
type EnumerableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EnumerableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EnumerableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EnumerableSession struct {
	Contract     *Enumerable       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EnumerableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EnumerableCallerSession struct {
	Contract *EnumerableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// EnumerableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EnumerableTransactorSession struct {
	Contract     *EnumerableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// EnumerableRaw is an auto generated low-level Go binding around an Ethereum contract.
type EnumerableRaw struct {
	Contract *Enumerable // Generic contract binding to access the raw methods on
}

// EnumerableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EnumerableCallerRaw struct {
	Contract *EnumerableCaller // Generic read-only contract binding to access the raw methods on
}

// EnumerableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EnumerableTransactorRaw struct {
	Contract *EnumerableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEnumerable creates a new instance of Enumerable, bound to a specific deployed contract.
func NewEnumerable(address common.Address, backend bind.ContractBackend) (*Enumerable, error) {
	contract, err := bindEnumerable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Enumerable{EnumerableCaller: EnumerableCaller{contract: contract}, EnumerableTransactor: EnumerableTransactor{contract: contract}, EnumerableFilterer: EnumerableFilterer{contract: contract}}, nil
}

// NewEnumerableCaller creates a new read-only instance of Enumerable, bound to a specific deployed contract.
func NewEnumerableCaller(address common.Address, caller bind.ContractCaller) (*EnumerableCaller, error) {
	contract, err := bindEnumerable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EnumerableCaller{contract: contract}, nil
}

// NewEnumerableTransactor creates a new write-only instance of Enumerable, bound to a specific deployed contract.
func NewEnumerableTransactor(address common.Address, transactor bind.ContractTransactor) (*EnumerableTransactor, error) {
	contract, err := bindEnumerable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EnumerableTransactor{contract: contract}, nil
}

// NewEnumerableFilterer creates a new log filterer instance of Enumerable, bound to a specific deployed contract.
func NewEnumerableFilterer(address common.Address, filterer bind.ContractFilterer) (*EnumerableFilterer, error) {
	contract, err := bindEnumerable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EnumerableFilterer{contract: contract}, nil
}

// bindEnumerable binds a generic wrapper to an already deployed contract.
func bindEnumerable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EnumerableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Enumerable *EnumerableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Enumerable.Contract.EnumerableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Enumerable *EnumerableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Enumerable.Contract.EnumerableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Enumerable *EnumerableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Enumerable.Contract.EnumerableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Enumerable *EnumerableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Enumerable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Enumerable *EnumerableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Enumerable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Enumerable *EnumerableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Enumerable.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Enumerable *EnumerableCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Enumerable.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Enumerable *EnumerableSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _Enumerable.Contract.BalanceOf(&_Enumerable.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Enumerable *EnumerableCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _Enumerable.Contract.BalanceOf(&_Enumerable.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Enumerable *EnumerableCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Enumerable.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Enumerable *EnumerableSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _Enumerable.Contract.GetApproved(&_Enumerable.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Enumerable *EnumerableCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _Enumerable.Contract.GetApproved(&_Enumerable.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Enumerable *EnumerableCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _Enumerable.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Enumerable *EnumerableSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _Enumerable.Contract.IsApprovedForAll(&_Enumerable.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Enumerable *EnumerableCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _Enumerable.Contract.IsApprovedForAll(&_Enumerable.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Enumerable *EnumerableCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Enumerable.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Enumerable *EnumerableSession) Name() (string, error) {
	return _Enumerable.Contract.Name(&_Enumerable.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Enumerable *EnumerableCallerSession) Name() (string, error) {
	return _Enumerable.Contract.Name(&_Enumerable.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Enumerable *EnumerableCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Enumerable.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Enumerable *EnumerableSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _Enumerable.Contract.OwnerOf(&_Enumerable.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Enumerable *EnumerableCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _Enumerable.Contract.OwnerOf(&_Enumerable.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Enumerable *EnumerableCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Enumerable.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Enumerable *EnumerableSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Enumerable.Contract.SupportsInterface(&_Enumerable.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Enumerable *EnumerableCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Enumerable.Contract.SupportsInterface(&_Enumerable.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Enumerable *EnumerableCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Enumerable.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Enumerable *EnumerableSession) Symbol() (string, error) {
	return _Enumerable.Contract.Symbol(&_Enumerable.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Enumerable *EnumerableCallerSession) Symbol() (string, error) {
	return _Enumerable.Contract.Symbol(&_Enumerable.CallOpts)
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_Enumerable *EnumerableCaller) TokenByIndex(opts *bind.CallOpts, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Enumerable.contract.Call(opts, &out, "tokenByIndex", index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_Enumerable *EnumerableSession) TokenByIndex(index *big.Int) (*big.Int, error) {
	return _Enumerable.Contract.TokenByIndex(&_Enumerable.CallOpts, index)
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_Enumerable *EnumerableCallerSession) TokenByIndex(index *big.Int) (*big.Int, error) {
	return _Enumerable.Contract.TokenByIndex(&_Enumerable.CallOpts, index)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_Enumerable *EnumerableCaller) TokenOfOwnerByIndex(opts *bind.CallOpts, owner common.Address, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Enumerable.contract.Call(opts, &out, "tokenOfOwnerByIndex", owner, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_Enumerable *EnumerableSession) TokenOfOwnerByIndex(owner common.Address, index *big.Int) (*big.Int, error) {
	return _Enumerable.Contract.TokenOfOwnerByIndex(&_Enumerable.CallOpts, owner, index)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_Enumerable *EnumerableCallerSession) TokenOfOwnerByIndex(owner common.Address, index *big.Int) (*big.Int, error) {
	return _Enumerable.Contract.TokenOfOwnerByIndex(&_Enumerable.CallOpts, owner, index)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Enumerable *EnumerableCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _Enumerable.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Enumerable *EnumerableSession) TokenURI(tokenId *big.Int) (string, error) {
	return _Enumerable.Contract.TokenURI(&_Enumerable.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Enumerable *EnumerableCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _Enumerable.Contract.TokenURI(&_Enumerable.CallOpts, tokenId)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Enumerable *EnumerableCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Enumerable.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Enumerable *EnumerableSession) TotalSupply() (*big.Int, error) {
	return _Enumerable.Contract.TotalSupply(&_Enumerable.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Enumerable *EnumerableCallerSession) TotalSupply() (*big.Int, error) {
	return _Enumerable.Contract.TotalSupply(&_Enumerable.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Enumerable *EnumerableTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Enumerable.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Enumerable *EnumerableSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Enumerable.Contract.Approve(&_Enumerable.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Enumerable *EnumerableTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Enumerable.Contract.Approve(&_Enumerable.TransactOpts, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Enumerable *EnumerableTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Enumerable.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Enumerable *EnumerableSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Enumerable.Contract.SafeTransferFrom(&_Enumerable.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Enumerable *EnumerableTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Enumerable.Contract.SafeTransferFrom(&_Enumerable.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Enumerable *EnumerableTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Enumerable.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Enumerable *EnumerableSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Enumerable.Contract.SafeTransferFrom0(&_Enumerable.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Enumerable *EnumerableTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Enumerable.Contract.SafeTransferFrom0(&_Enumerable.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Enumerable *EnumerableTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _Enumerable.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Enumerable *EnumerableSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Enumerable.Contract.SetApprovalForAll(&_Enumerable.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Enumerable *EnumerableTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Enumerable.Contract.SetApprovalForAll(&_Enumerable.TransactOpts, operator, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Enumerable *EnumerableTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Enumerable.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Enumerable *EnumerableSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Enumerable.Contract.TransferFrom(&_Enumerable.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Enumerable *EnumerableTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Enumerable.Contract.TransferFrom(&_Enumerable.TransactOpts, from, to, tokenId)
}

// EnumerableApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Enumerable contract.
type EnumerableApprovalIterator struct {
	Event *EnumerableApproval // Event containing the contract specifics and raw log

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
func (it *EnumerableApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EnumerableApproval)
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
		it.Event = new(EnumerableApproval)
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
func (it *EnumerableApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EnumerableApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EnumerableApproval represents a Approval event raised by the Enumerable contract.
type EnumerableApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_Enumerable *EnumerableFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*EnumerableApprovalIterator, error) {

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

	logs, sub, err := _Enumerable.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &EnumerableApprovalIterator{contract: _Enumerable.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_Enumerable *EnumerableFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *EnumerableApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _Enumerable.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EnumerableApproval)
				if err := _Enumerable.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_Enumerable *EnumerableFilterer) ParseApproval(log types.Log) (*EnumerableApproval, error) {
	event := new(EnumerableApproval)
	if err := _Enumerable.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EnumerableApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the Enumerable contract.
type EnumerableApprovalForAllIterator struct {
	Event *EnumerableApprovalForAll // Event containing the contract specifics and raw log

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
func (it *EnumerableApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EnumerableApprovalForAll)
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
		it.Event = new(EnumerableApprovalForAll)
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
func (it *EnumerableApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EnumerableApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EnumerableApprovalForAll represents a ApprovalForAll event raised by the Enumerable contract.
type EnumerableApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_Enumerable *EnumerableFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*EnumerableApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Enumerable.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &EnumerableApprovalForAllIterator{contract: _Enumerable.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_Enumerable *EnumerableFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *EnumerableApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Enumerable.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EnumerableApprovalForAll)
				if err := _Enumerable.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_Enumerable *EnumerableFilterer) ParseApprovalForAll(log types.Log) (*EnumerableApprovalForAll, error) {
	event := new(EnumerableApprovalForAll)
	if err := _Enumerable.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EnumerableTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Enumerable contract.
type EnumerableTransferIterator struct {
	Event *EnumerableTransfer // Event containing the contract specifics and raw log

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
func (it *EnumerableTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EnumerableTransfer)
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
		it.Event = new(EnumerableTransfer)
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
func (it *EnumerableTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EnumerableTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EnumerableTransfer represents a Transfer event raised by the Enumerable contract.
type EnumerableTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_Enumerable *EnumerableFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*EnumerableTransferIterator, error) {

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

	logs, sub, err := _Enumerable.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &EnumerableTransferIterator{contract: _Enumerable.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_Enumerable *EnumerableFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *EnumerableTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _Enumerable.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EnumerableTransfer)
				if err := _Enumerable.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_Enumerable *EnumerableFilterer) ParseTransfer(log types.Log) (*EnumerableTransfer, error) {
	event := new(EnumerableTransfer)
	if err := _Enumerable.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
