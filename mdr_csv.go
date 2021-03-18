// mdr_csv.go

// not sure that os.Exit(-1) is the answer to bad formats.
// add error return and pass the buck?

package mdr

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func fieldScrub(s string) string {
	s = strings.Replace(s, ",", "", -1)
	s = strings.Replace(s, "\"", "", -1)
	return strings.Trim(s, `/r/t/n "`)
}

// ExtractFloat32 assumes convention of period for decimal
func ExtractFloat32(s string) float32 {
	rv, err := strconv.ParseFloat(fieldScrub(s), 32)
	if err != nil {
		fmt.Printf("Cant extract Float32 from |%s|\n", s)
		os.Exit(-1)
	}
	return float32(rv)
}

// ExtractFloat64 assumes convention of period for decimal
func ExtractFloat64(s string) float64 {
	rv, err := strconv.ParseFloat(fieldScrub(s), 64)
	if err != nil {
		log.Printf("Cant extract Float64 from |%s|\n", s)
		os.Exit(-1)
	}
	return rv
}

func ExtractUint8(s string) uint8 {
	rv, err := strconv.ParseUint(fieldScrub(s), 10, 8)
	if err != nil {
		fmt.Printf("Cant extract Uint8 from |%s|\n", s)
		os.Exit(-1)
	}
	return uint8(rv)
}

func ExtractUint16(s string) uint16 {
	rv, err := strconv.ParseUint(fieldScrub(s), 10, 16)
	if err != nil {
		fmt.Printf("Cant extract Uint16 from |%s|\n", s)
		os.Exit(-1)
	}
	return uint16(rv)
}

func ExtractUint32(s string) uint32 {
	rv, err := strconv.ParseUint(fieldScrub(s), 10, 32)
	if err != nil {
		fmt.Printf("Cant extract Uint32 from |%s|\n", s)
		os.Exit(-1)
	}
	return uint32(rv)
}

func ExtractUint64(s string) uint64 {
	rv, err := strconv.ParseUint(fieldScrub(s), 10, 64)
	if err != nil {
		fmt.Printf("Cant extract Uint64 from |%s|\n", s)
		os.Exit(-1)
	}
	return uint64(rv)
}

func ExtractInt8(s string) int8 {
	rv, err := strconv.ParseInt(fieldScrub(s), 10, 8)
	if err != nil {
		fmt.Printf("Cant extract Int8 from |%s|\n", s)
		os.Exit(-1)
	}
	return int8(rv)
}

func ExtractInt16(s string) int16 {
	rv, err := strconv.ParseInt(fieldScrub(s), 10, 16)
	if err != nil {
		fmt.Printf("Cant extract Int16 from |%s|\n", s)
		os.Exit(-1)
	}
	return int16(rv)
}

func ExtractInt32(s string) int32 {
	rv, err := strconv.ParseInt(fieldScrub(s), 10, 32)
	if err != nil {
		fmt.Printf("Cant extract Int32 from |%s|\n", s)
		os.Exit(-1)
	}
	return int32(rv)
}

func ExtractInt64(s string) int64 {
	rv, err := strconv.ParseInt(fieldScrub(s), 10, 64)
	if err != nil {
		fmt.Printf("Cant extract Int64 from |%s|\n", s)
		os.Exit(-1)
	}
	return int64(rv)
}
