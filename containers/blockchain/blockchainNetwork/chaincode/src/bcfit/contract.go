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
	"math/rand"
	"strconv"

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
	contract.Id = randomString(10)
	contract.SellerId = seller_id
	contract.UserId = user_id
	contract.ProductId = product_id
	contract.Quantity = quantity

	//find product price from product in seller
	//get all sellers
	sellersBytes, err := stub.GetState("sellers")
	if err != nil {
		return shim.Error("Unable to get sellers.")
	}
	sellers := make(map[string]Seller)
	json.Unmarshal(sellersBytes, &sellers)

	//get seller
	seller := sellers[seller_id]
	if seller.Id != seller_id {
		return shim.Error("Seller not found")
	}

	var product Product
	//find the seller
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

	//get all contracts
	contractsBytes, err := stub.GetState(CONTRACTS_KEY)
	if err != nil {
		return shim.Error("Unable to get contracts.")
	}
	contracts := make(map[string]Contract)
	json.Unmarshal(contractsBytes, &contracts)

	//append contract to contracts and update the contracts state
	contracts[contract.Id] = contract
	updatedContractsBytes, _ := json.Marshal(contracts)
	err = stub.PutState(CONTRACTS_KEY, updatedContractsBytes)

	return shim.Success(nil)

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

	//get all contracts
	contractBytes, err := stub.GetState(CONTRACTS_KEY)
	if err != nil {
		return shim.Error("Unable to get contracts.")
	}
	contracts := make(map[string]Contract)
	json.Unmarshal(contractBytes, &contracts)

	// find and update contract state
	contract := contracts[contract_id]
	if contract.Id != contract_id {
		return shim.Error("Contract not found")
	}

	seller_id = contract.SellerId
	user_id = contract.UserId
	product_id = contract.ProductId
	quantity = contract.Quantity
	cost = contract.Cost

	//if newState is 'complete', updates user account, seller's account and product inventory
	if newState == STATE_COMPLETE {

		//get all users
		usersBytes, err := stub.GetState(USERS_KEY)
		if err != nil {
			return shim.Error("Unable to get users.")
		}
		users := make(map[string]User)
		json.Unmarshal(usersBytes, &users)

		//find the user in users
		user := users[user_id]
		if user.Id != user_id {
			return shim.Error("User not found")
		}

		//update user's FitcoinsBalance
		if (user.FitcoinsBalance - cost) >= 0 {
			user.FitcoinsBalance = user.FitcoinsBalance - cost
		} else {
			return shim.Error("Insufficient fitcoins")
		}

		//get all sellers
		sellersBytes, err := stub.GetState(SELLERS_KEY)
		if err != nil {
			return shim.Error("Unable to get sellers.")
		}
		sellers := make(map[string]Seller)
		json.Unmarshal(sellersBytes, &sellers)

		//get seller
		seller := sellers[seller_id]
		if seller.Id != seller_id {
			return shim.Error("Seller not found")
		}

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
		users[user_id] = user
		updatedUsersBytes, _ := json.Marshal(users)
		err = stub.PutState(USERS_KEY, updatedUsersBytes)

		//update sellers state
		sellers[seller_id] = seller
		updatedSellersBytes, _ := json.Marshal(sellers)
		err = stub.PutState(SELLERS_KEY, updatedSellersBytes)

	}

	// update contract state
	contract.State = newState

	//update contracts state
	contracts[contract_id] = contract
	updatedContractsBytes, _ := json.Marshal(contracts)
	err = stub.PutState(CONTRACTS_KEY, updatedContractsBytes)

	return shim.Success(nil)

}

// ============================================================================================================================
// Get contracts by ID
// Inputs - contractID
// ============================================================================================================================
func (t *SimpleChaincode) getContractByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}
	var err error

	//get contractID from args
	contract_id := args[0]

	//gets all contracts
	contractsBytes, err := stub.GetState(CONTRACTS_KEY)
	if err != nil {
		return shim.Error("Unable to get contracts.")
	}
	contracts := make(map[string]Contract)
	json.Unmarshal(contractsBytes, &contracts)

	//find and return the contract
	contract := contracts[contract_id]
	if contract.Id != contract_id {
		return shim.Error("Contract not found")
	}

	//return user info
	contractBytes, _ := json.Marshal(contract)
	return shim.Success(contractBytes)

}

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// Generate a random string of A-Z chars with len = l
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}
