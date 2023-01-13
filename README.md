## gh0.st
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fvoxelin%2Fgh0.st.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fvoxelin%2Fgh0.st?ref=badge_shield)


## Current instance: [https://files.ieljit.lol/](https://files.ieljit.lol/)

### DESCRIPTION

This is a powerful file-sharing server built in Go that improves upon the original 0x0.st server. It comes with a range of features, such as configurable blocklists, blocking of TOR exit nodes, native gzip compression, and native SSL support. All of these features are included under the CC0 1.0 Universal license, which can be found in the LICENSE file.

## USAGE

1. Visit the [releases page](https://github.com/voxelin/gh0.st/releases/latest) and download the latest release.
2. Make the binary executable by running `chmod +x ghost`.
3. Start the server on `localhost:3000` by running `./ghost`. Alternatively, you can specify a configuration file with the `-c` flag.

### CONFIGURATION

Configuration is done through a JSON file. By default, this file is located at `config.jsonc` in the current directory, however, you can specify a different location with the `-c` flag. The structure of the configuration file is outlined below:

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

Thanks to [joaoofreitas](https://github.com/joaoofreitas) for the great idea, which was further developed by 🇺🇦 [voxelin](https://github.com/voxelin).

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fvoxelin%2Fgh0.st.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fvoxelin%2Fgh0.st?ref=badge_large)
