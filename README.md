## gh0.st

## Current instance: [https://files.ieljit.lol/](https://files.ieljit.lol/)

### DESCRIPTION

This is a simple file sharing server. This is the Go version of the original [0x0.st](https://0x0.st) server. This project also includes a few more features, such as configurable blocklists, TOR exit nodes blocking, native gzip compression, native SSL support, and a few more. For more information, see the [original project](https://git.0x0.st/mia/0x0). This project is licensed under the CC0 1.0 Universal license. See the [LICENSE](/LICENSE) file for more information.

## USAGE

1. Download the latest release from the [releases page](https://github.com/voxelin/gh0.st/releases/latest).
2. Run `chmod +x ghost` to make the binary executable.
3. Run `./ghost` to start the server on `localhost:3000`. You can also specify a configuration file with the `-c` flag.

### CONFIGURATION

Configuration is done through a JSON file. The default location is `config.jsonc` in the current directory, but you can specify a different location with the `-c` flag. The configuration file is structured as follows:

```jsonc
{
    "address": "127.0.0.1", // Use 0.0.0.0 to listen on all interfaces. Default: 127.0.0.1(localhost)
    "port": 3000, // Port to listen on. Default: 3000
    "size_limit": 100, // Maximum size of a file in MB. Default: 100
    "db_path": "files.sqlite", // Path to the database file. Default: db.json
    "blocklist_path": "blocklist.txt", // Path to the blocklist file. Default: blocklist.json
    "index_page_path": "index.html", // Path to the index page. Default: index.html
    "block_tor": true, // Block TOR users. Default: true
    "fake_ssl": false, // Use this if you are using a reverse proxy with SSL enabled. There is no need to specify cert and key files. Default: false
    "ssl_": false, // Use this to use real SSL on this executable. Default: false
    "ssl_cert": "cert.pem", // Path to the SSL certificate. Default: cert.pem
    "ssl_key": "key.pem", // Path to the SSL key. Default: key.pem
    "gzip_": true // Use gzip compression. Default: true
}
```

### LICENSE

```
Creative Commons Legal Code
CC0 1.0 Universal
```

Thanks to [joaoofreitas](https://github.com/joaoofreitas) for idea. Made and advanced by 🇺🇦 [voxelin](https://github.com/voxelin)
