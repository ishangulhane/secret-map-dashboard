## Chaincode Functions

This document goes through all the functions in the chaincode. All chaincode function calls are through an input json. The json consists of the following attributes:

* type - enroll, invoke or query
* userID - id to call the function
* fcn - function name
* args - array of string


### Enroll

Enroll member call
```
var input = {
  type: enroll,
  params: {}
}
```


### User invoke calls

The invoke calls from user's iOS app which update the blockchain state.

#### Create user
```
var input = {
  type: invoke,
  params: {
    userId: userId,
    fcn: createMember
    args: userID, user
  }
}
```

#### Generate fitcoins
```
input = {
  type: invoke,
  params: {
    userId: userId,
    fcn: generateFitcoins
    args: userID, totalSteps
  }
}
```

#### Make purchase
```
input = {
  type: invoke,
  params: {
    userId: userId,
    fcn: makePurchase
    args: sellerID, userID, productID, quantity
  }
}
```

#### args def
- userID - the user ID returned from enroll
- totalSteps - the total steps walked by user
- sellerID - the seller's ID
- productID - the id of product with seller, picked by user through interface
- quantity - picked by user through interface


### Seller invoke calls

The invoke calls from seller dashboard which update the blockchain state.

#### Create seller
```
var input = {
  type: invoke,
  params: {
    userId: sellerID
    fcn: createMember
    args: sellerID, seller
  }
}
```

#### Create product inventory
```
var input = {
  type: invoke,
  params: {
    userId: sellerID
    fcn: createProduct
    args: sellerID, productID, productName, productCount productPrice
  }
}
```

#### Update product inventory
```
var input = {
  type: invoke,
  params: {
    userId: sellerID
    fcn: updateProduct
    args: sellerID, productID, productName, productCount productPrice
  }
}
```

#### Transact purchase
```
var input = {
  type: invoke,
  params: {
    userId: sellerID
    fcn: transactPurchase
    args: contractID, newState(complete or declined)
  }
}
```

#### args def
- sellerID - the seller's ID returned from enroll
- productID - product property: the id of product with seller
- productName - product property: the name of product
- productCount - product property: the count of product
- productPrice - product price: the price of product
- contractID - the contract ID generated when user perform 'makePurchase'


### Query calls

The calls that read data from blockchain state database.

#### Get State
Gets state with userId, sellerID or contractID as args
```
var input = {
  type: query,
  params: {
    userId: userID
    fcn: getState
    args: id (userId, sellerID or contractID)
  }
}
```

#### Get products for sale
Gets array of products available with sellerID
```
var input = {
  type: query,
  params: {
    userId: userID
    fcn: getProductsForSale
    args: (none)
  }
}
```

#### Get all user's contracts
```
var input = {
  type: query,
  params: {
    userId: userID,
    fcn: getAllContracts
    args: userID
  }
}
```

#### Get all contracts
```
var input = {
  type: query,
  params: {
    userId: userID,
    fcn: getAllContracts
    args: (none)
  }
}
```

#### Get product by Id
```
var input = {
  type: query,
  params: {
    userId: userID,
    fcn: getProductByID
    args: userID, productId
  }
}
```

#### args def
- userID - the user's ID
