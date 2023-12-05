This is a tool that generally does three things:

* Detect if a UTF-8 BOM (byte order mark) is at the start of file.

* If a UTF-8 BOM is detected, it is removed.

* If no UTF-8 BOM is found, this will add one to the start of the file.

A BOM (byte order mark) is a Unicode character that can be added to a text file to specify the byte order or encoding. If it's not added to a UTF-8—encoded file, many softwares won't be able to decode or display the file's content correctly. However, if it's added to a file encoded UTF-8, many softwares won't be able to decode or display the character correctly.

Therefore, the best way to work with the byte order mark isn't to totally accept or prevent it, but to use it when needed and prevent it if possible (or to use it if possible and prevent it if necessary). That's why I made this tool.

### Requirements

+ Go 1.16+

### Build & install

1. Clone the repository.
```sh
git clone https://github.com/sls1005/bom
```

2. `cd` to the `bom` directory.
```sh
cd bom
```

3. Build the executable.
```sh
go build
```
Or install:
```sh
go build
```

### Usage

```sh
./bom FILE
```

See `./bom --help` for more information.

### Note

+ This can only be compiled with Go compiler v1.16 or later. Some changes are required to make it to work with an older version of compiler.

+ This does not check if the file is UTF-8—encoded. Please check it with another tool.

+ This cannot be used to detect other UTFs' byte order marks (such as those for UTF-16 or UTF-32), nor can it be used to add a BOM for an encoding other than UTF-8.

+ Although this cannot be used to detect the encoding, a file with a UTF-8 BOM detected by this program is likely to be UTF-8—encoded.

### References

+ <https://en.m.wikipedia.org/wiki/Byte_order_mark>
