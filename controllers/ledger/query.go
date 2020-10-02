package ledger

import (
	"encoding/json"
	"fmt"
	"github.com/Akachain/akc-admin-go/common"
	"github.com/Akachain/akc-admin-go/constants"
	"github.com/gin-gonic/gin"
	fabCommon "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"io/ioutil"
	"os"
	"path/filepath"
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
	var msg string

	var ledgerRequest LedgerRequest
	c.BindJSON(&ledgerRequest)

	// Load client context
	ledger, err := common.LoadFabricLedger(ledgerRequest.OrgName)
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	var block *fabCommon.Block
	switch ledgerRequest.QueryType {
	case constants.ByBlockNumber:
		blockNumber, err := strconv.ParseUint(ledgerRequest.QueryValue, 10, 64)
		if err != nil {
			c.JSON(200, common.RequestResponse(false, err.Error()))
			return
		}
		block, err = ledger.QueryBlock(blockNumber)
	case constants.ByTxID:
		var txID = fab.TransactionID(ledgerRequest.QueryValue)
		block, err = ledger.QueryBlockByTxID(txID)
	case constants.ByBlockHash:
		block, err = ledger.QueryBlockByHash([]byte(ledgerRequest.QueryType))
	default:
		msg = fmt.Sprintf("Type '%s' currently doesn't supported.", ledgerRequest.QueryType)
		c.JSON(200, common.RequestResponse(false, msg))
		return
	}

	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	if block == nil {
		c.JSON(200, common.RequestResponseData(true, "NULL"))
		return
	}

	data, _ := json.Marshal(block)

	fileName := fmt.Sprintf("block_%s_%s_%s.block", ledgerRequest.QueryValue, ledgerRequest.OrgName, ledgerRequest.ChannelName)
	filePath := filepath.Join("tmp", "blocks", fileName)
	_ = os.MkdirAll(filepath.Join("tmp", "blocks"), os.ModePerm)
	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}
	c.JSON(200, common.RequestResponse(true, "Query successfully and file block saved!"))
}
