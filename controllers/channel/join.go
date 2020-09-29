package channel

import (
	"fmt"
	"github.com/Akachain/akc-admin-go/common"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"strings"
)

type JoinChannelRequest struct {
	OrgName     string `json:"orgName"`
	Peer        string `json:"peer"`
	ChannelName string `json:"channelName"`
}

func JoinChannel(c *gin.Context) {
	var err error
	var msg string

	var joinRequest JoinChannelRequest
	c.BindJSON(&joinRequest)

	// Load client context
	context, resourceManagement, err := common.GetResources()
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	var peers []string
	if joinRequest.Peer != "" {
		peers = append(peers, joinRequest.Peer)
	} else {
		peers = context.Peers
	}

	options := []resmgmt.RequestOption{
		resmgmt.WithTargetEndpoints(peers...),
	}

	if err := resourceManagement.JoinChannel(joinRequest.ChannelName, options...); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			msg = "Join channel fail: channel in peer already exists"
		} else {
			msg = err.Error()
		}
		c.JSON(200, common.RequestResponse(false, msg))
		return
	}

	msg = fmt.Sprintf("Successfully join channel '%s' to '%s'.", joinRequest.ChannelName, joinRequest.Peer)
	c.JSON(200, common.RequestResponse(true, msg))
}
