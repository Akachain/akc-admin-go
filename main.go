package main

import (
	"github.com/Akachain/akc-admin-go/controllers/chaincode"
	"github.com/Akachain/akc-admin-go/controllers/channel"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	// r.Static("/public", "./public")

	client := r.Group("/api")
	{
		client.POST("/chaincode/install", chaincode.InstallChaincode)
		client.POST("/chaincode/approve", chaincode.ApproveChaincode)
		client.POST("/channel/create", channel.CreateChannel)
	}

	return r
}

func main() {
	r := setupRouter()
	r.Run(":4001")
}
