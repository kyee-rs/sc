## gh0.st

### DESCRIPTION

This service is an updated version of the original `0xg0.st`, which in turn was a fork of `0x0.st` (lmao). The original `0x0.st` was written in Python and used a MySQL database. This version is written in Go and uses a SQLite database. The original `0xg0.st` was also a bit of a mess, so I decided to rewrite it from scratch. This version does not provide storage/* folder to store files, it uses GORM with SQLite driver to save the content of the files in the database. This version also provides a file size limit of 10MB. The original `0x0.st` had a file size limit of 100MB.
### LICENSE

```
Creative Commons Legal Code
CC0 1.0 Universal
```

Check LICENSE file for more information about this software license.
