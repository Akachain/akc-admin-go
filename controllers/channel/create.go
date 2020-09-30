package channel

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Akachain/akc-admin-go/common"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
)

type ChannelRequest struct {
	OrgName       string `json:"orgName"`
	ChannelName   string `json:"channelName"`
	ChannelConfig string `json:"channelConfig"`
}

func CreateChannel(c *gin.Context) {
	var err error
	var msg string

	var channelRequest ChannelRequest
	c.BindJSON(&channelRequest)

	// Load client context
	_, resourceManagement, err := common.GetResources()
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	projectPath, _ := os.Getwd()
	channelConfigTx := filepath.Join(projectPath, "configs", channelRequest.ChannelConfig)

	channelConfigTxRead, err := os.Open(channelConfigTx)
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	// Close read file config channel after application return.
	defer channelConfigTxRead.Close()

	if _, err := resourceManagement.SaveChannel(resmgmt.SaveChannelRequest{
		ChannelID:     channelRequest.ChannelName,
		ChannelConfig: channelConfigTxRead,
	}); err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	msg = fmt.Sprintf("Successfully created channel '%s'", channelRequest.ChannelName)
	c.JSON(200, common.RequestResponse(true, msg))
}
