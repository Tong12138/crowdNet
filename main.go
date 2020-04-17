//main.go

package main

import(
	"os"
	"fmt"
	"github.com/hyperledger/Demo/sdkInit"
	"github.com/hyperledger/Demo/service"
	"github.com/hyperledger/Demo/web"
	"github.com/hyperledger/Demo/web/controller"

)

const(
	configFile = "config.yaml"
	initialized = false
	SimpleCC = "simplecc"

)



func main(){
	initInfo := &sdkInit.InitInfo{
		ChannelID: "mychannel",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/hyperledger/Demo/fixtures/channel-artifacts/channel.tx",

		OrgAdmin: "Admin",
		OrgName:"OrgReq",
		OrdererOrgName: "orderer.crowd.com",

		ChaincodeID: SimpleCC,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath: "github.com/hyperledger/Demo/chaincode/",
		UserName:"User1",

	}
	var environment map[string] *sdkInit.Environ
	environment = make(map[string]*sdkInit.Environ)
	environment[initInfo.ChannelID]= &sdkInit.Environ{}

	
    var err error
	//实例化SDK
	environment[initInfo.ChannelID].Sdk, err = sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
    
    for key:= range environment{
    	defer environment[key].Sdk.Close()
    }
	

    //资源客户管理端

    err = sdkInit.CreateSourceClient(environment[initInfo.ChannelID], initInfo)

	err = sdkInit.CreateChannel(environment[initInfo.ChannelID], initInfo)
	if err != nil{
		fmt.Println(err.Error())
		return 
	}

	err = sdkInit.JoinChannel(environment[initInfo.ChannelID], initInfo)

	environment[initInfo.ChannelID].ChannelClient, err= sdkInit.InstallAndInstantiateCC(environment[initInfo.ChannelID], initInfo)
	if err != nil{
		fmt.Println(err.Error())
		return
	}


	 // =================SETINFO
	serviceSetup := service.ServiceSetup{

		Environment: &environment,
		ChaincodeID: SimpleCC,
		Info:initInfo,
	}

	msg, err := serviceSetup.SetInfo("yyt", "blockchain",initInfo.ChannelID)
	if err!= nil{
		fmt.Println(err)
	}else{
		fmt.Println(msg)
	}
	

	// / GETINFO
	msg, err = serviceSetup.GetInfo("yyt", initInfo.ChannelID)
	if err!= nil{
		fmt.Println(err)
	}else{
		fmt.Println(msg)
	}
		
	app := controller.Application{
		Fabric: &serviceSetup,
	}
	web.WebStart(&app)


}