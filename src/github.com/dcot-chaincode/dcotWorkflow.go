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

package main

import (
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/rs/xid"
)

//var logger = shim.NewLogger("dcot-chaincode")

// DcotWorkflowChaincode implementation
type DcotWorkflowChaincode struct {
	testMode bool
}

func (t *DcotWorkflowChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Initializing Chain of Custody")
	_, args := stub.GetFunctionAndParameters()
	//var err error

	// Upgrade Mode 1: leave ledger state as it was
	if len(args) == 0 {
		fmt.Print("Args correctly!!!")
		return shim.Success(nil)
	}

	// Upgrade mode 2: change all the names and account balances
	/*
		if len(args) != 8 {
			err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 8: {"+
				"Exporter, "+
				"Exporter's Bank, "+
				"Exporter's Account Balance, "+
				"Importer, "+
				"Importer's Bank, "+
				"Importer's Account Balance, "+
				"Carrier, "+
				"Regulatory Authority"+
				"}. Found %d", len(args)))
			return shim.Error(err.Error())
		}

		// Type checks
		_, err = strconv.Atoi(string(args[2]))
		if err != nil {
			fmt.Printf("Exporter's account balance must be an integer. Found %s\n", args[2])
			return shim.Error(err.Error())
		}
		_, err = strconv.Atoi(string(args[5]))
		if err != nil {
			fmt.Printf("Importer's account balance must be an integer. Found %s\n", args[5])
			return shim.Error(err.Error())
		}

		fmt.Printf("Exporter: %s\n", args[0])
		fmt.Printf("Exporter's Bank: %s\n", args[1])
		fmt.Printf("Exporter's Account Balance: %s\n", args[2])
		fmt.Printf("Importer: %s\n", args[3])
		fmt.Printf("Importer's Bank: %s\n", args[4])
		fmt.Printf("Importer's Account Balance: %s\n", args[5])
		fmt.Printf("Carrier: %s\n", args[6])
		fmt.Printf("Regulatory Authority: %s\n", args[7])

		// Map participant identities to their roles on the ledger
		roleKeys := []string{expKey, ebKey, expBalKey, impKey, ibKey, impBalKey, carKey, raKey}
		for i, roleKey := range roleKeys {
			err = stub.PutState(roleKey, []byte(args[i]))
			if err != nil {
				fmt.Errorf("Error recording key %s: %s\n", roleKey, err.Error())
				return shim.Error(err.Error())
			}
		}
	*/
	return shim.Success(nil)
}

func (t *DcotWorkflowChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	var creatorOrg, creatorCertIssuer string
	//var attrValue string
	var err error
	var isEnabled bool

	fmt.Println("DcotWorkflow Invoke")

	if !t.testMode {
		creatorOrg, creatorCertIssuer, err = getTxCreatorInfo(stub)
		if err != nil {
			fmt.Errorf("Error extracting creator identity info: %s\n", err.Error())
			return shim.Error(err.Error())
		}
		fmt.Printf("DcotWorkflow Invoke by '%s', '%s'\n", creatorOrg, creatorCertIssuer)

		isEnabled, _, err = isInvokerOperator(stub, "dcot-operator")
		if err != nil {
			fmt.Errorf("Error getting attribute info: %s\n", err.Error())
			return shim.Error(err.Error())
		}
	}

	function, args := stub.GetFunctionAndParameters()
	/*
		if function == "requestTrade" {
			// Importer requests a trade
			return t.requestTrade(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "acceptTrade" {
			// Exporter accepts a trade
			return t.acceptTrade(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "requestLC" {
			// Importer requests an L/C
			return t.requestLC(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "issueLC" {
			// Importer's Bank issues an L/C
			return t.issueLC(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "acceptLC" {
			// Exporter's Bank accepts an L/C
			return t.acceptLC(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "requestEL" {
			// Exporter requests an E/L
			return t.requestEL(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "issueEL" {
			// Regulatory Authority issues an E/L
			return t.issueEL(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "prepareShipment" {
			// Exporter prepares a shipment
			return t.prepareShipment(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "acceptShipmentAndIssueBL" {
			// Carrier validates the shipment and issues a B/L
			return t.acceptShipmentAndIssueBL(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "requestPayment" {
			// Exporter's Bank requests a payment
			return t.requestPayment(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "makePayment" {
			// Importer's Bank makes a payment
			return t.makePayment(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "updateShipmentLocation" {
			// Carrier updates the shipment location
			return t.updateShipmentLocation(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "getTradeStatus" {
			// Get status of trade agreement
			return t.getTradeStatus(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "getLCStatus" {
			// Get the L/C status
			return t.getLCStatus(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "getELStatus" {
			// Get the E/L status
			return t.getELStatus(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "getShipmentLocation" {
			// Get the shipment location
			return t.getShipmentLocation(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "getBillOfLading" {
			// Get the bill of lading
			return t.getBillOfLading(stub, creatorOrg, creatorCertIssuer, args)
		} else if function == "getAccountBalance" {
			// Get account balance: Exporter/Importer
			return t.getAccountBalance(stub, creatorOrg, creatorCertIssuer, args)
			/*} else if function == "delete" {
			// Deletes an entity from its state
			return t.delete(stub, creatorOrg, creatorCertIssuer, args)
		}
	*/
	if function == "initNewChain" {
		return t.initNewChain(stub, isEnabled, args)
	} else if function == "startTransfer" {
		return t.startTransfer(stub, isEnabled, args)
	} else if function == "completeTrasfer" {
		return t.completeTrasfer(stub, isEnabled, args)
	} else if function == "commentChain" {
		return t.commentChain(stub, isEnabled, args)
	} else if function == "cancelTrasfer" {
		return t.cancelTrasfer(stub, isEnabled, args)
	} else if function == "terminateChain" {
		return t.terminateChain(stub, isEnabled, args)
	} else if function == "updateDocument" {
		return t.updateDocument(stub, isEnabled, args)
	} else if function == "getAssetDetails" {
		return t.getAssetDetails(stub, isEnabled, args)
	} else if function == "getChainOfEvents" {
		return t.getChainOfEvents(stub, isEnabled, args)
	}
	return shim.Error("Invalid invoke function name")
}

/*
// Request a trade agreement
func (t *DcotWorkflowChaincode) requestTrade(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var tradeKey string
	var tradeAgreement *TradeAgreement
	var tradeAgreementBytes []byte
	var amount int
	var err error

	// Access control: Only an DCOT operatorcan invoke this transaction
	if !t.testMode && !authenticateImporterOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Importer Org. Access denied.")
	}

	if len(args) != 3 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 3: {ID, Amount, Description of Goods}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	amount, err = strconv.Atoi(string(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}

	tradeAgreement = &TradeAgreement{amount, args[2], REQUESTED, 0}
	tradeAgreementBytes, err = json.Marshal(tradeAgreement)
	if err != nil {
		return shim.Error("Error marshaling trade agreement structure")
	}

	// Write the state to the ledger
	tradeKey, err = getTradeKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(tradeKey, tradeAgreementBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Trade %s request recorded\n", args[0])

	return shim.Success(nil)
}

// Accept a trade agreement
func (t *DcotWorkflowChaincode) acceptTrade(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var tradeKey string
	var tradeAgreement *TradeAgreement
	var tradeAgreementBytes []byte
	var err error

	// Access control: Only an Exporting Entity Org member can invoke this transaction
	if !t.testMode && !authenticateExportingEntityOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Exporting Entity Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {ID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	// Get the state from the ledger
	tradeKey, err = getTradeKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	tradeAgreementBytes, err = stub.GetState(tradeKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(tradeAgreementBytes) == 0 {
		err = errors.New(fmt.Sprintf("No record found for trade ID %s", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(tradeAgreementBytes, &tradeAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	if tradeAgreement.Status == ACCEPTED {
		fmt.Printf("Trade %s already accepted", args[0])
	} else {
		tradeAgreement.Status = ACCEPTED
		tradeAgreementBytes, err = json.Marshal(tradeAgreement)
		if err != nil {
			return shim.Error("Error marshaling trade agreement structure")
		}
		// Write the state to the ledger
		err = stub.PutState(tradeKey, tradeAgreementBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	fmt.Printf("Trade %s acceptance recorded\n", args[0])

	return shim.Success(nil)
}

// Request an L/C
func (t *DcotWorkflowChaincode) requestLC(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var tradeKey, lcKey string
	var tradeAgreementBytes, letterOfCreditBytes, exporterBytes []byte
	var tradeAgreement *TradeAgreement
	var letterOfCredit *LetterOfCredit
	var err error

	// Access control: Only an DCOT operatorcan invoke this transaction
	if !t.testMode && !authenticateImporterOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Importer Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Trade ID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	// Lookup trade agreement from the ledger
	tradeKey, err = getTradeKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	tradeAgreementBytes, err = stub.GetState(tradeKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(tradeAgreementBytes) == 0 {
		err = errors.New(fmt.Sprintf("No record found for trade ID %s", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(tradeAgreementBytes, &tradeAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify that the trade has been agreed to
	if tradeAgreement.Status != ACCEPTED {
		return shim.Error("Trade has not been accepted by the parties")
	}

	// Lookup exporter (L/C beneficiary)
	exporterBytes, err = stub.GetState(expKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	letterOfCredit = &LetterOfCredit{"", "", string(exporterBytes), tradeAgreement.Amount, []string{}, REQUESTED}
	letterOfCreditBytes, err = json.Marshal(letterOfCredit)
	if err != nil {
		return shim.Error("Error marshaling letter of credit structure")
	}

	// Write the state to the ledger
	lcKey, err = getLCKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(lcKey, letterOfCreditBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Letter of Credit request for trade %s recorded\n", args[0])

	return shim.Success(nil)
}

// Issue an L/C
// We don't need to check the trade status if the L/C request has already been recorded
func (t *DcotWorkflowChaincode) issueLC(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var lcKey string
	var letterOfCreditBytes []byte
	var letterOfCredit *LetterOfCredit
	var err error

	// Access control: Only an DCOT operatorcan invoke this transaction
	if !t.testMode && !authenticateImporterOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Importer Org. Access denied.")
	}

	if len(args) < 3 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting at least 3: {Trade ID, L/C ID, Expiry Date} [List of Documents]. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	// Lookup L/C from the ledger
	lcKey, err = getLCKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	letterOfCreditBytes, err = stub.GetState(lcKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(letterOfCreditBytes, &letterOfCredit)
	if err != nil {
		return shim.Error(err.Error())
	}

	if letterOfCredit.Status == ISSUED {
		fmt.Printf("L/C for trade %s already issued", args[0])
	} else if letterOfCredit.Status == ACCEPTED {
		fmt.Printf("L/C for trade %s already accepted", args[0])
	} else {
		letterOfCredit.Id = args[1]
		letterOfCredit.ExpirationDate = args[2]
		letterOfCredit.Documents = args[3:]
		letterOfCredit.Status = ISSUED
		letterOfCreditBytes, err = json.Marshal(letterOfCredit)
		if err != nil {
			return shim.Error("Error marshaling L/C structure")
		}
		// Write the state to the ledger
		err = stub.PutState(lcKey, letterOfCreditBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	fmt.Printf("L/C issuance for trade %s recorded\n", args[0])

	return shim.Success(nil)
}

// Accept an L/C
func (t *DcotWorkflowChaincode) acceptLC(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var lcKey string
	var letterOfCreditBytes []byte
	var letterOfCredit *LetterOfCredit
	var err error

	// Access control: Only an Exporter Org member can invoke this transaction
	if !t.testMode && !authenticateExporterOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Exporter Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Trade ID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	// Lookup L/C from the ledger
	lcKey, err = getLCKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	letterOfCreditBytes, err = stub.GetState(lcKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(letterOfCreditBytes, &letterOfCredit)
	if err != nil {
		return shim.Error(err.Error())
	}

	if letterOfCredit.Status == ACCEPTED {
		fmt.Printf("L/C for trade %s already accepted", args[0])
	} else if letterOfCredit.Status == REQUESTED {
		fmt.Printf("L/C for trade %s has not been issued", args[0])
		return shim.Error("L/C not issued yet")
	} else {
		letterOfCredit.Status = ACCEPTED
		letterOfCreditBytes, err = json.Marshal(letterOfCredit)
		if err != nil {
			return shim.Error("Error marshaling L/C structure")
		}
		// Write the state to the ledger
		err = stub.PutState(lcKey, letterOfCreditBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	fmt.Printf("L/C acceptance for trade %s recorded\n", args[0])

	return shim.Success(nil)
}

// Request an E/L
func (t *DcotWorkflowChaincode) requestEL(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var tradeKey, lcKey, elKey string
	var tradeAgreementBytes, letterOfCreditBytes, exportLicenseBytes, exporterBytes, carrierBytes, approverBytes []byte
	var tradeAgreement *TradeAgreement
	var letterOfCredit *LetterOfCredit
	var exportLicense *ExportLicense
	var err error

	// Access control: Only an Exporting Entity Org member can invoke this transaction
	if !t.testMode && !authenticateExportingEntityOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Exporting Entity Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Trade ID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	// Lookup L/C from the ledger
	lcKey, err = getLCKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	letterOfCreditBytes, err = stub.GetState(lcKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(letterOfCreditBytes, &letterOfCredit)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify that the L/C has already been accepted
	if letterOfCredit.Status != ACCEPTED {
		fmt.Printf("L/C for trade %s has not been accepted", args[0])
		return shim.Error("L/C not accepted yet")
	}

	// Lookup trade agreement from the ledger
	tradeKey, err = getTradeKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	tradeAgreementBytes, err = stub.GetState(tradeKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(tradeAgreementBytes) == 0 {
		err = errors.New(fmt.Sprintf("No record found for trade ID %s", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(tradeAgreementBytes, &tradeAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Record the E/L request
	exporterBytes, err = stub.GetState(expKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Lookup exporter
	exporterBytes, err = stub.GetState(expKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Lookup carrier
	carrierBytes, err = stub.GetState(carKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Lookup regulatory authority (license approver)
	approverBytes, err = stub.GetState(raKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	exportLicense = &ExportLicense{"", "", string(exporterBytes), string(carrierBytes), tradeAgreement.DescriptionOfGoods, string(approverBytes), REQUESTED}
	exportLicenseBytes, err = json.Marshal(exportLicense)
	if err != nil {
		return shim.Error("Error marshaling export license structure")
	}

	// Write the state to the ledger
	elKey, err = getELKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(elKey, exportLicenseBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Export License request for trade %s recorded\n", args[0])

	return shim.Success(nil)
}

// Issue an E/L
func (t *DcotWorkflowChaincode) issueEL(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var elKey string
	var exportLicenseBytes []byte
	var exportLicense *ExportLicense
	var err error

	// Access control: Only a Regulator Org member can invoke this transaction
	if !t.testMode && !authenticateRegulatorOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Regulator Org. Access denied.")
	}

	if len(args) != 3 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 3: {Trade ID, L/C ID, Expiry Date}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	// Lookup E/L from the ledger
	elKey, err = getELKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	exportLicenseBytes, err = stub.GetState(elKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(exportLicenseBytes, &exportLicense)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify that the E/L has not already been issued
	if exportLicense.Status == ISSUED {
		fmt.Printf("E/L for trade %s has already been issued", args[0])
	} else {
		exportLicense.Id = args[1]
		exportLicense.ExpirationDate = args[2]
		exportLicense.Status = ISSUED
		exportLicenseBytes, err = json.Marshal(exportLicense)
		if err != nil {
			return shim.Error("Error marshaling E/L structure")
		}
		// Write the state to the ledger
		err = stub.PutState(elKey, exportLicenseBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	fmt.Printf("Export License issuance for trade %s recorded\n", args[0])

	return shim.Success(nil)
}

// Prepare a shipment; preparation is indicated by setting the location as SOURCE
func (t *DcotWorkflowChaincode) prepareShipment(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var elKey, shipmentLocationKey string
	var shipmentLocationBytes, exportLicenseBytes []byte
	var exportLicense *ExportLicense
	var err error

	// Access control: Only an Exporting Entity Org member can invoke this transaction
	if !t.testMode && !authenticateExportingEntityOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Exporting Entity Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Trade ID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	// Lookup shipment location from the ledger
	shipmentLocationKey, err = getShipmentLocationKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	shipmentLocationBytes, err = stub.GetState(shipmentLocationKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(shipmentLocationBytes) != 0 {
		if string(shipmentLocationBytes) == SOURCE {
			fmt.Printf("Shipment for trade %s has already been prepared", args[0])
			return shim.Success(nil)
		} else {
			fmt.Printf("Shipment for trade %s has passed the preparation stage", args[0])
			return shim.Error("Shipment past the preparation stage")
		}
	}

	// Lookup E/L from the ledger
	elKey, err = getELKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	exportLicenseBytes, err = stub.GetState(elKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(exportLicenseBytes, &exportLicense)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify that the E/L has already been issued
	if exportLicense.Status != ISSUED {
		fmt.Printf("E/L for trade %s has not been issued", args[0])
		return shim.Error("E/L not issued yet")
	}

	shipmentLocationKey, err = getShipmentLocationKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(shipmentLocationKey, []byte(SOURCE))
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Shipment preparation for trade %s recorded\n", args[0])

	return shim.Success(nil)
}

// Accept a shipment and issue a B/L
func (t *DcotWorkflowChaincode) acceptShipmentAndIssueBL(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var shipmentLocationKey, blKey, tradeKey string
	var shipmentLocationBytes, tradeAgreementBytes, billOfLadingBytes, exporterBytes, carrierBytes, beneficiaryBytes []byte
	var billOfLading *BillOfLading
	var tradeAgreement *TradeAgreement
	var err error

	// Access control: Only an Carrier Org member can invoke this transaction
	if !t.testMode && !authenticateCarrierOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Carrier Org. Access denied.")
	}

	if len(args) != 5 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 5: {Trade ID, B/L ID, Expiration Date, Source Port, Destination Port}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	// Lookup shipment location from the ledger
	shipmentLocationKey, err = getShipmentLocationKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	shipmentLocationBytes, err = stub.GetState(shipmentLocationKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(shipmentLocationBytes) == 0 {
		fmt.Printf("Shipment for trade %s has not been prepared yet", args[0])
		return shim.Error("Shipment not prepared yet")
	}
	if string(shipmentLocationBytes) != SOURCE {
		fmt.Printf("Shipment for trade %s has passed the preparation stage", args[0])
		return shim.Error("Shipment past the preparation stage")
	}

	// Lookup trade agreement from the ledger
	tradeKey, err = getTradeKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	tradeAgreementBytes, err = stub.GetState(tradeKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(tradeAgreementBytes) == 0 {
		err = errors.New(fmt.Sprintf("No record found for trade ID %s", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(tradeAgreementBytes, &tradeAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Lookup exporter
	exporterBytes, err = stub.GetState(expKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Lookup carrier
	carrierBytes, err = stub.GetState(carKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Lookup importer's bank (beneficiary of the title to goods after paymen tis made)
	beneficiaryBytes, err = stub.GetState(ibKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Create and record a B/L
	billOfLading = &BillOfLading{args[1], args[2], string(exporterBytes), string(carrierBytes), tradeAgreement.DescriptionOfGoods,
				     tradeAgreement.Amount, string(beneficiaryBytes), args[3], args[4]}
	billOfLadingBytes, err = json.Marshal(billOfLading)
	if err != nil {
		return shim.Error("Error marshaling bill of lading structure")
	}

	// Write the state to the ledger
	blKey, err = getBLKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(blKey, billOfLadingBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Bill of Lading for trade %s recorded\n", args[0])

	return shim.Success(nil)
}

// Request a payment
func (t *DcotWorkflowChaincode) requestPayment(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var shipmentLocationKey, paymentKey, tradeKey string
	var shipmentLocationBytes, paymentBytes, tradeAgreementBytes []byte
	var tradeAgreement *TradeAgreement
	var err error

	// Access control: Only an Exporting Entity Org member can invoke this transaction
	if !t.testMode && !authenticateExportingEntityOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Exporting Entity Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Trade ID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	// Lookup trade agreement from the ledger
	tradeKey, err = getTradeKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	tradeAgreementBytes, err = stub.GetState(tradeKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(tradeAgreementBytes) == 0 {
		err = errors.New(fmt.Sprintf("No record found for trade ID %s", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(tradeAgreementBytes, &tradeAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Lookup shipment location from the ledger
	shipmentLocationKey, err = getShipmentLocationKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	shipmentLocationBytes, err = stub.GetState(shipmentLocationKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(shipmentLocationBytes) == 0 {
		fmt.Printf("Shipment for trade %s has not been prepared yet", args[0])
		return shim.Error("Shipment not prepared yet")
	}

	// Check if there's already a pending payment request
	paymentKey, err = getPaymentKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	paymentBytes, err = stub.GetState(paymentKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(paymentBytes) != 0 {	// The value doesn't matter as this is a temporary key used as a marker
		fmt.Printf("Payment request already pending for trade %s\n", args[0])
	} else {
		// Check what has been paid up to this point
		fmt.Printf("Amount paid thus far for trade %s = %d; total required = %d\n", args[0], tradeAgreement.Payment, tradeAgreement.Amount)
		if tradeAgreement.Amount == tradeAgreement.Payment {	// Payment has already been settled
			fmt.Printf("Payment already settled for trade %s\n", args[0])
			return shim.Error("Payment already settled")
		}
		if string(shipmentLocationBytes) == SOURCE && tradeAgreement.Payment != 0 {	// Suppress duplicate requests for partial payment
			fmt.Printf("Partial payment already made for trade %s\n", args[0])
			return shim.Error("Partial payment already made")
		}

		// Record request on ledger
		err = stub.PutState(paymentKey, []byte(REQUESTED))
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Printf("Payment request for trade %s recorded\n", args[0])
	}
	return shim.Success(nil)
}

// Make a payment
func (t *DcotWorkflowChaincode) makePayment(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var shipmentLocationKey, paymentKey, tradeKey string
	var paymentAmount, expBal, impBal int
	var shipmentLocationBytes, paymentBytes, tradeAgreementBytes, impBalBytes, expBalBytes []byte
	var tradeAgreement *TradeAgreement
	var err error

	// Access control: Only an DCOT operatorcan invoke this transaction
	if !t.testMode && !authenticateImporterOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Importer Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Trade ID}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	// Check if there's already a pending payment request
	paymentKey, err = getPaymentKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	paymentBytes, err = stub.GetState(paymentKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(paymentBytes) == 0 {
		fmt.Printf("No payment request found for trade %s", args[0])
		return shim.Error("No payment request found")
	}

	// Lookup trade agreement from the ledger
	tradeKey, err = getTradeKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	tradeAgreementBytes, err = stub.GetState(tradeKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(tradeAgreementBytes) == 0 {
		err = errors.New(fmt.Sprintf("No record found for trade ID %s", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(tradeAgreementBytes, &tradeAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Lookup shipment location from the ledger
	shipmentLocationKey, err = getShipmentLocationKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	shipmentLocationBytes, err = stub.GetState(shipmentLocationKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(shipmentLocationBytes) == 0 {
		fmt.Printf("Shipment for trade %s has not been prepared yet", args[0])
		return shim.Error("Shipment not prepared yet")
	}

	// Lookup account balances
	expBalBytes, err = stub.GetState(expBalKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	expBal, err = strconv.Atoi(string(expBalBytes))
	if err != nil {
		return shim.Error(err.Error())
	}
	impBalBytes, err = stub.GetState(impBalKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	impBal, err = strconv.Atoi(string(impBalBytes))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Record transfer of funds
	if string(shipmentLocationBytes) == SOURCE {
		paymentAmount = tradeAgreement.Amount/2
	} else {
		paymentAmount = tradeAgreement.Amount - tradeAgreement.Payment
	}
	tradeAgreement.Payment += paymentAmount
	expBal += paymentAmount
	if impBal < paymentAmount {
		fmt.Printf("Importer's bank balance %d is insufficient to cover payment amount %d\n", impBal, paymentAmount)
	}
	impBal -= paymentAmount

	// Update ledger state
	tradeAgreementBytes, err = json.Marshal(tradeAgreement)
	if err != nil {
		return shim.Error("Error marshaling trade agreement structure")
	}
	err = stub.PutState(tradeKey, tradeAgreementBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(expBalKey, []byte(strconv.Itoa(expBal)))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(impBalKey, []byte(strconv.Itoa(impBal)))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Delete request key from ledger
	err = stub.DelState(paymentKey)
	if err != nil {
		fmt.Println(err.Error())
		return shim.Error("Failed to delete payment request from ledger")
	}

	return shim.Success(nil)
}

// Update shipment location; we will only allow SOURCE and DESTINATION as valid locations for this contract
func (t *DcotWorkflowChaincode) updateShipmentLocation(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var shipmentLocationKey string
	var shipmentLocationBytes []byte
	var err error

	// Access control: Only a Carrier Org member can invoke this transaction
	if !t.testMode && !authenticateCarrierOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Carrier Org. Access denied.")
	}

	if len(args) != 2 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Trade ID, Location}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	// Lookup shipment location from the ledger
	shipmentLocationKey, err = getShipmentLocationKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	shipmentLocationBytes, err = stub.GetState(shipmentLocationKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(shipmentLocationBytes) == 0 {
		fmt.Printf("Shipment for trade %s has not been prepared yet", args[0])
		return shim.Error("Shipment not prepared yet")
	}
	if string(shipmentLocationBytes) == args[1] {
		fmt.Printf("Shipment for trade %s is already in location %s", args[0], args[1])
	}

	// Write the state to the ledger
	err = stub.PutState(shipmentLocationKey, []byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Shipment location for trade %s recorded\n", args[0])

	return shim.Success(nil)
}

/*
// Deletes an entity from state
func (t *DcotWorkflowChaincode) delete(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var key string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: <key name>")
	}

	key = args[0]

	// Delete the key from the state in ledger
	err = stub.DelState(key)
	if err != nil {
		fmt.Println(err.Error())
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// Get current state of a trade agreement
func (t *DcotWorkflowChaincode) getTradeStatus(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var tradeKey, jsonResp string
	var tradeAgreement TradeAgreement
	var tradeAgreementBytes []byte
	var err error

	// Access control: Only an Importer or Exporter or Exporting Entity Org member can invoke this transaction
	if !t.testMode && !(authenticateImporterOrg(creatorOrg, creatorCertIssuer) || authenticateExporterOrg(creatorOrg, creatorCertIssuer) || authenticateExportingEntityOrg(creatorOrg, creatorCertIssuer)) {
		return shim.Error("Caller not a member of Importer or Exporter or Exporting Entity Org. Access denied.")
	}

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: <trade ID>")
	}

	// Get the state from the ledger
	tradeKey, err = getTradeKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	tradeAgreementBytes, err = stub.GetState(tradeKey)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + tradeKey + "\"}"
		return shim.Error(jsonResp)
	}

	if len(tradeAgreementBytes) == 0 {
		jsonResp = "{\"Error\":\"No record found for " + tradeKey + "\"}"
		return shim.Error(jsonResp)
	}

	// Unmarshal the JSON
	err = json.Unmarshal(tradeAgreementBytes, &tradeAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	jsonResp = "{\"Status\":\"" + tradeAgreement.Status + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(jsonResp))
}

// Get current state of a Letter of Credit
func (t *DcotWorkflowChaincode) getLCStatus(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var lcKey, jsonResp string
	var letterOfCredit LetterOfCredit
	var letterOfCreditBytes []byte
	var err error

	// Access control: Only an Importer or Exporter or Exporting Entity Org member can invoke this transaction
	if !t.testMode && !(authenticateImporterOrg(creatorOrg, creatorCertIssuer) || authenticateExporterOrg(creatorOrg, creatorCertIssuer) || authenticateExportingEntityOrg(creatorOrg, creatorCertIssuer)) {
		return shim.Error("Caller not a member of Importer or Exporter or Exporting Entity Org. Access denied.")
	}

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: <trade ID>")
	}

	// Get the state from the ledger
	lcKey, err = getLCKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	letterOfCreditBytes, err = stub.GetState(lcKey)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + lcKey + "\"}"
		return shim.Error(jsonResp)
	}

	if len(letterOfCreditBytes) == 0 {
		jsonResp = "{\"Error\":\"No record found for " + lcKey + "\"}"
		return shim.Error(jsonResp)
	}

	// Unmarshal the JSON
	err = json.Unmarshal(letterOfCreditBytes, &letterOfCredit)
	if err != nil {
		return shim.Error(err.Error())
	}

	jsonResp = "{\"Status\":\"" + letterOfCredit.Status + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(jsonResp))
}

// Get current state of an Export License
func (t *DcotWorkflowChaincode) getELStatus(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var elKey, jsonResp string
	var exportLicense ExportLicense
	var exportLicenseBytes []byte
	var err error

	// Access control: Only an Exporting Entity or Regulator Org member can invoke this transaction
	if !t.testMode && !(authenticateExportingEntityOrg(creatorOrg, creatorCertIssuer) || authenticateRegulatorOrg(creatorOrg, creatorCertIssuer)) {
		return shim.Error("Caller not a member of Exporting Entity or Regulator Org. Access denied.")
	}

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: <trade ID>")
	}

	// Get the state from the ledger
	elKey, err = getELKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	exportLicenseBytes, err = stub.GetState(elKey)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + elKey + "\"}"
		return shim.Error(jsonResp)
	}

	if len(exportLicenseBytes) == 0 {
		jsonResp = "{\"Error\":\"No record found for " + elKey + "\"}"
		return shim.Error(jsonResp)
	}

	// Unmarshal the JSON
	err = json.Unmarshal(exportLicenseBytes, &exportLicense)
	if err != nil {
		return shim.Error(err.Error())
	}

	jsonResp = "{\"Status\":\"" + exportLicense.Status + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(jsonResp))
}

// Get current location of a shipment
func (t *DcotWorkflowChaincode) getShipmentLocation(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var slKey, jsonResp string
	var shipmentLocationBytes []byte
	var err error

	// Access control: Only an Importer or Exporter or Exporting Entity or Carrier Org member can invoke this transaction
	if !t.testMode && !(authenticateImporterOrg(creatorOrg, creatorCertIssuer) || authenticateExporterOrg(creatorOrg, creatorCertIssuer) || authenticateExportingEntityOrg(creatorOrg, creatorCertIssuer) || authenticateCarrierOrg(creatorOrg, creatorCertIssuer)) {
		return shim.Error("Caller not a member of Importer or Exporter or Exporting Entity or Carrier Org. Access denied.")
	}

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: <trade ID>")
	}

	// Get the state from the ledger
	slKey, err = getShipmentLocationKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	shipmentLocationBytes, err = stub.GetState(slKey)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + slKey + "\"}"
		return shim.Error(jsonResp)
	}

	if len(shipmentLocationBytes) == 0 {
		jsonResp = "{\"Error\":\"No record found for " + slKey + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp = "{\"Location\":\"" + string(shipmentLocationBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(jsonResp))
}

// Get Bill of Lading
func (t *DcotWorkflowChaincode) getBillOfLading(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var blKey, jsonResp string
	var billOfLadingBytes []byte
	var err error

	// Access control: Only an Importer or Exporter or Exporting Entity or Carrier Org member can invoke this transaction
	if !t.testMode && !(authenticateImporterOrg(creatorOrg, creatorCertIssuer) || authenticateExporterOrg(creatorOrg, creatorCertIssuer) || authenticateExportingEntityOrg(creatorOrg, creatorCertIssuer) || authenticateCarrierOrg(creatorOrg, creatorCertIssuer)) {
		return shim.Error("Caller not a member of Importer or Exporter or Exporting Entity or Carrier Org. Access denied.")
	}

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: <trade ID>")
	}

	// Get the state from the ledger
	blKey, err = getBLKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	billOfLadingBytes, err = stub.GetState(blKey)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + blKey + "\"}"
		return shim.Error(jsonResp)
	}

	if len(billOfLadingBytes) == 0 {
		jsonResp = "{\"Error\":\"No record found for " + blKey + "\"}"
		return shim.Error(jsonResp)
	}
	fmt.Printf("Query Response:%s\n", string(billOfLadingBytes))
	return shim.Success(billOfLadingBytes)
}

// Get current account balance for a given participant
func (t *DcotWorkflowChaincode) getAccountBalance(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var entity, balanceKey, jsonResp string
	var balanceBytes []byte
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2: {Trade ID, Entity}")
	}

	entity = strings.ToLower(args[1])
	if entity == "exporter" {
		// Access control: Only an Exporter or Exporting Entity Org member can invoke this transaction
		if !t.testMode && !(authenticateExporterOrg(creatorOrg, creatorCertIssuer) || authenticateExportingEntityOrg(creatorOrg, creatorCertIssuer)) {
			return shim.Error("Caller not a member of Exporter or Exporting Entity Org. Access denied.")
		}
		balanceKey = expBalKey
	} else if entity == "importer" {
		// Access control: Only an DCOT operatorcan invoke this transaction
		if !t.testMode && !authenticateImporterOrg(creatorOrg, creatorCertIssuer) {
			return shim.Error("Caller not a member of Importer Org. Access denied.")
		}
		balanceKey = impBalKey
	} else {
		err = errors.New(fmt.Sprintf("Invalid entity %s; Permissible values: {exporter, importer}", args[1]))
		return shim.Error(err.Error())
	}

	// Get the account balances from the ledger
	balanceBytes, err = stub.GetState(balanceKey)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + balanceKey + "\"}"
		return shim.Error(jsonResp)
	}

	if len(balanceBytes) == 0 {
		jsonResp = "{\"Error\":\"No record found for " + balanceKey + "\"}"
		return shim.Error(jsonResp)
	}
	jsonResp = "{\"Balance\":\"" + string(balanceBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(jsonResp))
}
*/

func (t *DcotWorkflowChaincode) initNewChain(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {
	//TODO
	//var callerID string
	var jsonResp string
	var chainOfCustody *ChainOfCustody
	var err error
	var jsonCOC []byte
	var COCKey string
	//var chaincodeStubInterface ChaincodeStubInterface
	// Access control: Only an DCOT operatorcan invoke this transaction
	//if !t.testMode && !isEnabled {
	//	return shim.Error("Caller is not a DCOT operator.")
	//}
	//Check Args size is correct!!!
	//var cocKey string
	
	guid := xid.New()
	COCKey, err = getCOCKey(stub, guid.String())
	if err != nil {
		return shim.Error(err.Error())
	}


	err = json.Unmarshal([]byte(args[0]), &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}
	if (chainOfCustody.DocumentId == "" || len(chainOfCustody.DocumentId) == 0 ){
		return shim.Error("initNewChain ERROR: Document ID must not be null or empty string!!")

	}
	chainOfCustody.Id = guid.String()
	chainOfCustody.Status = IN_CUSTODY

	//_, callerID, err = getTxCreatorInfo(stub)
	//if err != nil {
	//	return shim.Error(err.Error())
	//}

	chainOfCustody.DeliveryMan = "admin"    //FIXME!!!! GESTIRE UTENZE PERMESSI!!!!!
	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(COCKey, jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}
	//TODO


	//jsonResp = "{\" **** initNewChain complete! ****\":\"" + string(jsonCOC) + "\"} "
	jsonResp = string(jsonCOC)
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(jsonResp))
}

func (t *DcotWorkflowChaincode) startTransfer(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {
	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte


	if len(args) != 2 {
		return shim.Error("startTransferAsset_ERROR: this method must want exactly two arguments!!")
	}

	COCKey, err = getCOCKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	// Access control: Only an DCOT operatorcan invoke this transaction
	//if !t.testMode && !isEnabled {
	//	return shim.Error("Caller is not a DCOT operator.")
	//}
	chainOfCustodyBytes, err = stub.GetState(COCKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal([]byte(chainOfCustodyBytes) , &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	if chainOfCustody.Status != IN_CUSTODY {
		return shim.Error("startTransferAsset ERROR : Asset have not status IN_CUSTODY!!")
	}

	//_, callerID, err = getTxCreatorInfo(stub)
	//if err != nil {
	//	return shim.Error(err.Error())
	//}
	//fmt.Printf("CALLER_ID :"+ string(callerID) +" . \n")

	//if callerID != chainOfCustody.DeliveryMan{
	//	return shim.Error("startTransferAsset ERROR : The caller must be the current custodian!!")
	//}


	//FIXME SEE UP!!!!
	if chainOfCustody.DeliveryMan != "admin"{
	return shim.Error("startTransferAsset ERROR : The caller must be the current custodian!!")
	}

	chainOfCustody.Status = TRANSFER_PENDING
	chainOfCustody.DeliveryMan = args[1]
	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(COCKey, jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}	

	fmt.Printf(string(jsonCOC))

	return shim.Success(nil)
}

func (t *DcotWorkflowChaincode) completeTrasfer(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {
	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	//TODO

	// Access control: Only an DCOT operatorcan invoke this transaction
	//if !t.testMode && !isEnabled {
	//	return shim.Error("Caller is not a DCOT operator.")
	//}
	//Check Args size is correct!!!
	if len(args) != 1 {
		return shim.Error("completeTrasfer ERROR: this method must want exactly one argument!!")
	}

	COCKey, err = getCOCKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	chainOfCustodyBytes, err = stub.GetState(COCKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal([]byte(chainOfCustodyBytes) , &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	if chainOfCustody.Status !=TRANSFER_PENDING {
		return shim.Error("completeTrasfer ERROR : Asset have not status TRANSFER_PENDING!!")
	}

	/*_, callerID, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	if callerID != chainOfCustody.DeliveryMan{
		return shim.Error("completeTrasfer ERROR : The caller must be the current custodian!!")
	}
	*/

	//FIXME SEE UP!!!
	if chainOfCustody.DeliveryMan != "admin"{
	return shim.Error("completeTrasfer ERROR : The caller must be the current custodian!!")
	}
	
	
	chainOfCustody.Status = IN_CUSTODY
	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(COCKey, jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}	

	return shim.Success(nil)
}

func (t *DcotWorkflowChaincode) commentChain(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {
	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	// Access control: Only an DCOT operatorcan invoke this transaction
	//if !t.testMode && !isEnabled {
	//	return shim.Error("Caller is not a DCOT operator.")
	//}

	//Check Args size is correct!!!

	if len(args) != 2{
		return shim.Error("commentChain ERROR: this method must want exactly two argument!!")
	}

	COCKey, err = getCOCKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	chainOfCustodyBytes, err = stub.GetState(COCKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal([]byte(chainOfCustodyBytes) , &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	chainOfCustody.Text = args[1]
	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	//TODO: CHECK ROLE OF CALLER!!!!

	err = shim.SetEvent("commentChain EVENT: ", jsonCOC)
	if err =! nil {
		return shim.Error(err.Error())
	}


	return shim.Success(nil)
}

func (t *DcotWorkflowChaincode) cancelTrasfer(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {
	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	//TODO

	// Access control: Only an DCOT operatorcan invoke this transaction
	//if !t.testMode && !isEnabled {
	//	return shim.Error("Caller is not a DCOT operator.")
	//}
	//Check Args size is correct!!!
	if len(args) != 1 {
		return shim.Error("cancelTrasfer ERROR: this method must want exactly one argument!!")
	}

	COCKey, err = getCOCKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	chainOfCustodyBytes, err = stub.GetState(COCKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal([]byte(chainOfCustodyBytes) , &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	if chainOfCustody.Status !=TRANSFER_PENDING {
		return shim.Error("cancelTrasfer ERROR : Asset have not status TRANSFER_PENDING!!")
	}

	/*
	_, callerID, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	if callerID != chainOfCustody.DeliveryMan{
		return shim.Error("cancelTrasfer ERROR : The caller must be the current custodian!!")
	} 
	*/

	if chainOfCustody.DeliveryMan != "admin"{
		return shim.Error("cancelTrasfer ERROR : The caller must be the current custodian!!")
	}
	chainOfCustody.Status = IN_CUSTODY
	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(COCKey, jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}	

	return shim.Success(nil)
}

func (t *DcotWorkflowChaincode) terminateChain(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {
	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	//TODO

	// Access control: Only an DCOT operatorcan invoke this transaction
	//if !t.testMode && !isEnabled {
	//	return shim.Error("Caller is not a DCOT operator.")
	//}
	//Check Args size is correct!!!
	if len(args) != 1 {
		return shim.Error("terminateChain ERROR: this method must want exactly one argument!!")
	}

	COCKey, err = getCOCKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	chainOfCustodyBytes, err = stub.GetState(COCKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal([]byte(chainOfCustodyBytes) , &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	if chainOfCustody.Status != IN_CUSTODY {
		return shim.Error("terminateChain ERROR : Asset have not status IN_CUSTODY!!")
	}

	/*
	_, callerID, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	if callerID != chainOfCustody.DeliveryMan{
		return shim.Error("terminateChain ERROR : The caller must be the current custodian!!")
	}
	*/

	//FIXME SEE UP
	if chainOfCustody.DeliveryMan != "admin"{
		return shim.Error("terminateChain ERROR : The caller must be the current custodian!!")
	}
	chainOfCustody.Status = RELEASED
	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(COCKey, jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}	

	return shim.Success(nil)
}

func (t *DcotWorkflowChaincode) updateDocument(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {
	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	var jsonResp string

	// Access control: Only an DCOT operatorcan invoke this transaction
	if !t.testMode && !isEnabled {
		return shim.Error("Caller is not a DCOT operator.")
	}

	//Check Args size is correct!!!
	if len(args) != 2 {
		return shim.Error("updateDocument ERROR: this method must want exactly two argument!!")
	}

	COCKey, err = getCOCKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	chainOfCustodyBytes, err = stub.GetState(COCKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal([]byte(chainOfCustodyBytes) , &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	if chainOfCustody.Status != IN_CUSTODY{
		return shim.Error("updateDocument ERROR: Asset's status is not IN_CUSTODY!!!")
	}
	chainOfCustody.DocumentId = args[1]

	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(COCKey, jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}

	//EVENT created	
	err = shim.SetEvent("updateDocument EVENT:", jsonCOC )
	if err =! nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
func (t *DcotWorkflowChaincode) getAssetDetails(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {
	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	var jsonResp string
	//TODO

	// Access control: Only an DCOT operatorcan invoke this transaction
	//if !t.testMode && !isEnabled {
	//	return shim.Error("Caller is not a DCOT operator.")
	//}
	//Check Args size is correct!!!
	if len(args) != 1 {
		return shim.Error("getAssetDetails ERROR: this method must want exactly one argument!!")
	}

	COCKey, err = getCOCKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	chainOfCustodyBytes, err = stub.GetState(COCKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal([]byte(chainOfCustodyBytes) , &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	/*
	_, callerID, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	*/

	//FIXME
	if chainOfCustody.DeliveryMan != "admin" {
		return shim.Error("getAssetDetails ERROR: The caller haven't a role compatible or admin privileges")
	}

	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}
	jsonResp = "{\" **** getAssetDetails complete! ****\":\"" + string(jsonCOC) + "\"} "
	fmt.Printf("Query Response:%s\n", jsonResp)
	
	return shim.Success([]byte(jsonResp))}

func (t *DcotWorkflowChaincode) getChainOfEvents(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {
	//TODO

	// Access control: Only an DCOT operatorcan invoke this transaction
	if !t.testMode && !isEnabled {
		return shim.Error("Caller is not a DCOT operator.")
	}

	//Check Args size is correct!!!

	return shim.Success(nil)
}

func main() {
	twc := new(DcotWorkflowChaincode)
	twc.testMode = true
	err := shim.Start(twc)
	if err != nil {
		fmt.Printf("Error starting Chain of Custody chaincode: %s", err)
	}
}
