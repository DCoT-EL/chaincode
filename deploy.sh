#!/bin/bash
cd src/github.com/dcot-chaincode
CORE_CHAINCODE_LOGLEVE=debug CORE_PEER_ADDRESS=127.0.0.1:7051 CORE_CHAINCODE_ID_NAME=dcot-chaincode:$1 ./dcot-chaincode
peer chaincode install -n dcot-chaincode -v $1 -p github.com/hyperledger/fabric/examples/chaincode/go/dcot-chaincode

