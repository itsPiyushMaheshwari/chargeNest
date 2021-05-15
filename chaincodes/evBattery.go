package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type BatteryDetails struct {
	//Battery Details
	BatteryType           string `json:"batterytype"`
	Status                string `json:"status"`
	ManufacturedDate      string `json:"manufactureddate"`
	ExpiryDate            string `json:"expirydate"`
	MaxChargeCycles       string `json:"maxchargecycles"`
	CurrentChargeCycle    string `json:"currentchargecycle"`
	IdealBatteryUsage     string `json:"idealbatteryusage"`
	ManufacturerName      string `json:"manufacturername"`
	ManufacturerAddress   string `json:"manufactureraddress"`
	SecurityCertification string `json:"securitycertification"`
	CurrentUser           string `json:"currentuser"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "addBatteryDetails" {
		return s.addBatteryDetails(APIstub, args)
	} else if function == "getBatteryDetails" {
		return s.getBatteryDetails(APIstub, args)
	} else if function == "checkAndUpdateBatteryDetails" {
		return s.checkAndUpdateBatteryDetails(APIstub, args)
	} else if function == "updateBatteryCurrentUser" {
		return s.updateBatteryCurrentUser(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	batteryDetails := []BatteryDetails{
		BatteryDetails{BatteryType: "type A", ManufacturedDate: "12/12/2020", ExpiryDate: "12/12/2028", MaxChargeCycles: "8000", CurrentChargeCycle: "1", IdealBatteryUsage: "20", ManufacturerName: "Surya Manufacturers", ManufacturerAddress: "Pune", SecurityCertification: "ISO", CurrentUser: "User1"},
		BatteryDetails{BatteryType: "type A", ManufacturedDate: "12/12/2020", ExpiryDate: "12/12/2028", MaxChargeCycles: "8000", CurrentChargeCycle: "1", IdealBatteryUsage: "20", ManufacturerName: "Surya Manufacturers", ManufacturerAddress: "Pune", SecurityCertification: "ISO", CurrentUser: "SS1"},
		BatteryDetails{BatteryType: "type B", ManufacturedDate: "12/12/2020", ExpiryDate: "12/12/2028", MaxChargeCycles: "8000", CurrentChargeCycle: "1", IdealBatteryUsage: "20", ManufacturerName: "Surya Manufacturers", ManufacturerAddress: "Pune", SecurityCertification: "ISO", CurrentUser: "User2"},
		BatteryDetails{BatteryType: "type B", ManufacturedDate: "12/12/2020", ExpiryDate: "12/12/2028", MaxChargeCycles: "8000", CurrentChargeCycle: "1", IdealBatteryUsage: "20", ManufacturerName: "Nila Manufacturers", ManufacturerAddress: "Lucknow", SecurityCertification: "ISO", CurrentUser: "SS2"},
		BatteryDetails{BatteryType: "type A", ManufacturedDate: "12/12/2020", ExpiryDate: "12/12/2028", MaxChargeCycles: "8000", CurrentChargeCycle: "1", IdealBatteryUsage: "20", ManufacturerName: "Nila Manufacturers", ManufacturerAddress: "Lucknow", SecurityCertification: "ISO", CurrentUser: "User3"},
	}

	i := 0
	for i < len(batteryDetails) {
		fmt.Println("i is ", i)
		carAsBytes, _ := json.Marshal(batteryDetails[i])
		APIstub.PutState("BAT"+strconv.Itoa(i), carAsBytes)
		fmt.Println("Added", batteryDetails[i])
		i = i + 1
	}
	return shim.Success(nil)
}

func (s *SmartContract) addBatteryDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 11 {
		return shim.Error("Incorrect no of arguments, Expected 11")
	}
	var batteryDetails = BatteryDetails{BatteryType: args[1], Status: "ACTIVE", ManufacturedDate: args[2], ExpiryDate: args[3], MaxChargeCycles: args[4], CurrentChargeCycle: args[5], IdealBatteryUsage: args[6], ManufacturerName: args[7], ManufacturerAddress: args[8], SecurityCertification: args[9], CurrentUser: args[10]}
	batdetsAsBytes, _ := json.Marshal(batteryDetails)
	APIstub.PutState(args[0], batdetsAsBytes)
	return shim.Success(nil)
}

func (s *SmartContract) getBatteryDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	batAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(batAsBytes)
}

func (s *SmartContract) checkAndUpdateBatteryDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//arg[0] -> Battery Id
	//arg[1] -> Swap Station Id

	batAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Invalid BatteryId")
	}
	batteryDets := BatteryDetails{}
	json.Unmarshal(batAsBytes, &batteryDets)
	batteryDets.CurrentUser = args[1]
	isExpired := checkBatteryExpired(batteryDets)
	if !isExpired {
		batteryDets.ExpiryDate = "Expired"
	}

	return shim.Success(nil)
}

func (s *SmartContract) updateBatteryCurrentUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//arg[0] -> Battery Id
	//arg[1] -> User Id
	batAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Invalid BatteryId")
	}
	batteryDets := BatteryDetails{}
	json.Unmarshal(batAsBytes, &batteryDets)
	batteryDets.CurrentUser = args[1]
	batdetsAsBytes, _ := json.Marshal(batteryDets)
	APIstub.PutState(args[0], batdetsAsBytes)
	return shim.Success(nil)
}

func checkBatteryExpired(batteryDets BatteryDetails) bool {
	now := time.Now()
	t, _ := time.Parse("2006-01-02 15:04:05.000000000 +0000 UTC", batteryDets.ExpiryDate)
	return now.Before(t)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
