package main

import (
	"debug/elf"
	"fmt"

	"github.com/ianlancetaylor/demangle"
)

// UtilConvHex converts a string to a C-like hexadecimal string literal
func UtilConvHex(buf string) string {
	var hexBuf string

	for i := 0; i < len(buf); i++ {
		hexBuf += fmt.Sprintf("\\x%X", buf[i])
	}

	return hexBuf
}

// UtilDemangle will demangle a symbol by string, this is
// simply just a friendly wrapped around the demangle package
func UtilDemangle(symbol *string) (string, error) {
	x, err := demangle.ToString(*symbol)
	if err != nil {
		return "", err
	}

	return x, nil
}

// UtilUniqueSlice will remove all duplicate instances
// from the offset slice.
func UtilUniqueSlice(s []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(s))
	j := 0

	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}

		seen[v] = struct{}{}
		s[j] = v

		j++
	}

	return s[:j]
}

// UtilIsNice will validate to make sure that the string in question
// is of 'human readable' format
func UtilIsNice(str string) bool {
	length := len(str)

	spaces := 0
	for i := 0; i < length; i++ {
		if str[i] == ' ' {
			spaces++
		}

		if str[i] < ' ' && !(str[i] == '\r' || str[i] == '\n') {
			return false
		}
	}

	if spaces == length {
		return false
	}

	return true
}

// UtilConvertMachine will convert from elf.Machine type to a string.
func UtilConvertMachine(mach elf.Machine) string {
	conv := map[elf.Machine]string{
		elf.EM_386:         "Intel 80386",
		elf.EM_68HC12:      "Motorola M68HC12",
		elf.EM_68K:         "Motorola 68000",
		elf.EM_860:         "Intel 80860",
		elf.EM_88K:         "Motorola 88000",
		elf.EM_960:         "Intel 80960",
		elf.EM_ALPHA:       "Digital Alpha",
		elf.EM_ARC:         "Argonaut RISC Core, Argonaut Technologies Inc.",
		elf.EM_ARM:         "Advanced RISC Machines ARM",
		elf.EM_COLDFIRE:    "Motorola ColdFire",
		elf.EM_FR20:        "Fujitsu FR20",
		elf.EM_H8_300:      "Hitachi H8/300",
		elf.EM_H8_300H:     "Hitachi H8/300H",
		elf.EM_H8S:         "Hitachi H8S",
		elf.EM_IA_64:       "Intel IA-64 processor architecture",
		elf.EM_M32:         "AT&T WE 32100",
		elf.EM_ME16:        "Toyota ME16 processor",
		elf.EM_MIPS:        "MIPS I Architecture",
		elf.EM_MIPS_RS3_LE: "MIPS RS3000 Little-endian",
		elf.EM_MIPS_X:      "Stanford MIPS-X",
		elf.EM_MMA:         "Fujitsu MMA Multimedia Accelerator",
		elf.EM_NCPU:        "Sony nCPU embedded RISC processor",
		elf.EM_NDR1:        "Denso NDR1 microprocessor",
		elf.EM_NONE:        "No machine",
		elf.EM_PARISC:      "Hewlett-Packard PA-RISC",
		elf.EM_PCP:         "Siemens PCP",
		elf.EM_PPC:         "PowerPC",
		elf.EM_PPC64:       "64-bit PowerPC",
		elf.EM_RCE:         "Motorola RCE",
		elf.EM_RH32:        "TRW RH-32",
		elf.EM_S370:        "IBM System/370 Processor",
		elf.EM_SH:          "Hitachi SH",
		elf.EM_SPARC:       "SPARC",
		elf.EM_SPARC32PLUS: "Enhanced instruction set SPARC",
		elf.EM_SPARCV9:     "SPARC Version 9",
		elf.EM_ST100:       "STMicroelectronics ST100 processor",
		elf.EM_STARCORE:    "Motorola Star*Core processor",
		elf.EM_TINYJ:       "Advanced Logic Corp. TinyJ embedded processor family",
		elf.EM_TRICORE:     "Siemens Tricore embedded processor",
		elf.EM_V800:        "NEC V800",
		elf.EM_VPP500:      "Fujitsu VPP500",
		elf.EM_X86_64:      "x86_64"}

	if arch, ok := conv[mach]; ok {
		return arch
	}

	return conv[elf.EM_NONE]
}
