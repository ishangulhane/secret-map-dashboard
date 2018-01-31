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
<<<<<<< HEAD
	var err error
=======
>>>>>>> Updated chaincode
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments")
	}

	//get id and type from args
	member_id := args[0]
	member_type := strings.ToLower(args[1])

	//check if type is 'user'
	if member_type == TYPE_USER {
<<<<<<< HEAD

		//create user
		var user User
		user.Id = member_id
		user.Type = TYPE_USER
		user.FitcoinsBalance = 0
		user.StepsUsedForConversion = 0
		user.TotalSteps = 0

		//store user
		userAsBytes, _ := json.Marshal(user)
		err = stub.PutState(user.Id, userAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
=======
		var user User
		var err error

		// build JSON string and decode into user structure instance
		id := "\"id\": \"" + member_id + "\", "
		memberType := "\"memberType\": \"" + member_type + "\", "
		fitcoinsBalance := "\"fitcoinsBalance\": 0, "
		transactionSteps := "\"transactionSteps\": 0"
		content := "{" + id + memberType + fitcoinsBalance + transactionSteps + "}"
		json.Unmarshal([]byte(content), &user)

		// get all users
		usersBytes, err := stub.GetState(USERS_KEY)
		if err != nil {
			return shim.Error("Unable to get users.")
		}
		var users []User
		json.Unmarshal(usersBytes, &users)

		//append user to users and update the users state
		users = append(users, user)
		updatedUsersBytes, _ := json.Marshal(users)
		err = stub.PutState(USERS_KEY, updatedUsersBytes)

	}

	//check if type is 'seller'
	if member_type == TYPE_SELLER {
		var seller Seller
		var err error

		// build JSON string and decode into seller structure instance
		id := "\"id\": \"" + member_id + "\", "
		memberType := "\"memberType\": \"" + member_type + "\", "
		fitcoinsBalance := "\"fitcoinsBalance\": 0, "
		sellerProducts := "\"products\": []"
		content := "{" + id + memberType + fitcoinsBalance + sellerProducts + "}"
		json.Unmarshal([]byte(content), &seller)

		// get all sellers
		sellersBytes, err := stub.GetState(SELLERS_KEY)
		if err != nil {
			return shim.Error("Unable to get users.")
		}
		var sellers []Seller
		json.Unmarshal(sellersBytes, &sellers)

		//append seller to sellers and update the users state
		sellers = append(sellers, seller)
		updatedSellersBytes, _ := json.Marshal(sellers)
		err = stub.PutState(SELLERS_KEY, updatedSellersBytes)

	}
>>>>>>> Updated chaincode

		//return user info
		return shim.Success(userAsBytes)

	} else if member_type == TYPE_SELLER {
		//check if type is 'seller'

		//create seller
		var seller Seller
		seller.Id = member_id
		seller.Type = TYPE_SELLER
		seller.FitcoinsBalance = 0

		// store seller
		sellerAsBytes, _ := json.Marshal(seller)
		err = stub.PutState(seller.Id, sellerAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		//get and update sellerIDs
		sellerIdsBytes, err := stub.GetState("sellerIds")
		if err != nil {
			return shim.Error("Unable to get users.")
		}
		var sellerIds []string
		// add sellerID to update sellers
		json.Unmarshal(sellerIdsBytes, &sellerIds)
		sellerIds = append(sellerIds, seller.Id)
		updatedSellerIdsBytes, _ := json.Marshal(sellerIds)
		err = stub.PutState("sellerIds", updatedSellerIdsBytes)

		//return seller info
		return shim.Success(sellerAsBytes)

	}

<<<<<<< HEAD
	return shim.Success(nil)

}

=======
	//get id and type from args
	member_id := args[0]
	member_type := strings.ToLower(args[1])

	if member_type == TYPE_USER {

		//get all users
		usersBytes, err := stub.GetState(USERS_KEY)
		if err != nil {
			return shim.Error("Unable to get users.")
		}
		var users []User
		json.Unmarshal(usersBytes, &users)

		//find the user in users
		for _, user := range users {
			if user.Id == member_id {
				userBytes, _ := json.Marshal(user)
				return shim.Success(userBytes)
			}
		}
		return shim.Error("User not found")

	}

	//check if type is 'seller'
	if member_type == TYPE_SELLER {

		//get all sellers
		sellersBytes, err := stub.GetState(SELLERS_KEY)
		if err != nil {
			return shim.Error("Unable to get sellers.")
		}
		var sellers []Seller
		json.Unmarshal(sellersBytes, &sellers)

		//find the seller in sellers
		for _, seller := range sellers {
			if seller.Id == member_id {
				sellerBytes, _ := json.Marshal(seller)
				return shim.Success(sellerBytes)
				break
			}
		}
		return shim.Error("Seller not found")

	}

	return shim.Error("Unknown type")
}

>>>>>>> Updated chaincode
// ============================================================================================================================
// Generate Fitcoins for the user
// Inputs - userId, transactionSteps
// ============================================================================================================================
func (t *SimpleChaincode) generateFitcoins(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments")
	}
	var err error

<<<<<<< HEAD
	//get user_id and newSteps from args
	user_id := args[0]
	newTransactionSteps, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	//get user
	var user User
	userAsBytes, err := stub.GetState(user_id)
	if err != nil {
		return shim.Error("Failed to get user")
	}
	json.Unmarshal(userAsBytes, &user)
	if user.Type != TYPE_USER {
		return shim.Error("Not user type")
	}

	//update user account
	var newSteps = newTransactionSteps - user.StepsUsedForConversion
	if newSteps > STEPS_TO_FITCOIN {
		var newFitcoins = newSteps / STEPS_TO_FITCOIN
		var remainderSteps = newSteps % STEPS_TO_FITCOIN
		user.FitcoinsBalance = user.FitcoinsBalance + newFitcoins
		user.StepsUsedForConversion = newTransactionSteps - remainderSteps
		user.TotalSteps = newTransactionSteps

		//update users state
		updatedUserAsBytes, _ := json.Marshal(user)
		err = stub.PutState(user_id, updatedUserAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		//return user info
		return shim.Success(updatedUserAsBytes)
	}

	return shim.Success(userAsBytes)
=======
	//get user_id and newTransactionSteps from args
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
	var users []User
	json.Unmarshal(usersBytes, &users)

	//find the user in users
	userFound := false
	for g := 0; g < len(users); g++ {
		if users[g].Id == user_id {
			userFound = true
			//check if steps walked since the last transaction meets the required number of steps
			if (newTransactionSteps - users[g].TransactionSteps) == STEPS_TO_FITCOIN {
				//update user's FitcoinsBalance
				users[g].FitcoinsBalance = users[g].FitcoinsBalance + 1
				users[g].TransactionSteps = newTransactionSteps
			} else {
				return shim.Error("Incorrect transactionSteps.")
			}
		}
	}

	//if user not found return error
	if userFound != true {
		return shim.Error("User not found")
	}

	//update users state
	updatedUsersBytes, _ := json.Marshal(users)
	err = stub.PutState(USERS_KEY, updatedUsersBytes)

	return shim.Success(nil)
>>>>>>> Updated chaincode

}
