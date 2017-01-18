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

type FishContract struct {
	ID       string
	Seller   string
	FishName string
	Price    uint64
	Weight   uint64
}

// Init method will be called during deployment.
// args: contractId(string), seller(string), fishName(uint), price(uint), weight(uint)
// TODO: confirm what seller value should be for later tracking. ECA/ TCA etc.
func (t *TestContractChainCode) Init(stub shim.ChaincodeStubInterface, methodName string, args []string) ([]byte, error) {
	myLogger.Debug("Init Chaincode...")
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	contractId := args[0]
	seller := args[1]
	fishName := args[2]
	price, err := strconv.ParseUint(args[3], 10, 64)
	if err != nil {
		return nil, errors.New("Failed to parse price.")
	}
	weight, err := strconv.ParseUint(args[4], 10, 64)
	if err != nil {
		return nil, errors.New("Failed to parse weight.")
	}
	contract := FishContract{
		ID:       contractId,
		Seller:   seller,
		FishName: fishName,
		Price:    price,
		Weight:   weight,
	}
	b, err := json.Marshal(contract)

	stub.PutState("contracts/"+contractId, b)

	myLogger.Debug("Init Chaincode...done")

	return nil, nil
}

func (t *TestContractChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// TODO: implement it.
	return nil, nil
}

// args: contractId(string)
func (t *TestContractChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	myLogger.Debug("Query Chaincode...")
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	contractId := args[0]
	state, err := stub.GetState("contracts/" + contractId)
	return state, err
}

func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(TestContractChainCode))
	if err != nil {
		fmt.Printf("Error starting TestContractChainCode: %s", err)
	}
}
