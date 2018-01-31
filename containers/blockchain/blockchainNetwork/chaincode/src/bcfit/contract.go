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
<<<<<<< HEAD
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
=======
	"math/rand"
	"strconv"
>>>>>>> Updated chaincode

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

	//creates contract struct with properties, and get sellerID, userID, productID, quantity from args
	var contract Contract
	contract.Id = "c" + randomInts(10)
	contract.SellerId = args[0]
	contract.UserId = args[1]
	contract.ProductId = args[2]
	quantity, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("4th argument 'quantity' must be a numeric string")
	}
	contract.Quantity = quantity

<<<<<<< HEAD
	//get seller
	sellerAsBytes, err := stub.GetState(contract.SellerId)
	if err != nil {
		return shim.Error("Failed to get seller")
	}
	var seller Seller
	json.Unmarshal(sellerAsBytes, &seller)
	if seller.Type != TYPE_SELLER {
		return shim.Error("Not seller type")
	}
=======
	//find product price from product in seller
	//get all sellers
	sellersBytes, err := stub.GetState("sellers")
	if err != nil {
		return shim.Error("Unable to get sellers.")
	}
	var sellers []Seller
	json.Unmarshal(sellersBytes, &sellers)
>>>>>>> Updated chaincode

	//find the product
	var product Product
	productFound := false
<<<<<<< HEAD
	for h := 0; h < len(seller.Products); h++ {
		if seller.Products[h].Id == contract.ProductId {
			productFound = true
			product = seller.Products[h]
			break
		}
	}

=======
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
>>>>>>> Updated chaincode
	//if product not found return error
	if productFound != true {
		return shim.Error("Product not found")
	}

	//calculates cost and assigns to contract
	contract.Cost = product.Price * contract.Quantity
	//gets product name
	contract.ProductName = product.Name
	//assign 'Pending' state
	contract.State = STATE_PENDING

<<<<<<< HEAD
	// get user's current state
	var user User
	userAsBytes, err := stub.GetState(contract.UserId)
	if err != nil {
		return shim.Error("Failed to get user")
	}
	json.Unmarshal(userAsBytes, &user)
	if user.Type != TYPE_USER {
		return shim.Error("Not user type")
	}

	//check if user has enough Fitcoinsbalance
	if user.FitcoinsBalance < contract.Cost {
		return shim.Error("Insufficient funds")
	}
=======
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
>>>>>>> Updated chaincode

	//store contract
	contractAsBytes, _ := json.Marshal(contract)
	err = stub.PutState(contract.Id, contractAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//append contractId
	user.ContractIds = append(user.ContractIds, contract.Id)

	//update user's state
	updatedUserAsBytes, _ := json.Marshal(user)
	err = stub.PutState(contract.UserId, updatedUserAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

<<<<<<< HEAD
	//return contract info
	return shim.Success(contractAsBytes)

}

=======
>>>>>>> Updated chaincode
// ============================================================================================================================
// Transact Purchase - update user account, update seller's account and product inventory, update contract state
// Inputs - contractID, newState(complete or declined)
// ============================================================================================================================
func (t *SimpleChaincode) transactPurchase(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments")
	}

	//get contractID args
	contract_id := args[0]
	newState := args[1]

<<<<<<< HEAD
	// Get contract from the ledger
	contractAsBytes, err := stub.GetState(contract_id)
	if err != nil {
		return shim.Error("Failed to get user")
	}
	var contract Contract
	json.Unmarshal(contractAsBytes, &contract)

	//if current contract state is pending, then execute transaction
	if contract.State == STATE_PENDING {

		if newState == STATE_COMPLETE {

			// get user's current state
			var user User
			userAsBytes, err := stub.GetState(contract.UserId)
			if err != nil {
				return shim.Error("Failed to get user")
			}
			json.Unmarshal(userAsBytes, &user)
			if user.Type != TYPE_USER {
				return shim.Error("Not user type")
			}

			//update user's FitcoinsBalance
			if (user.FitcoinsBalance - contract.Cost) >= 0 {
				user.FitcoinsBalance = user.FitcoinsBalance - contract.Cost
			} else {
				return shim.Error("Insufficient fitcoins")
			}

			// get seller's current state
			var seller Seller
			sellerAsBytes, err := stub.GetState(contract.SellerId)
			if err != nil {
				return shim.Error("Failed to get user")
			}
			json.Unmarshal(sellerAsBytes, &seller)
			if seller.Type != TYPE_SELLER {
				return shim.Error("Not seller type")
			}

			//update seller's FitcoinsBalance
			seller.FitcoinsBalance = seller.FitcoinsBalance + contract.Cost

			//update seller's product count
			productFound := false
			for h := 0; h < len(seller.Products); h++ {
				if seller.Products[h].Id == contract.ProductId {
					productFound = true
					seller.Products[h].Count = seller.Products[h].Count - contract.Quantity
					break
				}
			}

			//if product not found return error
			if productFound != true {
				return shim.Error("Product not found")
			}

			//update users state
			updatedUserAsBytes, _ := json.Marshal(user)
			err = stub.PutState(contract.UserId, updatedUserAsBytes)
			if err != nil {
				return shim.Error(err.Error())
			}

			//update seller's state
			updatedSellerAsBytes, _ := json.Marshal(seller)
			err = stub.PutState(contract.SellerId, updatedSellerAsBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
		}
=======
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
>>>>>>> Updated chaincode

		// update contract state
		contract.State = STATE_COMPLETE

<<<<<<< HEAD
		// update contract state on ledger
		updatedContractAsBytes, _ := json.Marshal(contract)
		err = stub.PutState(contract.Id, updatedContractAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		//return contract info
		return shim.Success(updatedContractAsBytes)

	} else {
		return shim.Error("Contract already Complete or Declined")
	}
=======
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
>>>>>>> Updated chaincode

}

// ============================================================================================================================
// Get all user contracts
// Inputs - userID
// ============================================================================================================================
func (t *SimpleChaincode) getAllUserContracts(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}
	var err error

	//get userID from args
	user_id := args[0]

<<<<<<< HEAD
	//get user
	userAsBytes, err := stub.GetState(user_id)
	if err != nil {
		return shim.Error("Failed to get user")
=======
	//gets all contracts
	contractBytes, err := stub.GetState(CONTRACTS_KEY)
	if err != nil {
		return shim.Error("Unable to get contracts.")
>>>>>>> Updated chaincode
	}
	var user User
	json.Unmarshal(userAsBytes, &user)
	if user.Type != TYPE_USER {
		return shim.Error("Not user type")
	}

	//get user contracts
	var contracts []Contract
<<<<<<< HEAD
	for h := 0; h < len(user.ContractIds); h++ {
		//get contract from the ledger
		contractAsBytes, err := stub.GetState(user.ContractIds[h])
		if err != nil {
			return shim.Error("Failed to get contract")
		}
		var contract Contract
		json.Unmarshal(contractAsBytes, &contract)
		contracts = append(contracts, contract)
	}
	//change to array of bytes
	contractsAsBytes, _ := json.Marshal(contracts)
	return shim.Success(contractsAsBytes)

}

// ============================================================================================================================
// Get all contracts
// Inputs - (none)
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
=======
	json.Unmarshal(contractBytes, &contracts)

	//find and return the contract
	for _, contract := range contracts {
		if contract.Id == contract_id {
			contractBytes, _ := json.Marshal(contract)
			return shim.Success(contractBytes)
			break
>>>>>>> Updated chaincode
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on contract id - ", queryKeyAsStr)
		var contract Contract
		json.Unmarshal(queryValAsBytes, &contract)
		contracts = append(contracts, contract)
	}
<<<<<<< HEAD

	//change to array of bytes
	contractsAsBytes, _ := json.Marshal(contracts)
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
=======
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
>>>>>>> Updated chaincode
}
