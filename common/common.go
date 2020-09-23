package common

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-cli/pkg/environment"
	"github.com/hyperledger/fabric-cli/pkg/fabric"
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
		"data": data,
	}
}

func GetResources() (*environment.Context, fabric.ResourceManagement, error) {
	var err error
	var config *environment.Config
	config = environment.NewConfig()

	err = config.LoadFromFile(filepath.Join("configs", "config.yaml"))
	if err != nil {
		return nil, nil, err
	}

	factory, err := fabric.NewFactory(config)
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

//func GetResourcesByOrg(orgName string) (*environment.Context, fabric.ResourceManagement, error) {
//	var err error
//	var config *environment.Config
//	config = environment.NewConfig()
//
//	err = config.LoadFromFile(filepath.Join("configs", "config.yaml"))
//	if err != nil {
//		return nil, nil, err
//	}
//
//	factory, err := fabric.NewFactory(config)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	resourceManagement, err := factory.ResourceManagement()
//	if err != nil {
//		return nil, nil, err
//	}
//
//	context := config.Contexts[orgName]
//
//	return context, resourceManagement, nil
//}

// func LoadFabricSDK() (*fabsdk.FabricSDK, error) {
// configPath := filepath.Join("configs", "networks", "harisato.yaml")
// backend, err := config.FromFile(configPath)()
// if err != nil {
// 	return nil, err
// }

// configProvider := func() ([]core.ConfigBackend, error) {
// 	return backend, nil
// }

// sdk, err := fabsdk.New(configProvider)
// if err != nil {
// 	return nil, err
// }
// }

func LoadFabricSDK() (fabric.SDK, error) {
	var err error
	var config *environment.Config
	config = environment.NewConfig()

	err = config.LoadFromFile(filepath.Join("configs", "config.yaml"))
	if err != nil {
		return nil, err
	}

	factory, err := fabric.NewFactory(config)
	if err != nil {
		return nil, err
	}

	sdk, err := factory.SDK()
	if err != nil {
		return nil, err
	}

	return sdk, nil
}
