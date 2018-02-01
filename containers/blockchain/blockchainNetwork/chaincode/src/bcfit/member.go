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
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// Create member
// Inputs - id, type(user or seller)
// ============================================================================================================================
func (t *SimpleChaincode) createMember(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments")
	}

	//get id and type from args
	member_id := args[0]
	member_type := strings.ToLower(args[1])

	//check if type is 'user'
	if member_type == TYPE_USER {
		var err error

		//create user
		var user User
		user.Id = member_id
		user.Type = member_type
		user.FitcoinsBalance = 0
		user.StepsUsedForConversion = 0
		user.TotalSteps = 0

		// get all users
		usersBytes, err := stub.GetState(USERS_KEY)
		if err != nil {
			return shim.Error("Unable to get users.")
		}
		users := make(map[string]User)
		json.Unmarshal(usersBytes, &users)

		//add user and update users state
		users[member_id] = user
		updatedUsersBytes, _ := json.Marshal(users)
		err = stub.PutState(USERS_KEY, updatedUsersBytes)

		//return user info
		userBytes, _ := json.Marshal(user)
		return shim.Success(userBytes)
	}

	//check if type is 'seller'
	if member_type == TYPE_SELLER {
		var err error

		//create seller
		var seller Seller
		seller.Id = member_id
		seller.Type = member_type
		seller.FitcoinsBalance = 0

		// get all sellers
		sellersBytes, err := stub.GetState(SELLERS_KEY)
		if err != nil {
			return shim.Error("Unable to get users.")
		}
		sellers := make(map[string]Seller)
		json.Unmarshal(sellersBytes, &sellers)

		//add seller and update sellers state
		sellers[member_id] = seller
		updatedSellersBytes, _ := json.Marshal(sellers)
		err = stub.PutState(SELLERS_KEY, updatedSellersBytes)

		//return seller info
		sellerBytes, _ := json.Marshal(seller)
		return shim.Success(sellerBytes)

	}

	return shim.Success(nil)

}

// ============================================================================================================================
// Get member
// Inputs - id, type(user or seller)
// ============================================================================================================================
func (t *SimpleChaincode) getMember(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments")
	}

	//get id and type from args
	member_id := args[0]
	member_type := strings.ToLower(args[1])

	if member_type == TYPE_USER {

		//get all users
		usersBytes, err := stub.GetState(USERS_KEY)
		if err != nil {
			return shim.Error("Unable to get users.")
		}
		users := make(map[string]User)
		json.Unmarshal(usersBytes, &users)

		//find the user in users
		user := users[member_id]
		if user.Id != member_id {
			return shim.Error("User not found")
		}

		//return user info
		userBytes, _ := json.Marshal(user)
		return shim.Success(userBytes)
	}

	//check if type is 'seller'
	if member_type == TYPE_SELLER {

		//get all sellers
		sellersBytes, err := stub.GetState(SELLERS_KEY)
		if err != nil {
			return shim.Error("Unable to get sellers.")
		}
		sellers := make(map[string]Seller)
		json.Unmarshal(sellersBytes, &sellers)

		//find the seller in sellers
		seller := sellers[member_id]
		if seller.Id != member_id {
			return shim.Error("Seller not found")
		}

		//return seller info
		sellerBytes, _ := json.Marshal(seller)
		return shim.Success(sellerBytes)

	}

	return shim.Error("Unknown type")
}

// ============================================================================================================================
// Generate Fitcoin for the user
// Inputs - userId, transactionSteps
// ============================================================================================================================
func (t *SimpleChaincode) generateFitcoin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments")
	}
	var err error

	//get user_id and newSteps from args
	user_id := args[0]
	newTransactionSteps, err := strconv.Atoi(args[1])
	if err != nil {
		// handle error
		return shim.Error(err.Error())
	}

	//get all users
	usersBytes, err := stub.GetState(USERS_KEY)
	if err != nil {
		return shim.Error("Unable to get users.")
	}
	users := make(map[string]User)
	json.Unmarshal(usersBytes, &users)

	//update the user in users
	user := users[user_id]
	if user.Id != user_id {
		return shim.Error("User not found")
	}

	var newSteps = newTransactionSteps - user.StepsUsedForConversion
	if newSteps > STEPS_TO_FITCOIN {
		var newFitcoins = newSteps / STEPS_TO_FITCOIN
		var remainderSteps = newSteps % STEPS_TO_FITCOIN
		user.FitcoinsBalance = user.FitcoinsBalance + newFitcoins
		user.StepsUsedForConversion = newTransactionSteps - remainderSteps
		user.TotalSteps = newTransactionSteps
		users[user_id] = user

		//update users state
		updatedUsersBytes, _ := json.Marshal(users)
		err = stub.PutState(USERS_KEY, updatedUsersBytes)
	}

	userBytes, _ := json.Marshal(user)
	return shim.Success(userBytes)

}
