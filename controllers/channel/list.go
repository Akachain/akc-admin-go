package channel

import (
	"encoding/json"
	"github.com/Akachain/akc-admin-go/common"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
)

type ListChannelRequest struct {
	Peer    string `json:"peer"`
	OrgName string `json:"orgName"`
}

func ListChannel(c *gin.Context) {
	var err error

	var listRequest ListChannelRequest
	c.BindJSON(&listRequest)

	// Load client context
	context, resourceManagement, err := common.GetResourcesByOrg(listRequest.OrgName)
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	var peers []string
	if listRequest.Peer != "" {
		peers = append(peers, listRequest.Peer)
	} else {
		peers = append(peers, context.Peers[0])
	}

	options := []resmgmt.RequestOption{
		resmgmt.WithTargetEndpoints(peers...),
	}

	resp, err := resourceManagement.QueryChannels(options...)
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	type DataResponse struct {
		Channels []string
	}

	var data DataResponse
	for _, channel := range resp.Channels {
		data.Channels = append(data.Channels, channel.ChannelId)
	}

	res, _ := json.Marshal(data)
	c.JSON(200, common.RequestResponseData(true, string(res)))
}
