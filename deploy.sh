#!/bin/bash

FABRIC_DIR='/home/ubuntu/configuration-network-fabric'


CHANNEL='ledgerchannel'
HOST='localhost'
OPERATION='instantiate'
if [ "$1" == "" ]; then
	echo "Insert the version of your chaincode ex: '1.0'"
	echo "Usage ./deploy <version of chaincode> <upgrade or instantiate>"
	exit 0
fi
if [ "$2" != "" ]; then
	OPERATION='upgrade'
fi
echo "Use $HOST as HOST for HLF"
cd src/github.com/dcot-chaincode
go get github.com/hyperledger/fabric/core/chaincode/lib/cid
go get github.com/hyperledger/fabric/core/chaincode/shim
go get github.com/rs/xid
go get github.com/hyperledger/fabric/protos/peer
go build
#CORE_CHAINCODE_LOGLEVEL=debug CORE_PEER_ADDRESS=$HOST:10051 CORE_CHAINCODE_ID_NAME=dcot-chaincode:$1 ./dcot-chaincode
cd .. && cp -fR dcot-chaincode $FABRIC_DIR/chaincode/go
echo "Install chaincode"
docker exec -it cli peer chaincode install -n dcot-chaincode -v $1 -p github.com/hyperledger/fabric/examples/chaincode/go/dcot-chaincode
sleep 30 
echo "Instantiate chaincode"
docker exec -it cli peer chaincode $OPERATION -n dcot-chaincode -c '{"Args":[""]}' -C $CHANNEL -v $1
