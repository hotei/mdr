// loadSHA256.go (c) 2014 David Rook - all rights reserved 

package mdr

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// LoadSHA256asList returns lines of file as []string , nil on success
//  each line is known to split without error
// comments are stripped
func LoadSHA256asList(fname string) ([]string, error) {
	fileBytes, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("error reading from %s: %v\n", fname, err)
	}
	list256 := make([]string, 0, 1000)
	//Verbose.Printf("Loading from %s\n", fname)

	lines := bytes.Split(fileBytes, []byte{'\n'})

	//Verbose.Printf("found %d lines in %s\n", len(lines), fname)
	// fmt.Printf("%v\n",lines)
	for ndx, line := range lines {
		_ = ndx
		//Verbose.Printf("line[%d] = %q\n", ndx, line)
		if len(line) < 3 { // must have at least 3 pipe symbols
			continue
		}
		if line[0] == '#' {
			continue
		}
		_, err := Split256(string(line))
		if err != nil {
			log.Fatalf("split failed\n")
		}
		list256 = append(list256, string(line))
	}
	fmt.Printf("added %d precalculated digests from %s\n", len(list256), fname)
	return list256, nil
}

// LoadSHA256asMap returns map[SHA256]path , nil on success
func LoadSHA256asMap(fname string) (map[string]string, error) {
	fileBytes, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("error reading from %s: %v\n", fname, err)
	}
	hashMap := make(map[string]string, 1000)
	//Verbose.Printf("Loading from %s\n", fname)

	lines := bytes.Split(fileBytes, []byte{'\n'})

	//Verbose.Printf("found %d lines in %s\n", len(lines), fname)
	// fmt.Printf("%v\n",lines)
	for ndx, line := range lines {
		sline := string(line)
		_ = ndx
		//Verbose.Printf("line[%d] = %q\n", ndx, sline)
		if len(sline) < 3 { // must have at least 3 pipe symbols
			continue
		}
		if sline[0] == '#' {
			continue
		}
		partRec, err := Split256(sline)
		if err != nil {
			log.Fatalf("Can't load line[%d] %q from %s\n", ndx, sline, fname)
		}
		hashMap[partRec.SHA] = sline
	}
	fmt.Printf("added %d precalculated digests from %s\n", len(hashMap), fname)
	return hashMap, nil
}

// LoadSHA256asDirMap returns *dirMap[directoryPath]*DirNode
//  uses scanner to allow large file.256 to be used as input 
// Used when its necessary to compute something with directory as input
func LoadSHA256asDirMap(fname string) (map[string]*DirNode,error) {
	var dirMap = make(map[string]*DirNode, 10000)
	fmt.Printf("Loading SHA256 list from %s may take a while...\n", fname)
	input, err := os.Open(fname)
	if err != nil {
		log.Fatalf("can't open file %s\n", fname)
	}
	scanner := bufio.NewScanner(input)
	lineCt := 0
	fileCt := 0
	for scanner.Scan() {
		if (lineCt % 100) == 0 {
			Spinner()
		}
		sline := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatalf("error reading %s\n", fname)
		}
		lineCt++
		//Verbose.Printf("line[%d] = %q\n", lineCt, sline)
		if len(sline) < 3 { // must have at least 3 pipe symbols
			continue
		}
		if sline[0] == '#' {
			continue
		}
		sr, err := Split256(sline)
		var fr FileRec = FileRec{R: sr, Touched: false}
		if err != nil {
			log.Fatalf("Can't load line[%d] %q from %s\n", lineCt, sline, fname)
		}
		dir := filepath.Dir(sr.Name)
		if dir != "/" {
			if dir[len(dir)-1] != '/' {
				dir = dir + "/"
			}
		}
		fr.R.Name = filepath.Base(sr.Name)
		_, ok := dirMap[dir]
		if ok { // add file to existing directory node
			dirMap[dir].AddFile(fr)
		} else { // add both directory node and file
			var dn DirNode
			dn.Pathname = dir
			dn.Files = []FileRec{fr}
			dirMap[dir] = &dn
		}
		fileCt++
	}
	fmt.Printf("added %d directory nodes with %d files from %s\n", len(dirMap), fileCt, fname)
	return dirMap, nil
}