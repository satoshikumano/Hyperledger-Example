/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	"github.com/op/go-logging"
)

var myLogger = logging.MustGetLogger("test_contract")

// TODO: review the explanation later.
// TestContractChainCode is simple chaincode implementing a basic Contract
// https://github.com/hyperledger/fabric/blob/master/docs/tech/application-ACL.md
type TestContractChainCode struct {
}

type Asset struct {
	ID             string
	Owner          string
	FishName       string
	Weight         uint64
	MinTemperature uint64
	MaxTemperature uint64
	Price          uint64
	Location       string
}

type Contract struct {
	AssetID                  string
	PreviousTxId             string
	AcceptableMinTemperature uint64
	AcceptableMaxTemperature uint64
	Location                 string
	Completed                bool
}

// Init method will be called during deployment.
// args: contractId(string), seller(string), fishName(uint), price(uint), weight(uint)
// TODO: confirm what seller value should be for later tracking. ECA/ TCA etc.
func (t *TestContractChainCode) Init(stub shim.ChaincodeStubInterface, methodName string, args []string) ([]byte, error) {
	myLogger.Debug("Init Chaincode...done")
	return nil, nil
}

func (t *TestContractChainCode) start_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// TODO: only the owner of asset can execute this function.
	assetId := args[0]
	previousTxId := args[1]
	acceptableMinTemp, err := strconv.ParseUint(args[3], 10, 64)
	if err != nil {
		return nil, errors.New("Failed to parse minTemp.")
	}
	acceptableMaxTemp, err := strconv.ParseUint(args[3], 10, 64)
	if err != nil {
		return nil, errors.New("Failed to parse maxTemp.")
	}
	contract := Contract{
		AssetID:                  assetId,
		PreviousTxId:             previousTxId,
		AcceptableMinTemperature: acceptableMinTemp,
		AcceptableMaxTemperature: acceptableMaxTemp,
		Completed:                false,
		Location:                 "",
	}
	b, err := json.Marshal(contract)

	txId := stub.GetTxID()
	stub.PutState("contract/"+txId, b)

	return nil, nil
}
func (t *TestContractChainCode) complete_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//TODO: send money to seller
	//TODO: get the owner from the caller
	txId := args[0]

	var contract Contract
	contract_json, err := stub.GetState("contract/" + txId)
	err = json.Unmarshal(contract_json, &contract)
	if err != nil {
		return nil, errors.New("contract id not valid")
	}

	var asset Asset
	asset_json, err1 := stub.GetState("asset/" + contract.AssetID)
	if err1 != nil {
		return nil, errors.New("asset not found")
	}
	err1 = json.Unmarshal(asset_json, &asset)
	if err1 != nil {
		return nil, errors.New("failed to unmarshal asset")
	}

	if asset.MaxTemperature > contract.AcceptableMaxTemperature {
		myLogger.Debug("WARNING: Acceptable max temperature is below the asset temperature")
		return nil, errors.New("max temperature warning")
	}
	if asset.MinTemperature < contract.AcceptableMinTemperature {
		myLogger.Debug("WARNING: Acceptable min temperature is above the asset temperature")
		return nil, errors.New("min temperature warning")
	}
	contract.Completed = true
	contract.Location = args[2]

	c, err := json.Marshal(contract)
	stub.PutState("contract/"+txId, c)

	asset.Owner = args[1]
	a, err := json.Marshal(asset)
	stub.PutState("asset/"+contract.AssetID, a)

	return nil, nil
}

func (t *TestContractChainCode) create_supply_chain(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 8")
	}

	assetId := args[0]
	owner := args[1]
	fishName := args[2]
	price, err := strconv.ParseUint(args[3], 10, 64)
	if err != nil {
		return nil, errors.New("Failed to parse price.")
	}
	weight, err := strconv.ParseUint(args[4], 10, 64)
	if err != nil {
		return nil, errors.New("Failed to parse weight.")
	}

	location := args[5]

	minTemp, err := strconv.ParseUint(args[6], 10, 64)
	if err != nil {
		return nil, errors.New("Failed to parse minTemp.")
	}
	maxTemp, err := strconv.ParseUint(args[7], 10, 64)
	if err != nil {
		return nil, errors.New("Failed to parse maxTemp.")
	}

	asset := Asset{
		ID:             assetId,
		Owner:          owner,
		FishName:       fishName,
		Price:          price,
		Weight:         weight,
		MinTemperature: minTemp,
		MaxTemperature: maxTemp,
		Location:       location,
	}
	b, err := json.Marshal(asset)

	stub.PutState("asset/"+assetId, b)

	return nil, nil
}

func (t *TestContractChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "create_supply_chain" {
		return t.create_supply_chain(stub, args)
	}

	if function == "start_trade" {
		return t.start_trade(stub, args)
	}

	if function == "complete_trade" {
		return t.complete_trade(stub, args)
	}
	return nil, nil
}

// args: transaction id

//result:
//print all the previous transactions data structure

func (t *TestContractChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	myLogger.Debug("Query Chaincode...")
	if function == "query_asset" {
		return t.query_asset(stub, args)
	}
	if function == "query_one_contract" {
		return t.query_one_contract(stub, args)
	}
	if function == "query_contract_ancestors" {
		return t.query_contract_ancestors(stub, args)
	}
	return nil, nil
}

func (t *TestContractChainCode) query_asset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	assetId := args[0]
	asset_bytes, err := stub.GetState("asset/" + assetId)
	if err != nil {
		return nil, errors.New("Asset not found.")
	}
	return asset_bytes, nil
}

func (t *TestContractChainCode) query_one_contract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	txId := args[0]
	contract_bytes, err := stub.GetState("contract/" + txId)
	if err != nil {
		return nil, errors.New("Contract not found.")
	}
	return contract_bytes, nil
}

func (t *TestContractChainCode) get_contract(stub shim.ChaincodeStubInterface, contractId string) (Contract, error) {
	var contract Contract
	contract_bytes, err := stub.GetState("contract/" + contractId)
	if err != nil {
		return contract, errors.New("Contract not found.")
	}
	err = json.Unmarshal(contract_bytes, &contract)
	if err != nil {
		return contract, errors.New("Failed to unmarshal contract")
	}
	return contract, nil
}

func (t *TestContractChainCode) query_contract_ancestors(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var contracts []Contract
	var contract Contract
	var prevTxId string
	contract, err := t.get_contract(stub, args[0])
	if err != nil {
		return nil, errors.New("Contract not found.")
	}
	contracts = append(contracts, contract)
	prevTxId = contract.PreviousTxId
	for prevTxId != "" {
		contract, err := t.get_contract(stub, prevTxId)
		if err != nil {
			break
		}
		prevTxId = contract.PreviousTxId
		contracts = append(contracts, contract)
	}
	b, err := json.Marshal(contracts)
	if err != nil {
		return nil, errors.New("Failed to marshal contracts")
	}
	return b, nil
}

func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(TestContractChainCode))
	if err != nil {
		fmt.Printf("Error starting TestContractChainCode: %s", err)
	}
}
