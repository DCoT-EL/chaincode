#!/bin/bash
if [ "$1" == "" ]; then
	echo "Insert the version of your chaincode ex: '1.0'"
	exit 0
fi
cd src/github.com/dcot-chaincode
go get github.com/hyperledger/fabric/core/chaincode/lib/cid
go get github.com/hyperledger/fabric/core/chaincode/shim
go get github.com/rs/xid
go get github.com/hyperledger/fabric/protos/peer
go build
CORE_CHAINCODE_LOGLEVEL=debug CORE_PEER_ADDRESS=127.0.0.1:7051 CORE_CHAINCODE_ID_NAME=dcot-chaincode:$1 ./dcot-chaincode
peer chaincode install -n dcot-chaincode -v $1 -p github.com/hyperledger/fabric/examples/chaincode/go/dcot-chaincode

