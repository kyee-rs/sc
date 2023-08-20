# 📂🗄️ Simple Cache: Temporary (or not?) file hosting
*Made with ❤️ by [@12subnet](https://github.com/12subnet). Project heavily inspired by [0x0.st](https://git.0x0.st/mia/0x0)*

> ⚠️ **WARNING**: This branch is used specially to integrate with TimescaleDB. If you have a regular Postgresql database, please refer to origin/main.

## Introduction
Simple Cache is a service for minimalistic temporary (or not?) file hosting.
SC allows you to upload any file via `curl`, `httpie`,
or any other HTTP client, and store it for a specified number of days (or forever).

## Routes
- **GET** `/` - Returns a 403 status code.
- **POST** `/` - Uploads a file from multipart/form-data (Field: `file` or `url`).
- **GET** `/:id` - Loads a file preview in the browser if possible. If not, download the file.

## Configuration
The configuration can be done either via the file or the environment variables. The priority for loading is as follows:
```text
┌──────────────────────┐     ┌─────────────────┐     ┌────────────────────────────┐
│ LOAD DEFAULT VALUES  │ ──> │ LOAD CONFIG.HCL │ ──> │ LOAD ENVIRONMENT VARIABLES │
└──────────────────────┘     └─────────────────┘     └────────────────────────────┘
```

If some keys exist in multiple places (e.g., server.port in `config.hcl` and SC_SERVER_PORT is given), then the latest load is prioritized. (Environment variables are superior.)

### config.hcl
Example `config.hcl`:
```hcl
server {
  serverName = "Simple Cache v1.2.4"                                         # Sets the `Server` header for each response.
  appName = "Simple Cache"                                                   # Sets the app name to display in the terminal.
  port = 8080                                                                # Server port
  databaseUrl = "postgres://user:password@localhost:5432/timescaledb"        # TimescaleDB Connection URL
  seed = 3719                                                                # Seed for ID generation [type: random number > 0]
}

logger {
  forceColors = false                                                        # Force colored output
  fullTimestamp = true
}

limits {
  maxSize = 10                                                               # Max uploaded body size in MB
  blockTor = false                                                           # Block TOR exit nodes
  ipBlacklist = []                                                           # ["168.0.2.3", "1.2.3.4"...]
}
```

### Environment Variables
Configuration can also be done through environment variables. See the table below.

| Environment             | Default Value         | Description                                |
|-------------------------|-----------------------|--------------------------------------------|
| SC_SERVER_SERVERNAME    | "Simple Cache v1.2.4" | `Server` header for each response          |
| SC_SERVER_APPNAME       | "Simple Cache"        | App name to display in the terminal        |
| SC_SERVER_PORT          | 8080                  | Server port                                |
| SC_SERVER_DATABASEURL   | [placeholder]         | TimescaleDB Connection URL                 |
| SC_SERVER_SEED          | 3719                  | Seed for ID generation [random number > 0] |
| SC_LOGGER_FORCECOLORS   | false                 | Force colored output                       |
| SC_LOGGER_FULLTIMESTAMP | true                  | -                                          |
| SC_LIMITS_MAXSIZE       | 10                    | Max uploaded body size in MB               |
| SC_LIMITS_BLOCKTOR      | false                 | Block TOR exit nodes                       |
