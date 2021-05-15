package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type SwapTransactions struct {
	BatteryID    string `json:"batteryid"`
	BatteryUsage string `json:"batteryusage"`
	EVOwnerID    string `json:"evownerid"`
	TotalBill    string `json:"totalbill"`
	SwapAmount   string `json:"swapamount"`
	Date         string `json:"date"`
	Status       string `json:"status"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "GetUserTransactions" {
		return s.getUserTransactions(APIstub, args)
	} else if function == "addTransaction" {
		return s.addTransaction(APIstub, args)
	} else if function == "updatePaidTransaction" {
		return s.updatePaidTransaction(APIstub, args)
	} else if function == "getTransactionsDetails" {
		return s.getTransactionsDetails(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	transDetails := []SwapTransactions{
		SwapTransactions{BatteryID: "BAT3", BatteryUsage: "30", EVOwnerID: "User1", TotalBill: "50", SwapAmount: "0", Date: "2020-11-15 15:04:05.000000000 +0000 UTC", Status: "Success"},
		SwapTransactions{BatteryID: "BAT5", BatteryUsage: "10", EVOwnerID: "User2", TotalBill: "50", SwapAmount: "10", Date: "2020-11-15 15:04:05.000000000 +0000 UTC", Status: "Success"},
	}

	i := 0
	for i < len(transDetails) {
		fmt.Println("i is ", i)
		transAsBytes, _ := json.Marshal(transDetails[i])
		APIstub.PutState("TR"+strconv.Itoa(i), transAsBytes)
		fmt.Println("Added", transDetails[i])
		i = i + 1
	}
	return shim.Success(nil)
}

func (s *SmartContract) getUserTransactions(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var queryStr = "{\"selector\":{\"evownerid\":" + args[0] + "}}"
	resultsIterator, err := APIstub.GetQueryResult(queryStr)
	defer resultsIterator.Close()
	if err != nil {
		return shim.Error("Invalid Query")
	}
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("Iterator error")
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) getTransactionsDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	evTDAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(evTDAsBytes)
}

func (s *SmartContract) addTransaction(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect no of arguments, Expected 4")
	}
	var swapTransactionDetails = SwapTransactions{BatteryID: args[1], BatteryUsage: args[2], EVOwnerID: args[3], TotalBill: "100", SwapAmount: "0", Date: time.Now().String(), Status: "Pending"}
	swpTrnAsBytes, _ := json.Marshal(swapTransactionDetails)
	APIstub.PutState(args[0], swpTrnAsBytes)
	return shim.Success(swpTrnAsBytes)
}

func (s *SmartContract) updatePaidTransaction(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect no of arguments, Expected 1")
	}

	evTDAsBytes, _ := APIstub.GetState(args[0])
	trans := SwapTransactions{}
	json.Unmarshal(evTDAsBytes, &trans)
	trans.Status = "Success"
	swpTrnAsBytes, _ := json.Marshal(trans)
	APIstub.PutState(args[0], swpTrnAsBytes)
	return shim.Success(swpTrnAsBytes)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
