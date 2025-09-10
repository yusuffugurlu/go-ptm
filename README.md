# Go PTM Project

This project is a **Personal Transaction Manager (PTM)** application developed using the **Go** programming language. It enables users to manage their financial transactions, track their balances, and view their transaction history.

## Features

- **User Management**: User registration, login, and authorization.
- **Balance Management**: View and update user balances.
- **Transaction Management**: Record financial transactions and view transaction history.
- **Caching**: Redis-based caching for performance optimization.
- **Logging and Monitoring**: Integration with Prometheus and Grafana for system monitoring and logging.

## Project Structure

The project is organized as follows:

- **cmd/**: Entry point of the application.
- **internal/**: Business logic, services, data access layer, and controllers.
- **pkg/**: Shared utility libraries.
- **build/**: Configuration files for monitoring and logging.
- **configs/**: Database and application configurations.
- **postgres-data/**: PostgreSQL database files.

## Requirements

- **Go** (1.20+)
- **Docker** and **Docker Compose**
- **Redis**
- **PostgreSQL**

## Setup

1. Install project dependencies:
   ```bash
   go mod tidy
   ```

2. Start Docker containers:
   ```bash
   docker-compose up -d
   ```

3. Run the application:
   ```bash
   go run cmd/app/main.go
   ```

## Monitoring and Logging

- **Prometheus** and **Grafana** configurations are located in the `build/` directory.
- Access the monitoring dashboards at `http://localhost:3000`.