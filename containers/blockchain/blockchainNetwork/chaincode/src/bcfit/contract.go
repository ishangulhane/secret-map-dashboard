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
	var sellers []Seller
	json.Unmarshal(sellersBytes, &sellers)

	//find the seller
	var product Product
	sellerFound := false
	productFound := false
	for g := 0; g < len(sellers); g++ {
		if sellers[g].Id == seller_id {
			sellerFound = true
			//find and retrieve the product
			for h := 0; h < len(sellers[g].Products); h++ {
				if sellers[g].Products[h].Id == product_id {
					productFound = true
					product = sellers[g].Products[h]
					break
				}
			}
		}
	}

	//if seller not found return error
	if sellerFound != true {
		return shim.Error("Seller not found")
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
	var contracts []Contract
	json.Unmarshal(contractsBytes, &contracts)

	//append contract to contracts and update the contracts state
	contracts = append(contracts, contract)
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
	var contracts []Contract
	json.Unmarshal(contractBytes, &contracts)

	// find and update contract state
	contractFound := false
	for _, contract := range contracts {
		if contract.Id == contract_id {
			contractFound = true
			seller_id = contract.SellerId
			user_id = contract.UserId
			product_id = contract.ProductId
			quantity = contract.Quantity
			cost = contract.Cost
		}
	}

	//if contract not found return error
	if contractFound != true {
		return shim.Error("Contract not found")
	}

	//if newState is 'complete', updates user account, seller's account and product inventory
	if newState == STATE_COMPLETE {

		//get all users
		usersBytes, err := stub.GetState(USERS_KEY)
		if err != nil {
			return shim.Error("Unable to get users.")
		}
		var users []User
		json.Unmarshal(usersBytes, &users)

		//find the user in users
		userFound := false
		for g := 0; g < len(users); g++ {
			//update user's FitcoinsBalance
			if users[g].Id == user_id {
				userFound = true
				if (users[g].FitcoinsBalance - cost) >= 0 {
					users[g].FitcoinsBalance = users[g].FitcoinsBalance - cost
					break
				} else {
					return shim.Error("Insufficient fitcoins")
				}
			}
		}

		//if contract not found return error
		if userFound != true {
			return shim.Error("User not found")
		}

		//get all sellers
		sellersBytes, err := stub.GetState(SELLERS_KEY)
		if err != nil {
			return shim.Error("Unable to get sellers.")
		}
		var sellers []Seller
		json.Unmarshal(sellersBytes, &sellers)

		//find the seller in sellers
		sellerFound := false
		productFound := false
		for g := 0; g < len(sellers); g++ {
			if sellers[g].Id == seller_id {
				sellerFound = true
				//update seller's FitcoinsBalance
				sellers[g].FitcoinsBalance = sellers[g].FitcoinsBalance + cost
				//find the product for seller and update Count
				for h := 0; h < len(sellers[g].Products); h++ {
					if sellers[g].Products[h].Id == product_id {
						productFound = true
						sellers[g].Products[h].Count = sellers[g].Products[h].Count - quantity
					}
				}
			}
		}

		//if seller not found return error
		if sellerFound != true {
			return shim.Error("Seller not found")
		}
		//if product not found return error
		if productFound != true {
			return shim.Error("Product not found")
		}

		//update users state
		updatedUsersBytes, _ := json.Marshal(users)
		err = stub.PutState(USERS_KEY, updatedUsersBytes)

		//update sellers state
		updatedSellersBytes, _ := json.Marshal(sellers)
		err = stub.PutState(SELLERS_KEY, updatedSellersBytes)

	}

	// update contract state
	for g := 0; g < len(contracts); g++ {
		if contracts[g].Id == contract_id {
			contracts[g].State = newState
			break
		}
	}

	//update contracts state
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
	contractBytes, err := stub.GetState(CONTRACTS_KEY)
	if err != nil {
		return shim.Error("Unable to get contracts.")
	}
	var contracts []Contract
	json.Unmarshal(contractBytes, &contracts)

	//find and return the contract
	for _, contract := range contracts {
		if contract.Id == contract_id {
			contractBytes, _ := json.Marshal(contract)
			return shim.Success(contractBytes)
			break
		}
	}
	//otherwise return nil
	nilBytes, _ := json.Marshal(nil)
	return shim.Success(nilBytes)

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
