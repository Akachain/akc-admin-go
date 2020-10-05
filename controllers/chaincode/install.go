package chaincode

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Akachain/akc-admin-go/common"

	"github.com/gin-gonic/gin"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
	lifecyclepkg "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
	// mspClient "github.com/hyperledger/fabric-sdk-go/internal/client/msp"
)

type ChaincodeRequest struct {
	OrgName          string `json:"orgName"`
	ChannelName      string `json:"channelName"`
	ChaincodeId      string `json:"chaincodeId"`
	ChaincodeVersion string `json:"chaincodeVersion"`
	ChaincodeType    string `json:"chaincodeType"`
	ChaincodePath    string `json:"chaincodePath"`
}

func InstallChaincode(c *gin.Context) {
	var err error
	var msg string

	var chaincodeRequest ChaincodeRequest
	c.BindJSON(&chaincodeRequest)

	// Load client context
	context, resourceManagement, err := common.GetResourcesByOrg(chaincodeRequest.OrgName)

	projectPath, _ := os.Getwd()
	adminGoPath := filepath.Join(projectPath, "artifacts")

	// Package chaincode
	// internal, err := gopackager.NewCCPackage(chaincodeRequest.ChaincodePath, adminGoPath)
	// if err != nil {
	// 	c.JSON(200, common.RequestResponse(false, err.Error()))
	// 	return
	// }

	pkg, err := lifecyclepkg.NewCCPackage(&lifecyclepkg.Descriptor{
		Path:  filepath.Join(adminGoPath, "src", chaincodeRequest.ChaincodePath),
		Type:  pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value[strings.ToUpper(chaincodeRequest.ChaincodeType)]),
		Label: chaincodeRequest.ChaincodeId,
	})
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	chaincodePackagedName := fmt.Sprintf("%s.%s.tgz", chaincodeRequest.ChaincodeId, chaincodeRequest.ChaincodeVersion)
	if err := ioutil.WriteFile(filepath.Join("artifacts", chaincodePackagedName), pkg, 0644); err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	fmt.Printf("Successfully packaged chaincode '%s'\n", chaincodeRequest.ChaincodeId)

	responses, err := resourceManagement.LifecycleInstallCC(
		resmgmt.LifecycleInstallCCRequest{
			Label:   chaincodeRequest.ChaincodeId,
			Package: pkg,
		},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithTargetEndpoints(context.Peers...),
	)
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	if len(responses) == 0 {
		packageID := lifecycle.ComputePackageID(chaincodeRequest.ChaincodeId, pkg)
		msg = fmt.Sprintf("Chaincode '%s' has already been installed on all peers. Package ID '%s'", chaincodeRequest.ChaincodeId, packageID)
	} else {
		msg = fmt.Sprintf("Successfully installed chaincode '%s'. Package ID '%s'", chaincodeRequest.ChaincodeId, responses[0].PackageID)
	}

	c.JSON(200, common.RequestResponse(true, msg))
}
