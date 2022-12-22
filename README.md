## gh0.st

## Current instance: [https://gofile.pp.ua/](https://gofile.pp.ua)

### DESCRIPTION

This is a simple file sharing server. It is designed to be used as a command line tool, and it is not meant to be used as a web server. This is the Go version of the original [0x0.st](https://0x0.st) server. This project also includes a few more features, such as configurable blocklists, TOR exit nodes blocking, native gzip compression, native SSL support, and a few more. For more information, see the [original project](https://git.0x0.st/mia/0x0). This project is licensed under the CC0 1.0 Universal license. See the [LICENSE](/LICENSE) file for more information.

### CONFIGURATION

Configuration is done through a JSON file. The default location is `config.json` in the current directory, but you can specify a different location with the `-c` flag. The configuration file is structured as follows:

```json
{
    "size_limit": 10,
    "db_path": "./ghost.files.sqlite",
    "blocklist_path": "./ghost.blocklist.txt",
    "index_page_path": "./ghost.index.html",
    "ssl_": false,
    "ssl_cert": "./cert.pem",
    "ssl_key": "./key.pem",
    "gzip_": true
}
```

### LICENSE

```
Creative Commons Legal Code
CC0 1.0 Universal
```

Thanks to [joaoofreitas](https://github.com/joaoofreitas) for idea. Made and advanced by 🇺🇦 [voxelin](https://github.com/voxelin)
