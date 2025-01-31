package models

import "fmt"

type ContainerPort int32

const DefaultContainerHttpPort ContainerPort = 8080
const DefaultContainerHttpsPort ContainerPort = 8443
const DefaultContainerGrpcPort ContainerPort = 5101

func (p ContainerPort) ToString() string {
	return fmt.Sprintf("%d", p)
}

func (p ContainerPort) ToBindAddress() string {
	return fmt.Sprintf(":%d", p)
}

func (p ContainerPort) ToDockerTcpExposedPort() string {
	return fmt.Sprintf("%d/tcp", p)
}
