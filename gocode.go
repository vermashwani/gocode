package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	//"github.com/op/go-logging"
)

//var myLogger = logging.MustGetLogger("poc_chaincode")

// SKH is a high level smart contract that
type SKH struct {

}

// ImbalanceDetails is for storing Imbalance Details
type ImbalanceDetails struct{
	Esco string `json:"esco"`
	UserId string `json:"userId"`
	TotalImbalance string `json:"totalImbalance"`
	LastUpdateDate string `json:"lastUpdateDate"`
}

// Transaction is for storing Transaction Details
type Transaction struct{	
	TransId string `json:"transId"`
	TransDate string `json:"transDate"`
	From string `json:"from"`
	To string `json:"to"`
	Quantity string `json:"quantity"`
	Type string `json:"type"`
	Status string `json:"status"`
	LastUpdateDate string `json:"lastUpdateDate"`
}

// to return the verify result
type VerifyU struct{	
	Result string `json:"result"`
}

// Init initializes the smart contracts
func (t *SKH) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Println("Init..")
	fmt.Errorf("Ashwani Init..")
	
	// Check if table already exists
	_, err := stub.GetTable("ImbalanceDetails")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}
	fmt.Println("Creating Table --> ImbalanceDetails")
	
	// Create ImbalanceDetails Table
	err = stub.CreateTable("ImbalanceDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "esco", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "userId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "totalImbalance", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "lastUpdateDate", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating ImbalanceDetails table.")
	}

	// Check if Transaction table already exists
	_, err = stub.GetTable("Transaction")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}
	fmt.Println("Creating Table --> Transaction")

	// Create Transaction Table
	err = stub.CreateTable("Transaction", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "transId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "transDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "from", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "to", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "quantity", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "type", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "status", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "lastUpdateDate", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating Transaction table.")
	}
	
	// setting up the users esco
	stub.PutState("user_type1_1", []byte("ESCO_A"))
	stub.PutState("user_type1_2", []byte("ESCO_B"))
	stub.PutState("user_type1_3", []byte("ESCO_C"))
	stub.PutState("user_type1_4", []byte("ESCO_D"))	
	
	fmt.Println("Init Done.")
	return nil, nil
}
	
//addImbalance to ESCO
func (t *SKH) addImbalance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	fmt.Println("function --> addImbalance()")
	
		if len(args) != 4 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 4 . Got: %d.", len(args))
		}
	
	esco:=args[0]
	userId:=args[1]
	totalImbalance:=args[2]
	lastUpdateDate:=args[3]
		
	/*assignerOrg1, err := stub.GetState(args[11])
	assignerOrg := string(assignerOrg1)
	createdBy:=assignerOrg*/

	// Get the row pertaining to this Esco
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: esco}}
	columns = append(columns, col1)

	row, err := stub.GetRow("ImbalanceDetails", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get the data for ESCO " + esco + "\"}"
		return nil, errors.New(jsonResp)
	}
	if len(row.Columns) == 0 {
		// Insert a row
		ok, err := stub.InsertRow("ImbalanceDetails", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: esco}},
				&shim.Column{Value: &shim.Column_String_{String_: userId}},
				&shim.Column{Value: &shim.Column_String_{String_: totalImbalance}},
				&shim.Column{Value: &shim.Column_String_{String_: lastUpdateDate}},
			}})

		if err != nil {
			return nil, fmt.Errorf("Failed inserting row [%s]", err)
		}
		if !ok {
			return nil, errors.New("Failed inserting row.")
		}
	}
	fmt.Println("function --> addImbalance() Exit.")
	return nil, nil
}

//acceptTransaction to ESCO
func (t *SKH) acceptTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("function --> acceptTransaction()")

	if len(args) != 2 {
		return nil, fmt.Errorf("Incorrect number of arguments. Expecting 2 . Got: %d.", len(args))
	}
	
	transId := args[0]
	lastUpdateDate := args[1]
	
		//Get the row for Transaction
		var columns []shim.Column
		col := shim.Column{Value: &shim.Column_String_{String_: transId}}
		columns = append(columns, col)

		row, err := stub.GetRow("Transaction", columns)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get the data for Transaction Id " + transId + "\"}"
			return nil, errors.New(jsonResp)
		}	
		
		escoFrom := row.Columns[2].GetString_()
		escoTo := row.Columns[3].GetString_()
		transQuantity, _:=strconv.ParseInt(row.Columns[4].GetString_(), 10, 0) 
		transType := row.Columns[5].GetString_()
		transStatus := row.Columns[6].GetString_()
	
		fmt.Println("function --> addTransaction() :: TransId [%s], escoFrom [%s], escoTo [%s], Quantity [%s], TransType [%s], Status [%s]", transId, escoFrom, escoTo, transQuantity, transType, transStatus)
	
    if(transStatus == "Pending") {
	
		fmt.Println("function --> acceptTransaction() :: transStatus condition TRUE.")
		
		// Get the row for escoFrom
		var columns1 []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: escoFrom}}
		columns1 = append(columns1, col1)

		row1, err1 := stub.GetRow("ImbalanceDetails", columns1)
		if err1 != nil {
			jsonResp := "{\"Error\":\"Failed to get the data for ESCO " + escoFrom + "\"}"
			return nil, errors.New(jsonResp)
		}

		// Get the row for escoTo
		var columns2 []shim.Column
		col2 := shim.Column{Value: &shim.Column_String_{String_: escoTo}}
		columns2 = append(columns2, col2)

		row2, err2 := stub.GetRow("ImbalanceDetails", columns2)
		if err2 != nil {
			jsonResp := "{\"Error\":\"Failed to get the data for ESCO " + escoTo + "\"}"
			return nil, errors.New(jsonResp)
		}
		
		fmt.Println("function --> acceptTransaction() :: Transaction Count [%d] FromCount[%d] ToCount[%d]", len(row.Columns), len(row1.Columns), len(row2.Columns))

		//Checking data availability for the Transaction
		if len(row.Columns) > 0 && len(row1.Columns) > 0 && len(row2.Columns) > 0{
		
		var ok3 bool = false
		var ok4 bool = false
		
		if (transType == "BUY") {
			//Update Quantity Transfer from
			fmt.Println("function --> acceptTransaction() :: Condition --> BUY")
			
			totalQuantity, _:=strconv.ParseInt(row1.Columns[2].GetString_(), 10, 0)
			updateQuantity :=  strconv.Itoa(int(totalQuantity) + int(transQuantity))
			
			fmt.Println("function --> acceptTransaction() :: Update EscoFrom Before [%d] After [%d] update",totalQuantity, updateQuantity)
			
			ok3, err3 := stub.ReplaceRow("ImbalanceDetails", shim.Row{
			Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: escoFrom}},
			&shim.Column{Value: &shim.Column_String_{String_: row1.Columns[1].GetString_()}},	
			&shim.Column{Value: &shim.Column_String_{String_: updateQuantity}},
			&shim.Column{Value: &shim.Column_String_{String_: lastUpdateDate}},
			}})
			
			if err3 != nil {
				return nil, fmt.Errorf("Failed replacing row [%s]", err3)
			}
			if !ok3 {
				return nil, errors.New("Failed replacing row.")
			}
			
			//Update Quantity Transfer to
			totalQuantity1, _:=strconv.ParseInt(row2.Columns[2].GetString_(), 10, 0)
			updateQuantity1 :=strconv.Itoa(int(totalQuantity1) - int(transQuantity))
			
			fmt.Println("function --> acceptTransaction() :: Update EscoTo Before [%d] After [%d] update",totalQuantity1, updateQuantity1)
			
			ok4, err4 := stub.ReplaceRow("ImbalanceDetails", shim.Row{
			Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: escoTo}},
			&shim.Column{Value: &shim.Column_String_{String_: row2.Columns[1].GetString_()}},			
			&shim.Column{Value: &shim.Column_String_{String_: updateQuantity1}},
			&shim.Column{Value: &shim.Column_String_{String_: lastUpdateDate}},
			}})
			
			if err4 != nil {
				return nil, fmt.Errorf("Failed replacing row [%s]", err4)
			}
			if !ok4 {
				return nil, errors.New("Failed replacing row.")
			}
		} else if(transType == "SELL") {
			//Update Quantity Transfer from
			fmt.Println("function --> acceptTransaction() :: Condition --> SELL")
			
			totalQuantity, _:=strconv.ParseInt(row1.Columns[2].GetString_(), 10, 0)
			updateQuantity :=  strconv.Itoa(int(totalQuantity) - int(transQuantity))
			
			fmt.Println("function --> acceptTransaction() :: Update EscoFrom Before [%d] After [%d] update",totalQuantity, updateQuantity)
			
			ok3, err3 := stub.ReplaceRow("ImbalanceDetails", shim.Row{
			Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: escoFrom}},
			&shim.Column{Value: &shim.Column_String_{String_: row1.Columns[1].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: updateQuantity}},
			&shim.Column{Value: &shim.Column_String_{String_: lastUpdateDate}},
			}})
			
			if err3 != nil {
				return nil, fmt.Errorf("Failed replacing row [%s]", err3)
			}
			if !ok3 {
				return nil, errors.New("Failed replacing row.")
			}
			
			//Update Quantity Transfer to
			totalQuantity1, _:=strconv.ParseInt(row2.Columns[2].GetString_(), 10, 0)
			updateQuantity1 :=  strconv.Itoa(int(totalQuantity1) + int(transQuantity))
			
			fmt.Println("function --> acceptTransaction() :: Update EscoTo Before [%d] After [%d] update",totalQuantity1, updateQuantity1)
			
			ok4, err4 := stub.ReplaceRow("ImbalanceDetails", shim.Row{
			Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: escoTo}},
			&shim.Column{Value: &shim.Column_String_{String_: row2.Columns[1].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: updateQuantity1}},
			&shim.Column{Value: &shim.Column_String_{String_: lastUpdateDate}},
			}})
			
			if err4 != nil {
						return nil, fmt.Errorf("Failed replacing row [%s]", err4)
					}
			if !ok4 {
						return nil, errors.New("Failed replacing row.")
				    }
			}//End Transaction Type Condition
			
			//Update Transaction Table on successful From and To ESCO Update	 
			fmt.Println("function --> acceptTransaction() :: Update Status ESCO From [%t] To [%t]", ok3, ok4)
			
			if(ok3 && ok4) {
				fmt.Println("function --> acceptTransaction() :: Updateing Transaction table status.")
					//Update Transaction Status on successful Imbalance Details updation
					ok5, err5 := stub.ReplaceRow("Transaction", shim.Row{
					Columns: []*shim.Column{
					&shim.Column{Value: &shim.Column_String_{String_: transId}},
					&shim.Column{Value: &shim.Column_String_{String_: row.Columns[1].GetString_()}},
					&shim.Column{Value: &shim.Column_String_{String_: row.Columns[2].GetString_()}},
					&shim.Column{Value: &shim.Column_String_{String_: row.Columns[3].GetString_()}},
					&shim.Column{Value: &shim.Column_String_{String_: row.Columns[4].GetString_()}},
					&shim.Column{Value: &shim.Column_String_{String_: row.Columns[5].GetString_()}},
					&shim.Column{Value: &shim.Column_String_{String_: "Accepted"}},
					&shim.Column{Value: &shim.Column_String_{String_: lastUpdateDate}},
					}})
			
					if err5 != nil {
							return nil, fmt.Errorf("Failed replacing row [%s]", err5)
						}
					if !ok5 {
							return nil, errors.New("Failed replacing row.")
						}
				}else{
					return nil, errors.New("Transaction Rollback code will be followed..<TBD>")
				 }
	    }else{ 
			return nil, fmt.Errorf("Column lengths -->> . Got: %d. %d.  %d.", len(row.Columns), len(row1.Columns), len(row2.Columns))
		   }
	}else{
			return nil, errors.New("Incorrect Status Type. Should be Pending")
	     }
	return nil, nil
}

//addTransaction - Add Imbalance Transaction
func (t *SKH) addTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("function --> addTransaction()")

		if len(args) != 6 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 6. Got: %d.", len(args))
		}
		transId:=args[0]
		transDate:=args[1]
		from:=args[2]
		to:=args[3]
		quantity:=args[4]
		transType:=args[5]

		fmt.Println("function --> addTransaction() :: TransId [%s], TransDate [%s], From [%s], To [%s], Quantity [%s], TransType [%s]", transId, transDate, from, to, quantity, transType)
		
		// Insert a row
		ok, err := stub.InsertRow("Transaction", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: transId}},
				&shim.Column{Value: &shim.Column_String_{String_: transDate}},
				&shim.Column{Value: &shim.Column_String_{String_: from}},
				&shim.Column{Value: &shim.Column_String_{String_: to}},
				&shim.Column{Value: &shim.Column_String_{String_: quantity}},
				&shim.Column{Value: &shim.Column_String_{String_: transType}},
				&shim.Column{Value: &shim.Column_String_{String_: "Pending"}},
				&shim.Column{Value: &shim.Column_String_{String_: transDate}},
			}})

		if err != nil {
			return nil, err 
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}
		fmt.Println("function --> addTransaction() Exit.")
	return nil, nil
}

//get All ESCO Imbalances 
func (t *SKH) getAllImbalances(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("function --> getAllImbalances()")
	
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect number of arguments. Expecting 1. Got: %d.", len(args))
	}

	fmt.Println("function --> getAllImbalances() :: Input [%s]", args[0])

	var columns []shim.Column

	rows, err := stub.GetRows("ImbalanceDetails", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
	
	res2E:= []*ImbalanceDetails{}	
	
	for row := range rows {		
		newApp:= new(ImbalanceDetails)
		newApp.Esco = row.Columns[0].GetString_()
		newApp.UserId = row.Columns[1].GetString_()
		newApp.TotalImbalance = row.Columns[2].GetString_()
		newApp.LastUpdateDate = row.Columns[3].GetString_()
		
		//if newApp.EmployeeId == EmployeeId && newApp.Source == assignerOrg{
		res2E=append(res2E,newApp)		
		//}				
	}
	
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	fmt.Println("function --> getAllImbalances() Exit.")
	
	return mapB, nil
}

// to get the imbalance deatils of an ESCO
func (t *SKH) getImbalance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("function --> getImbalance()")
	
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect number of arguments. Expecting 1. Got: %d.", len(args))
	}

	Esco := args[0]
	
	fmt.Println("function --> getImbalance() :: ESCO [%s]", Esco)
	
	// Get the row pertaining to this Esco
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: Esco}}
	columns = append(columns, col1)

	row, err := stub.GetRow("ImbalanceDetails", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get the data for ESCO " + Esco + "\"}"
		return nil, errors.New(jsonResp)
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"No Data available for ESCO " + Esco + "\"}"
		return nil, errors.New(jsonResp)
	}
	res2E := ImbalanceDetails{}
	
	res2E.Esco = row.Columns[0].GetString_()
	res2E.UserId = row.Columns[1].GetString_()
	res2E.TotalImbalance = row.Columns[2].GetString_()
	res2E.LastUpdateDate = row.Columns[3].GetString_()

    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	fmt.Println("function --> getImbalance() Exit.")
	
	return mapB, nil
}

// to get the Transaction deatils
func (t *SKH) getTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("function --> getTransaction()")
	
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect number of arguments. Expecting 1. Got: %d.", len(args))
	}

	TransId := args[0]
	
	fmt.Println("function --> getTransaction() :: TransId [%s]", TransId)
	
	// Get the row pertaining to this Transaction
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: TransId}}
	columns = append(columns, col1)

	row, err := stub.GetRow("Transaction", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get the data for Transaction " + TransId + "\"}"
		return nil, errors.New(jsonResp)
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"No Data available for Transaction Id " + TransId + "\"}"
		return nil, errors.New(jsonResp)
	}
	res2E := Transaction{}
	
	res2E.TransId = row.Columns[0].GetString_()
	res2E.TransDate = row.Columns[1].GetString_()
	res2E.From = row.Columns[2].GetString_()
	res2E.To = row.Columns[3].GetString_()
	res2E.Quantity = row.Columns[4].GetString_()
	res2E.Type = row.Columns[5].GetString_()
	res2E.Status = row.Columns[6].GetString_()
	res2E.LastUpdateDate = row.Columns[7].GetString_()
		
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	fmt.Println("function --> getTransaction() Exit.")
	
	return mapB, nil
}

//get All Transaction that are sent for Approval
func (t *SKH) getTransactionSent(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("function --> getTransactionSent()")
	
	if len(args) != 2 {
		return nil, fmt.Errorf("Incorrect number of arguments. Expecting 2. Got: %d.", len(args))
	}

	Esco := args[0]
	Status := args[1]
	
	fmt.Println("function --> getTransactionSent() :: ESCO [%s], Status [%s]", Esco, Status)

	var columns []shim.Column

	rows, err := stub.GetRows("Transaction", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
		
	res2E:= []*Transaction{}	
	
	for row := range rows {		
		newApp:= new(Transaction)
		newApp.TransId = row.Columns[0].GetString_()
		newApp.TransDate = row.Columns[1].GetString_()
		newApp.From = row.Columns[2].GetString_()
		newApp.To = row.Columns[3].GetString_()
		newApp.Quantity = row.Columns[4].GetString_()
		newApp.Type = row.Columns[5].GetString_()
		newApp.Status = row.Columns[6].GetString_()
		newApp.LastUpdateDate = row.Columns[7].GetString_()
		
		if((newApp.From == Esco && newApp.To != Esco) && (newApp.Status == Status)){
		res2E=append(res2E,newApp)		
		}				
	}
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	fmt.Println("function --> getTransactionSent() Exit.")
	
	return mapB, nil
}

//get All Transaction that are received for Approval
func (t *SKH) getTransactionReceived(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("function --> getTransactionReceived()")
	
	if len(args) != 2 {
		return nil, fmt.Errorf("Incorrect number of arguments. Expecting 2. Got: %d.", len(args))
	}

	Esco := args[0]
	Status := args[1]

	fmt.Println("function --> getTransactionReceived() :: ESCO [%s], Status [%s]", Esco, Status)

	var columns []shim.Column

	rows, err := stub.GetRows("Transaction", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
		
	res2E:= []*Transaction{}	
	
	for row := range rows {		
		newApp:= new(Transaction)
		newApp.TransId = row.Columns[0].GetString_()
		newApp.TransDate = row.Columns[1].GetString_()
		newApp.From = row.Columns[2].GetString_()
		newApp.To = row.Columns[3].GetString_()
		newApp.Quantity = row.Columns[4].GetString_()
		newApp.Type = row.Columns[5].GetString_()
		newApp.Status = row.Columns[6].GetString_()
		newApp.LastUpdateDate = row.Columns[7].GetString_()
		
		if((newApp.To == Esco && newApp.From != Esco) && (newApp.Status == Status)){
		res2E=append(res2E,newApp)		
		}				
	}
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	fmt.Println("function --> getTransactionReceived() Exit.")
	
	return mapB, nil
}

//get All Accepted Transactions
func (t *SKH) getTransactionAccepted(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("function --> getTransactionAccepted()")

	if len(args) != 2 {
		return nil, fmt.Errorf("Incorrect number of arguments. Expecting 2. Got: %d.", len(args))
	}
	Esco := args[0]
	Status := args[1]
	
	fmt.Println("function --> getTransactionAccepted() :: ESCO [%s], Status [%s]", Esco, Status)
	
	var columns []shim.Column

	rows, err := stub.GetRows("Transaction", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
	
	res2E:= []*Transaction{}	
	
	for row := range rows {		
		newApp:= new(Transaction)
		newApp.TransId = row.Columns[0].GetString_()
		newApp.TransDate = row.Columns[1].GetString_()
		newApp.From = row.Columns[2].GetString_()
		newApp.To = row.Columns[3].GetString_()
		newApp.Quantity = row.Columns[4].GetString_()
		newApp.Type = row.Columns[5].GetString_()
		newApp.Status = row.Columns[6].GetString_()
		newApp.LastUpdateDate = row.Columns[7].GetString_()
		
		if((newApp.From == Esco || newApp.To == Esco) && (newApp.Status == Status)){
		res2E=append(res2E,newApp)	
		}				
	}
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	fmt.Println("function --> getTransactionAccepted() Exit.")
	return mapB, nil
}

// Invoke invokes the chaincode
func (t *SKH) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "addImbalance" {
		t := SKH{}
		return t.addImbalance(stub, args)
	}else if function == "addTransaction" { 
		t := SKH{}
		return t.addTransaction(stub, args)
	}else if function == "acceptTransaction" { 
		t := SKH{}
		return t.acceptTransaction(stub, args)
	}

	return nil, fmt.Errorf("Received unknown function invocation [%s]", function)
}

// query queries the chaincode
func (t *SKH) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "getImbalance" {
		t := SKH{}
		return t.getImbalance(stub, args)
	}else if function == "getAllImbalances" {
			 t := SKH{}
			 return t.getAllImbalances(stub, args)
	}else if function == "getTransaction" { 
			 t := SKH{}
			 return t.getTransaction(stub, args)
	}else if function == "getTransactionSent" { 
			 t := SKH{}
			 return t.getTransactionSent(stub, args)
	}else if function == "getTransactionReceived" { 
			 t := SKH{}
			 return t.getTransactionReceived(stub, args)
	}else if function == "getTransactionAccepted" { 
			 t := SKH{}
			 return t.getTransactionAccepted(stub, args)
	}
	return nil, fmt.Errorf("Received unknown function invocation [%s]", function)
}

func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(SKH))
	if err != nil {
		 fmt.Printf("Error starting SKH: %s", err)
	}
}