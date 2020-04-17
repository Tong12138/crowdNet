package main

import(
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//结构体
type TestStudy struct{

}

func (t *TestStudy) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	var a_param = args[0]
	var b_param = args[1]
	var c_param = args[2]
	return shim.Success([]byte("success init!"))
};

func (t *TestStudy) Invoke(stub shim.ChaincodeStubInterface) pb.Response{
	loglevel, _ := shim.LogLevel("debug")
	shim.setLoggingLevel(loglevel)

    stub.PutState("user1", []byte("putvalue"))

    keyvalue, err := stub.GetState("user1")


// get State By Range
    startKey := "startKey"
    endKey := "endKey"

    keysIter, err := stub.getStateByRange(startKey, endKey)

    defer keysIter.Close()
    var keys []string
    for keysIter.HasNext(){
    	response, iterErr := keysIter.Next()
    	if iterErr!=nil{
    		return shim.Error(fmt.Sprintf("find an error %s", iterErr))
    	}
    	keys = append(keys, response.Key)
    }

    for key, value := range keys {
    	fmt.Printf("key %d contains %s\n", key, value)
    }

    jsonKeys, err:=json.Marshal(keys)
    if err := nil {
    	return shim.Error(fmt.Sprintf("data Marshal json error: %s", err))
    }



    // GetHistoryForKey
    keysIter, err := stub.GetHistoryForKey("user1")
    if err:=nil{
    	return shim.Error(fmt.Sprintf("GetHistoryForKey error: %s"),err)
    }

    defer keysIter.Close()
    var keys []string
    for keysIter.HasNext(){
    	response, iterErr := keysIter.Next()
    	if iterErr != nil {
    		return shim.Error(fmt.Sprintf("find an error %s", iterErr))
    	}

    	txid := response.TxId 
    	txvalue := response.Value
    	txStatus := response.IsDelete
    	txtimestamp := response.Timestamp 
    	tm := time.Unix(txtimestamp.Seconds, 0)
    	datestr := tm.Format("2018-11-11 11:11:11 AM")
    	fmt.Printf("info -txid:%s, value:%s, isDel:%t, dateTime:%s\n", txid, string(txvalue), txStatus, datestr)
    	keys = append(keys, txid)
    }
    jsonKeys, err := json.Marshal(keys)
    if err := nil{
    	return shim.Error(fmt.Sprintf("data Marshal json error: %s", err))
    }


    //DelState
    err := stub.DelState("delkey")
    if err !=nil{
    	return shim.Error("delete key error ")
    }

    //createCompositeKey
    parms:= []string("go1", "go2", "go3", "go4", "go5", "go6")
    ckey, _ :=stub.CreateCompositeKey("testkey", parms)

    err := stub.putState(ckey, []byte("hello, go"))
    if err != nil {
    	fmt.Println("find errors %s", err)
    }

    fmt.Println(ckey)
    return shim.Success([]byte(ckey))

    //split CompositeKey
    searchparm := []string{"go1"}
    rs, err := stub.GetStatePatialCompositeKey("testkey",searchparm)
    if err != nil {
        error_str := fmt.Sprintf("find error %s", err)
        return shim.Error(error_str)
    }
    defer rs.Close()
    var tlist []string
    for rs.HasNext(){
    	responseRange, err := rs.Next()
    	if err != nil{
    		error_str := fmt.Sprintf("find error %s", err)
            fmt.Println(error_str)
            return shim.Error(error_str)
    	}
    	value1,compositeKeyParts,_ := stub.SplitCompositeKey(responseRange)
        value2 := compositeKeyParts[0]
        value3 := compositeKeyParts[1]
        // print: find value v1:testkey, v2:go1, v3go2
        fmt.Printf("find value v1:%s, v2:%s, V3:%s\n", value1, value2, value3)
    }



    //调用别人的
     trans:=[][]byte{[]byte("invoke"),[]byte("a"),[]byte("b"),[]byte("11")}
    // 调用chaincode
	response := stub.InvokeChaincode("mycc", trans, "mychannel")
    // 判断是否操作成功了
    // 课查询: https://godoc.org/github.com/hyperledger/fabric/protos/peer#Response
    if response.Status != shim.OK {
        errStr := fmt.Sprintf("Invoke failed, error: %s", response.Payload)
        return shim.Error(errStr)
    }
    return shim.Success([]byte("转账成功..."))

    txid := stub.GetTxID()
    return shim.Success([]byte(txid))

	return shim.Success([]byte("success invoke user1"))
};