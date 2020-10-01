package common

import (
	"path/filepath"

	"github.com/Akachain/akc-admin-go/internal/environment"
	"github.com/Akachain/akc-admin-go/internal/fabric"
	"github.com/gin-gonic/gin"
)

func RequestResponse(success bool, message string) gin.H {
	return gin.H{
		"success": success,
		"message": message,
	}
}

func RequestResponseData(success bool, data string) gin.H {
	return gin.H{
		"success": success,
		"data":    data,
	}
}

func GetResources() (*environment.Context, fabric.ResourceManagement, error) {
	var err error

	config, err := InitConfig()
	if err != nil {
		return nil, nil, err
	}

	factory, err := fabric.NewFactoryDefault(config)
	if err != nil {
		return nil, nil, err
	}

	resourceManagement, err := factory.ResourceManagement()
	if err != nil {
		return nil, nil, err
	}

	context, err := config.GetCurrentContext()
	if err != nil {
		return nil, nil, err
	}

	return context, resourceManagement, nil
}

func InitConfig() (*environment.Config, error) {
	var err error
	var config *environment.Config
	config = environment.NewConfig()

	err = config.LoadFromFile(filepath.Join("configs", "config.yaml"))
	if err != nil {
		return nil, err
	}

	return config, nil
}

func GetResourcesByOrg(orgName string) (*environment.Context, fabric.ResourceManagement, error) {
	var err error

	config, err := InitConfig()
	if err != nil {
		return nil, nil, err
	}

	factory, err := fabric.NewFactory(config, orgName)
	if err != nil {
		return nil, nil, err
	}

	resourceManagement, err := factory.ResourceManagement()
	if err != nil {
		return nil, nil, err
	}

	context, err := config.GetContextByName(orgName)
	if err != nil {
		return nil, nil, err
	}

	return context, resourceManagement, nil
}

func LoadFabricLedger(orgName string) (fabric.Ledger, error) {
	var err error

	config, err := InitConfig()
	if err != nil {
		return nil, err
	}

	factory, err := fabric.NewFactory(config, orgName)
	if err != nil {
		return nil, err
	}

	ledger, err := factory.Ledger()
	if err != nil {
		return nil, err
	}

	return ledger, nil
}

func LoadFabricSDK() (fabric.SDK, error) {
	var err error

	config, err := InitConfig()
	if err != nil {
		return nil, err
	}

	factory, err := fabric.NewFactoryDefault(config)
	if err != nil {
		return nil, err
	}

	sdk, err := factory.SDK()
	if err != nil {
		return nil, err
	}

	return sdk, nil
}
