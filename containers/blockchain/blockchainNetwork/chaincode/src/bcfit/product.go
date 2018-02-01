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

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// Create product inventory for seller
// Inputs - sellerId, productID, productName, productCount, productPrice
// ============================================================================================================================
func (t *SimpleChaincode) createProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments")
	}

	//get sellerID from args
	seller_id := args[0]
	//create new product object from args
	var product Product
	product.Id = args[1]
	product.Name = args[2]
	count, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("4th argument must be a numeric string")
	}
	product.Count = count
	price, err := strconv.Atoi(args[4])
	if err != nil {
		return shim.Error("5th argument must be a numeric string")
	}
	product.Price = price

	//get all sellers
	sellersBytes, err := stub.GetState(SELLERS_KEY)
	if err != nil {
		return shim.Error("Unable to get sellers.")
	}
	sellers := make(map[string]Seller)
	json.Unmarshal(sellersBytes, &sellers)

	//find the seller, and append product
	seller := sellers[seller_id]
	if seller.Id != seller_id {
		return shim.Error("Seller not found")
	}
	seller.Products = append(seller.Products, product)
	sellers[seller_id] = seller

	//update seller state
	updatedSellersBytes, _ := json.Marshal(sellers)
	err = stub.PutState(SELLERS_KEY, updatedSellersBytes)

	return shim.Success(nil)

}

// ============================================================================================================================
// Update product inventory for seller
// Inputs - sellerId, productID, newProductName, newProductCount, newProductPrice
// ============================================================================================================================
func (t *SimpleChaincode) updateProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments")
	}

	//get sellerID from args
	seller_id := args[0]
	//get productID from args
	product_id := args[1]

	//get new product properties from args
	newProductName := args[2]
	newProductCount, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("3rd argument must be a numeric string")
	}
	newProductPrice, err := strconv.Atoi(args[4])
	if err != nil {
		return shim.Error("4th argument must be a numeric string")
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

	//find the product and update the properties
	productFound := false
	for h := 0; h < len(seller.Products); h++ {
		if seller.Products[h].Id == product_id {
			productFound = true
			seller.Products[h].Name = newProductName
			seller.Products[h].Count = newProductCount
			seller.Products[h].Price = newProductPrice
			break
		}
	}

	//if product not found return error
	if productFound != true {
		return shim.Error("Product not found")
	}

	//update sellers state
	sellers[seller_id] = seller
	updatedSellersBytes, _ := json.Marshal(sellers)
	err = stub.PutState(SELLERS_KEY, updatedSellersBytes)

	return shim.Success(nil)

}

// ============================================================================================================================
// Get product inventory for seller
// Inputs - sellerId, productID
// ============================================================================================================================
func (t *SimpleChaincode) getProductByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments")
	}

	//get sellerID, productID from args
	seller_id := args[0]
	product_id := args[1]

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

	//return product type
	productAsBytes, _ := json.Marshal(product)
	return shim.Success(productAsBytes)

}
