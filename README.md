# Puff

![go workflow](https://github.com/chronohq/puff/actions/workflows/go.yml/badge.svg)

Puff is a command-line tool for quickly generating random values in
different formats. At some point, software developers find themselves
needing random values for a variety of purposes. Puff aims to make this
process as intuitive as possible.

## Installation

**macOS**

```bash
brew install chronohq/tap/puff
```

**Linux**

1. Download the [latest binary release](https://github.com/chronohq/puff/releases/latest) for your system
2. Extract the binary into `/usr/local/bin`. For example:

```bash
rm /usr/local/bin/puff && tar -C /usr/local/bin xvzf puff-1.2.3.linux-amd64.tar.gz
```

## Quickstart

**Print Random Hexadecimal Strings (Default: 16 Bytes)**

Use `puff hex --help` to explore command options.

```bash
# Print one value
puff hex
2cc84ba90e5277f6733aa71386a4de3b

## Print two values
puff hex -n 2
f906e2fa87fbf6e9f0b0b44e2fc81993
5d7347c2fe7fda44097604c06ae4f25f

## Print a value with custom byte length
puff hex --bytes 32
8c5955a659c59d4414072b45bac872964a8a8077ffbd0f0083ffad47e5b33c66
```

**Print UUID (Version 7) as Hexadecimal Strings**

Use `puff uuid --help` to explore command options.

```bash
# Print two UUIDs in standard dashed format
puff uuid -n 2
019191ac-e84b-7e58-8ca9-d44d68ecd15f
019191ac-e84b-7e92-b277-07c67a7db551

# Print two UUIDs in compact dash-less format
puff uuid -n 2 --compact
019191ace35c71bfb517a4c1269716b8
019191ace35c71fa8d14ddf9e56c1faf
```

**Generate Random Bytes and Write It to a File**

Similar to using `dd` for generating test data, you can use `puff` to create a binary file with random bytes:

```bash
puff binary --bytes 10485760 -o /tmp/puff-10mb.bin
stat -c %s /tmp/puff-10mb.bin
10485760
```

## Contributing

Contributions of any kind are welcome.
If you're submitting a PR, please follow [Go's commit message structure](https://go.dev/wiki/CommitMessage).

## License

Puff is released under the [MIT license](https://opensource.org/license/MIT).
See the [LICENSE](LICENSE) file for details.
