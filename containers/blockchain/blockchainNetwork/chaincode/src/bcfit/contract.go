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
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// Make Purchase - creates purchase Contract
// Inputs - sellerID, userID, productID, quantity
// ============================================================================================================================
func (t *SimpleChaincode) makePurchase(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments")
	}
	var err error

	//gets sellerID, userID, productID, quantity from args
	seller_id := args[0]
	user_id := args[1]
	product_id := args[2]
	quantity, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("4th argument must be a numeric string")
	}

	//creates contract struct with properties
	var contract Contract
	contract.Id = "c" + randomInts(10)
	contract.SellerId = seller_id
	contract.UserId = user_id
	contract.ProductId = product_id
	contract.Quantity = quantity

	//get seller
	sellerAsBytes, err := stub.GetState(seller_id)
	if err != nil {
		return shim.Error("Failed to get seller")
	}
	seller := Seller{}
	json.Unmarshal(sellerAsBytes, &seller)

	//find the product
	var product Product
	productFound := false
	for h := 0; h < len(seller.Products); h++ {
		if seller.Products[h].Id == product_id {
			productFound = true
			product = seller.Products[h]
			break
		}
	}

	//if product not found return error
	if productFound != true {
		return shim.Error("Product not found")
	}

	//calculates cost and assigns to contract
	contract.Cost = product.Price * quantity
	//assign 'Pending' state
	contract.State = STATE_PENDING

	//store contract
	contractAsBytes, _ := json.Marshal(contract)      //convert to array of bytes
	err = stub.PutState(contract.Id, contractAsBytes) //store owner by its Id
	if err != nil {
		return shim.Error(err.Error())
	}

	// get user's current state
	var user User
	userAsBytes, err := stub.GetState(user_id)
	if err != nil {
		return shim.Error("Failed to get user")
	}
	json.Unmarshal(userAsBytes, &user)

	//append contractId
	user.ContractIds = append(user.ContractIds, contract.Id)

	//update seller's state
	updatedUserAsBytes, _ := json.Marshal(user)
	err = stub.PutState(user_id, updatedUserAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//return contract info
	return shim.Success(contractAsBytes)

}

// ============================================================================================================================
// Transact Purchase - update user account, update seller's account and product inventory, update contract state
// Inputs - contractID, contractState
// ============================================================================================================================
func (t *SimpleChaincode) transactPurchase(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments")
	}
	var err error
	var seller_id string
	var user_id string
	var product_id string
	var quantity int
	var cost int

	//get contractID and state from args
	contract_id := args[0]
	newState := args[1]

	// Get contract from the ledger
	contractAsBytes, err := stub.GetState(contract_id)
	if err != nil {
		return shim.Error("Failed to get user")
	}

	var contract Contract
	json.Unmarshal(contractAsBytes, &contract)

	seller_id = contract.SellerId
	user_id = contract.UserId
	product_id = contract.ProductId
	quantity = contract.Quantity
	cost = contract.Cost

	//if newState is 'complete', updates user account, seller's account and product inventory
	if newState == STATE_COMPLETE {

		// get user's current state
		var user User
		userAsBytes, err := stub.GetState(user_id)
		if err != nil {
			return shim.Error("Failed to get user")
		}
		json.Unmarshal(userAsBytes, &user)

		//update user's FitcoinsBalance
		if (user.FitcoinsBalance - cost) >= 0 {
			user.FitcoinsBalance = user.FitcoinsBalance - cost
		} else {
			return shim.Error("Insufficient fitcoins")
		}

		// get seller's current state
		var seller Seller
		sellerAsBytes, err := stub.GetState(seller_id)
		if err != nil {
			return shim.Error("Failed to get user")
		}
		json.Unmarshal(sellerAsBytes, &seller)

		//update seller's FitcoinsBalance
		seller.FitcoinsBalance = seller.FitcoinsBalance + cost

		//update seller's product count
		productFound := false
		for h := 0; h < len(seller.Products); h++ {
			if seller.Products[h].Id == product_id {
				productFound = true
				seller.Products[h].Count = seller.Products[h].Count - quantity
				break
			}
		}

		//if product not found return error
		if productFound != true {
			return shim.Error("Product not found")
		}

		//update users state
		updatedUserAsBytes, _ := json.Marshal(user)
		err = stub.PutState(user_id, updatedUserAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		//update seller's state
		updatedSellerAsBytes, _ := json.Marshal(seller)
		err = stub.PutState(seller_id, updatedSellerAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

	}

	// update contract state
	contract.State = newState

	//store contract
	updatedContractAsBytes, _ := json.Marshal(contract)
	err = stub.PutState(contract.Id, updatedContractAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//return contract info
	return shim.Success(updatedContractAsBytes)
}

// ============================================================================================================================
// Get all contracts
// Inputs -
// ============================================================================================================================
func (t *SimpleChaincode) getAllContracts(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var contracts []Contract

	// ---- Get All Contracts ---- //
	resultsIterator, err := stub.GetStateByRange("c0", "c9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		aKeyValue, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on contract id - ", queryKeyAsStr)
		var contract Contract
		json.Unmarshal(queryValAsBytes, &contract)
		contracts = append(contracts, contract)
	}

	//change to array of bytes
	contractsAsBytes, _ := json.Marshal(contracts) //convert to array of bytes
	return shim.Success(contractsAsBytes)

}

//generate an array of random ints
func randomArray(len int) []int {
	a := make([]int, len)
	for i := 0; i <= len-1; i++ {
		a[i] = rand.Intn(10)
	}
	return a
}

// Generate a random string of ints with length len
func randomInts(len int) string {
	rand.Seed(time.Now().UnixNano())
	intArray := randomArray(len)
	var stringInt []string
	for _, i := range intArray {
		stringInt = append(stringInt, strconv.Itoa(i))
	}
	return strings.Join(stringInt, "")
}
