## gh0.st

## Current instance: [https://ghost.blackvoxel.space/](https://ghost.blackvoxel.space/)

### DESCRIPTION

This is a powerful file-sharing server built in Go that improves upon the original 0x0.st server. It comes with a range of features, such as configurable blocklists, blocking of TOR exit nodes, native gzip compression, and native SSL support. All of these features are included under the CC0 1.0 Universal license, which can be found in the LICENSE file.

## USAGE

1. Visit the [releases page](https://github.com/voxelin/gh0.st/releases/latest) and download the latest release.
2. Make the binary executable by running `chmod +x ghost`.
3. Start the server on `localhost:8080` by running `./ghost`.

### CONFIGURATION

Configuration is done through a YAML (`config.yml`) file in some of this directories:

1. /etc/ghost/config/config.yml
2. ./config/config.yml
3. ./config.yml

Alternatively you can specify configuration flags through the environment variables from the list: [Env](#Environment). The structure of the configuration file is outlined below:

```yaml
host: 0.0.0.0 # or localhost (127.0.0.1)
port: 8080
size_limit: 10 # in MB
db_path: files.db
blocklist_path: blocklist.txt
index_path: index.html
block_tor: true # Block TOR exit nodes.
fake_ssl: true # Fake SSL. Specify this if you're using a reverse proxy with SSL. No need to specify ssl_cert and ssl_key.
enable_ssl: false # Real SSL. Specify this if you're using ghost standalone with SSL. Specify ssl_cert and ssl_key.
ssl_cert: cert.pem
ssl_key: key.pem
enable_gzip: true # Enable gzip compression for files.
```

### ENVIRONMENT

| Variable               | Description                       | Default       |
| ---------------------- | --------------------------------- | ------------- |
| `GHOST_HOST`           | Host to listen on                 | 0.0.0.0       |
| `GHOST_PORT`           | Port to listen on                 | 8080          |
| `GHOST_SIZE_LIMIT`     | Maximum file size in MB           | 10            |
| `GHOST_DB_PATH`        | Path to the database file         | files.db      |
| `GHOST_BLOCKLIST_PATH` | Path to the blocklist file        | blocklist.txt |
| `GHOST_INDEX_PATH`     | Path to the index file            | index.html    |
| `GHOST_BLOCK_TOR`      | Block TOR exit nodes              | true          |
| `GHOST_FAKE_SSL`       | Fake SSL                          | true          |
| `GHOST_ENABLE_SSL`     | Real SSL                          | false         |
| `GHOST_SSL_CERT`       | Path to the SSL certificate       | cert.pem      |
| `GHOST_SSL_KEY`        | Path to the SSL key               | key.pem       |
| `GHOST_ENABLE_GZIP`    | Enable gzip compression for files | true          |

### LICENSE

```
Creative Commons Legal Code
CC0 1.0 Universal
```

Thanks to [joaoofreitas](https://github.com/joaoofreitas) for the great idea, which was further developed by 🇺🇦 [voxelin](https://github.com/voxelin).
