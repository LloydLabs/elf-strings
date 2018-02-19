# elf-strings
elf-strings will programmatically read an ELF binary's string sections within a given binary. This is meant to be much like the `strings` UNIX utility, however is purpose built for ELF binaries. 

This means that you can get suitable information about the strings within the binary, such as the section they reside in, the offset in the section, etc.. This utility also has the functionality to 'demangle' C++ symbols, iterate linked libraries and print basic information about the ELF.

This can prove extremely useful for quickly grabbing strings when analysing a binary.

# Output
![alt text](https://i.imgur.com/plIdQCF.png "example of demangled strings")

# Building
```
git clone https://github.com/LloydLabs/elf-strings
cd elf-strings
go build
```

# Arguments
```
  -binary string
        the path to the ELF you wish to parse
  -demangle
        demangle C++ symbols into their original source identifiers, prettify found C++ symbols (optional)
  -hex
        output the strings as a hexadecimal literal (optional)
  -libs
        show the linked libraries in the binary (optional)
  -max uint
        the maximum amount of strings that you wish to be output (optional)
  -offset
        show the offset of the string in the section (default, recommended) (default true)
  -output-file string
        the path of the output file that you want to output to (optional)
  -output-format string
        the format you want to output as (optional, plain/json/xml) (default "plain")
```

# Example

An example grabbing the strings from the `echo` utility.

```
./elf-strings --binary=/bin/echo --min=4 --max-count=10

[+] Size: 31 kB
[+] Arch: x86_64
[+] Entry point: 0x401800
[+] Class: ELFCLASS64
[+] Byte order: LittleEndian

[.dynstr+0x0]: libc.so.6
[.dynstr+0xa]: fflush
[.dynstr+0x11]: __printf_chk
[.dynstr+0x1e]: setlocale
[.dynstr+0x28]: mbrtowc
[.dynstr+0x30]: strncmp
[.dynstr+0x38]: strrchr
[.dynstr+0x40]: dcgettext
[.dynstr+0x4a]: error
[.dynstr+0x50]: __stack_chk_fail
```
