package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type SwapStation struct {
	StationName    string `json:"stationName"`
	Address        string `json:"address"`
	LicenseDetails string `json:"licensedetails"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "addSwapStationDetails" {
		return s.addSwapStationDetails(APIstub, args)
	} else if function == "getStationDetails" {
		return s.getStationDetails(APIstub, args)
	} else if function == "GetSwappingStationInventory" {
		return s.GetSwappingStationInventory(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	ssDetails := []SwapStation{
		SwapStation{StationName: "ABC", Address: "street, area, city,state", LicenseDetails: "Active"},
		SwapStation{StationName: "XYZ", Address: "street, area, city,state", LicenseDetails: "Active"},
	}

	i := 0
	for i < len(ssDetails) {
		fmt.Println("i is ", i)
		ssAsBytes, _ := json.Marshal(ssDetails[i])
		APIstub.PutState("SS"+strconv.Itoa(i), ssAsBytes)
		fmt.Println("Added", ssDetails[i])
		i = i + 1
	}
	return shim.Success(nil)
}

func (s *SmartContract) addSwapStationDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect no of arguments, Expected 4")
	}
	var swapStationDetails = SwapStation{StationName: args[1], Address: args[2], LicenseDetails: args[3]}
	swpStnAsBytes, _ := json.Marshal(swapStationDetails)
	APIstub.PutState(args[0], swpStnAsBytes)
	return shim.Success(nil)
}

func (s *SmartContract) getStationDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	evSSAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(evSSAsBytes)
}

func (s *SmartContract) GetSwappingStationInventory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	return shim.Success(nil)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
