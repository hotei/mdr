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
	ArchiveJarType // not quite Zip but similar
	ArchiveCpioType
	ArchiveMailboxType
)

type dispatch struct {
	pattern string
	arch    ArchiveType
}

type patternType struct {
	pattern string
	count   int64
}

var (
	dispatcher         = []dispatch{}
	g_strictZ          = false // .z will return as if .Z unless this is set true
	collectionPatterns = []patternType{}
)

// examine name to determine archive method used
// should match magic numbers but that's a different function
func WhichArchiveType(s string) ArchiveType {
	ext := strings.ToLower(filepath.Ext(s))
	// ----------------------------------------- ArchiveZipType
	if ext == ".zip" {
		return ArchiveZipType
	}

	// ----------------------------------------- ArchiveJarType
	if ext == ".jar" {
		return ArchiveJarType
	}

	// ----------------------------------------- ArchiveTarType
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
	if ext == ".tbz2" {
		return ArchiveTarType
	}

	// ----------------------------------------- ArchiveCpioType
	if ext == ".cpio" {
		return ArchiveCpioType
	}

	// ----------------------------------------- ArchiveMailboxType
	if ext == ".mbx" {
		return ArchiveMailboxType
	}
	if len(dispatcher) == 0 {
		Verbose.Printf("Building dispatcher\n")
		// compressed tar collection
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.gz\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.Z\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.z\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.bz\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.bz2\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.bzip2\\E$", ArchiveTarType})
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

const ( // note: add new types at bottom only!
	CompressNoMatchType CompressType = iota
	CompressZipType
	CompressGzipType
	CompressBz2Type
	CompressZcompressType // .Z not common except possibly in Japan
	CompresszPackType     // .z deprecated - very rare in last 20 years
	CompressBz1Type       // .bz deprecated - very rare in last 20+ years
	CompressShrinkType    // found in zip files of old
	CompressCompactType   // ?? deprecated
	CompressFreezeType    // ?? deprecated
	Compress7zType        // new kid on the block
	CompressXvType        // even newer kid on the block
)

func init() {
	Verbose.Printf("mdr.archives.go init() entry\n")
	defer Verbose.Printf("mdr.archives.go init() exit\n")

	// note that filenames are converted to lowercase before matching takes place
	// paterns should have longest match possiblity listed first below  (tar.z before .z)
	///////////////////////////////////////////////////////////////////////////
	// N E E E D E D

	// fmt.Printf("5\n")
	// file.log
	//	pat := ".*\\Q.\\Elog$" // just a placeholder
	//	collectionPatterns = append(collectionPatterns, patternType{ pat,0})

	/////////////////////////////////////////////////////////////////////////////
	// W O R K - I N - P R O G R E S S
	//  file.z -  recognized by dispatcher
	//	pat = ".*\\Q.\\E\\Qz\\E$" // Unix "pack" or compress
	//	collectionPatterns = append(collectionPatterns, patternType{ pat,0})

	//  file.tar.z - recognized by dispatcher
	pat := ".*\\Q.tar.z\\E$" // Unix "pack" or compress
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})

	// file.jar - recognized by dispatcher
	pat = ".*\\Q.\\Ejar$" // java jar is somewhat similar to zip
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})
	////////////////////////////////////////////////////////////////////////////
	// D O N E

	// file.shar - needs more testing
	pat = ".*\\Q.\\Eshar$" // Shell Archive not compressed
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})

	// file.shar.Z(z) - needs more testing
	pat = ".*\\Q.\\Eshar\\Q.\\Ez$" // Shell Archive not compressed
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})

	// file.zip - testOK
	pat = ".*\\Q.zip\\E$" // PK-Zip usually compressed
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})

	// file.tar - testOK
	pat = ".*\\Q.tar\\E$" // uncompressed tar
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})

	// file.tgz - testOK
	pat = ".*\\Q.\\Etgz$" // tar compressed with gz
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})

	// file.taz -
	pat = ".*\\Q.\\Etaz$" // tar compressed with gz
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})

	// file.tar.gz - testOK
	pat = ".*\\Q.tar.gz\\E$" // tar compressed with gz
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})

	// file.tbz - testOK
	pat = ".*\\Q.\\Etbz$" // tar compressed with bz2
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})

	// file.tar.bz
	pat = ".*\\Q.tar.bz\\E$" // tar compressed with bz2
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})

	// file.bz2
	pat = ".*\\Q.bz2\\E$" // file compressed with bz2
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})

	// file.bzip2
	pat = ".*\\Q.\\Ebzip2$" // file compressed with bz2
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})

	// file.tar.bz2
	pat = ".*\\Q.tar.bz2\\E$" // tar compressed with bz2
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})

	// file.tbz2
	pat = ".*\\Q.\\Etbz2$" // tar compressed with bz2
	collectionPatterns = append(collectionPatterns, patternType{pat, 0})

}

// examine name to determine compression method used
// should match magic numbers but that's a different function
func WhichCompressType(s string) CompressType {
	origExt := filepath.Ext(s)

	// most .z files we have found are simply .Z that have been mislabeled, thus strictZ defaults false
	// the caller can also test magic header bytes in file body to see which is correct
	//   if the CompressZCompressType is returned and decompress fails
	if origExt == ".z" {
		if g_strictZ {
			return CompresszPackType
		}
		return CompressZcompressType
	}
	if origExt == ".Z" {
		return CompressZcompressType
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

func FileNameIsCollection(fname string) (bool, string) {
	//fmt.Printf("testing %s\n",fname)
	nameBytes := []byte(strings.ToLower(fname))
	maxNdx := len(collectionPatterns)
	for i := 0; i < maxNdx; i++ {
		pat := collectionPatterns[i]
		//		fun := dispatcher[i].Function
		isMatch, err := regexp.Match(pat.pattern, nameBytes)
		if err != nil {
			fmt.Printf("!Err-> ?pattern error in re2 %v\n", pat)
			log.Panicf(fmt.Sprintf("re2 pattern error %v", pat))
		}
		if isMatch {
			collectionPatterns[i].count++
			return true, pat.pattern
		} // only need the first match
	}
	// zipLog(fmt.Sprintf("!Err-> %s isn't a collection\n", fname))
	return false, "<no-match>"
}
