//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -o=./pkg/models/app.gen.go --config=../../configs/oapi-codegen/models.cfg.yaml ../../api/app-rest.model.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -o=./pkg/client/client.gen.go --config=../../configs/oapi-codegen/client.cfg.yaml ../../api/app-rest.api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -o=./pkg/server/server.gen.go --config=../../configs/oapi-codegen/server.cfg.yaml ../../api/app-rest.api.yaml

// Package restapis generator root
package restapis

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)
