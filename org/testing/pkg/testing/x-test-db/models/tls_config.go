package models

type DatabaseTLSConfig struct {
	ClientKeyPath    string
	ClientCertPath   string
	ServerRootCAPath string
	Mode             string
}

func NewDatabaseTLSConfig(
	clientKeyPath string,
	clientCertPath string,
	serverRootCAPath string,
	mode string,
) DatabaseTLSConfig {
	return DatabaseTLSConfig{
		ClientKeyPath:    clientKeyPath,
		ClientCertPath:   clientCertPath,
		ServerRootCAPath: serverRootCAPath,
		Mode:             mode,
	}
}
