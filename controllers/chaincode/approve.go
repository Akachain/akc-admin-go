package chaincode

import (
	"fmt"
	"strconv"

	"github.com/Akachain/akc-admin-go/common"

	"github.com/gin-gonic/gin"
	cliCommon "github.com/hyperledger/fabric-cli/cmd/commands/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

type ApproveRequest struct {
	OrgName             string `json:"orgName"`
	ChaincodeId         string `json:"chaincodeId"`
	ChaincodeVersion    string `json:"chaincodeVersion"`
	PackageID           string `json:"packageId"`
	SignaturePolicy     string `json:"signaturePolicy"`
	ChannelConfigPolicy string `json:"channelConfigPolicy"`
	CollectionsConfig   string `json:"collectionsConfig"`
	Sequence            string `json:"sequence"`
	InitRequired        bool   `json:"initRequired"`
	EndorsementPlugin   string `json:"endorsementPlugin"`
	ValidationPlugin    string `json:"validationPlugin"`
}

func ApproveChaincode(c *gin.Context) {
	var err error
	var msg string

	var approveRequest ApproveRequest
	c.BindJSON(&approveRequest)

	// Load client context
	context, resourceManagement, err := common.GetResourcesByOrg(approveRequest.OrgName)
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	signaturePolicy, err := cliCommon.GetChaincodePolicy(approveRequest.SignaturePolicy)
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	collectionsConfig, err := cliCommon.GetCollectionConfigFromFile(approveRequest.CollectionsConfig)
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	sequence, err := strconv.ParseInt(approveRequest.Sequence, 10, 64)
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	req := resmgmt.LifecycleApproveCCRequest{
		Name:                approveRequest.ChaincodeId,
		Version:             approveRequest.ChaincodeVersion,
		PackageID:           approveRequest.PackageID,
		Sequence:            sequence,
		SignaturePolicy:     signaturePolicy,
		ChannelConfigPolicy: approveRequest.ChannelConfigPolicy,
		CollectionConfig:    collectionsConfig,
		InitRequired:        approveRequest.InitRequired,
		EndorsementPlugin:   approveRequest.EndorsementPlugin,
		ValidationPlugin:    approveRequest.ValidationPlugin,
	}

	options := []resmgmt.RequestOption{
		resmgmt.WithTargetEndpoints(context.Peers...),
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	}

	if _, err := resourceManagement.LifecycleApproveCC(context.Channel, req, options...); err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	msg = fmt.Sprintf("Successfully approved chaincode '%s'", approveRequest.ChaincodeId)

	c.JSON(200, common.RequestResponse(true, msg))
}
