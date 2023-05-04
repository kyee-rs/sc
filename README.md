# Ghost - Zero bullshit file hosting

Ghost is a simple and lightweight file sharing service written in Go, designed to make file sharing quick and easy without the need for any complicated setup. With Ghost, you can share files with anyone, anywhere in the world, without having to worry about file size limitations, annoying ads, or any other kind of bullshit.

## Installation
### Docker
Ghost is available as a Docker image on [GHCR.io](https://ghcr.io/voxelin/ghost). To run Ghost using Docker, simply run the following command:

```bash
$ docker run -d -p 8080:8080 ghcr.io/voxelin/ghost:latest
```

### Binary

Ghost is also available as a standalone binary for Linux, macOS, and Windows. To install Ghost using the binary, download the latest release from the [releases page](https://github.com/voxelin/ghost/releases) and open the binary file.

### Source

To install Ghost from source, you will need to have Go 1.16 or higher installed on your system. To install Ghost from source, run the following commands:

```bash
$ git clone https://github.com/voxelin/ghost
$ cd ghost
$ go build -o ghost ./internal
$ ./ghost
```

## Usage

Ghost is designed to be as simple and easy to use as possible. It needs no additional configuration by default, but if you want to extend the [default values](#environment-variables), please, use the [configuration file](#configuration-file) or [environment variables](#environment-variables) specified below.

### Environment variables
| Variable | Description | Default |
| --- | --- | --- |
| `GS_PORT` | The port to run Ghost on | `8080` |
| `GS_HOST` | The host to run Ghost on | `127.0.0.1` |
| `GS_DB_PATH` | The path to the database file | `./ghost.db` |
| `GS_BLOCK_TOR` | Whether or not to block Tor users | `false` |
| `GS_GZIP` | Whether or not to compress files using Gzip | `true` |
| `GS_AUTO_CLEANUP` | The number of days to keep files for | `7` |
| `GS_MAX_SIZE` | The maximum file size in megabytes | `100` |
| `GS_LANGUAGE` | The language to use [available: en, uk] | `en` |

### Configuration file
Ghost also supports configuration files. To use a configuration file, create a file named `cfg.json` in the same directory as the Ghost binary and add the following:

```json
{
    "host": "127.0.0.1",
    "port": 8080,
    "db_path": "./ghost.db",
    "block_tor": false,
    "gzip": true,
    "auto_cleanup": 7,
    "max_size": 100,
    "language": "en"
}
```
Also, there are some additional paths that Ghosts looks for configuration files in:
- `./cfg.json`
- `./cfg/cfg.json`
- `./ghost/cfg.json`
- `/etc/ghost/cfg.json`

## Contributing

Ghost is completely open-source, and contributions are always welcome! If you're interested in contributing to Ghost, check out the official GitHub repository and feel free to submit pull requests or open issues.

## License

Ghost is released under the CC0 1.0 Universal License. See the [`LICENSE`](LICENSE) file for more information.
