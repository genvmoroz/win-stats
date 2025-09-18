# Win Stats Picker Project Documentation

## Overview

The Win Stats Picker project is a Go-based application designed to collect and serve hardware statistics. It includes an API for retrieving these stats and a service for running the application.

## Project Structure

- `picker/api/openapi.yml`: OpenAPI specification for the Win Stats Picker API.
- `picker/api/openapi.cfg.yml`: Configuration for generating Go code from the OpenAPI specification.
- `picker/cmd/service/main.go`: Main entry point for the service.
- `prometheus-collector/deployment/grafana/provisioning/dashboards/win_stats.json`: Grafana dashboard configuration for visualizing the collected stats.

## API Documentation

The API is documented using OpenAPI 3.1.0. Below are the main endpoints:

### `/stats`

- **Method**: GET
- **Summary**: Returns a map of hardware stats.
- **Responses**:
    - `200 OK`: Returns a list of hardware stats.
    - `500 Internal Server Error`: Returns an error response.

### `/health`

- **Method**: GET
- **Summary**: Returns the health status of the API.
- **Responses**:
    - `200 OK`: Returns a plain text status.

## Components

### Schemas

- **Stats**: Describes a response to the GetStats endpoint.
- **Hardware**: Describes a hardware component.
- **Sensor**: Describes a sensor.
- **SensorValue**: Describes a sensor value.
- **Error**: Describes an error response.

## Service

The service is implemented in Go and is responsible for running the HTTP server. The main entry point is `picker/cmd/service/main.go`.

### Main Function

The `main` function initializes the service and handles graceful shutdown on system signals.

### Start Function

The `start` function sets up signal handling, recovers from panics, and starts the HTTP server using an error group for concurrency.

## Dependencies

Dependencies are managed using Go modules (`go.mod`). The project imports several packages, including:

- `context`
- `fmt`
- `log`
- `os/signal`
- `runtime/debug`
- `syscall`
- `github.com/genvmoroz/win-stats/picker/internal/dependency`
- `golang.org/x/sync/errgroup`

## Grafana Dashboard

The Grafana dashboard configuration is located in `prometheus-collector/deployment/grafana/provisioning/dashboards/win_stats.json`. It visualizes the hardware stats collected by the service.

## Configuration

The OpenAPI configuration file (`picker/api/openapi.cfg.yml`) specifies the package name, generation options, and output path for the generated Go code.

## Running the Service

To run the service, execute the `main.go` file:

```sh
go run picker/cmd/service/main.go
