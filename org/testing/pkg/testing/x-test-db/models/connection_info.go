package models

import "fmt"

type DatabaseConnectionInfo struct {
	HostName       DatabaseHostName
	Port           DatabaseHostPort
	HostNameDocker DatabaseHostNameDocker
	PortDocker     DatabaseHostPortDocker
	DatabaseName   DatabaseName
	PoolMaxConn    int32
	OwnerUser      string
	OwnerPassword  string
	TLSConfig      DatabaseTLSConfig
}

func NewDatabaseConnectionInfo(
	hostName DatabaseHostName,
	port DatabaseHostPort,
	hostNameDocker DatabaseHostNameDocker,
	portDocker DatabaseHostPortDocker,
	databaseName DatabaseName,
	poolMaxConn int32,
	ownerUser string,
	ownerPassword string,
	tlsConfig DatabaseTLSConfig,
) DatabaseConnectionInfo {
	return DatabaseConnectionInfo{
		HostName:       hostName,
		Port:           port,
		HostNameDocker: hostNameDocker,
		PortDocker:     portDocker,
		DatabaseName:   databaseName,
		PoolMaxConn:    poolMaxConn,
		OwnerUser:      ownerUser,
		OwnerPassword:  ownerPassword,
		TLSConfig:      tlsConfig,
	}
}

func (c *DatabaseConnectionInfo) PreparePostgreConnectionString() string {
	dbURI := fmt.Sprintf("host=%s user=%s password=%s port=%d database=%s sslmode=%s",
		c.HostName, c.OwnerUser, c.OwnerPassword, c.Port, c.DatabaseName, c.TLSConfig.Mode)

	if len(c.TLSConfig.ServerRootCAPath) > 0 {
		dbURI += fmt.Sprintf(" sslrootcert=%s", c.TLSConfig.ServerRootCAPath)
	}
	if len(c.TLSConfig.ServerRootCAPath) > 0 {
		dbURI += fmt.Sprintf(" sslcert=%s", c.TLSConfig.ClientCertPath)
	}
	if len(c.TLSConfig.ServerRootCAPath) > 0 {
		dbURI += fmt.Sprintf(" sslkey=%s", c.TLSConfig.ClientKeyPath)
	}
	if c.PoolMaxConn > 0 {
		dbURI += fmt.Sprintf(" pool_max_conns=%d", c.PoolMaxConn)
	}
	return dbURI
}
