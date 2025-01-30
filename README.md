# Project 2112
# Getting Started

![2112 Architecture Overview](docs/img/architecture.png)

1. Clone this repository `git clone git@github.com:elbujito/2112.git`
2. Run `cd 2112`
3. Run `go get`
4. Run `go run . db migrate`
5. Run `go run . db seed`
6. Run `go run . start` to start the server, you should see the following:
```
⇨ http server started on [::]:8081
⇨ http server started on [::]:8080
```
7. List available routes using `go run . info protected-api-routes` and use your favourite API client to test. or use the following to get started and make sure you're up and running.
```bash
curl -H "Accept: application/json" http://127.0.0.1:8081/health/alive
curl -H "Accept: application/json" http://127.0.0.1:8081/health/ready
```

> Recommended: run `go run .` and explore all available options, it should be straightforward.

For more details on running and using the service, scroll down to "[Operations](#operations)" section. 


   ```bash
   git clone git@github.com:elbujito/2112.git
   cd 2112

# GraphQL Gateway Setup

This repository sets up a **GraphQL Gateway** service implemented in Go, which interacts with a **Redis** service for Pub/Sub messaging. The setup is containerized using **Docker** and managed with **Docker Compose**.

## Project Structure

- **GraphQL Gateway (Go)**: A GraphQL API for querying satellite position data.
- **Redis Service**: Used for Pub/Sub messaging between services.
- **Docker Compose**: Manages and orchestrates the services.

## Prerequisites

- **Docker** and **Docker Compose** installed on your machine.
- **Go** installed (if you plan to modify the Go code locally).
