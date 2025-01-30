## API Repo

### Introduction

This repo shall hold all the API definition files (OpenAPI, AsyncAPI, Proto, Kafka Topic) and shall be responsible for generating and publishing code for API Models, Client, Server implementations for Golang, Typescript and others.

### Output

#### Golang

The codegen folder shall contain a go module. With go.mod and pkg folder that contains all the generated code

#### Typescript

The codegen shall create a npm package that contains all the models and client/server code

### Folder Structure

```
+-- <api-name>
|   +-- api                     # Directory where the api definitions themselves go
|   |   +-- openapi <optional>
|   |   |   ...
|   |   +-- asyncapi <optional>
|   |   |   ...
|   |   +-- kafka <optional>
|   |   |   ...
|   |   +-- proto <optional>
|   |   |   ...
|   +-- codegen                 # Directory where the generated code of the api goes
|   |   +-- golang
|   |   |   go.mod
|   |   |   pkg
|   |   |   ...
|   |   +-- typescript
|   |   |   client
|   |   |   server
|   |   |   ...
|   |   CHANGELOG.md            # Changelog for tracking the API changes
+-- ...                         # Repeat for other APIs
Makefile                        # Commands to run all generate, lint, publish, build commands
CHANGELOG.md                    # Changelog for tracking overall repo changes
...                             # Configuration files for generators, linters, pipelines etc.
```

### Linters

- OpenAPI - [Redocly CLI](https://redocly.com/docs/cli/)
- AsyncAPI - [Redocly CLI](https://redocly.com/docs/cli/)
- Kafka - TODO
- Proto - TODO

### Generators

- OpenAPI to TS (React) [Orval](https://orval.dev/overview)
- OpenAPI to Golang [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
