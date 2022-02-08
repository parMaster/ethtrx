## Ethereum Transactions Explorer
Educational Project - Junior Golang Developer Test Problem

## Sandboxed demo

Demo server is deployed, API endpoint is available at `http://cdns.com.ua:8088/getTxList`:

`curl -X POST http://cdns.com.ua:8088/getTxList -H 'Content-Type: application/json' -d '{\"date\":\"2022-02-07\"}'`

## Prerequisites
- Access to Mongo (database=etc, collection=transactions). Sandboxed demo is running docker container with Mongo 4.2 like this:
	`docker run --name mongodb -p27018:27017 mongo:4.2`
- Set correct Mongo connection URI in `configs/apiserver.toml` 
- Set correct Etherscan API Key in `configs/apiserver.toml`

## Using API

`getTxList` API provides a json-encoded transactions list

### Basic request:

`POST /getTxList HTTP/1.1`

```
{
    "transactions": [
        {
            "hash": "0x4f65b92e705c450988557995b60b77978a9fde743a22efba07dd3206357e3c1d",
            "blockHash": "0xd1452f4505b203cbcfb897e03c5ad134f824fbb35ace4e83bb0a3722cd9e2c8b",
            "blockNumber": "0xd80d4f",
            "blockHeight": 14159183,
            "blockTime": "0x62011e67",
            "blockDate": "2022-02-07",
            "from": "0x00dee1f836998bcc736022f314df906588d44808",
            "to": "0x4a137fd5e7a256ef08a7de531a17d0be0cc7b6b6",
            "value": "0x10dd27a770322c",
            "valueNumber": "0.004746762009981484",
            "confirmations": 356,
            "type": "0x2",
            "chainId": "0x1"
        },
        {
            "hash": "0xd6e86d336ae817a85a9492f88d493edc9a9bd509fcd0fbdb4bea686df4905809",
            "blockHash": "0xd1452f4505b203cbcfb897e03c5ad134f824fbb35ace4e83bb0a3722cd9e2c8b",
            "blockNumber": "0xd80d4f",
            "blockHeight": 14159183,
            "blockTime": "0x62011e67",
            "blockDate": "2022-02-07",
            "from": "0xe92f359e6f05564849afa933ce8f62b8007a1d5d",
            "to": "0x9008d19f58aabd9ed0d60971565aa8510560ab41",
            "value": "0x0",
            "valueNumber": "0.",
            "confirmations": 356,
            "type": "0x2",
            "chainId": "0x1"
        }
    ]
}
```

Identical to the following request with `page` and `limit` parameters:

`POST /getTxList -d '{"page":1, "limit":2}' HTTP/1.1`

By default `page` and `limit` are set to `1` and `10` accordingly

### Additional filters:

#### Filter transactions by hash

`POST /getTxList -d '{"txhash":"0x90db94375b973f4046ba991e8a326968f3035ff12c53f4891d010ab2a3850c50"}' HTTP/1.1`

```
{
    "transactions": [
        {
            "hash": "0x90db94375b973f4046ba991e8a326968f3035ff12c53f4891d010ab2a3850c50",
            "blockHash": "0xd1452f4505b203cbcfb897e03c5ad134f824fbb35ace4e83bb0a3722cd9e2c8b",
            "blockNumber": "0xd80d4f",
            "blockHeight": 14159183,
            "blockTime": "0x62011e67",
            "blockDate": "2022-02-07",
            "from": "0x9162c9a00a8e4c3654a4e3f236e8e2bf097c9040",
            "to": "0x9162c9a00a8e4c3654a4e3f236e8e2bf097c9040",
            "value": "0x0",
            "valueNumber": "0.",
            "confirmations": 356,
            "type": "0x2",
            "chainId": "0x1"
        }
    ]
}
```

#### Filter transactions by block number in decimal

`POST /getTxList -d '{"blockheight":14159183}' HTTP/1.1`

Identical to:

`POST /getTxList -d '{"blocknumber":"0x62011e67"}' HTTP/1.1`


```
{
    "transactions": [
        {
            "hash": "0x4f65b92e705c450988557995b60b77978a9fde743a22efba07dd3206357e3c1d",
            "blockHash": "0xd1452f4505b203cbcfb897e03c5ad134f824fbb35ace4e83bb0a3722cd9e2c8b",
            "blockNumber": "0xd80d4f",
            "blockHeight": 14159183,
            "blockTime": "0x62011e67",
            "blockDate": "2022-02-07",
            "from": "0x00dee1f836998bcc736022f314df906588d44808",
            "to": "0x4a137fd5e7a256ef08a7de531a17d0be0cc7b6b6",
            "value": "0x10dd27a770322c",
            "valueNumber": "0.004746762009981484",
            "confirmations": 356,
            "type": "0x2",
            "chainId": "0x1"
        },
... rest of response body ommited
}
```

#### Filter transactions by receiver

`POST /getTxList -d '{"to":"0x9162c9a00a8e4c3654a4e3f236e8e2bf097c9040"}' HTTP/1.1`

#### Filter transactions by sender

`POST /getTxList -d '{"from":"0x9162c9a00a8e4c3654a4e3f236e8e2bf097c9040"}' HTTP/1.1`


```
{
    "transactions": [
        {
            "hash": "0x90db94375b973f4046ba991e8a326968f3035ff12c53f4891d010ab2a3850c50",
            "blockHash": "0xd1452f4505b203cbcfb897e03c5ad134f824fbb35ace4e83bb0a3722cd9e2c8b",
            "blockNumber": "0xd80d4f",
            "blockHeight": 14159183,
            "blockTime": "0x62011e67",
            "blockDate": "2022-02-07",
            "from": "0x9162c9a00a8e4c3654a4e3f236e8e2bf097c9040",
            "to": "0x9162c9a00a8e4c3654a4e3f236e8e2bf097c9040",
            "value": "0x0",
            "valueNumber": "0.",
            "confirmations": 356,
            "type": "0x2",
            "chainId": "0x1"
        }
    ]
}
```

#### Filter transactions by date

`POST /getTxList -d '{"date":"2022-02-07"}' HTTP/1.1`


```
{
    "transactions": [
        {
            "hash": "0x90db94375b973f4046ba991e8a326968f3035ff12c53f4891d010ab2a3850c50",
            "blockHash": "0xd1452f4505b203cbcfb897e03c5ad134f824fbb35ace4e83bb0a3722cd9e2c8b",
            "blockNumber": "0xd80d4f",
            "blockHeight": 14159183,
            "blockTime": "0x62011e67",
            "blockDate": "2022-02-07",
            "from": "0x9162c9a00a8e4c3654a4e3f236e8e2bf097c9040",
            "to": "0x9162c9a00a8e4c3654a4e3f236e8e2bf097c9040",
            "value": "0x0",
            "valueNumber": "0.",
            "confirmations": 356,
            "type": "0x2",
            "chainId": "0x1"
        }
    ]
}
```

#### Filters can be combined

`POST /getTxList -d '{"blockheight":14159183, "from":"0x9162c9a00a8e4c3654a4e3f236e8e2bf097c9040"}' HTTP/1.1`

Will apply both conditions `blockheight == 14159183 && from == "0x9162c9a00a8e4c3654a4e3f236e8e2bf097c9040"`

#### Empty result
Looks like this
```
{
    "transactions": null
}
```

#### Errors

`POST /getTxList -d '{"txhash":"0x9"}' HTTP/1.1`

```
{
    "error": "Invalid 'txhash' parameter: hex string of odd length"
}
```


`POST /getTxList -d '{"page":"q"}' HTTP/1.1`

```
{
     "error": "Error decoding request: json: cannot unmarshal string into Go struct field request.page of type int"
}
```



`POST /getTxList -d '{"to":"0xdb298285fe"}' HTTP/1.1`

`POST /getTxList -d '{"from":"0xdb298285fe"}' HTTP/1.1`

```
{
    "error": "'to' parameter is not a valid hex address"
}
```

etc.


### (!) Please, note:
Calculating transaction _comission_ requires the `Gas User` value, which is provided by _eth_getTransactionReceipt_ API for each transaction. Considering the fact Ethereum blockchain serves more that 1 million transactions per day, it would require to call the receipt API more than 11 times/sec, which is imposible within 5 calls/sec using a Free _etherscan_ plan.

Anyway, commission formula is:

`(Block Base Fee Per Gas + Max Priority Fee Per Gas) * Gas Used`