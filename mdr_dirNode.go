// DirNode.go (c) 2014 David Rook - all rights reserved

package mdr

import (
	// standard lib go 1.2 pkgs
	"crypto/sha256"
	"fmt"
	"hash"
	"log"
	"os"
	"sort"
)

const (
	paranoid      = false
	G_datetimeFmt = "2006-01-02:15_04_05"
)

type FileRec struct {
	R Rec256
	//	fname   string // needed to sort on
	Touched bool
}

func (fr FileRec) Dump() {
	fr.R.Dump()
	fmt.Printf("touched = %v\n", fr.Touched)
}

type DirNode struct {
	Pathname       string
	IsSortedByName bool
	Files          []FileRec
	Size           int64
	SHA256         string
}

type ByName []FileRec

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].R.Name < a[j].R.Name }

// Finalize computes the total size and SHA256 of the leaf
func (d *DirNode) Finalize() (string, int64) {
	d.Size = 0
	d.SortFilesByName()
	var h hash.Hash = sha256.New()
	for i := 0; i < len(d.Files); i++ {
		d.Size += d.Files[i].R.Size
		h.Write([]byte(d.Files[i].R.SHA))
	}
	digest := h.Sum(nil)
	rv := ""
	for _, hex := range digest {
		rv = rv + fmt.Sprintf("%02x", hex)
	}
	d.SHA256 = rv
	return rv, d.Size
}

func (d *DirNode) UnTouchAll() {
	for i := 0; i < len(d.Files); i++ {
		d.Files[i].Touched = false
	}
}

func (d *DirNode) AddFile(f FileRec) {
	//Verbose.Printf("Adding file %s to %d others\n", f.r.Name, len(d.files))
	d.Files = append(d.Files, f)
	flen := len(d.Files)
	if flen == 1 {
		d.IsSortedByName = true
		return
	}
	// remains sorted if addition is > last element
	if d.IsSortedByName && (d.Files[flen-2].R.Name < d.Files[flen-1].R.Name) {
		d.IsSortedByName = true
	} else {
		d.IsSortedByName = false
	}
}

func (d *DirNode) SortFilesByName() {
	//Verbose.Printf("SortFilesByName in %s\n", d.pathname)
	if len(d.Files) <= 1 {
		d.IsSortedByName = true
		return
	}
	// magic happens here
	//fmt.Printf("Before sort\n")
	//d.Dump()
	sort.Sort(ByName(d.Files))
	//if sort.IsSorted(ByName(d.files)) == false {
	//	log.Fatalf("houston we have a problem\n")
	//}
	//fmt.Printf("After sort\n")
	d.IsSortedByName = true
	//d.Dump()
}

func (d *DirNode) DeleteFile(fname string) {
	//Verbose.Printf("Deleting %s\n", fname)
	log.Fatalf("here")
	// it is an error to delete a file that's not in list
	ndx := d.IndexOf(fname)
	if ndx < 0 {
		d.Dump()
		log.Fatalf("!Err---> cant find file %s to delete it\n", fname)
	}
	// otherwise we found it somewhere
	if len(d.Files) == 1 { // special case if it's the only file
		d.Files = []FileRec{}
		d.IsSortedByName = false // calling zero length unsorted
		return
	}
	if ndx == 0 { // first one is a special case, doesn't change sorted status
		d.Files = d.Files[1:]
		return
	}
	// if its the last one thats a special case, doesn't change sorted status
	lastNdx := len(d.Files) - 1
	if ndx == lastNdx {
		d.Files = d.Files[:lastNdx-1] // truncate the list
		return
	}
	// if possible concat slices and preserve sorted status (expensive but cheaper than a resort)
	// otherwise swap current and last item and truncate the list
	if d.IsSortedByName {
		d.Files = append(d.Files[:ndx], d.Files[ndx+1:]...)
	} else {
		d.Files[ndx] = d.Files[lastNdx]
		d.Files = d.Files[:lastNdx-1]
	}
}

// brute force works
func (d *DirNode) IndexOfBrute(fname string) int {
	for ndx, f := range d.Files {
		//if ndx > 0 {
		//	if d.files[ndx].r.Name < d.files[ndx-1].r.Name {
		//		fmt.Printf("+++ WARNING +++ file not sorted\n")
		//	}
		//}
		if f.R.Name == fname {
			return ndx
		}
	}
	return -1
}

func (d *DirNode) IsSortedProperly() bool {
	flen := len(d.Files)
	if flen <= 1 {
		return true
	}
	for i := 1; i < len(d.Files)-1; i++ {
		if d.Files[i].R.Name < d.Files[i-1].R.Name {
			fmt.Printf("out of order %s %s\n",
				d.Files[i].R.Name,
				d.Files[i-1].R.Name)
			log.Panicf("crashed here")
		}
	}
	return true
}

// special case speedup if we check last returned value +1 first?
// better handles case where we're stepping through list ???
// Returns -1 if not found
func (d *DirNode) IndexOf(fname string) int {
	//return d.IndexOfBrute(fname)
	if d.IsSortedByName == false {
		d.SortFilesByName()
	}
	// use binary search here
	length := len(d.Files)
	if length == 0 {
		return -1
	}
	if length == 1 {
		if fname == d.Files[0].R.Name {
			return 0
		} else {
			return -1
		}
	}
	lo := 0
	hi := length - 1
	var loopCt int
	foundAt := -1
	for {
		/*
			if paranoid {
				if d.files[lo].r.Name > d.files[hi].r.Name {
					d.Dump()
					fmt.Printf("g_dirMap has %d items\n", len(g_dirMap))
					fmt.Printf("refreshed[%d] errs[%d] ok[%d] added[%d]\n",
						g_refreshCount, g_errCount, g_okCount, g_addCount)
					fmt.Printf("sort called %d times\n", g_sortCt)
					fmt.Printf("lo[%d]%s > hi[%d]%s\n",
						lo, d.files[lo].r.Name, hi, d.files[hi].r.Name)
					fmt.Printf("lo value > hi value so\n")
					log.Panicf("something is broken in DirNode's binary search")
				}
			}
		*/
		if fname < d.Files[lo].R.Name || fname > d.Files[hi].R.Name {
			// not found
			return -1
		}

		mid := (hi + lo) / 2
		//Verbose.Printf("lo[%d] hi[%d]  mid = %d loopct=%d\n", lo, hi, mid, loopCt)
		if fname == d.Files[lo].R.Name {
			foundAt = lo
			break
		}
		if fname == d.Files[hi].R.Name {
			foundAt = hi
			break
		}
		if fname == d.Files[mid].R.Name {
			foundAt = mid
			break
		}

		if fname < d.Files[lo].R.Name { // not found
			return -1
		}
		if fname > d.Files[hi].R.Name { // not found
			return -1
		}

		if (hi - lo) < 2 {
			break
		}
		if true {
			loopCt++
			if loopCt > 64 { // implies > 2^64 elements in list
				fmt.Printf("lo[%d] hi[%d]\n", lo, hi)
				log.Panicf("something is broken in DirNode's binary search")
			}
		}
		if fname < d.Files[mid].R.Name {
			hi = mid
			continue
		}
		if fname > d.Files[mid].R.Name {
			lo = mid
			continue
		}
	}
	return foundAt
}

// calls AddFile if file not currently in DirNode
// name should be just the base part of path
func (d *DirNode) Update(fname string) error {
	// BUG(mdr): ? be paranoid and check if we have any slashes in fname?

	var fr FileRec
	fullPath := d.Pathname + fname
	info, err := os.Stat(fullPath)
	if err != nil {
		fmt.Printf("!Err---> %v\n", err)
		return err
	}
	fr.R.Size = info.Size()
	fr.R.Date = info.ModTime().Format(G_datetimeFmt)
	tmpsha, err := FileSHA256(fullPath)
	if err != nil {
		fmt.Printf("!Err---> %v\n", err)
		return err
	}
	fr.R.SHA = tmpsha
	fr.R.Name = fname
	ndx := d.IndexOf(fname)
	if ndx >= 0 {
		d.Files[ndx].R = fr.R
	} else {
		d.AddFile(fr)
	}
	//Verbose.Printf("Updated %s\n", fullPath)
	return nil
}

func (d *DirNode) Dump() {
	var sortStatus string
	if d.IsSortedByName {
		sortStatus = "is"
	} else {
		sortStatus = "is NOT"
	}
	fmt.Printf("Dir path = %s %s sorted\n", d.Pathname, sortStatus)
	//	if d.isSortedByName == false {
	//		d.SortFilesByName()
	//	}
	for i := 0; i < len(d.Files)-1; i++ {
		fmt.Printf("%d %s\n", i, d.Files[i].R.Name)
		d.Files[i].Dump()
	}
}
