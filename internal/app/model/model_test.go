package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Making models, playing with models and unmarshalling
func Test_Unmarshal(t *testing.T) {
	shortBlock := `
			{
				"jsonrpc":"2.0",
				"id":1,
				"result":{
				   "baseFeePerGas":"0x5cfe76044",
				   "difficulty":"0x1b4ac252b8a531",
				   "extraData":"0xd883010a06846765746888676f312e31362e36856c696e7578",
				   "timestamp":"0x6110bab2",
				   "totalDifficulty":"0x612789b0aba90e580f8",
				   "transactions":[
					  "0x40330c87750aa1ba1908a787b9a42d0828e53d73100ef61ae8a4d925329587b5",
					  "0x6fa2208790f1154b81fc805dd7565679d8a8cc26112812ba1767e1af44c35dd4",
					  "0xe31d8a1f28d4ba5a794e877d65f83032e3393809686f53fa805383ab5c2d3a3c",
					  "0xa6a83df3ca7b01c5138ec05be48ff52c7293ba60c839daa55613f6f1c41fdace",
					  "0x4e46edeb68a62dde4ed081fae5efffc1fb5f84957b5b3b558cdf2aa5c2621e17",
					  "0x356ee444241ae2bb4ce9f77cdbf98cda9ffd6da244217f55465716300c425e82",
					  "0x1a4ec2019a3f8b1934069fceff431e1370dcc13f7b2561fe0550cc50ab5f4bbc",
					  "0xad7994bc966aed17be5d0b6252babef3f56e0b3f35833e9ac414b45ed80dac93"
				   ],
				   "transactionsRoot":"0xaceb14fcf363e67d6cdcec0d7808091b764b4428f5fd7e25fb18d222898ef779",
				   "uncles":[
					  "0x9e8622c7bf742bdeaf96c700c07151c1203edaf17a38ea8315b658c2e6d873cd"
				   ]
				}
			}
			`

	var cont BlockResponse
	err := json.Unmarshal([]byte(shortBlock), &cont)
	assert.Error(t, err) // not a full block, expecting error

	// fmt.Printf("%+v\n", cont)
	// fmt.Printf("%+v\n", cont.Transactions)
	// fmt.Printf("\n\n\n")

	jStr2 := `
	{
		"jsonrpc":"2.0",
		"id":1,
		"result":{
			"blockHash":"0xf850331061196b8f2b67e1f43aaa9e69504c059d3d3fb9547b04f9ed4d141ab7",
			"blockNumber":"0xcf2420",
			"from":"0x00192fb10df37c9fb26829eb2cc623cd1bf599e8",
			"gas":"0x5208",
			"gasPrice":"0x19f017ef49",
			"maxFeePerGas":"0x1f6ea08600",
			"maxPriorityFeePerGas":"0x3b9aca00",
			"hash":"0xbc78ab8a9e9a0bca7d0321a27b2c03addeae08ba81ea98b03cd3dd237eabed44",
			"input":"0x",
			"nonce":"0x33b79d",
			"to":"0xc67f4e626ee4d3f272c2fb31bad60761ab55ed9f",
			"transactionIndex":"0x5b",
			"value":"0x19755d4ce12c00",
			"type":"0x2",
			"accessList":[

			],
			"chainId":"0x1",
			"v":"0x0",
			"r":"0xa681faea68ff81d191169010888bbbe90ec3eb903e31b0572cd34f13dae281b9",
			"s":"0x3f59b0fa5ce6cf38aff2cfeb68e7a503ceda2a72b4442c7e2844d63544383e3"
		}
	}
	`

	var tx TransactionResponse
	err = json.Unmarshal([]byte(jStr2), &tx)
	assert.NoError(t, err)

	// Block with transactions
	jStr3 := `
		{
			"jsonrpc": "2.0",
			"id": 1,
			"result": {
			  "baseFeePerGas": "0x37ff479c0a",
			  "difficulty": "0x2c97d85bf4b231",
			  "extraData": "0x4555312f706f6f6c2e62696e616e63652e636f6d2f00126cf5",
			  "gasLimit": "0x1c9c346",
			  "gasUsed": "0x1c9a0c1",
			  "hash": "0x9a412290aeca44bdbcf4fc7bd541afdf22225c7c570873cc0e032770ac079846",
			  "logsBloom": "0xef3eb1c9ffcf77b798ef7effd7ec3f3f06f9f985dffd3b58dacb7ffbbbfa7d6dadf73f7feefb67b456f07cfdddbd13ed7ef3ef86ddeafbea8797f8cffa75df757e4d6fd65f9fedf9099ff3bb33bf5d3e3e86fbfcfb67bb6c746af5d3abe520491bf0f1a3ef3cd9ef7db762bfefed5ec9e65be27947fead4f57e0f793db8e0bf93ceadbaf7ffff7dfb7cf579063d48afffd8fee93fd695ffeb67bfcd7aef8377fc3afbf9ffbc875fbeefccfe17eefd933d67f4be5afb9ffa17eeef7f43deedc7679373b735bbfb696fbc5d8f4fbb1b9ffbb7e9aa7dd5abed93f3f3b8751f8beb5fabf6bffef7f41ae8ee607ae775bedb2e3ed4e68f3ddc75f66cff333e853dae5",
			  "miner": "0xc365c3315cf926351ccaf13fa7d19c8c4058c8e1",
			  "mixHash": "0x43aeb90efa71579c5fe7c02553b19c7700aea42b69094b68b0df1036d21880ad",
			  "nonce": "0xd8bd5ad05e94b70c",
			  "number": "0xd7c548",
			  "parentHash": "0x1627f062a52e16ed675425455b1b59cbfa8a61c00b4d9c916387a62149f23d37",
			  "receiptsRoot": "0x1f6ea43a0554c4da0e8c87c466cc75465aa87fa1a58ab2bb7945f49b950022cc",
			  "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			  "size": "0x1427b",
			  "stateRoot": "0x01af94f188862f78750abfc5784899624ed7809ab1f1713db0e875c788a0efdb",
			  "timestamp": "0x61fd5bd7",
			  "totalDifficulty": "0x89b92da2384a27d4e6c",
			  "transactions": [
				{
				  "blockHash": "0x9a412290aeca44bdbcf4fc7bd541afdf22225c7c570873cc0e032770ac079846",
				  "blockNumber": "0xd7c548",
				  "from": "0xf07704777d6bc182bf2c67fbda48913169b84983",
				  "gas": "0x33450",
				  "gasPrice": "0x4b3a7a7339",
				  "maxFeePerGas": "0x681346b600",
				  "maxPriorityFeePerGas": "0x133b32d72f",
				  "hash": "0xa71509e636639a1bf7094ec059ca3c8a536592b111648f341cf6b589cb4ca442",
				  "input": "0x8803dbee00000000000000000000000000000000000000000000000089c96f1771d9d2280000000000000000000000000000000000000000000394de5cfa99899f1da17c00000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000f07704777d6bc182bf2c67fbda48913169b849830000000000000000000000000000000000000000000000000000000061fd5c2a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000090185f2135308bad17527004364ebcc2d37e5f6000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
				  "nonce": "0x9c83",
				  "to": "0xd9e1ce17f2641f24ae83637ab66a2cca9c378b9f",
				  "transactionIndex": "0x0",
				  "value": "0x0",
				  "type": "0x2",
				  "accessList": [],
				  "chainId": "0x1",
				  "v": "0x0",
				  "r": "0xfd4e458233c47302ef2d43a7866416b2f61365c027f0a2c0930064c1c3bed273",
				  "s": "0x226b16db51b0e97308503e0e12a6e42326cc4def99e097955405aaf7fb3a465d"
				},
				{
				  "blockHash": "0x9a412290aeca44bdbcf4fc7bd541afdf22225c7c570873cc0e032770ac079846",
				  "blockNumber": "0xd7c548",
				  "from": "0xcfdfe45b1d80c318d6e7d77c0cc6a7f529ed0230",
				  "gas": "0x388e5",
				  "gasPrice": "0x8bb2c97000",
				  "maxFeePerGas": "0x8bb2c97000",
				  "maxPriorityFeePerGas": "0x8bb2c97000",
				  "hash": "0xfac31f5bbc2a4a271e543ff39e2a292ac917a9b7c32db1237146322d5f600193",
				  "input": "0xa0712d680000000000000000000000000000000000000000000000000000000000000006",
				  "nonce": "0x4d",
				  "to": "0x2dddae5c2c27ae6c1751cac72adacfe82a60f8a4",
				  "transactionIndex": "0x1",
				  "value": "0x4fefa17b7240000",
				  "type": "0x2",
				  "accessList": [],
				  "chainId": "0x1",
				  "v": "0x0",
				  "r": "0x8b72fbdc2c48af58d8c62702ad1990f6af29f9474b37a20eb8b5e9f7bad3452b",
				  "s": "0x6a0945c6630e085734c8e938c6156897430e9793a554dbf8b94cf7a71c4f5ff1"
				},
				{
				  "blockHash": "0x9a412290aeca44bdbcf4fc7bd541afdf22225c7c570873cc0e032770ac079846",
				  "blockNumber": "0xd7c548",
				  "from": "0x76d4ec21be60d595bf80c945861afd6b384de40e",
				  "gas": "0x388e5",
				  "gasPrice": "0x8587a05e0a",
				  "maxFeePerGas": "0x9b10b18400",
				  "maxPriorityFeePerGas": "0x4d8858c200",
				  "hash": "0x979bc875b0f9d0ed227573f4ae5fefc314d23b0b86a13e2290eddda507def5b1",
				  "input": "0xa0712d680000000000000000000000000000000000000000000000000000000000000006",
				  "nonce": "0x1c",
				  "to": "0x2dddae5c2c27ae6c1751cac72adacfe82a60f8a4",
				  "transactionIndex": "0x2",
				  "value": "0x4fefa17b7240000",
				  "type": "0x2",
				  "accessList": [],
				  "chainId": "0x1",
				  "v": "0x0",
				  "r": "0xc68b7992c7977689f08d3a2744efc024dbdf21422fcb5bf4cdec965d54e1b0e2",
				  "s": "0x476d41ca9d40ff3bd96dc89bc02bcfc93b26fd7080c68e6a881cfc4a12dc452c"
				},
				{
				  "blockHash": "0x9a412290aeca44bdbcf4fc7bd541afdf22225c7c570873cc0e032770ac079846",
				  "blockNumber": "0xd7c548",
				  "from": "0x5da54d44a193979243bf73dde653d8ead953c680",
				  "gas": "0x3db29",
				  "gasPrice": "0x8587a05e0a",
				  "maxFeePerGas": "0xcec0ecb000",
				  "maxPriorityFeePerGas": "0x4d8858c200",
				  "hash": "0xdfc71a47607abefc32875261368844aa85afd1c19a54cf7f8d8be88e826fdec1",
				  "input": "0xa0712d680000000000000000000000000000000000000000000000000000000000000006",
				  "nonce": "0x3e",
				  "to": "0x2dddae5c2c27ae6c1751cac72adacfe82a60f8a4",
				  "transactionIndex": "0x3",
				  "value": "0x4fefa17b7240000",
				  "type": "0x2",
				  "accessList": [],
				  "chainId": "0x1",
				  "v": "0x1",
				  "r": "0x6dd0f5660379cb262bfc2ef072fac2c2053df6d139d30897ef206140c6fb940d",
				  "s": "0x506187624fef7b34ae518743ec794a2f66d498625bb995735338e2ee8e84e7f1"
				},
				{
				  "blockHash": "0x9a412290aeca44bdbcf4fc7bd541afdf22225c7c570873cc0e032770ac079846",
				  "blockNumber": "0xd7c548",
				  "from": "0x53c787bbb6c6e7649828c79bb237b8128fd92a6a",
				  "gas": "0x3db29",
				  "gasPrice": "0x8587a05e0a",
				  "maxFeePerGas": "0xcec0ecb000",
				  "maxPriorityFeePerGas": "0x4d8858c200",
				  "hash": "0xb5ea42c5b2e441e1486bd97086e2d748b5f84ed82afd0faba24b98870a9bdf88",
				  "input": "0xa0712d680000000000000000000000000000000000000000000000000000000000000006",
				  "nonce": "0x3b",
				  "to": "0x2dddae5c2c27ae6c1751cac72adacfe82a60f8a4",
				  "transactionIndex": "0x4",
				  "value": "0x4fefa17b7240000",
				  "type": "0x2",
				  "accessList": [],
				  "chainId": "0x1",
				  "v": "0x1",
				  "r": "0x54925223be6df99648b2308b4bc206da5432c71e1327f1f05ac32260aff04a34",
				  "s": "0x781763900dbee8470bebc76a24dfb22cf0fc14b7cf3a63f012dc50164a9830f5"
				}
			  ],
			  "transactionsRoot": "0x738602f360f1b4a472a2744eda59de0b13c9aa96be867a347533894ccfd8546a",
			  "uncles": []
			}
		  }
		`

	var contFull BlockResponse
	err = json.Unmarshal([]byte(jStr3), &contFull)
	assert.NoError(t, err)
	// fmt.Printf("%+v\n", contFull)
	// fmt.Printf("%+v\n", contFull.Transactions)

	for _, v := range contFull.Transactions {
		// fmt.Printf("%+v\n\n", v)
		assert.NotNil(t, v)
	}

}
