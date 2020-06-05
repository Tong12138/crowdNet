//fabricInitInfo.go

package sdkInit

import(
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"

)

type InitInfo struct{
	ChannelID string
	ChannelConfig string
	OrgAdmin string
	OrgName string
	OrdererOrgName string
    AdminIdentity msp.SigningIdentity


	ChaincodeID string
	ChaincodeGoPath string
	ChaincodePath string
	UserName string

}

type Environ struct{
	Sdk *fabsdk.FabricSDK
	OrgResMgmt *resmgmt.Client
	ChannelClient *channel.Client 
}