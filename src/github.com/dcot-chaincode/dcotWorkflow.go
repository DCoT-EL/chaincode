
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

	return shim.Success(nil)
}

func (t *DcotWorkflowChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	var creatorOrg, creatorCertIssuer string
	//var attrValue string
	var err error
	var isEnabled bool
	var callerRole string

	logger.Debug("DcotWorkflow Invoke\n")

	if !t.testMode {
		creatorOrg, creatorCertIssuer, err = getTxCreatorInfo(stub)
		if err != nil {
			logger.Error("Error extracting creator identity info: \n", err.Error())
			return shim.Error(err.Error())
		}
		logger.Info("DcotWorkflow Invoke by '', ''\n", creatorOrg, creatorCertIssuer)
		callerRole, _, err = getTxCreatorInfo(stub)
		if err != nil {
			return shim.Error(err.Error())
		}

		isEnabled, _, err = isInvokerOperator(stub, callerRole)
		if err != nil {
			logger.Error("Error getting attribute info: \n", err.Error())
			return shim.Error(err.Error())
		}
	}

	function, args := stub.GetFunctionAndParameters()
	
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

func (t *DcotWorkflowChaincode) initNewChain(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {

	logger.Debug("***start initNewChain***")
	//TODO
	//var callerID string
	var jsonResp string
	var chainOfCustody ChainOfCustody
	var err error
	var jsonCOC []byte
	var COCKey string
	var callerRole, callerUID string
	var operation string
	var event Event
	//moment := time.Now()
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
	operation = "initNewChain"
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
	
	event, err = createEvent(stub, callerUID, callerRole, operation)
	if err !=nil {
		return shim.Error(err.Error())
	}
	chainOfCustody.Event = event
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
	var chainOfCustody ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	var callerRole, callerUID string
	var operation string
	var event Event

	if len(args) != 2 {
		return shim.Error("startTransfer ERROR: this method must want exactly two arguments!!\n")
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
		return shim.Error("startTransferAsset ERROR : Asset have not status IN_CUSTODY!!\n")
	}
	callerRole, callerUID, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	if callerUID != chainOfCustody.DeliveryMan {
		return shim.Error("startTransferAsset ERROR : The caller must be the current custodian!!\n")
	}
	//logger.Info("startTransferAsset: Ok! Caller confirmed!!\n")
	operation = "startTransfer"
	chainOfCustody.Status = TRANSFER_PENDING
	chainOfCustody.DeliveryMan = args[1]
	event, err = createEvent(stub, callerUID, callerRole, operation)
	if err !=nil {
		return shim.Error(err.Error())
	}
	chainOfCustody.Event = event
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
	var operation string
	var event Event

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
	//logger.Info("caller_ROLE :"+ string(callerRole) +" . \n")
	//logger.Info("DeliveryMan :"+ string(chainOfCustody.DeliveryMan) +" . \n")

	if callerUID != chainOfCustody.DeliveryMan{
		return shim.Error("completeTrasfer ERROR : The caller must be the current custodian!!\n")
	}
	logger.Info("completeTrasfer: Ok! Caller confirmed!!\n")
	operation = "completeTrasfer"
	chainOfCustody.Status = IN_CUSTODY
	event, err = createEvent(stub, callerUID, callerRole, operation)
	if err !=nil {
		return shim.Error(err.Error())
	}
	chainOfCustody.Event = event
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
	var callerUID string
	var  callerRole string
	var operation string
	var event Event

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
	callerRole, callerUID, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("caller_ROLE :"+ string(callerRole) +" . \n")

	if (callerRole == CALLER_ROLE_1|| callerRole == CALLER_ROLE_2){
	
	logger.Info("commentChain: Ok! Caller confirmed!!\n")

	operation = "commentChain"
	chainOfCustody.Text = args[1]
	event, err = createEvent(stub, callerUID, callerRole, operation)
	if err !=nil {
		return shim.Error(err.Error())
	}
	chainOfCustody.Event = event
	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(COCKey, jsonCOC)
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
	return shim.Error("completeTrasfer ERROR : The caller must be a dcot-operator or network admin!!\n")

}

func (t *DcotWorkflowChaincode) cancelTrasfer(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {

	logger.Debug("***start cancelTrasfer***")

	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	var callerUID, callerRole string
	var operation string
	var event Event

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

	if (callerUID == chainOfCustody.DeliveryMan ||  callerRole == CALLER_ROLE_1 ||callerRole != CALLER_ROLE_2){
	logger.Info("cancelTrasfer: Ok! Caller confirmed!!\n")
	operation = "cancelTrasfer"
	chainOfCustody.Status = IN_CUSTODY
	event, err = createEvent(stub, callerUID, callerRole, operation)
	if err !=nil {
		return shim.Error(err.Error())
	}
	chainOfCustody.Event = event
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

	return shim.Error("cancelTrasfer ERROR : The caller must be the current custodian or have a admin/operator role!!\n")

}

func (t *DcotWorkflowChaincode) terminateChain(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {

	logger.Debug("***start terminateChain***")

	var COCKey string
	var err error
	var chainOfCustody *ChainOfCustody
	var chainOfCustodyBytes []byte
	var jsonCOC []byte
	var callerUID, callerRole string
	var operation string
	var event Event

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

	operation = "terminateChain"	
	callerRole, callerUID, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	//logger.Info("caller_UID :"+ string(callerUID) +" . \n")

	if callerUID != chainOfCustody.DeliveryMan {
		return shim.Error("terminateChain ERROR : The caller must be the current!!\n")
	}

	logger.Info("terminateChain: Ok! Caller confirmed!!\n")

	chainOfCustody.Status = RELEASED
	event, err = createEvent(stub, callerUID, callerRole, operation)
	if err !=nil {
		return shim.Error(err.Error())
	}
	chainOfCustody.Event = event
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
	var operation string
	var event Event

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

	logger.Info("caller_ROLE :"+ string(callerRole) +" . \n")

	if (callerUID == chainOfCustody.DeliveryMan ||  callerRole ==CALLER_ROLE_1 || callerRole != CALLER_ROLE_2){
	logger.Info("updateDocument: Ok! Caller confirmed!!\n")

	if chainOfCustody.Status != IN_CUSTODY {
		return shim.Error("updateDocument ERROR: Asset's status is not IN_CUSTODY!!!")
	}
	operation = "updateDocument"
	chainOfCustody.DocumentId = args[1]
	event, err = createEvent(stub, callerUID, callerRole, operation)
	if err !=nil {
		return shim.Error(err.Error())
	}
	chainOfCustody.Event = event
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
return shim.Error("cancelTrasfer ERROR : The caller must be the current custodian or dcot-operator/admin!!\n")

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

	if (callerRole == CALLER_ROLE_1 || callerRole == CALLER_ROLE_2){
	
	logger.Info("getAssetDetails: Ok! Caller confirmed!!\n")


	jsonCOC, err = json.Marshal(&chainOfCustody)
	if err != nil {
		return shim.Error(err.Error())
	}
	jsonResp = string(jsonCOC)
	logger.Info("Query Response:\n", jsonResp)

	logger.Debug("***end getAssetDetails***")

	return shim.Success([]byte(jsonResp))}
	return shim.Error("getAssetDetails ERROR : The caller must be a dcot-operator or network admin!!\n")

}

func (t *DcotWorkflowChaincode) getChainOfEvents(stub shim.ChaincodeStubInterface, isEnabled bool, args []string) pb.Response {

	logger.Debug("***start getChainOfEvents***")

	var COCKey string
	var err2 error
	var chainOfCustody *ChainOfCustody
	var jsonCOC []byte
	var jsonResp, jsonResponse string
	var callerRole string
	var err error
	
	if len(args) != 1 {
		return shim.Error("getChainOfEvents ERROR: this method must want exactly one argument!!")
	}    
	callerRole, _, err = getTxCreatorInfo(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("caller_ROLE :"+ string(callerRole) +" . \n")

	if (callerRole == CALLER_ROLE_1 || callerRole == CALLER_ROLE_2){
	
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
	return shim.Error("getChainOfEvents ERROR : The caller must be a dcot-operator or dcot-admin!!\n")

}

func main() {
	twc := new(DcotWorkflowChaincode)
	twc.testMode = true
	err := shim.Start(twc)
	if err != nil {
		logger.Error("Error starting Chain of Custody chaincode: ", err)
	}
}
