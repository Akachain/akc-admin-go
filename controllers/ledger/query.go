package ledger

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Akachain/akc-admin-go/common"
	"github.com/Akachain/akc-admin-go/constants"
	"github.com/gin-gonic/gin"
	fabCommon "github.com/hyperledger/fabric-protos-go/common"
	"strconv"
)

// LedgerRequest parameters
type LedgerRequest struct {
	Peer        string `json:"peer"`
	OrgName     string `json:"orgName"`
	ChannelName string `json:"channelName"`
	QueryType   string `json:"queryType"`
	QueryValue  string `json:"queryValue"`
}

func QueryBlock(c *gin.Context) {
	var err error

	var ledgerRequest LedgerRequest
	c.BindJSON(&ledgerRequest)

	// Load client context
	ledger, err := common.LoadFabricLedger(ledgerRequest.OrgName)
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	var block *fabCommon.Block
	if ledgerRequest.QueryType == constants.BY_BLOCK_NUMBER {
		blockNumber, err := strconv.ParseUint(ledgerRequest.QueryValue, 10, 64)
		block, err = ledger.QueryBlock(blockNumber)
		if err != nil {
			c.JSON(200, common.RequestResponse(false, err.Error()))
			return
		}
	}

	data, _ := json.Marshal(block.Data.Data[0])

	decoded, err := base64.StdEncoding.DecodeString(string(data))
	print(decoded)


	c.JSON(200, common.RequestResponseData(true, string(data)))
}
