//start.go
package sdkInit

import(
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"

	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
)

const ChaincodeVersion = "1.0"

func Getkeys(username string){
	//得到私钥
	privateKey,_:=rsa.GenerateKey(rand.Reader, 2048)
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	x509_Privatekey:=x509.MarshalPKCS1PrivateKey(privateKey)
	//创建一个用来保存私钥的以.pem结尾的文件
	fp,_:=os.Create(username + "_PrivateKey.pem")
	defer fp.Close()
	//将私钥字符串设置到pem格式块中
	pem_block:=pem.Block{
		Type:"csdn_privateKey",
		Bytes:x509_Privatekey,
	}
	//转码为pem并输出到文件中
	pem.Encode(fp,&pem_block)

	//处理公钥,公钥包含在私钥中
	publickKey:=privateKey.PublicKey
	//接下来的处理方法同私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	x509_PublicKey,_:=x509.MarshalPKIXPublicKey(&publickKey)
	pem_PublickKey:=pem.Block{
		Type:"csdn_PublicKey",
		Bytes:x509_PublicKey,
	}
	file,_:=os.Create(username + "_PublicKey.pem")
	defer file.Close()
	//转码为pem并输出到文件中
	pem.Encode(file,&pem_PublickKey)
}

func SetupSDK(ConfigFile string, initialized bool) (*fabsdk.FabricSDK, error){

	if initialized{
		return nil, fmt.Errorf("Fabric SDK 已被实例化")
	}

	sdk, err := fabsdk.New(config.FromFile(ConfigFile))
	if err != nil {
		return nil, fmt.Errorf("实例化Fabric SDK失败: %v", err)
	}

	fmt.Println("Fabric SDK 初始化成功")
	return sdk, nil
}

func CreateSourceClient(envir *Environ, info *InitInfo) error{
	clientContext :=envir.Sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName))
	if clientContext == nil{
		return fmt.Errorf("根据指定的组织名称与管理员创建资源管理客户端Context失败")
	}

	resMgmtClient, err := resmgmt.New(clientContext)
	if err!=nil{
		return fmt.Errorf("根据指定的资源管理客户端Context创建通道管理客户端失败：%v", err)
	}

	envir.OrgResMgmt = resMgmtClient
	fmt.Println("创建资源客户端成功")
	return nil
}

func Register(envir *Environ, username, password string) error{

 	ctx := envir.Sdk.Context()
    c, err := mspclient.New(ctx)
    if err != nil{
    	fmt.Println("failed to create msp client")
    	return err
    }

    department := "org1.department1"

	request := &mspclient.RegistrationRequest{
		Name:        username,
		Type:        "client",
		Affiliation: department,
		Secret:      password,
	}

    secret, err := c.Register(request)
	if err != nil {
		fmt.Printf("register %s [%s]\n", username, err)
		return err
	}
	Getkeys(username)
	fmt.Printf("register %s successfully,with password %s\n", username, secret)
	return nil
}

func Enroll(envir *Environ, username, password string) error{

 	ctx := envir.Sdk.Context()
    c, err := mspclient.New(ctx)
    if err != nil{
    	fmt.Println("failed to create msp client")
    	return err
    }

    err = c.Enroll(username, mspclient.WithSecret(password))

    if err!= nil{
    	fmt.Println("fail to enroll user: %s\n", username)
    	return err
    }

	fmt.Printf("User %s successfully enrolled.\n", username)
	return nil
	// Create new identity
	// newIdentity, err := c.CreateIdentity(req)
	// if err != nil {
	// 	fmt.Println("Create identity failed: %s", err)
	// }

	// if newIdentity.Secret == "" {
	// 	fmt.Println("Secret should have been generated")
	// }

	// identity, err := c.GetIdentity(username)
	// if err != nil {
	// 	fmt.Println("get identity failed: %s", err)
	// }

	// fmt.Println("Get Identity: [%v]:", identity)

	// if !verifyIdentity(req, identity) {
	// 	fmt.Println("verify identity failed req=[%v]; resp=[%v] ", req, identity)
	// }

	// return err

}

func Getidentity(envir *Environ, name string) error{

 	ctx := envir.Sdk.Context()
    c, err := mspclient.New(ctx)
    if err != nil{
    	fmt.Println("failed to create msp client")
    	return err
    }

	_, err = c.GetSigningIdentity(name)
	if err != nil {
	    fmt.Printf("Get identitie return error %s\n", err)
	    return err
	}
    // fmt.Println("get identity successfully !")
    // fmt.Println(identity)	
    return nil
}

func CreateChannel(envir *Environ, info *InitInfo) error{
	
	mspClient, err := mspclient.New(envir.Sdk.Context(), mspclient.WithOrg(info.OrgName))
	if err!=nil{
		return fmt.Errorf("根据指定的Orgname创建的 ORG MSP客户端实例失败：%v", err)
	}

	adminIdentity, err:= mspClient.GetSigningIdentity(info.OrgAdmin)
	info.AdminIdentity = adminIdentity
	if err!=nil{
		return fmt.Errorf("获取指定id的签名标识失败%v", err)
	}

	channelReq := resmgmt.SaveChannelRequest{ChannelID:info.ChannelID, ChannelConfigPath:info.ChannelConfig, SigningIdentities:[]msp.SigningIdentity{adminIdentity}}

	_, err = envir.OrgResMgmt.SaveChannel(channelReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err!=nil{
		return fmt.Errorf("创建应用通道失败：%v", err)
	}

	fmt.Println("通道已成功创建")

	return nil
} 

func JoinChannel(envir *Environ, info *InitInfo) error{
	err := envir.OrgResMgmt.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err != nil{
		return fmt.Errorf("Peers加入通道失败：%v",err)
	}

	fmt.Println("peers 已成功加入通道")
	return nil
}  

func InstallAndInstantiateCC (envir *Environ, info *InitInfo) error{
	fmt.Println("开始安装链码...")
	ccPkg, err := gopackager.NewCCPackage(info.ChaincodePath, info.ChaincodeGoPath)
	if err!=nil{
		return fmt.Errorf("创建链码包失败：%v", err)
	}

	installCCReq := resmgmt.InstallCCRequest{Name: info.ChaincodeID, Path: info.ChaincodePath, Version: ChaincodeVersion, Package: ccPkg}
    
    _, err = envir.OrgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
    if err!=nil{
    	return fmt.Errorf("安装链码失败: %v", err)
    }
    fmt.Println("指定的链码安装成功")
    fmt.Println("开始实例化链码....")

    ccPolicy := cauthdsl.SignedByAnyMember([]string{"OrgReqMSP"})

    instantiateCCReq := resmgmt.InstantiateCCRequest{Name: info.ChaincodeID, Path: info.ChaincodePath, Version: ChaincodeVersion, Args: [][]byte{[]byte("init")}, Policy: ccPolicy}
    
    _, err= envir.OrgResMgmt.InstantiateCC(info.ChannelID, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))

    if err != nil{
    	return fmt.Errorf("实例化链码失败：%v", err)
    }

    fmt.Println("实例化链码成功")
    // return nil

    clientChannelContext := envir.Sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(info.UserName), fabsdk.WithOrg(info.OrgName))

    envir.ChannelClient, err = channel.New(clientChannelContext)
    if err!=nil{
    	return fmt.Errorf("创建应用通道客户端失败: %v", err)
    }

    fmt.Println("通道客户端创建成功， 可用此客户端调用链码进行查询或执行事务")
    return  nil
}

func CreateChannelClient(envir *Environ, info *InitInfo, name string)error{
	var err error
	clientChannelContext := envir.Sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(name), fabsdk.WithOrg(info.OrgName))
    envir.ChannelClient, err = channel.New(clientChannelContext)

    // environment[initInfo.ChannelID].ChannelClient, err := channel.New(clientChannelContext)
    if err!=nil{
    	return fmt.Errorf("创建应用通道客户端失败: %v", err)
    }

    fmt.Println("通道客户端创建成功， 可用此客户端调用链码进行查询或执行事务")
    return nil
}