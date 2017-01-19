
# Hyperledger POC


## Commands

### Registrar
```
POST localhost:7050/registrar
{
  "enrollId": "lukas",
  "enrollSecret": "NPKYL39uKbkj"
}
```

### Deploy chaincode
```
POST localhost:7050/chaincode
{
  "jsonrpc": "2.0",
  "method": "deploy",
  "params": {
    "chaincodeID":{
        "name":"mycc"
    },
"ctorMsg": {
        "args":["init"]
    },
    "secureContext": "lukas"
  },
  "id": "1"  
}
```

### Create supply chain

```
POST localhost:7050/chaincode
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "chaincodeID":{
        "name":"mycc"
    },
"ctorMsg": {
        "args":["create_supply_chain", "assetid1", "@satoshi", "tuna", "3", "1000", "Ohma", "0", "5" ]
    },
    "secureContext": "lukas"
  },
  "id": "2"  
}
```

### Start trade

```
POST localhost:7050/chaincode
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "chaincodeID":{
        "name":"mycc"
    },
"ctorMsg": {
        "args":["start_trade", "assetid1", "", "4", "10" ]
    },
    "secureContext": "lukas"
  },
  "id": "4"  
}
```

### Complete trade

```
POST localhost:7050/chaincode
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "chaincodeID":{
        "name":"mycc2"
    },
"ctorMsg": {
        "args":["complete_trade", "c5d13807-e9f4-4037-a026-ef82b737a5aa", "@miguel", "Tokio"]
    },
    "secureContext": "lukas"
  },
  "id": "5"  
}

### Query assets

```
POST localhost:7050/chaincode
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
      "type": 1,
      "chaincodeID":{
          "name":"mycc2"
      },
      "ctorMsg": {
         "args":["query_asset", "assetid1"]
      },
      "secureContext": "lukas"
  },
  "id": 3
}
```

### Query one contract 

```
POST localhost:7050/chaincode
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
      "type": 1,
      "chaincodeID":{
          "name":"mycc2"
      },
      "ctorMsg": {
         "args":["query_one_contract", "{transactionID}"]
      },
      "secureContext": "lukas"
  },
  "id": 3
}
```

