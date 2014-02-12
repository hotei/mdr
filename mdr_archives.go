// mdr_archives.go (c) 2013,2014 David Rook - see LICENSE.md

package mdr

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

type ArchiveType int

const (
	ArchiveNoMatchType ArchiveType = iota
	ArchiveTarType
	ArchiveZipType
	ArchiveArkType
	ArchiveCpioType
)

type dispatch struct {
	pattern string
	arch    ArchiveType
}

var (
	dispatcher = []dispatch{}
	g_strictZ  = false
)

// examine name to determine archive method used
// should match magic numbers but that's a different function
func WhichArchiveType(s string) ArchiveType {
	ext := strings.ToLower(filepath.Ext(s))
	// ArchiveZipType -----------------------------------------
	if ext == ".zip" {
		return ArchiveZipType
	}
	// ArchiveTarType -----------------------------------------
	if ext == ".tar" {
		return ArchiveTarType
	}
	if ext == ".tgz" {
		return ArchiveTarType
	}
	if ext == ".taz" {
		return ArchiveTarType
	}
	if ext == ".tbz" {
		return ArchiveTarType
	}
	// ArchiveCpioType -----------------------------------------
	if ext == ".cpio" {
		return ArchiveCpioType
	}
	// ArchiveArkType -----------------------------------------
	if ext == ".ark" {
		return ArchiveArkType
	}

	if len(dispatcher) == 0 {
		if Verbose {
			fmt.Printf("Building dispatcher\n")
		}
		// compressed tar collection
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.gz\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.Z\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.z\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.bz\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.bz2\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.bzip2\\E$", ArchiveTarType})
		// compressed ark collection
		dispatcher = append(dispatcher, dispatch{".*\\Q.ark.Z\\E$", ArchiveArkType})
	}

	nameBytes := []byte(s)
	for _, d := range dispatcher {
		isMatch, err := regexp.Match(d.pattern, nameBytes)
		if err != nil {
			fmt.Printf("!Err-> ?pattern error in %s\n", d.pattern)
			log.Panicf("better fix it now")
		}
		if isMatch {
			return d.arch
		} // only need the first match
	}

	// NO MATCH -----------------------------------------
	return ArchiveNoMatchType
}

type CompressType int

const (
	CompressNoMatchType CompressType = iota
	CompressZipType
	CompressGzipType
	CompressBz2Type
	CompressZcompressType // .Z not common except possibly in Japan
	CompressPackType      // .z deprecated - very rare in last 20 years
	CompressBz1Type       // .bz deprecated - very rare in last 20+ years
)

// examine name to determine compression method used
// should match magic numbers but that's a different function
func WhichCompressType(s string) CompressType {
	origExt := filepath.Ext(s)
	// note - we don't do the next step because a lot of the .z files we have
	// found are simply .Z that have been mislabeled, thus strictZ defaults false
	// if you have the opposite, then set the global true instead
	// the caller  should test magic header bytes in file body to see which is correct
	if g_strictZ {
		if origExt == ".z" {
			return CompressPackType
		}
	}
	// unless we have a good reason not to we fold case for the following tests
	ext := strings.ToLower(origExt)
	// CompressZipType ---------------------------------
	if ext == ".zip" {
		return CompressZipType
	}
	// CompressGzipType ---------------------------------
	if ext == ".gz" {
		return CompressGzipType
	}
	if ext == ".gzip" {
		return CompressGzipType
	}
	if ext == ".tgz" {
		return CompressGzipType
	}
	// CompressZcompressType ---------------------------------
	if ext == ".z" { // .Z converted to .z by tolower - see above - usually not a bug
		return CompressZcompressType
	}
	if ext == ".taz" {
		return CompressZcompressType
	}
	// CompressBz2Type ---------------------------------
	if ext == ".bz" {
		return CompressBz2Type
	}
	if ext == ".tbz" {
		return CompressBz2Type
	}
	if ext == ".bz2" {
		return CompressBz2Type
	}
	if ext == ".tbz2" {
		return CompressBz2Type
	}
	if ext == ".bzip2" {
		return CompressBz2Type
	}
	// NO MATCH -----------------------------------------
	return CompressNoMatchType
}
