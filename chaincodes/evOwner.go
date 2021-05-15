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

type EVOwner struct {
	Name                string `json:"name"`
	SubscriptionDetails string `json:"subscriptiondetails"`
	SubscriptionStatus  string `json:"subscriptionstatus"`
	CurrentBattery      string `json:"currentbattery"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "addEVOwnerDetails" {
		return s.addEVOwnerDetails(APIstub, args)
	} else if function == "getEVOwnerDetails" {
		return s.getEVOwnerDetails(APIstub, args)
	} else if function == "AuthenticateUser" {
		return s.AuthenticateUser(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	evOwnerDetails := []EVOwner{
		EVOwner{Name: "Karthik", SubscriptionDetails: "Starter Pack", SubscriptionStatus: "Active", CurrentBattery: "BAT1"},
		EVOwner{Name: "Vidhya", SubscriptionDetails: "Starter Pack", SubscriptionStatus: "Active", CurrentBattery: "BAT3"},
		EVOwner{Name: "Adhi", SubscriptionDetails: "Starter Pack", SubscriptionStatus: "Active", CurrentBattery: "BAT5"},
	}

	i := 0
	for i < len(evOwnerDetails) {
		fmt.Println("i is ", i)
		evoAsBytes, _ := json.Marshal(evOwnerDetails[i])
		APIstub.PutState("User"+strconv.Itoa(i), evoAsBytes)
		fmt.Println("Added", evOwnerDetails[i])
		i = i + 1
	}
	return shim.Success(nil)
}

func (s *SmartContract) addEVOwnerDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect no of arguments, Expected 5")
	}
	var evOwnerDetails = EVOwner{Name: args[1], SubscriptionDetails: args[2], SubscriptionStatus: args[3], CurrentBattery: args[4]}
	evOwnerDetsAsBytes, _ := json.Marshal(evOwnerDetails)
	APIstub.PutState(args[0], evOwnerDetsAsBytes)
	return shim.Success(nil)
}

func (s *SmartContract) getEVOwnerDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	evOAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(evOAsBytes)
}

//Authenticate User
func (s *SmartContract) AuthenticateUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	_, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Invalid User")
	} else {
		return shim.Success([]byte("True"))
	}
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
