## gh0.st

## Current instance: [https://ghost.blackvoxel.space/](https://ghost.blackvoxel.space/)

### DESCRIPTION

This is a powerful file-sharing server built in Go that improves upon the original 0x0.st server. It comes with a range of features, such as configurable blacklists, blocking of TOR exit nodes, native gzip compression, and native SSL support. All of these features are included under the CC0 1.0 Universal license, which can be found in the LICENSE file.

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
host: 0.0.0.0 # or 127.0.0.1 (localhost)
port: 8080 # use 443 for SSL
size_limit: 10 # in MB
db_path: files.db # path to the database file
blacklist_path: blacklist.txt # path to the blacklist file
index_path: index.html # path to the index file
block_tor: true # block TOR exit nodes
fake_ssl: false # fake SSL
enable_ssl: false # real SSL (requires ssl_cert and ssl_key)
# ssl_cert: cert.pem # path to the SSL certificate
# ssl_key: key.pem # path to the SSL key
enable_gzip: true # enable gzip compression for files (recommended)
# trusted_platform: "" # trusted platform to take IP from. When using other, specify the Real IP Header (e.g. X-CDN-IP)
# allowed_ips: [] # comma-separated list of allowed reverse proxy IPs (recommended by GIN Documentation)
```

### ENVIRONMENT

| Variable                 | Description                                                                                                                                | Default       |
| ------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------ | ------------- |
| `GHOST_HOST`             | Host to listen on                                                                                                                          | 0.0.0.0       |
| `GHOST_PORT`             | Port to listen on                                                                                                                          | 8080          |
| `GHOST_SIZE_LIMIT`       | Maximum file size in MB                                                                                                                    | 10            |
| `GHOST_DB_PATH`          | Path to the database file                                                                                                                  | files.db      |
| `GHOST_BLACKLIST_PATH`   | Path to the blacklist file                                                                                                                 | blacklist.txt |
| `GHOST_INDEX_PATH`       | Path to the index file                                                                                                                     | index.html    |
| `GHOST_BLOCK_TOR`        | Block TOR exit nodes                                                                                                                       | true          |
| `GHOST_FAKE_SSL`         | Fake SSL                                                                                                                                   | false         |
| `GHOST_ENABLE_SSL`       | Real SSL                                                                                                                                   | false         |
| `GHOST_SSL_CERT`         | Path to the SSL certificate                                                                                                                | `nil`         |
| `GHOST_SSL_KEY`          | Path to the SSL key                                                                                                                        | `nil`         |
| `GHOST_ENABLE_GZIP`      | Enable gzip compression for files                                                                                                          | true          |
| `GHOST_TRUSTED_PLATFORM` | Trusted platform to take IP from. When using other, specify the Real IP Header (e.g. `X-CDN-IP`) [cloudflare \| google \| other (specify)] | `nil`         |
| `GHOST_ALLOWED_IPS`      | Comma-separated list of allowed reverse proxy IPs (Recommended by GIN Documentation) [e.g. `1.2.3.4,2.2.2.2`]                              | `nil`         |
| `GIN_MODE`               | Set to `release` to disable debug mode                                                                                                     | `debug`       |

### LICENSE

```
Creative Commons Legal Code
CC0 1.0 Universal
```

Thanks to [joaoofreitas](https://github.com/joaoofreitas) for the great idea, which was further developed by 🇺🇦 [voxelin](https://github.com/voxelin).
