package msp

import (
	"fmt"
	"strings"

	"github.com/Akachain/akc-admin-go/common"
	"github.com/gin-gonic/gin"

	mspClient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
)

type UserRequest struct {
	OrgName    string                `json:"orgName"`
	UserName   string                `json:"userName"`
	Type       string                `json:"type"`
	Attributes []mspClient.Attribute `json:"attrs"`
}

func RegisterUser(c *gin.Context) {
	var err error
	var msg string

	var regRequest UserRequest
	c.BindJSON(&regRequest)

	sdk, err := common.LoadFabricSDK()

	ctxProvider := sdk.Context()
	msp, err := mspClient.New(ctxProvider)
	enrollSecret, err := msp.Register(&mspClient.RegistrationRequest{
		Name:        regRequest.UserName,
		Type:        regRequest.Type,
		Affiliation: regRequest.OrgName,
		Attributes:  regRequest.Attributes,
	})

	if err != nil {
		if strings.Contains(err.Error(), "already registered") {
			msg = fmt.Sprintf("Identity '%s' is already registered.", regRequest.UserName)
			c.JSON(200, common.RequestResponse(true, msg))
			return
		}
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	// Enroll the new user
	err = msp.Enroll(regRequest.UserName, mspClient.WithSecret(enrollSecret))
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	// Get the new user's signing identity
	// si, err := msp.GetSigningIdentity(regRequest.UserName)
	// if err != nil {
	// 	c.JSON(200, common.RequestResponse(false, err.Error()))
	// 	return
	// }

	// println(si)

	msg = fmt.Sprintf("Register user '%s' is completed! Enrolled with secret '%s'.", regRequest.UserName, enrollSecret)

	c.JSON(200, common.RequestResponse(true, msg))
}

func RevokeUser(c *gin.Context) {
	var err error
	var msg string

	var revokeRequest UserRequest
	c.BindJSON(&revokeRequest)

	sdk, err := common.LoadFabricSDK()

	ctxProvider := sdk.Context()
	msp, err := mspClient.New(ctxProvider)

	revokeResponse, err := msp.Revoke(&mspClient.RevocationRequest{
		Name:   revokeRequest.UserName,
		GenCRL: true,
	})
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	if revokeResponse.RevokedCerts == nil {
		msg = fmt.Sprintf("User '%s' has been revoked.", revokeRequest.UserName)
		c.JSON(200, common.RequestResponse(true, msg))
		return
	}

	msg = fmt.Sprintf("Successfully revoke user '%s'", revokeRequest.UserName)
	c.JSON(200, common.RequestResponse(true, msg))
}

type EnrollRequest struct {
	OrgName  string `json:"orgName"`
	UserName string `json:"userName"`
	Secret   string `json:"enrollSecret"`
}

func EnrollUser(c *gin.Context) {
	var err error
	var msg string

	var enrollUser EnrollRequest
	c.BindJSON(&enrollUser)

	sdk, err := common.LoadFabricSDK()

	ctxProvider := sdk.Context()
	msp, err := mspClient.New(ctxProvider)

	// Enroll the new user
	err = msp.Enroll(enrollUser.UserName, mspClient.WithSecret(enrollUser.Secret))
	if err != nil {
		c.JSON(200, common.RequestResponse(false, err.Error()))
		return
	}

	msg = fmt.Sprintf("Successfully enroll user '%s'", enrollUser.UserName)
	c.JSON(200, common.RequestResponse(true, msg))
}
