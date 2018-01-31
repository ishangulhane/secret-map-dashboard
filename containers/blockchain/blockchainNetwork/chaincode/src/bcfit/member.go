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
    var user User
    var err error

    // build JSON string and decode into user structure instance
    id := "\"id\": \"" + member_id + "\", "
    memberType := "\"memberType\": \"" + member_type + "\", "
    fitcoinsBalance := "\"fitcoinsBalance\": 0, "
    transactionSteps := "\"transactionSteps\": 0"
    content := "{" + id + memberType + fitcoinsBalance + transactionSteps + "}"
    json.Unmarshal( []byte(content), &user )

		// get all users
    usersBytes, err := stub.GetState( USERS_KEY )
    if err != nil {
      return shim.Error( "Unable to get users." )
    }
    var users []User
    json.Unmarshal( usersBytes, &users )

		//append user to users and update the users state
		users = append( users, user )
    updatedUsersBytes, _ := json.Marshal( users )
    err = stub.PutState( USERS_KEY, updatedUsersBytes )

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
    json.Unmarshal( []byte(content), &seller )

		// get all sellers
    sellersBytes, err := stub.GetState( SELLERS_KEY )
    if err != nil {
      return shim.Error( "Unable to get users." )
    }
    var sellers []Seller
    json.Unmarshal( sellersBytes, &sellers )

		//append seller to sellers and update the users state
    sellers = append( sellers, seller )
    updatedSellersBytes, _ := json.Marshal( sellers )
    err = stub.PutState( SELLERS_KEY, updatedSellersBytes )

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
    usersBytes, err := stub.GetState( USERS_KEY )
    if err != nil {
      return shim.Error( "Unable to get users." )
    }
    var users []User
    json.Unmarshal( usersBytes, &users )

		//find the user in users
    for _, user := range users {
      if user.Id == member_id {
        userBytes, _ := json.Marshal( user )
        return shim.Success(userBytes)
      }
    }
    return shim.Error("User not found")

  }

	//check if type is 'seller'
  if member_type == TYPE_SELLER {

		//get all sellers
    sellersBytes, err := stub.GetState( SELLERS_KEY)
    if err != nil {
      return shim.Error( "Unable to get sellers." )
    }
    var sellers []Seller
    json.Unmarshal( sellersBytes, &sellers )

    //find the seller in sellers
    for _, seller := range sellers {
      if seller.Id == member_id {
        sellerBytes, _ := json.Marshal( seller )
        return shim.Success(sellerBytes)
        break
      }
    }
		return shim.Error("Seller not found")

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

	//get user_id and newTransactionSteps from args
  user_id := args[0]
	newTransactionSteps, err := strconv.Atoi(args[1])
	if err != nil {
	    // handle error
			return shim.Error(err.Error())
	}

	//get all users
  usersBytes, err := stub.GetState( USERS_KEY )
  if err != nil {
    return shim.Error( "Unable to get users." )
  }
  var users []User
  json.Unmarshal( usersBytes, &users )

  //find the user in users
	userFound := false
  for g := 0; g < len( users ); g++ {
		if users[g].Id == user_id {
			userFound = true
			//check if steps walked since the last transaction meets the required number of steps
			if (newTransactionSteps - users[g].TransactionSteps) == STEPS_TO_FITCOIN {
				//update user's FitcoinsBalance
	      users[g].FitcoinsBalance = users[g].FitcoinsBalance + 1
				users[g].TransactionSteps = newTransactionSteps
			} else {
				return shim.Error( "Incorrect transactionSteps." )
			}
    }
  }

	//if user not found return error
	if userFound != true {
    return shim.Error( "User not found" )
  }

  //update users state
  updatedUsersBytes, _ := json.Marshal( users )
  err = stub.PutState( USERS_KEY, updatedUsersBytes )

  return shim.Success(nil)

}
