package service

import(
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"os"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
	"os/exec"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/crowdsourcing/Demo/sdkInit"
	"time"

)

const (
	ChaincodeVersion = "1.0"
	ConfigFile = "Demo/config.yaml"
)

func (t *ServiceSetup) RegisterChain(name, info, channelname string)(string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "userRegister", Args: [][]byte{[]byte(name), []byte(info)}}
	response, err:= (*(t.Environment))[channelname].ChannelClient.Execute(req)
	if err != nil{
		return "", err
	}
	return string(response.TransactionID), nil
}
func (t *ServiceSetup) Recharge(num, channelname string)(string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "recharge", Args: [][]byte{[]byte(num)}}
	response, err:= (*(t.Environment))[channelname].ChannelClient.Execute(req)
	if err != nil{
		return "", err
	}
	return string(response.TransactionID), nil
}

func (t *ServiceSetup) Addskills(skills, channelname string)(string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "userAddSkill", Args: [][]byte{[]byte(skills)}}
	response, err:= (*(t.Environment))[channelname].ChannelClient.Execute(req)
	if err != nil{
		return "", err
	}
	return string(response.TransactionID), nil
}

func(t *ServiceSetup) Getuser(channelname string)(string, error){
	req := channel.Request{ChaincodeID:t.ChaincodeID, Fcn:"userQuery", Args:[][]byte{}}
	response,err := (*(t.Environment))[channelname].ChannelClient.Query(req)
	if err != nil{
		return "", err
	}
	return string(response.Payload), nil
}

func(t * ServiceSetup) Getusers(channelname string)(string, error){

	req := channel.Request{ChaincodeID:t.ChaincodeID, Fcn:"alluserQuery", Args:[][]byte{}}
	response,err := (*(t.Environment))[channelname].ChannelClient.Query(req)
	if err != nil{
		return "", err
	}
	return string(response.Payload), nil
}


func (t *ServiceSetup) Posttask(Name, idd, Type, detail, reward,  requirement, channelname  string, posttime, receivetime,  deadline time.Time)(string, error){

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "taskPost", Args: [][]byte{[]byte(Name), []byte(idd), []byte(Type), []byte(detail),[]byte(reward), []byte(posttime.Format("2006-01-02 15:04:05")), []byte(receivetime.Format("2006-01-02 15:04:05")), []byte(deadline.Format("2006-01-02 15:04:05")), []byte(requirement)}}
	response, err:= (*(t.Environment))[channelname].ChannelClient.Execute(req)

	if err != nil{
		return "", err
	}

	return string(response.TransactionID), nil
}

func (t *ServiceSetup) Recievetask(taskid, channelname string, recievetime time.Time)(string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "taskReceive", Args: [][]byte{[]byte(taskid), []byte(recievetime.Format("2006-01-02 15:04:05"))}}
	response, err:= (*(t.Environment))[channelname].ChannelClient.Execute(req)

	if err != nil{
		return "", err
	}

	return string(response.TransactionID), nil
}

func (t *ServiceSetup) Alloreward(taskid, workerid, rate, channelname string, allotime time.Time)(string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "rewardAllocate", Args: [][]byte{[]byte(taskid), []byte(workerid), []byte(rate), []byte(allotime.Format("2006-01-02 15:04:05"))}}
	response, err:= (*(t.Environment))[channelname].ChannelClient.Execute(req)

	if err != nil{
		return "", err
	}

	return string(response.TransactionID), nil
}

func (t *ServiceSetup) Committask(taskid, solution, channelname string, submittime time.Time)(string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "taskCommit", Args: [][]byte{[]byte(taskid), []byte(solution),  []byte(submittime.Format("2006-01-02 15:04:05"))}}
	response, err:= (*(t.Environment))[channelname].ChannelClient.Execute(req)

	if err != nil{
		return "", err
	}

	return string(response.TransactionID), nil
}

func(t *ServiceSetup) Gettask(taskid, channelname string)(string, error){
	req := channel.Request{ChaincodeID:t.ChaincodeID, Fcn:"taskQuery", Args:[][]byte{[]byte(taskid)}}
	response,err := (*(t.Environment))[channelname].ChannelClient.Query(req)
	if err != nil{
		return "", err
	}
	return string(response.Payload), nil
}

func(t *ServiceSetup) Getrecord(taskid, channelname string)(string, error){
	req := channel.Request{ChaincodeID:t.ChaincodeID, Fcn:"recordQuery", Args:[][]byte{[]byte(taskid)}}
	response,err := (*(t.Environment))[channelname].ChannelClient.Query(req)
	if err != nil{
		return "", err
	}
	return string(response.Payload), nil
}

func(t * ServiceSetup) Gettasks(channelname string)(string, error){

	req := channel.Request{ChaincodeID:t.ChaincodeID, Fcn:"alltaskQuery", Args:[][]byte{}}
	response,err := (*(t.Environment))[channelname].ChannelClient.Query(req)
	if err != nil{
		return "", err
	}
	return string(response.Payload), nil
}


func (t *ServiceSetup) SetInfo(name, num, channelname string)(string, error){
	eventID := "eventSetInfo"
	reg, notifier := registerEvent((*(t.Environment))[channelname].ChannelClient, t.ChaincodeID, eventID)

	// reg, notifier := registerEvent((*(t.Environment))[channelname].channelClient, t.ChaincodeID, eventID)
	defer (*(t.Environment))[channelname].ChannelClient.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "set", Args: [][]byte{[]byte(name), []byte(num), []byte(eventID)}}
	response, err:= (*(t.Environment))[channelname].ChannelClient.Execute(req)

	if err != nil{
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err!= nil{
		return "", err
	}
	return string(response.TransactionID), nil
}

func(t * ServiceSetup) GetInfo(name, channelname string)(string, error){

	req := channel.Request{ChaincodeID:t.ChaincodeID, Fcn:"get", Args:[][]byte{[]byte(name)}}
	response,err := (*(t.Environment))[channelname].ChannelClient.Query(req)
	if err != nil{
		return "", err
	}
	return string(response.Payload), nil
}

func(t * ServiceSetup) CreateNewChannel(newchannel string)(error){

    //create channel.tx
    fmt.Println("newchannel is %v", newchannel)
    command := exec.Command("make","newchannel", "channelname="+newchannel)
    command.Dir = "./Demo/"

    err:= command.Run()
	if err!= nil{
		return fmt.Errorf("创建tx文件失败%v", err)
    }

    //更改配置文件
    configdata, err := ioutil.ReadFile(ConfigFile)
    if err!= nil{
    	return fmt.Errorf("读取配置文件失败%v", err)
    }

    resultMap:=make(map[string]map[string]interface{})
    err = yaml.Unmarshal(configdata, resultMap)
    if err!= nil{
    	return fmt.Errorf("反序列化文件失败%v", err)
    }
    resultMap["channels"][newchannel]= resultMap["channels"]["mychannel"]
    out, err := yaml.Marshal(resultMap)
    if err!= nil{
    	return fmt.Errorf("序列化文件失败%v", err)
    }
    err = ioutil.WriteFile(ConfigFile, out, 0655)
    if err!= nil{
    	return fmt.Errorf("写配置文件失败%v", err)
    }
    (*(t.Environment))[newchannel]= &sdkInit.Environ{}

	(*(t.Environment))[newchannel].Sdk, err = fabsdk.New(config.FromFile(ConfigFile))
       if err!=nil{
    	return fmt.Errorf("更新sdk失败: %v", err)
    }
    // defer (*(t.Environment))[newchannel].Sdk.Close()

	clientContext :=(*(t.Environment))[newchannel].Sdk.Context(fabsdk.WithUser(t.Info.OrgAdmin), fabsdk.WithOrg(t.Info.OrgName))
	if clientContext == nil{
		return fmt.Errorf("根据指定的组织名称与管理员创建资源管理客户端Context失败")
	}

	resMgmtClient, err := resmgmt.New(clientContext)
	if err!=nil{
		return fmt.Errorf("根据指定的资源管理客户端Context创建通道管理客户端失败：%v", err)
	}

	(*(t.Environment))[newchannel].OrgResMgmt = resMgmtClient

    
    //创建通道
	ChannelConfig := os.Getenv("GOPATH") + "/src/github.com/hyperledger/crowdsourcing/Demo/fixtures/channel-artifacts/" +newchannel+".tx"
	fmt.Println("SimpleService createchannel")
	mspClient, err := mspclient.New((*(t.Environment))[newchannel].Sdk.Context(), mspclient.WithOrg(t.Info.OrgName))
	if err!=nil{
		return fmt.Errorf("根据指定的Orgname创建的 ORG MSP客户端实例失败：%v", err)
	}

	adminIdentity, err:= mspClient.GetSigningIdentity(t.Info.OrgAdmin)
	// fmt.Println(t.Info.AdminIdentity)
	if err!=nil{
		return fmt.Errorf("获取指定id的签名标识失败%v", err)
	}
	channelReq := resmgmt.SaveChannelRequest{ChannelID:newchannel, ChannelConfigPath:ChannelConfig, SigningIdentities:[]msp.SigningIdentity{adminIdentity}}

	_, err = (*(t.Environment))[newchannel].OrgResMgmt.SaveChannel(channelReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(t.Info.OrdererOrgName))
	if err!=nil{
		return fmt.Errorf("创建应用通道失败：%v", err)
	}

	fmt.Println("新通道已成功创建")
    //加入通道
    err = (*(t.Environment))[newchannel].OrgResMgmt.JoinChannel(newchannel, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(t.Info.OrdererOrgName))
	if err != nil{
		return fmt.Errorf("Peers加入通道失败：%v",err)
	}

	fmt.Println("peers 已成功加入通道")
    
    //链码安装实例化
    fmt.Println("开始安装链码...")
	ccPkg, err := gopackager.NewCCPackage(t.Info.ChaincodePath, t.Info.ChaincodeGoPath)
	if err!=nil{
		return fmt.Errorf("创建链码包失败：%v", err)
	}

	installCCReq := resmgmt.InstallCCRequest{Name: t.Info.ChaincodeID, Path: t.Info.ChaincodePath, Version: ChaincodeVersion, Package: ccPkg}
    
    _, err = (*(t.Environment))[newchannel].OrgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
    if err!=nil{
    	return fmt.Errorf("安装链码失败: %v", err)
    }
    fmt.Println("指定的链码安装成功")
    


    fmt.Println("开始实例化链码....")

    ccPolicy := cauthdsl.SignedByAnyMember([]string{"OrgReqMSP"})
    // fmt.Println("ccpolicy :%v", ccPolicy)

    instantiateCCReq := resmgmt.InstantiateCCRequest{Name: t.Info.ChaincodeID, Path: t.Info.ChaincodePath, Version: ChaincodeVersion, Args: [][]byte{[]byte("init")}, Policy: ccPolicy}
    
    // fmt.Println("test ::: %v", instantiateCCReq)
    // fmt.Println("info.ChannelID %v", info.ChannelID)
    _, err= (*(t.Environment))[newchannel].OrgResMgmt.InstantiateCC(newchannel, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))

    if err != nil{
    	return fmt.Errorf("实例化链码失败：%v", err)
    }

    fmt.Println("实例化链码成功")

 //    clientChannelContext := (*(t.Environment))[newchannel].Sdk.ChannelContext(newchannel, fabsdk.WithUser(t.Info.UserName), fabsdk.WithOrg(t.Info.OrgName))

 //    channelClient, err := channel.New(clientChannelContext)
 //    if err!=nil{
 //    	return fmt.Errorf("创建应用通道客户端失败: %v", err)
 //    }

 //    fmt.Println("通道客户端创建成功， 可用此客户端调用链码进行查询或执行事务 newchannel")
 //    (*(t.Environment))[newchannel].ChannelClient = channelClient

	return nil
}