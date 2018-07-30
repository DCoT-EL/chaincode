/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
//Reference for cid library here https://github.com/hyperledger/fabric/blob/master/core/chaincode/lib/cid/README.md
package main

import (
	//"crypto/x509"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func getTxCreatorInfo(stub shim.ChaincodeStubInterface) (string, string, error) {
	var mspid string
	var err error
	var attrValue string
	var found bool
	mspid, err = cid.GetMSPID(stub)
	if err != nil {
		fmt.Printf("Error getting MSP identity: %s\n", err.Error())
		return "", "", err
	}

	attrValue, found, err = cid.GetAttributeValue(stub, TOKEN)
	if err != nil {
		fmt.Printf("Error getting Attribute Value: %s\n", err.Error())
<<<<<<< HEAD
=======
		return "", "", err
	}
	if found == false {
		fmt.Printf("Error getting Attribute Value NOT FOUND!!!")
		err.Error()
>>>>>>> 5e2b8b8ea2722174c932f4f4c868cc45fb052a64
		return "", "", err
	}
	if found == false {
		fmt.Printf("Error getting Attribute Value NOT FOUND!!!")
	//	err.Error()
	//	return "", "", err
	}

	return mspid, attrValue , nil
}

/*
// For now, just hardcode an ACL
// We will support attribute checks in an upgrade

func authenticateExportingEntityOrg(mspID string, certCN string) bool {
	return (mspID == "ExportingEntityOrgMSP") && (certCN == "ca.exportingentityorg.trade.com")
}

func authenticateExporterOrg(mspID string, certCN string) bool {
	return (mspID == "ExporterOrgMSP") && (certCN == "ca.exporterorg.trade.com")
}

func authenticateImporterOrg(mspID string, certCN string) bool {
	return (mspID == "ImporterOrgMSP") && (certCN == "ca.importerorg.trade.com")
}

func authenticateCarrierOrg(mspID string, certCN string) bool {
	return (mspID == "CarrierOrgMSP") && (certCN == "ca.carrierorg.trade.com")
}

func authenticateRegulatorOrg(mspID string, certCN string) bool {
	return (mspID == "RegulatorOrgMSP") && (certCN == "ca.regulatororg.trade.com")
}
*/

func isInvokerOperator(stub shim.ChaincodeStubInterface, attrName string) (bool, string, error) {
	var found bool
	var attrValue string
	var err error

	attrValue, found, err = cid.GetAttributeValue(stub, attrName)
	if err != nil {
		fmt.Printf("Error getting Attribute Value: %s\n", err.Error())
		return false, "", err
	}
	return found, attrValue, nil
}
