# Ghost - Zero bullshit file hosting

Ghost is a simple and lightweight file sharing service written in Go, designed to make file sharing quick and easy without the need for any complicated setup. With Ghost, you can share files with anyone, anywhere in the world, without having to worry about file size limitations, annoying ads, or any other kind of bullshit.

## Features

- Zero configuration: Ghost is designed to be as easy to use as possible, with no complicated setup required.
- Fast and lightweight: Ghost is built with Go, making it fast and efficient, even on low-end hardware.
- Simple and intuitive: Ghost is designed to be as simple and intuitive as possible, with no complicated menus or settings.
- Serviceless: Ghost is completely serviceless, meaning that you don't need to sign up for an account or pay any fees to use it. No external databases or services are required.
- Open-source: Ghost is completely free and open-source, with no hidden fees or restrictions.

## Getting started

To get started with Ghost, simply download the latest release from the official GitHub repository, extract the files to a directory on your computer, and run the `ghost` executable. Ghost will automatically start listening for incoming file uploads on port 8080.

Once Ghost is running, you can upload files to it by using a `curl` command like the following:
```bash
$ curl -F "file=@yourFile.txt" http://localhost:8080
```

## Configuration

Ghost requires no configuration by default, but you can customize its behavior by using an environment variables. The following environment variables are supported:

- `GHOST_HOST`: The host that Ghost should listen on. Defaults to `0.0.0.0`.
- `GHOST_PORT`: The port that Ghost should listen on. Defaults to `8080`.
- `GHOST_DB_PATH`: The path to the database file that Ghost should use. Defaults to `ghost.db`.
- `GHOST_BLOCK_TOR`: Whether or not Ghost should block Tor users. Defaults to `false`.
- `GHOST_GZIP`: Whether or not Ghost should gzip compress responses. Defaults to `true`.
- `GHOST_AUTOCLEANUP`: Interval (in days) at which Ghost should automatically delete expired files. Defaults to `0` (disabled).
- `GHOST_MAXSIZE`: Maximum file size (in megabytes) that Ghost should allow. Defaults to `0` (unlimited).

## Contributing

Ghost is completely open-source, and contributions are always welcome! If you're interested in contributing to Ghost, check out the official GitHub repository and feel free to submit pull requests or open issues.

## License

Ghost is released under the CC0 1.0 Universal license. See the [`LICENSE`](LICENSE) file for more information.
