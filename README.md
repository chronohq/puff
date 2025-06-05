# Puff

[![go workflow](https://github.com/chronohq/puff/actions/workflows/go.yml/badge.svg)](https://github.com/chronohq/puff/actions/workflows/go.yml)
[![mit license](https://img.shields.io/badge/license-MIT-green)](/LICENSE)

Puff is a command-line tool designed to quickly generate random values
in various formats such as hexadecimal strings, UUIDs, and binary blobs.
Whether for testing, cryptographic purposes, or data seeding, Puff makes
the process straightforward.

## Installation

**macOS**

```bash
brew install chronohq/tap/puff
```

**Linux**

1. Download the [latest binary release](https://github.com/chronohq/puff/releases/latest) for your system
2. Extract the binary into `/usr/local/bin`, for example:

```bash
rm /usr/local/bin/puff && tar -C /usr/local/bin xvzf puff-1.2.3.linux-amd64.tar.gz
```

**Windows**

1. Download the [latest binary release](https://github.com/chronohq/puff/releases/latest) for your system
2. Extract the binary using the built-in Windows file explorer or a tool like 7-Zip
3. Move `puff.exe` to a directory in your PATH
4. Open Windows Terminal and run `puff.exe` from any directory

## Quickstart

### Print Hexadecimal Values (Default: 16 Bytes)

Use `puff hex --help` to view available command options.

```bash
# Print a single hex-encoded value
puff hex
2cc84ba90e5277f6733aa71386a4de3b

# Print two hex-encoded values
puff hex -n 2
f906e2fa87fbf6e9f0b0b44e2fc81993
5d7347c2fe7fda44097604c06ae4f25f

# Print a hex-encoded value with custom byte length (32 bytes)
puff hex --bytes 32
8c5955a659c59d4414072b45bac872964a8a8077ffbd0f0083ffad47e5b33c66

# Print hex-encoded values separated by a comma
puff hex -n 2 --delimiter ","
d4bc48da024a728fee985a6257e88611,63100951b7ff67de3f7e9c1d0b98101d

# Print hex-encoded values with a custom suffix
puff hex -n 2 --suffix ".png"
2259c58f9a774fe171720349cd715bed.png
314fb96973ccf36ef579fcaa9303ddcc.png
```

### Print UUIDs as Hexadecimal Values

By default, Puff generates version 4 UUIDs, but version 7 is also supported.
Use `puff uuid --help` to view available command options.

```bash
# Print two version 4 UUIDs (default)
puff uuid -n 2
da8456b6-81b6-4440-b875-0613d7183e4d
b7aa0977-0acb-4983-8c4c-36265a8dd60b

# Print two compact version 4 UUIDs
puff uuid -n 2 --compact
3cfe31da3d444d7baf0903066c864a41
d97ab07f0f564643be5d47ddea6a95af

# Print a compact version 7 UUID (time-ordered)
puff uuid --version 7 --compact
019741e3fd837da48a946baaa14db4ce

# Print UUIDs with custom suffix
puff uuid -n 2 --compact --suffix ".png"
8659b17b768a4ec5b48d1bac08279750.png
805fc8e5b2aa4de0a7f37deecbf9a606.png
```

### Print Base64 Values (Default: 16 Bytes)

Use `puff base64 --help` to view available command options.

```bash
# Print a single base64-encoded value
puff base64
4YamGT7lPlJakROuu7zN5w==

# Print a base64-encoded value with custom byte length (32 bytes)
puff base64 --bytes 32
YPjRDxaUYFkurRAhz0MUzzl8Hh3Y0Z79DZcrJX5R/4g=

# Print a base64-encoded value using URL-safe encoding
puff base64 --url-safe
7KtUKGFAvVHoDdsQDIuRtQ

# Print a URL-safe base64-encoded with custom suffix
puff base64 --url-safe --suffix ".png"
BvOqihvrdg2kmtPQaxKx6A.png
```

### Create a Binary File with Random Bytes

Similar to using `dd` for generating test data, you can use `puff` to create a binary file with random bytes:

```bash
# Create a 10MB binary file with random bytes
puff binary --bytes 10485760 -o /tmp/puff-10mb.bin

# Verify the file size
stat -c %s /tmp/puff-10mb.bin
10485760
```

## Contributing

Contributions of any kind are welcome.
If you're submitting a PR, please follow [Go's commit message structure](https://go.dev/wiki/CommitMessage).

## License

Puff is available under the [MIT license](https://opensource.org/license/MIT).
See the [LICENSE](LICENSE) file for details.
