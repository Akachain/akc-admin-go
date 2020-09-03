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
