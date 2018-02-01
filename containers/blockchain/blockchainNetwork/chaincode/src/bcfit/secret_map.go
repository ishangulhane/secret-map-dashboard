/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//steps to fitcoin constant
const STEPS_TO_FITCOIN = 100

//contract state
const STATE_COMPLETE = "complete"
const STATE_PENDING = "pending"
const STATE_DECLINED = "declined"

//member type
const TYPE_USER = "user"
const TYPE_SELLER = "seller"

//keys for key-value store
const USERS_KEY = "users"
const SELLERS_KEY = "sellers"
const CONTRACTS_KEY = "contracts"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Member object for participants
type Member struct {
	Id              string `json:"id"`
	Type            string `json:"memberType"`
	FitcoinsBalance int    `json:"fitcoinsBalance"`
}

// User
type User struct {
	Member
	TotalSteps             int `json:"totalSteps"`
	StepsUsedForConversion int `json:"stepsUsedForConversion"`
}

// Seller
type Seller struct {
	Member
	Products []Product `json:"products"`
}

// Product
type Product struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
	Price int    `json:"price"`
}

// Contract
type Contract struct {
	Id        string `json:"id"`
	SellerId  string `json:"sellerId"`
	UserId    string `json:"userId"`
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
	Cost      int    `json:"price"`
	State     string `json:"state"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - initialize the chaincode
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	//initalize users key-value store
	users := make(map[string]User)
	usersBytes, err := json.Marshal(users)
	if err != nil {
		return shim.Error("Error initializing users.")
	}
	err = stub.PutState(USERS_KEY, usersBytes)

	//initalize sellers key-value store
	sellers := make(map[string]Seller)
	sellersBytes, err := json.Marshal(sellers)
	if err != nil {
		return shim.Error("Error initializing sellers.")
	}
	err = stub.PutState(SELLERS_KEY, sellersBytes)

	//initalize contracts key-value store
	contracts := make(map[string]Contract)
	contractBytes, err := json.Marshal(contracts)
	if err != nil {
		return shim.Error("Error initializing contracts.")
	}
	err = stub.PutState(CONTRACTS_KEY, contractBytes)

	return shim.Success(nil)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)

	//call functions
	if function == "createMember" {
		return t.createMember(stub, args)
	} else if function == "getMember" {
		return t.getMember(stub, args)
	} else if function == "generateFitcoin" {
		return t.generateFitcoin(stub, args)
	} else if function == "getDataByKey" {
		return t.getDataByKey(stub, args)
	} else if function == "createProduct" {
		return t.createProduct(stub, args)
	} else if function == "updateProduct" {
		return t.updateProduct(stub, args)
	} else if function == "getProductByID" {
		return t.getProductByID(stub, args)
	} else if function == "makePurchase" {
		return t.makePurchase(stub, args)
	} else if function == "transactPurchase" {
		return t.transactPurchase(stub, args)
	} else if function == "getContractByID" {
		return t.getContractByID(stub, args)
	}

	return shim.Error("Function with the name " + function + " does not exist.")
}

// ============================================================================================================================
// Get all data for a key
// Inputs - key
// ============================================================================================================================
func (t *SimpleChaincode) getDataByKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}

	//get key from args
	key := args[0]

	//get data by key
	dataBytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Unable to get data by key - check key")
	}

	return shim.Success(dataBytes)

}
