package main

import (
	"github.com/Akachain/akc-admin-go/controllers/chaincode"
	"github.com/Akachain/akc-admin-go/controllers/channel"
	"github.com/Akachain/akc-admin-go/controllers/msp"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	// r.Static("/public", "./public")

	client := r.Group("/api")
	{
		client.POST("/chaincode/install", chaincode.InstallChaincode)
		client.POST("/chaincode/approve", chaincode.ApproveChaincode)
		client.POST("/chaincode/commit", chaincode.CommitChaincode)
		client.POST("/channel/create", channel.CreateChannel)
		client.POST("/msp/registerUser", msp.RegisterUser)
		client.POST("/msp/revokeUser", msp.RevokeUser)
		client.POST("/msp/enrollUser", msp.EnrollUser)
	}

	return r
}

func main() {
	r := setupRouter()
	r.Run(":4001")
}
