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

	return t.updateProduct(stub, args)
}

// ============================================================================================================================
// Update product inventory for seller
// Inputs - sellerId, productID, newProductName, newProductCount, newProductPrice
// ============================================================================================================================
func (t *SimpleChaincode) updateProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments")
	}
	var err error

	//get sellerID from args
	seller_id := args[0]
	//get productID from args
	product_id := args[1]

	//get new product properties from args
	newProductName := args[2]
	newProductCount, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("3rd argument 'productCount' must be a numeric string")
	}
	newProductPrice, err := strconv.Atoi(args[4])
	if err != nil {
		return shim.Error("4th argument 'productPrice' must be a numeric string")
	}

	//get seller
	sellersBytes, err := stub.GetState("sellers")
	if err != nil {
		return shim.Error("Unable to get sellers.")
	}
	var sellers []Seller
	json.Unmarshal(sellersBytes, &sellers)

	//find seller and product
	sellerFound := false
	productFound := false

	//find seller in sellers array
	for g := 0; g < len(sellers); g++ {
		if sellers[g].Id == seller_id {
			sellerFound = true
			//find the product and update the properties
			for h := 0; h < len(sellers[g].Products); h++ {
				if sellers[g].Products[h].Id == product_id {
					productFound = true
					sellers[g].Products[h].Name = newProductName
					sellers[g].Products[h].Count = newProductCount
					sellers[g].Products[h].Price = newProductPrice
					break
				}
			}
			//if product not found, append product
			if productFound != true {
				var product Product
				product.Id = product_id
				product.Name = newProductName
				product.Count = newProductCount
				product.Price = newProductPrice
				//append product
				sellers[g].Products = append(sellers[g].Products, product)
			}
			break
		}
	}
	//if product or seller not found return error
	if sellerFound != true {
		return shim.Error("Seller not found")
	}

	updatedSellerAsBytes, _ := json.Marshal(sellers)
	err = stub.PutState("sellers", updatedSellerAsBytes)

	//return seller info
	return shim.Success(updatedSellerAsBytes)

}

// ============================================================================================================================
// Get product inventory for seller
// Inputs - sellerId, productID
// ============================================================================================================================
func (t *SimpleChaincode) getProductByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments")
	}
	var err error

	//get sellerID, productID from args
	seller_id := args[0]
	product_id := args[1]

	//get sellers array
	sellersBytes, err := stub.GetState("sellers")
	if err != nil {
		return shim.Error("Unable to get sellers.")
	}
	var sellers []Seller
	json.Unmarshal(sellersBytes, &sellers)

	//find seller and product
	var product Product
	sellerFound := false
	productFound := false

	//find seller in sellers array
	for g := 0; g < len(sellers); g++ {
		if sellers[g].Id == seller_id {
			sellerFound = true
			//find the product and update the properties
			for h := 0; h < len(sellers[g].Products); h++ {
				if sellers[g].Products[h].Id == product_id {
					productFound = true
					product = sellers[g].Products[h]
					break
				}
			}
			break
		}
	}
	//if product or seller not found return error
	if productFound != true {
		return shim.Error("Product not found")
	}
	if sellerFound != true {
		return shim.Error("Seller not found")
	}

	//return product
	productAsBytes, _ := json.Marshal(product)
	return shim.Success(productAsBytes)
}

// ============================================================================================================================
// Get all products for sale
// Inputs - (none)
// ============================================================================================================================
func (t *SimpleChaincode) getProductsForSale(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//get sellers array
	sellersBytes, err := stub.GetState("sellers")
	if err != nil {
		return shim.Error("Unable to get sellers.")
	}
	var sellers []Seller
	json.Unmarshal(sellersBytes, &sellers)

	// create return object array
	type ReturnProductSale struct {
		SellerID  string `json:"sellerid"`
		ProductId string `json:"productid"`
		Name      string `json:"name"`
		Count     int    `json:"count"`
		Price     int    `json:"price"`
	}
	var returnProducts []ReturnProductSale

	//go through all sellers
	for g := 0; g < len(sellers); g++ {
		//go through all products
		for h := 0; h < len(sellers[g].Products); h++ {
			if sellers[g].Products[h].Count > 0 {
				var returnProduct ReturnProductSale
				returnProduct.SellerID = sellers[g].Id
				returnProduct.ProductId = sellers[g].Products[h].Id
				returnProduct.Name = sellers[g].Products[h].Name
				returnProduct.Count = sellers[g].Products[h].Count
				returnProduct.Price = sellers[g].Products[h].Price
				//append to array
				returnProducts = append(returnProducts, returnProduct)
			}
		}
	}

	//return products for sole
	returnProductsBytes, _ := json.Marshal(returnProducts)
	return shim.Success(returnProductsBytes)
}
