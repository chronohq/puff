# Puff

Puff is a command-line tool for quickly generating random values in
different formats. At some point, software developers find themselves
needing random values for a variety of purposes. Puff aims to make this
process as intuitive as possible.

## Quickstart

### Print a random hexadecimal string (Default: 16 bytes).

```bash
$ puff hex
2cc84ba90e5277f6733aa71386a4de3b
```

### Print a random hexadecimal string with a custom byte length.

```bash
$ puff hex --bytes 32
8c5955a659c59d4414072b45bac872964a8a8077ffbd0f0083ffad47e5b33c66
```

### Print four random hexadecimal strings.

```bash
$ puff hex -n 4
f906e2fa87fbf6e9f0b0b44e2fc81993
5d7347c2fe7fda44097604c06ae4f25f
55ea03cb94847b52b6911d4b517bcc01
2b6ce246d8eeb1b546ed56e8ae927437
```

### Generate random bytes and write it to file.

Similar to using `dd` for generating test data, you can use `puff` to create a binary file with random bytes:

```bash
$ puff binary --bytes 10485760 -o /tmp/puff-10mb.bin
$ stat -c %s /tmp/puff-10mb.bin
10485760
```

## Contributing

Contributions of any kind are welcome.
If you're submitting a PR, please follow [Go's commit message structure](https://go.dev/wiki/CommitMessage).

## License

Puff is released under the [MIT license](https://opensource.org/license/MIT).
See the [LICENSE](LICENSE) file for details.