/*******************************************************************************
Copyright (c) 2016 IBM Corporation and other Contributors.


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and limitations under the License.


Contributors:

Sumabala Nair - Initial Contribution
Kim Letkeman - Initial Contribution
Sumabala Nair - Updated for hyperledger May 2016
Sumabala Nair - Partial updates added May 2016
******************************************************************************/
//SN: March 2016

// IoT Blockchain Simple Smart Contract v 1.0

// This is a simple contract that creates a CRUD interface to 
// create, read, update and delete an asset

package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "strings"
     "reflect"
     "net/http"
     "bytes"
     "io/ioutil"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

const CONTRACTSTATEKEY string = "ContractStateKey"  
// store contract state - only version in this example
const MYVERSION string = "1.0"

// ************************************
// asset and contract state 
// ************************************

type ContractState struct {
    Version      string                        `json:"version"`
}

 type Response struct{
     Result  Result  `json:"result,omitempty"` 
     Id      int64      `json:"id,omitempty"` 
 }
type Result struct{
   Status   string `json:"status,omitempty"`
   Message string `json:"message,omitempty"`
}
type Geolocation struct {
    Latitude    *string `json:"latitude,omitempty"`
    Longitude   *string `json:"longitude,omitempty"`
}

type AssetState struct {
    AssetID         		*string       `json:"assetID,omitempty"`        // all assets must have an ID, primary key of contract
    Location        		*Geolocation  `json:"location,omitempty"`       // current asset location
    Status          		*string       `json:"kitstatus,omitempty"`        // the name of the carrier
	Role            		*string       `json:"role,omitempty"`
	Lastowner				*string 	  `json:"lastowner,omitempty"`
	Ownername            	*string       `json:"ownername,omitempty"`
	Ownerid            		*string       `json:"ownerid,omitempty"`
	Overallstatus           *string       `json:"overallstatus,omitempty"`
    AlertType               *string       `json:"alerttype,omitempty"`
    Alert                   *string       `json:"alert,omitempty"`

	
}
type Message123 struct {
    AssetID         		*string       `json:"assetID,omitempty"`        // all assets must have an ID, primary key of contract
    Location        		*Geolocation  `json:"location,omitempty"`       // current asset location
    Status          		*string       `json:"kitstatus,omitempty"`        // the name of the carrier
	Role            		*string       `json:"role,omitempty"`
	Lastowner				*string 	  `json:"lastowner,omitempty"`
	Ownername            	*string       `json:"ownername,omitempty"`
	Ownerid            		*string       `json:"ownerid,omitempty"`
	Overallstatus           *string       `json:"overallstatus,omitempty"`
	
}


type AutoGenerated struct {
	Jsonrpc string 
	Method string 
	Params Params  `json:"params"`
	Id int 
}

type Params struct {
		Type int 
		ChaincodeID ChaincodeID `json:"ChaincodeId"`
		CtorMsg CtorMsg `json:"CtorMsg"`
		SecureContext string 
	}

type ChaincodeID struct {
			Name string 
		}
type CtorMsg struct {
			Function string 
			Args []string 
		}
 
var contractState = ContractState{MYVERSION}


// ************************************
// deploy callback mode 
// ************************************
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    var stateArg ContractState
    var err error
    if len(args) != 1 {
        return nil, errors.New("init expects one argument, a JSON string with tagged version string")
    }
    err = json.Unmarshal([]byte(args[0]), &stateArg)
    if err != nil {
        return nil, errors.New("Version argument unmarshal failed: " + fmt.Sprint(err))
    }
    if stateArg.Version != MYVERSION {
        return nil, errors.New("Contract version " + MYVERSION + " must match version argument: " + stateArg.Version)
    }
    contractStateJSON, err := json.Marshal(stateArg)
    fmt.Println("contractStateJSON--------%+s",string(contractStateJSON))
    if err != nil {
        return nil, errors.New("Marshal failed for contract state" + fmt.Sprint(err))
    }
    err = stub.PutState(CONTRACTSTATEKEY, contractStateJSON)
    if err != nil {
        return nil, errors.New("Contract state failed PUT to ledger: " + fmt.Sprint(err))
    }
    return nil, nil
}

// ************************************
// deploy and invoke callback mode 
// ************************************
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    // Handle different functions
    if function == "createAsset" {
        // create assetID
        return t.createAsset(stub, args)
    } else if function == "updateAsset" {
        // create assetID
        return t.updateAsset(stub, args)
    } else if function == "deleteAsset" {
        // Deletes an asset by ID from the ledger
        return t.deleteAsset(stub, args)
    }
    return nil, errors.New("Received unknown invocation: " + function)
}

// ************************************
// query callback mode 
// ************************************
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    // Handle different functions
    if function == "readAsset" {
        // gets the state for an assetID as a JSON struct
        return t.readAsset(stub, args)
    } else if function =="readAssetObjectModel" {
        return t.readAssetObjectModel(stub, args)
    }  else if function == "readAssetSamples" {
		// returns selected sample objects 
		return t.readAssetSamples(stub, args)
	} else if function == "readAssetSchemas" {
		// returns selected sample objects 
		return t.readAssetSchemas(stub, args)
	}
    return nil, errors.New("Received unknown invocation: " + function)
}

/**********main implementation *************/

func main() {
    err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting Simple Chaincode: %s", err)
    }
}

/*****************ASSET CRUD INTERFACE starts here************/

/****************** 'deploy' methods *****************/

/******************** createAsset ********************/

func (t *SimpleChaincode) createAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    _,erval:=t. createOrUpdateAsset(stub, args)
    return nil, erval
}

//******************** updateAsset ********************/

func (t *SimpleChaincode) updateAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
     _,erval:=t. createOrUpdateAsset(stub, args)
      fmt.Println("UpdateAsset method called ------->>>>>>", erval);
    return nil, erval
}


//******************** deleteAsset ********************/

func (t *SimpleChaincode) deleteAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var assetID string // asset ID
    var err error
    var stateIn AssetState
	
    // validate input data for number of args, Unmarshaling to asset state and obtain asset id
    stateIn, err = t.validateInput(args)
    if err != nil {
        return nil, err
    }
    assetID = *stateIn.AssetID
    // Delete the key / asset from the ledger
    err = stub.DelState(assetID)
    if err != nil {
        err = errors.New("DELSTATE failed! : "+ fmt.Sprint(err))
       return nil, err
    }
    return nil, nil
}

/******************* Query Methods ***************/

//********************readAsset********************/

func (t *SimpleChaincode) readAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var assetID string // asset ID
    var err error
    var state AssetState

     // validate input data for number of args, Unmarshaling to asset state and obtain asset id
    stateIn, err:= t.validateInput(args)
    if err != nil {
        return nil, errors.New("Asset does not exist!")
    }
    assetID = *stateIn.AssetID
        // Get the state from the ledger
    assetBytes, err:= stub.GetState(assetID)
    if err != nil  || len(assetBytes) ==0{
        err = errors.New("Unable to get asset state from ledger")
        return nil, err
    } 
    err = json.Unmarshal(assetBytes, &state)
    if err != nil {
         err = errors.New("Unable to unmarshal state data obtained from ledger")
        return nil, err
    }
    return assetBytes, nil
}

//*************readAssetObjectModel*****************/

func (t *SimpleChaincode) readAssetObjectModel(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var state AssetState = AssetState{}

    // Marshal and return
    stateJSON, err := json.Marshal(state)
    if err != nil {
        return nil, err
    }
    return stateJSON, nil
}
//*************readAssetSamples*******************/

func (t *SimpleChaincode) readAssetSamples(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return []byte(samples), nil
}
//*************readAssetSchemas*******************/

func (t *SimpleChaincode) readAssetSchemas(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return []byte(schemas), nil
}

// ************************************
// validate input data : common method called by the CRUD functions
// ************************************
func (t *SimpleChaincode) validateInput(args []string) (stateIn AssetState, err error) {
    var assetID string // asset ID
    var state AssetState = AssetState{} // The calling function is expecting an object of type AssetState

    if len(args) !=1 {
        err = errors.New("Incorrect number of arguments. Expecting a JSON strings with mandatory assetID")
        return state, err
    }
    jsonData:=args[0]
    assetID = ""
    stateJSON := []byte(jsonData)
    err = json.Unmarshal(stateJSON, &stateIn)
    if err != nil {
        err = errors.New("Unable to unmarshal input JSON data")
        return state, err
        // state is an empty instance of asset state
    }      
    // was assetID present?
    // The nil check is required because the asset id is a pointer. 
    // If no value comes in from the json input string, the values are set to nil
    
    if stateIn.AssetID !=nil { 
        assetID = strings.TrimSpace(*stateIn.AssetID)
        if assetID==""{
            err = errors.New("AssetID not passed")
            return state, err
        }
    } else {
        err = errors.New("Asset id is mandatory in the input JSON data")
        return state, err
    }
    
    
    stateIn.AssetID = &assetID
    return stateIn, nil
}
//******************** createOrUpdateAsset ********************/

func (t *SimpleChaincode) createOrUpdateAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var assetID string                 // asset ID                    // used when looking in map
    var err error
    var stateIn AssetState
    var stateStub AssetState
   

    // validate input data for number of args, Unmarshaling to asset state and obtain asset id

    stateIn, err = t.validateInput(args)
    if err != nil {
        fmt.Println("error validateInput inside createOrUpdateAsset------",err);
        return nil, err
    }
    // fmt.Println("stateIn inside createOrUpdateAsset------",string(stateIn));
    assetID = *stateIn.AssetID
    fmt.Println("assetId inside createOrUpdateAsset------",assetID);
    
    // getting the chaincodeid from cloudant
   jsonStringchaincode :=getchaincodeid()
    //fmt.Println("----getchaincodeid Response------------------",jsonStringchaincode)
                
   
     
    // Partial updates introduced here
    // Check if asset record existed in stub
    assetBytes, err:= stub.GetState(assetID)
    if err != nil || len(assetBytes)==0{
        // This implies that this is a 'create' scenario
         stateStub = stateIn // The record that goes into the stub is the one that cme in
           fmt.Println("stateStub inside createOrUpdateAsset--------",stateStub);
    } else {
        // This is an update scenario
        fmt.Println("assetBytes inside createOrUpdateAsset---updaet scenario-----",string(assetBytes));
        err = json.Unmarshal(assetBytes, &stateStub)
        if err != nil {
            err = errors.New("Unable to unmarshal JSON data from stub")
            return nil, err
            // state is an empty instance of asset state
        }
        var name string="Anzal Beg"
        stateStub.Ownername=&name
          // Merge partial state updates
        stateStub, err =t.mergePartialState(stateStub,stateIn)
        if err != nil {
            err = errors.New("Unable to merge state")
            return nil,err
        }
                             
    // post called 
    jsonString := getcurrentKitOwner(jsonStringchaincode)
    fmt.Println("----------------------",jsonString)
    var pro Response 
    var kitownername string
   // var msg Message123
     var data map[string]interface{}
    err := json.Unmarshal([]byte(jsonString), &pro)
    if err == nil {
    fmt.Printf("%+v\n", pro.Result.Status)
    message_unquoted:= strings.Replace(pro.Result.Message,"\"{", "`{", 2)
    fmt.Printf("message_unquoted-------%+s\n", string(message_unquoted))
                err1 := json.Unmarshal([]byte(message_unquoted), &data)
                if err1 == nil{
                        val, ok := data["ownername"]
                        if ok {
                            //do something here
                             kitownername:=val
                              fmt.Printf("kitownername----%+s\n", kitownername)
                               
                        }else{
                            kitownername:="null"
                              fmt.Printf("kitownername----%+s\n", kitownername)
                        }
                       
                } else{
                        fmt.Println(err1)
        }
    } else {
    fmt.Println(err)
    }
     fmt.Printf("kitownername--global--%+s\n", kitownername)
    }
    var name string="Anzal Beg"
        stateStub.Ownername=&name
    stateJSON, err := json.Marshal(stateStub)
     fmt.Println("stateJSON inside createOrUpdateAsset---updaet scenario-Marshal----",string(stateJSON));
    if err != nil {
        return nil, errors.New("Marshal failed for contract state" + fmt.Sprint(err))  }
    // Get existing state from the stub

  

    // Write the new state to the ledger
    err = stub.PutState(assetID, stateJSON)
     fmt.Println("putstate stateJSON in createOrUpdateAsset-------",string(stateJSON));
    if err != nil {
        err = errors.New("PUT ledger state failed: "+ fmt.Sprint(err))            
        return nil, err
    } 
    return nil, nil
}
/*********************************  internal: mergePartialState ****************************/	
 func (t *SimpleChaincode) mergePartialState(oldState AssetState, newState AssetState) (AssetState,  error) {
     
    old := reflect.ValueOf(&oldState).Elem()
    new := reflect.ValueOf(&newState).Elem()
    for i := 0; i < old.NumField(); i++ {
        oldOne:=old.Field(i)
        newOne:=new.Field(i)
        if ! reflect.ValueOf(newOne.Interface()).IsNil() {
            oldOne.Set(reflect.Value(newOne))
        } 
    }
    return oldState, nil
 }

    func getchaincodeid() string {
	url := "https://f2b6a79d-29d6-4df7-9db3-79f61eaef462-bluemix:efc09df85252cc5c7c2fa29f478895930b7f043ccdbe95ec8dc57ee2f41bd77a@f2b6a79d-29d6-4df7-9db3-79f61eaef462-bluemix.cloudant.com/chaincodeid/_find"
	fmt.Println("URL getchaincodeid---:>", url)

	var jsonStr = []byte(`{
                            "selector": {
                                "_id": {
                                "$gt": 0
                                }
                            },
                            "fields": [
                                "chaincodeid"
                                
                            ],
                            "sort": [
                                {
                                "_id": "asc"
                                }
                            ]}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(" getchaincodeid----response Status:", resp.Status)
	fmt.Println(" getchaincodeid ----response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("getchaincodeid-------response Body:", string(body))
     var rootMap map[string]interface{}
    // unmarshal data into map
    json.Unmarshal([]byte(string(body)), &rootMap)
    // access and cast until desired data is retrieved
    docs := rootMap["docs"].([]interface{})
    firstRoute := docs[0].(map[string]interface{})
    chaincodeid := firstRoute["chaincodeid"].(string)
    fmt.Println("URL chaincodeid---:>", chaincodeid)
	return chaincodeid
}                  
func getcurrentKitOwner(chaincodeidvalue string) string {
    fmt.Println("chaincodeid-------inside getcurrentKitOwner---:>", chaincodeidvalue)
	url := "https://c40925427abd45b08c5c38b04e20a438-vp1.us.blockchain.ibm.com:5002/chaincode"
	fmt.Println("URL getcurrentKitOwner---:>", url)
    readcurrentstatus:=AutoGenerated{
                                                            Jsonrpc: "2.0",
                                                            Method: "query",
                                                            Params: Params{
                                                                Type: 1,
                                                                ChaincodeID: ChaincodeID{
                                                                Name: chaincodeidvalue,
                                                                },
                                                                CtorMsg: CtorMsg{
                                                                Function: "readAsset",
                                                                Args: []string{"{\"assetID\":\"KIT10009\"}"},
                                                                },
                                                                SecureContext: "user_type1_1",
                                                            },
                                                            Id: 1,
                                      }
                                        
    js,_:=json.Marshal(readcurrentstatus)
	var jsonStr = []byte(js)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(" getcurrentKitOwner----response Status:", resp.Status)
	fmt.Println(" getcurrentKitOwner ----response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("getcurrentKitOwner-------response Body:", body)
	return string(body)
}

