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
	"bytes"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/rs/xid"
)

var logger = shim.NewLogger("dcot-chaincode-log")

//var logger = shim.NewLogger("dcot-chaincode")

// DcotWorkflowChaincode implementation
type DcotWorkflowChaincode struct {
	testMode bool
}

func (t *DcotWorkflowChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	logger.Info("Initializing Chain of Custody")
	logger.SetLevel(shim.LogDebug)
	_, args := stub.GetFunctionAndParameters()
	//var err error

	// Upgrade Mode 1: leave ledger state as it was
	if len(args) == 0 {
		//logger.Info("Args correctly!!!")
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
			logger.Info("Exporter's account balance must be an integer. Found \n", args[2])
			return shim.Error(err.Error())
		}
		_, err = strconv.Atoi(string(args[5]))
		if err != nil {
			logger.Info("Importer's account balance must be an integer. Found \n", args[5])
			return shim.Error(err.Error())
		}

		logger.Info("Exporter: \n", args[0])
		logger.Info("Exporter's Bank: \n", args[1])
		logger.Info("Exporter's Account Balance: \n", args[2])
		logger.Info("Importer: \n", args[3])
		logger.Info("Importer's Bank: \n", args[4])
		logger.Info("Importer's Account Balance: \n", args[5])
		logger.Info("Carrier: \n", args[6])
		logger.Info("Regulatory Authority: \n", args[7])

		// Map participant identities to their roles on the ledger
		roleKeys := []string{expKey, ebKey, expBalKey, impKey, ibKey, impBalKey, carKey, raKey}
		for i, roleKey := range roleKeys {
			err = stub.PutState(roleKey, []byte(args[i]))
			if err != nil {
				logger.Error("Error recording key : \n", roleKey, err.Error())
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

	logger.Debug("DcotWorkflow Invoke\n")

	if !t.testMode {
		creatorOrg, creatorCertIssuer, err = getTxCreatorInfo(stub)
		if err != nil {
			logger.Error("Error extracting creator identity info: \n", err.Error())
			return shim.Error(err.Error())
		}
		logger.Info("DcotWorkflow Invoke by '', ''\n", creatorOrg, creatorCertIssuer)

		isEnabled, _, err = isInvokerOperator(stub, "dcot-operator")
		if err != nil {
			logger.Error("Error getting attribute info: \n", err.Error())
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
	logger.Info("Trade  request recorded\n", args[0])

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
		err = errors.New(fmt.Sprintf("No record found for trade ID ", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(tradeAgreementBytes, &tradeAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	if tradeAgreement.Status == ACCEPTED {
		logger.Info("Trade  already accepted", args[0])
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
	logger.Info("Trade  acceptance recorded\n", args[0])

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
		err = errors.New(fmt.Sprintf("No record found for trade ID ", args[0]))
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
	logger.Info("Letter of Credit request for trade  recorded\n", args[0])

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
		logger.Info("L/C for trade  already issued", args[0])
	} else if letterOfCredit.Status == ACCEPTED {
		logger.Info("L/C for trade  already accepted", args[0])
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
	logger.Info("L/C issuance for trade  recorded\n", args[0])

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
		logger.Info("L/C for trade  already accepted", args[0])
	} else if letterOfCredit.Status == REQUESTED {
		logger.Info("L/C for trade  has not been issued", args[0])
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
	logger.Info("L/C acceptance for trade  recorded\n", args[0])

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
		logger.Info("L/C for trade  has not been accepted", args[0])
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
		err = errors.New(fmt.Sprintf("No record found for trade ID ", args[0]))
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
	logger.Info("Export License request for trade  recorded\n", args[0])

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
		logger.Info("E/L for trade  has already been issued", args[0])
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
	logger.Info("Export License issuance for trade  recorded\n", args[0])

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
			logger.Info("Shipment for trade  has already been prepared", args[0])
			return shim.Success(nil)
		} else {
			logger.Info("Shipment for trade  has passed the preparation stage", args[0])
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
		logger.Info("E/L for trade  has not been issued", args[0])
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
	logger.Info("Shipment preparation for trade  recorded\n", args[0])

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
		logger.Info("Shipment for trade  has not been prepared yet", args[0])
		return shim.Error("Shipment not prepared yet")
	}
	if string(shipmentLocationBytes) != SOURCE {
		logger.Info("Shipment for trade  has passed the preparation stage", args[0])
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
		err = errors.New(fmt.Sprintf("No record found for trade ID ", args[0]))
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
	logger.Info("Bill of Lading for trade  recorded\n", args[0])

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
		err = errors.New(fmt.Sprintf("No record found for trade ID ", args[0]))
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
		logger.Info("Shipment for trade  has not been prepared yet", args[0])
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
		logger.Info("Payment request already pending for trade \n", args[0])
	} else {
		// Check what has been paid up to this point
		logger.Info("Amount paid thus far for trade  = %d; total required = %d\n", args[0], tradeAgreement.Payment, tradeAgreement.Amount)
		if tradeAgreement.Amount == tradeAgreement.Payment {	// Payment has already been settled
			logger.Info("Payment already settled for trade \n", args[0])
			return shim.Error("Payment already settled")
		}
		if string(shipmentLocationBytes) == SOURCE && tradeAgreement.Payment != 0 {	// Suppress duplicate requests for partial payment
			logger.Info("Partial payment already made for trade \n", args[0])
			return shim.Error("Partial payment already made")
		}

		// Record request on ledger
		err = stub.PutState(paymentKey, []byte(REQUESTED))
		if err != nil {
			return shim.Error(err.Error())
		}
		logger.Info("Payment request for trade  recorded\n", args[0])
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
		logger.Info("No payment request found for trade ", args[0])
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
		err = errors.New(fmt.Sprintf("No record found for trade ID ", args[0]))
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
		logger.Info("Shipment for trade  has not been prepared yet", args[0])
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
		logger.Info("Importer's bank balance %d is insufficient to cover payment amount %d\n", impBal, paymentAmount)
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
		logger.Info(err.Error())
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
		logger.Info("Shipment for trade  has not been prepared yet", args[0])
		return shim.Error("Shipment not prepared yet")
	}
	if string(shipmentLocationBytes) == args[1] {
		logger.Info("Shipment for trade  is already in location ", args[0], args[1])
	}

	// Write the state to the ledger
	err = stub.PutState(shipmentLocationKey, []byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("Shipment location for trade  recorded\n", args[0])

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
		logger.Info(err.Error())
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
	logger.Info("Query Response:\n", jsonResp)
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
	logger.Info("Query Response:\n", jsonResp)
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
	logger.Info("Query Response:\n", jsonResp)
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
	logger.Info("Query Response:\n", jsonResp)
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
	logger.Info("Query Response:\n", string(billOfLadingBytes))
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
		err = errors.New(fmt.Sprintf("Invalid entity ; Permissible values: {exporter, importer}", args[1]))
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
	logger.Info("Query Response:\n", jsonResp)
	return shim.Success([]byte(jsonResp))
}
*/

func (t *DcotWorkflowChaincode) initNewChain(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {

	logger.Debug("***start initNewChain***")
	//TODO
	//var callerID string
	var jsonResp string
	var chainOfCustody *ChainOfCustody
	var err error
	var jsonCOC []byte
	var COCKey string
	var callerRole, callerUID string

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
	if chainOfCustody.DocumentId == "" || len(chainOfCustody.DocumentId) == 0 {
		return shim.Error("initNewChain ERROR: Document ID must not be null or empty string!!\n")

	}
	chainOfCustody.Id = guid.String()
	chainOfCustody.Status = IN_CUSTODY

	callerRole, callerUID, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("caller_UID :" + string(callerUID) + "\n")
	logger.Info("caller_ROLE :" + string(callerRole) + "\n")

	if len(callerUID) == 0 {
		return shim.Error("initNewChain ERROR: caller_UID is empty!!!\n")
	}
	chainOfCustody.DeliveryMan = string(callerUID)
	logger.Info("initNewChain: field DelivaryMan: " + string(callerUID) + "\n")
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
	logger.Info("Query Response:\n", jsonResp)

	err = stub.SetEvent("initNewChain EVENT: ", jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}

	logger.Info("initNewChain EVENT: ", string(jsonCOC))

	logger.Debug("***end initNewChain***")

	return shim.Success([]byte(jsonResp))
}

func (t *DcotWorkflowChaincode) startTransfer(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {

	logger.Debug("***start startTransfer***")

	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	var callerRole, callerUID string

	if len(args) != 2 {
		return shim.Error("startTransfer ERROR: this method must want exactly two arguments!!\n")
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

	err = json.Unmarshal([]byte(chainOfCustodyBytes), &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	if chainOfCustody.Status != IN_CUSTODY {
		return shim.Error("startTransferAsset ERROR : Asset have not status IN_CUSTODY!!\n")
	}

	callerRole, callerUID, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	//logger.Info("caller_UID :"+ string(callerUID) +" . \n")
	logger.Info("caller_ROLE :" + string(callerRole) + " . \n")

	if callerUID != chainOfCustody.DeliveryMan {
		return shim.Error("startTransferAsset ERROR : The caller must be the current custodian!!\n")
	}
	logger.Info("startTransferAsset: Ok! Caller confirmed!!\n")

	//FIXME SEE UP!!!!
	//if chainOfCustody.DeliveryMan != "5a9654f5-ff72-49dd-9be3-b3b524228556" {
	//	return shim.Error("startTransferAsset ERROR : The caller must be the current custodian!!")
	//}

	chainOfCustody.Status = TRANSFER_PENDING
	chainOfCustody.DeliveryMan = args[1]
	logger.Info("startTransferAsset: New DeliveryMan: \n", chainOfCustody.DeliveryMan)
	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(COCKey, jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("startTransfer EVENT: ", jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("startTransfer EVENT: ", string(jsonCOC))

	logger.Debug("***end startTransfer***")

	return shim.Success(nil)
}

func (t *DcotWorkflowChaincode) completeTrasfer(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {

	logger.Debug("***start completeTrasfer***")

	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	var callerRole, callerUID string

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

	err = json.Unmarshal([]byte(chainOfCustodyBytes), &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	if chainOfCustody.Status != TRANSFER_PENDING {
		return shim.Error("completeTrasfer ERROR : Asset have not status TRANSFER_PENDING!!")
	}

	callerRole, callerUID, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	//logger.Info("caller_UID :"+ string(callerUID) +" . \n")
	logger.Info("caller_ROLE :"+ string(callerRole) +" . \n")

	if callerUID != chainOfCustody.DeliveryMan{
		return shim.Error("completeTrasfer ERROR : The caller must be the current custodian!!\n")
	}
	logger.Info("completeTrasfer: Ok! Caller confirmed!!\n")

	chainOfCustody.Status = IN_CUSTODY
	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(COCKey, jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("completeTrasfer EVENT: ", jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("*** completeTrasfer EVENT: ", string(jsonCOC))

	return shim.Success(nil)
}

func (t *DcotWorkflowChaincode) commentChain(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {

	logger.Debug("***start commentChain***")

	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	var  callerRole string
	// Access control: Only an DCOT operatorcan invoke this transaction
	//if !t.testMode && !isEnabled {
	//	return shim.Error("Caller is not a DCOT operator.")
	//}

	//Check Args size is correct!!!

	if len(args) != 2 {
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

	err = json.Unmarshal([]byte(chainOfCustodyBytes), &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(COCKey, chainOfCustodyBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	callerRole, _, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("caller_ROLE :"+ string(callerRole) +" . \n")

	if callerRole !="dcot-operator"{
		return shim.Error("completeTrasfer ERROR : The caller must be a dcot-operator or network admin!!\n")
	}
	logger.Info("commentChain: Ok! Caller confirmed!!\n")


	chainOfCustody.Text = args[1]
	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("commentChain EVENT: ", jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("commentChain EVENT: ", string(jsonCOC))

	logger.Debug("***end commentChain***")

	return shim.Success(nil)
}

func (t *DcotWorkflowChaincode) cancelTrasfer(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {

	logger.Debug("***start cancelTrasfer***")

	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	var callerUID, callerRole string
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

	err = json.Unmarshal([]byte(chainOfCustodyBytes), &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	if chainOfCustody.Status != TRANSFER_PENDING {
		return shim.Error("cancelTrasfer ERROR : Asset have not status TRANSFER_PENDING!!")
	}

	callerRole, callerUID, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("caller_UID :"+ string(callerUID) +" . \n")
	logger.Info("caller_ROLE :"+ string(callerRole) +" . \n")

	if callerUID != chainOfCustody.DeliveryMan ||  callerRole != "dcot-operator"{
		return shim.Error("cancelTrasfer ERROR : The caller must be the current custodian or dcot-operator!!\n")
	}

	logger.Info("cancelTrasfer: Ok! Caller confirmed!!\n")

	chainOfCustody.Status = IN_CUSTODY
	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(COCKey, jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.SetEvent("cancelTrasfer EVENT: ", jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("cancelTrasfer EVENT: ", string(jsonCOC))

	logger.Debug("***end cancelTrasfer***")

	return shim.Success(nil)
}

func (t *DcotWorkflowChaincode) terminateChain(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {

	logger.Debug("***start terminateChain***")

	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	var callerUID string
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

	err = json.Unmarshal([]byte(chainOfCustodyBytes), &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	if chainOfCustody.Status != IN_CUSTODY {
		return shim.Error("terminateChain ERROR : Asset have not status IN_CUSTODY!!")
	}

	
	_, callerUID, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("caller_UID :"+ string(callerUID) +" . \n")

	if callerUID != chainOfCustody.DeliveryMan {
		return shim.Error("terminateChain ERROR : The caller must be the current!!\n")
	}

	logger.Info("terminateChain: Ok! Caller confirmed!!\n")

	chainOfCustody.Status = RELEASED
	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(COCKey, jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("terminateChain EVENT: ", jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("terminateChain EVENT: ", string(jsonCOC))

	logger.Debug("***end terminateChain***")

	return shim.Success(nil)
}

func (t *DcotWorkflowChaincode) updateDocument(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {

	logger.Debug("***start updateDocument***")

	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	var jsonResp string
	var callerUID, callerRole string

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

	err = json.Unmarshal([]byte(chainOfCustodyBytes), &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}

	callerRole, callerUID, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("caller_UID :"+ string(callerUID) +" . \n")
	logger.Info("caller_ROLE :"+ string(callerRole) +" . \n")

	if callerUID != chainOfCustody.DeliveryMan || callerRole != "dcot-operator"{
		return shim.Error("updateDocument ERROR : The caller must be the current custodian or dcot-operator!!\n")
	}

	logger.Info("updateDocument: Ok! Caller confirmed!!\n")

	if chainOfCustody.Status != IN_CUSTODY {
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
	err = stub.SetEvent("updateDocument EVENT:", jsonCOC)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("updateDocument EVENT: ", string(jsonCOC))

	jsonResp = string(jsonCOC)
	logger.Info("Query Response:\n", jsonResp)

	logger.Debug("***end updateDocument***")

	return shim.Success([]byte(jsonResp))
}

func (t *DcotWorkflowChaincode) getAssetDetails(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {

	logger.Debug("***start getAssetDetails***")

	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	var jsonResp string
	var callerRole string
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

	err = json.Unmarshal([]byte(chainOfCustodyBytes), &chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}


	callerRole, _, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("caller_ROLE :"+ string(callerRole) +" . \n")

	if callerRole !="dcot-operator"{
		return shim.Error("getAssetDetails ERROR : The caller must be a dcot-operator or network admin!!\n")
	}
	logger.Info("getAssetDetails: Ok! Caller confirmed!!\n")


	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}
	//jsonResp = "{\" **** getAssetDetails complete! ****\":\"" + string(jsonCOC) + "\"} "
	jsonResp = string(jsonCOC)
	logger.Info("Query Response:\n", jsonResp)

	logger.Debug("***end getAssetDetails***")

	return shim.Success([]byte(jsonResp))
}

func (t *DcotWorkflowChaincode) getChainOfEvents(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {

	logger.Debug("***start getChainOfEvents***")

	var COCKey string
	var err2 error
	var chainOfCustody *ChainOfCustody
	//var COCarray []*ChainOfCustody
	//var chainOfCustodyBytes []byte
	var jsonCOC []byte
	var jsonResp, jsonResponse string
	var callerRole string
	var err error
	//var history *HistoryQueryIteratorInterface
	// Access control: Only an DCOT operatorcan invoke this transaction
	//if !t.testMode && !isEnabled {
	//	return shim.Error("Caller is not a DCOT operator.")
	//}

	if len(args) != 1 {
		return shim.Error("getChainOfEvents ERROR: this method must want exactly one argument!!")
	}

	callerRole, _, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("caller_ROLE :"+ string(callerRole) +" . \n")

	if callerRole !="dcot-operator"{
		return shim.Error("getChainOfEvents ERROR : The caller must be a dcot-operator or network admin!!\n")
	}
	logger.Info("getChainOfEvents: Ok! Caller confirmed!!\n")


	COCKey, err2 = getCOCKey(stub, args[0])
	if err2 != nil {
		return shim.Error(err2.Error())
	}

	historyResponse, err3 := stub.GetHistoryForKey(COCKey)
	if err3 != nil {
		return shim.Error(err3.Error())
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")

	for historyResponse.HasNext() {
		COCarray, err1 := historyResponse.Next()
		if err1 != nil {
			return shim.Error(err1.Error())
		}
		//logger.Debug("COCarray :", string(COCarray))
		//jsonCOC, err2 = json.Marshal(&COCarray.Value)
		err = json.Unmarshal([]byte(COCarray.Value), &chainOfCustody)
		if err != nil {
			return shim.Error(err.Error())
		}
		jsonCOC, err2 = json.Marshal(&chainOfCustody)
		if err2 != nil {
			return shim.Error(err2.Error())
		}
		logger.Debug("jsonCOC :", string(jsonCOC))
		buffer.WriteString(string(jsonCOC))
		buffer.WriteString(",")

	}
	jsonResp = buffer.String()
	subString := jsonResp[0 : len(jsonResp)-1]
	jsonResponse = subString + "]"
	logger.Debug("Query Response:\n" + jsonResponse)

	logger.Debug("***end getChainOfEvents***")

	return shim.Success([]byte(jsonResponse))
}

func main() {
	twc := new(DcotWorkflowChaincode)
	twc.testMode = true
	err := shim.Start(twc)
	if err != nil {
		logger.Error("Error starting Chain of Custody chaincode: ", err)
	}
}
