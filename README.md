# Akachain admin

The Akachain Admin Tool provides RESTful API for an administrator to interact with a Hyperledger Fabric network. The list of supported functions are:

## Table of Contents
1. [MSP](#msp)
    1. [Enroll User](#enroll-user)
    2. [Register User](#register-user)
    3. [Revoke User](#revoke-user)
2. [Channel](#channel)
    1. [Create Channel](#create-channel)
    2. [Join Channel](#join-channel)
    3. [List Channel](#list-channel)
3. [Chaincode](#chaincode)
    1. [Install Chaincode](#install-chaincode)
    2. [Approve Chaincode](#approve-chaincode)
    3. [Commit Chaincode](#commit-chaincode)

## MSP

##### Enroll User
```bash
curl --location --request POST 'http://localhost:4001/api/msp/enrollUser' \
--header 'Content-Type: application/json' \
--data-raw '{
  "userName": "admin",
  "enrollSecret": "adminpw",
  "orgName": "Org1"
}'
```

##### Register User
```bash
curl --location --request POST 'http://localhost:4001/api/msp/registerUser' \
--header 'Content-Type: application/json' \
--data-raw '{
  "orgName": 	"Org1",
  "affiliation": "Org1.affiliation1",
  "userName": "appUser",
  "type": "client",
  "attrs": [{ "name": "role", "value": "SuperAdmin", "ecert": true }, {...}]
}'
```

##### Revoke User
```bash
curl --location --request POST 'http://localhost:4001/api/msp/registerUser' \
--header 'Content-Type: application/json' \
--data-raw '{
  "orgName": 	"Org1",
  "userName": "appUser",
}'
```

## Channel
##### Create Channel
```bash
curl --location --request POST 'http://localhost:4001/api/channel/create' \
--header 'Content-Type: application/json' \
--data-raw '{
    "orgName": "Org1",
	"channelName": "mychannel",
	"channelConfig": "mychannel.tx"
}'
```

##### Join Channel
```bash
curl --location --request POST 'http://localhost:4001/api/channel/join' \
--header 'Content-Type: application/json' \
--data-raw '{
    "orgName": "Org1",
	"channelName": "mychannel",
	"peer": "peer0.example.com"
}'
```

##### List Channel
```bash
curl --location --request POST 'http://localhost:4001/api/channel/list' \
--header 'Content-Type: application/json' \
--data-raw '{
    "orgName": "Org1",
	"peer": "peer0.example.com"
}'
```

## Chaincode

##### Install chaincode
```bash
curl --location --request POST 'http://localhost:4001/api/chaincode/install' \
--header 'Content-Type: application/json' \
--data-raw '{
    "orgName": "Org1",
    "chaincodeId": "abstore",
    "chaincodeVersion": "v1.0.0.0",
    "chaincodeType": "golang",
    "chaincodePath": "chaincodes/abstore/go"
}'
```

##### Approve chaincode
```bash
curl --location --request POST 'http://localhost:4001/api/chaincode/approve' \
--header 'Content-Type: application/json' \
--data-raw '{
    "orgName": "Org1",
    "chaincodeId": "abstore",
    "chaincodeVersion": "2",
    "packageId": "abstore:36d51e930bee55f6ca59e825263f7b7ff279de1dc6a3884f3d3876c658c114e9",
    "sequence": "1",
    "initRequired": false
}'
```

##### Commit chaincode
```bash
curl --location --request POST 'http://localhost:4001/api/chaincode/commit' \
--header 'Content-Type: application/json' \
--data-raw '{
    "chaincodeId": "abstore",
    "chaincodeVersion": "2",
    "sequence": "2",
    "initRequired": false,
    "orgs": ["Org1", "Org2"]
}'
```