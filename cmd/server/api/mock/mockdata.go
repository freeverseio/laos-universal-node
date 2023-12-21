package mock

import "encoding/json"

var MockResponseBlock = json.RawMessage(`{
    "difficulty": "0x4ea3f27bc",
    "extraData": "0xd883010a08846765746888676f312e31352e31856c696e7578",
    "gasLimit": "0x47e7c4",
    "gasUsed": "0x0",
    "hash": "0x7c2a5a3e2e7a6d3e2a1f9d6d8e0c7d2a9a6f6f0d6d2f3e3d9d1d4d0d3c3f5f7",
    "logsBloom": "0x000000",
    "miner": "0x0000000000000000000000000000000000000000",
    "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "nonce": "0x0000000000000000",
    "number": "0x29b8ef5",
    "parentHash": "0x879892876fc8821b4e37058a3a7445e5b16b4e344223d1f7c69321ab25b855ea",
        "receiptsRoot": "0xc9f18a6c6b66c1e4873939d1263a943cfbdbf90ef3afa1515b9cb5ae8c58d08e",
        "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
        "size": "0x361",
        "stateRoot": "0x48779a67043e53ad6acac5ec4f46a94fb18b017fb2b7d2b29bb8e654e86bd5e3",
        "timestamp": "0x658148b2",
        "totalDifficulty": "0xff8d4a7",
        "transactions": [
            {
                "blockHash": "0x782a46ff73d68a58153d9df0d309d4a4c5dd48f9802ded44efc6bd6e16677aac",
                "blockNumber": "0x29b8ef5",
                "from": "0xf429a23b511b23a0ddb4415bd7afb327afca3c7d",
                "gas": "0x3f6c0",
                "gasPrice": "0x9c2a0700",
                "maxFeePerGas": "0x9c2a0700",
                "maxPriorityFeePerGas": "0x9c2a0700",
                "hash": "0xeebb76b4f36fe5d1e368dbf6fb62c6f90d0a548d28b9c335f75e646f84377bfc",
                "input": "0xe8eda9df000000000000000000000000169929497847235c5569de26daa042b643f3d3a300000000000000000000000000000000000000000000000821ab0d4414980000000000000000000000000000f429a23b511b23a0ddb4415bd7afb327afca3c7d0000000000000000000000000000000000000000000000000000000000000000",
                "nonce": "0x81",
                "to": "0x14cac28d3b8ce63f362ca6a7354c4b13d6b561ec",
                "transactionIndex": "0x0",
                "value": "0x0",
                "type": "0x2",
                "accessList": [],
                "chainId": "0x13881",
                "v": "0x0",
                "r": "0x3412f2ed3a7876380e43624cb983ba71e631cd1647752f74c44e98ddb539d163",
                "s": "0x51be04e78fe2c4c7c13c90349ad40286a88c8be297c096c4bb6a31e59b87800f"
            }
        ],
        "transactionsRoot": "0x6ea17a28615cd164720317abce24e1ab140400fc42dc7afbb8cf836e1d9550f0",
        "uncles": []
        }`)

var MockResponseTransaction = json.RawMessage(`{
          "blockHash": "0x782a46ff73d68a58153d9df0d309d4a4c5dd48f9802ded44efc6bd6e16677aac",
            "blockNumber": "0x29b8ef5",
            "from": "0xf429a23b511b23a0ddb4415bd7afb327afca3c7d",
            "gas": "0x3f6c0",
            "gasPrice": "0x9c2a0700",
            "maxFeePerGas": "0x9c2a0700",
            "maxPriorityFeePerGas": "0x9c2a0700",
            "hash": "0xeebb76b4f36fe5d1e368dbf6fb62c6f90d0a548d28b9c335f75e646f84377bfc",
            "input": "0xe8eda9df000000000000000000000000169929497847235c5569de26daa042b643f3d3a300000000000000000000000000000000000000000000000821ab0d4414980000000000000000000000000000f429a23b511b23a0ddb4415bd7afb327afca3c7d0000000000000000000000000000000000000000000000000000000000000000",
            "nonce": "0x81",
            "to": "0x14cac28d3b8ce63f362ca6a7354c4b13d6b561ec",
            "transactionIndex": "0x0",
            "value": "0x0",
            "type": "0x2",
            "accessList": [],
            "chainId": "0x13881",
            "v": "0x0",
            "r": "0x3412f2ed3a7876380e43624cb983ba71e631cd1647752f74c44e98ddb539d163",
            "s": "0x51be04e78fe2c4c7c13c90349ad40286a88c8be297c096c4bb6a31e59b87800f"
        }`)

func GetFilterRequest(fromBlock, toBlock string) json.RawMessage {
	return json.RawMessage(`{
  "fromBlock": "` + fromBlock + `",
  "toBlock": "` + toBlock + `",
  "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
  "topics": [
    ["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"],
    null,
    [
      "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
      "0x0000000000000000000000000aff3454fce5edbc8cca8697c15331677e6ebccc"
    ]
  ]
}`)
}

func GetLogsRequest(fromBlock, toBlock string) json.RawMessage {
	return json.RawMessage(`{
  "fromBlock": "` + fromBlock + `",
  "toBlock": "` + toBlock + `",
  "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
  "address": "0x8888f1f195afa192cfee860698584c030f4c9db1",
  "topics": [
          "0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"
  ]
  }`)
}
