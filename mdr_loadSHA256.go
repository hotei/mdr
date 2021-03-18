// loadSHA256.go (c) 2014 David Rook - all rights reserved

package mdr

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	//"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Split256(line string) (Rec256, error) {
	var rec Rec256
	var err error
	x := strings.Split(line, "|")
	if len(x) != 4 {
		Verbose.Printf("!ERR--->  found other than 4 parts during split of %q\n", line)
		return rec, CantCreateRec
	}
	rec.Size, err = strconv.ParseInt(strings.Trim(x[0], " \n\t\r"), 10, 64)
	if err != nil {
		Verbose.Printf("!ERR--->  Cant parse to int64: %s\n", x[0])
		return rec, CantCreateRec
	}
	rec.SHA = strings.Trim(x[1], " \n\t\r")
	if !ValidHexString(rec.SHA) {
		Verbose.Printf("!ERR--->  %q is not valid hex \n", rec.SHA)
		return rec, CantCreateRec
	}
	rec.Date = strings.Trim(x[2], " \n\t\r")
	rec.Name = strings.Trim(x[3], " \n\t\r")
	if strings.ContainsAny(rec.Name, "\n\t\r") {
		fmt.Printf("%s name contains unusual character\n", rec.Name)
		return rec, CantCreateRec
	}
	return rec, nil
}

// LoadSHA256asList returns lines of file as []string , nil on success
//  each line is known to split without error
// comments are stripped
func LoadSHA256asList(fname string) ([]string, error) {
	fileBytes, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("error reading from %s: %v\n", fname, err)
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
			return nil, fmt.Errorf("split failed on line[%d] = %q", ndx, line)
		}
		list256 = append(list256, string(line))
	}
	Verbose.Printf("added %d precalculated digests from %s\n", len(list256), fname)
	return list256, nil
}

// LoadSHA256asMap returns map[SHA256]path , nil on success
func LoadSHA256asMap(fname string) (map[string]string, error) {
	fileBytes, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("error reading from %s: %v\n", fname, err)
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
			return nil, fmt.Errorf("Can't load line[%d] %q from %s\n", ndx, sline, fname)
		}
		hashMap[partRec.SHA] = sline
	}
	Verbose.Printf("added %d precalculated digests from %s\n", len(hashMap), fname)
	return hashMap, nil
}

// LoadSHA256asDirMap returns *dirMap[directoryPath]*DirNode
//  uses scanner to allow large file.256 to be used as input
// Used when its necessary to compute something with directory as input
func LoadSHA256asDirMap(fname string) (map[string]*DirNode, error) {
	var dirMap = make(map[string]*DirNode, 10000)
	fmt.Printf("Loading SHA256 list from %s may take a while...\n", fname)
	input, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("can't open file %s\n", fname)
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
			return nil, fmt.Errorf("error reading %s\n", fname)
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
			return nil, fmt.Errorf("Can't load line[%d] %q from %s\n", lineCt, sline, fname)
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
	Verbose.Printf("added %d directory nodes with %d files from %s\n", len(dirMap), fileCt, fname)
	return dirMap, nil
}

// LoadSHA256Names returns the array of pathnames
func LoadSHA256Names(fname string) ([]string, error) {
	Verbose.Printf("Loading from %s\n", fname)
	listOfNames := make([]string, 0, 100000)
	input, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("can't open file %s\n", fname)
	}
	scanner := bufio.NewScanner(input)
	lineCt := 0
	fileCt := 0
	for scanner.Scan() {
		sline := scanner.Text()
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("LoadSHA256Names(): error reading from %s\n", fname)
		}
		lineCt++
		//Verbose.Printf("line[%d] = %q\n", line)
		if len(sline) < 3 { // must have at least 3 pipe symbols
			continue
		}
		if sline[0] == '#' {
			continue
		}
		r, err := Split256(string(sline))
		if err != nil {
			return nil, fmt.Errorf("split failed\n")
		}
		listOfNames = append(listOfNames, string(r.Name))
		fileCt++
	}
	Verbose.Printf("added %d precalculated digests from %s\n", len(listOfNames), fname)
	return listOfNames, nil
}
