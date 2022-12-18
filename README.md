## gh0.st

### DESCRIPTION

This service is an updated version of the original `0xg0.st`, which in turn was a fork of `0x0.st` (lmao). The original `0x0.st` was written in Python and used a MySQL database. This version is written in Go and uses a SQLite database. The original `0xg0.st` was also a bit of a mess, so I decided to rewrite it from scratch. This version does not provide storage/\* folder to store files, it uses GORM with SQLite driver to save the content of the files in the database. This version also provides a file size limit control, gzip compression support, and a few other features like SSL support.

### CONFIGURATION

The configuration file is located in `config.go`. The file is self-explanatory, but here's a quick overview:

```go
var size_limit int64 = 10 * 1024 * 1024 // 10MB file size limit
var db_path string = "./db/files.sqlite" // SQLite database path
var ssl_ bool = false // SSL support
var ssl_cert string = "./cert.pem" // SSL certificate path
var ssl_key string = "./key.pem" // SSL key path
var gzip_ bool = false // gzip compression support
```

### LICENSE

```
Creative Commons Legal Code
CC0 1.0 Universal
```

This software was made by [joaoofreitas](https://github.com/joaoofreitas) and advanced by 🇺🇦 [voxelin](https://github.com/voxelin)
