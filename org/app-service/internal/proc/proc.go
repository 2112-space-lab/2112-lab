package proc

import (
	clientsPkg "github.com/org/2112-space-lab/org/app-service/internal/clients"
	"github.com/org/2112-space-lab/org/app-service/internal/clients/dbc"
	"github.com/org/2112-space-lab/org/app-service/internal/clients/service"
	"github.com/org/2112-space-lab/org/app-service/internal/config"
	"github.com/org/2112-space-lab/org/app-service/internal/data/models"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
)

func InitServiceEnv(serviceName string, version string) {
	config.SetServiceName(serviceName)
	config.SetServiceVersion(version)
	config.InitEnv()
	config.ResolveDevMode()
	config.InitFeatures()
	config.ResolveFlags()
	config.PrintEnvInEnvMode()
}

var clients []clientsPkg.IClient

func InitClients() {
	InitServiceClient()
	InitDbClient()
}

func ConfigureClients() {
	log.Debug("Configuring clients ...")
	for _, c := range clients {
		feature := config.Feature(c.Name())
		if feature.IsEnabled() {
			log.Debugf("Configuring %s client ...", c.Name())
			c.Configure(feature.Config)
			continue
		}
		log.Warnf("Client: '%s' is disabled, This may cause runtime errors if this client is used.", c.Name())
	}
}

func addClient(client clientsPkg.IClient) {
	clients = append(clients, client)
}

func InitServiceClient() {
	client := service.GetClient()
	log.Debugf("Activating %s client ...", client.Name())
	addClient(client)
}

func InitDbClient() {
	client := dbc.GetDBClient()
	log.Debugf("Activating %s client ...", client.Name())
	client.SetSilent(!config.DevModeFlag)
	addClient(client)
}

func InitDbConnection() {
	log.Debug("Initializing database connection ...")
	dbc.GetDBClient().InitDBConnection()
}

func InitModels() {
	log.Debug("Activating models ...")
	models.Init(dbc.GetDBClient().DB)
}
