package main

import "github.com/hyperledger/fabric/core/chaincode/shim"

func getCOCKey(stub shim.ChaincodeStubInterface, custodyId string) (string, error) {
	cocKey, err := stub.CreateCompositeKey("DCoT_ChainOfCustodyKey", []string{custodyId})
	if err != nil {
		return "", err
	} else {
		return cocKey, nil
	}
}
