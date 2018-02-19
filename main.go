package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/shawnsmithdev/zermelo"
)

var (
	demangleOpt = flag.Bool("demangle", false, "demangle C++ symbols into their original source identifiers, prettify found C++ symbols (optional)")
	hexOpt      = flag.Bool("hex", false, "output the strings as a hexadecimal literal (optional)")
	offsetOpt   = flag.Bool("offset", true, "show the offset of the string in the section (default, recommended)")
	binaryOpt   = flag.String("binary", "", "the path to the ELF you wish to parse")
	formatOpt   = flag.String("output-format", "plain", "the format you want to output as (optional, plain/json/xml)")
	outputOpt   = flag.String("output-file", "", "the path of the output file that you want to output to (optional)")
	maxOpt      = flag.Uint64("max-count", 0, "the maximum amount of strings that you wish to be output (optional)")
	libOpt      = flag.Bool("libs", false, "show the linked libraries in the binary (optional)")
	infoOpt     = flag.Bool("no-info", false, "don't show any information about the binary")
	minOpt      = flag.Uint64("min", 0, "the minimum length of the string")
	colorOpt    = flag.Bool("no-color", false, "disable color output in the results")
	trimOpt     = flag.Bool("no-trim", false, "disable triming whitespace and trailing newlines")
	humanOpt    = flag.Bool("no-human", false, "don't validate that its a human readable string, this could increase the amount of junk.")
)

// ReadSection is the main logic here
// it combines all of the modules, etc.
func ReadSection(reader *ElfReader, section string) {
	var err error
	var writer *OutWriter
	var count uint64

	sect := reader.ReaderParseSection(section)

	if *outputOpt != "" {
		writer, err = NewOutWriter(*outputOpt, OutParseTypeStr(*formatOpt))
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if sect != nil {
		nodes := reader.ReaderParseStrings(sect)

		// Since maps in Go are unsorted, we're going to have to make
		// a slice of keys, then iterate over this and just use the index
		// from the map.
		keys := make([]uint64, len(nodes))
		for k, _ := range nodes {
			keys = append(keys, k)
		}

		err = zermelo.Sort(keys)
		if err != nil {
			return
		}

		keys = UtilUniqueSlice(keys)

		for _, off := range keys {
			if *maxOpt != 0 {
				if count == *maxOpt {
					break
				}
			}

			str := string(nodes[off])
			if uint64(len(str)) < *minOpt {
				continue
			}

			if !*humanOpt {
				if !UtilIsNice(str) {
					continue
				}
			}

			str = strings.TrimSpace(str)

			if !*trimOpt {
				bad := []string{"\n", "\r"}
				for _, char := range bad {
					str = strings.Replace(str, char, "", -1)
				}
			}

			if *demangleOpt {
				demangled, err := UtilDemangle(&str)
				if err == nil {
					str = demangled
				}
			}

			if *hexOpt {
				str = UtilConvHex(str)
			}

			if *offsetOpt {
				if os.Getenv("NO_COLOR") != "" || *colorOpt {
					fmt.Printf("[%s+%#x]: %s\n",
						section,
						off,
						str)
				} else {
					fmt.Printf("[%s%s]: %s\n",
						color.BlueString(section),
						color.GreenString("+%#x", off),
						str)
				}
			} else {
				fmt.Println(str)
			}

			if writer != nil {
				writer.WriteResult(str, off)
			}

			count++
		}
	}
}

// ReadBasic will read the basic information
// about the ELF
func ReadBasic(reader *ElfReader) {
	stat, err := reader.File.Stat()
	if err != nil {
		return
	}

	size := humanize.Bytes(uint64(stat.Size()))

	fmt.Printf(
		"[+] Size: %s\n"+
			"[+] Arch: %s\n"+
			"[+] Entry point: %#x\n"+
			"[+] Class: %s\n"+
			"[+] Byte order: %s\n",
		size,
		UtilConvertMachine(reader.ExecReader.Machine),
		reader.ExecReader.Entry,
		reader.ExecReader.Class.String(),
		reader.ExecReader.ByteOrder.String(),
	)

	if *libOpt {
		fmt.Println("[+] Libraries:")
		libs, err := reader.ExecReader.ImportedLibraries()
		if err == nil {
			for _, lib := range libs {
				fmt.Printf("\t [!] %s\n", lib)
			}
		}
	}

	fmt.Println(strings.Repeat("-", 16))
}

// main is the entrypoint for this program
func main() {
	flag.Parse()

	if *binaryOpt == "" {
		flag.PrintDefaults()
		return
	}

	r, err := NewELFReader(*binaryOpt)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer r.Close()

	ReadBasic(r)

	sections := []string{".dynstr", ".rodata", ".rdata",
		".strtab", ".comment", ".note",
		".stab", ".stabstr", ".note.ABI-tag", ".note.gnu.build-id"}

	for _, section := range sections {
		ReadSection(r, section)
	}
}
