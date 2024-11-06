# Queuer

**Queuer** is a simple, configurable task queue designed for developers and system administrators to manage tasks that interact with services and databases efficiently. It reads tasks from a queue, processes them through defined services, and logs the results to maintain a robust transaction history. Queuer is written in Go and supports PostgreSQL as its database backend.

## Features

- **Queue Processing**: Reads from a queue database and processes each task by matching it to a configured service.
- **Service-Oriented Configuration**: Configure different services for each queue in a single `config.json` file.
- **Logging and Transaction Tracking**: Logs each processed task to a transaction database for tracking, and writes diagnostic information to a log database.
- **Flexible Verbosity**: Adjust logging verbosity using command-line flags for detailed insights into execution.

## Installation

To run Queuer, you need to have Go installed. Simply clone the repository and run Queuer with:

```sh
go run .
```

## Usage

### Command-Line Options

```sh
Usage:
  queuer [flags]

Flags:
  -f, --file string       Path to the configuration file (default "config.json")
  -h, --help              help for queuer
  -v, --verbosity count   Set the verbosity level (e.g., -v for verbose, -vv for very verbose)
```

### Examples

```sh
# Run with default config.json
queuer

# Run with a custom configuration file and verbose logging
queuer -f custom_config.json -vv
```

## Configuration

Queuer uses a JSON configuration file to define each queue’s service, environment, and database connections. Each task is matched to a service specified in this file and processed accordingly. Below is an example configuration.

### Example `config.json`

```json
[
  {
    "name": "Addition Service",
    "environment": "development",
    "service": "adder",
    "timeout": 1000,
    "retries": 3,

    "queueDatabaseHost": "localhost",
    "queueDatabasePort": "5430",
    "queueDatabaseName": "queue",

    "targetDatabaseHost": "localhost",
    "targetDatabasePort": "5431",
    "targetDatabaseName": "target",

    "logDatabaseHost": "localhost",
    "logDatabasePort": "5432",
    "logDatabaseName": "logs"
  },
  {
    "name": "Square Service",
    "environment": "prod",
    "service": "squarer",
    "timeout": 1000,
    "retries": 3,

    "queueDatabaseHost": "localhost",
    "queueDatabasePort": "5430",
    "queueDatabaseName": "queue",

    "targetDatabaseHost": "localhost",
    "targetDatabasePort": "5431",
    "targetDatabaseName": "target",

    "logDatabaseHost": "localhost",
    "logDatabasePort": "5432",
    "logDatabaseName": "logs"
  }
]
```

#### Configuration Fields

- `name`: Name of the service (for reference).
- `environment`: Environment in which the service is running (e.g., `development`, `prod`).
- `service`: Type of service; tasks in the queue are matched by this value (e.g., `adder`, `squarer`).
- `timeout`: Maximum duration (in milliseconds) allowed for a task to run.
- `retries`: Number of retries for a task in case of failure.
- `queueDatabaseHost`, `queueDatabasePort`, `queueDatabaseName`: Connection details for the queue database.
- `targetDatabaseHost`, `targetDatabasePort`, `targetDatabaseName`: Connection details for the target database where results are stored.
- `logDatabaseHost`, `logDatabasePort`, `logDatabaseName`: Connection details for the log database.

## Logging Levels

Queuer provides flexible logging options with adjustable verbosity:

- Default (no `-v` flag): Logs only errors.
- `-v`: Logs warnings and informational messages.
- `-vv`: Adds debug messages.
- `-vvv`: Adds full path details for deep diagnostics.

### Environment Configuration

Database connection parameters for each service should be set via environment variables. Only PostgreSQL is currently supported.

## Limitations

- **Database Compatibility**: Currently, only PostgreSQL is supported for the queue, target, and log databases.
- **Single Instance Execution**: Each instance reads from and processes a single queue at a time.
