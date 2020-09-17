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

type CommitRequest struct {
	ChaincodeId         string `json:"chaincodeId"`
	ChaincodeVersion    string `json:"chaincodeVersion"`
	SignaturePolicy     string `json:"signaturePolicy"`
	ChannelConfigPolicy string `json:"channelConfigPolicy"`
	CollectionsConfig   string `json:"collectionsConfig"`
	Sequence            string `json:"sequence"`
	InitRequired        bool   `json:"initRequired"`
	EndorsementPlugin   string `json:"endorsementPlugin"`
	ValidationPlugin    string `json:"validationPlugin"`
}

func CommitChaincode(c *gin.Context) {
	var err error
	var msg string

	var approveRequest CommitRequest
	c.BindJSON(&approveRequest)

	// Load client context
	context, resourceManagement, err := common.GetResources()
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

	req := resmgmt.LifecycleCommitCCRequest{
		Name:                approveRequest.ChaincodeId,
		Version:             approveRequest.ChaincodeVersion,
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

	if _, err := resourceManagement.LifecycleCommitCC(context.Channel, req, options...); err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	msg = fmt.Sprintf("Successfully commit chaincode '%s'", approveRequest.ChaincodeId)

	c.JSON(200, common.RequestResponse(true, msg))
}
