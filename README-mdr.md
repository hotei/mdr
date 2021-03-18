<center>
mdr
===
</center>

<h3>   <a href="http://godoc.org/github.com/hotei/mdr">
<img src="https://godoc.org/github.com/hotei/mdr?status.png" alt="mdr" /><br>
<p>
</a><a href="http://travis-ci.org/hotei/mdr">
<img src="https://secure.travis-ci.org/hotei/mdr.png" alt="Build Status" /></a>
Travis build status.
</h3>


License details are at the end of this document. 
This document is (c) 2012-2015 David Rook.

Comments can be sent to <hotei1352@gmail.com> .

This is a package of "utility" code I wrote.  I use it frequently.  If you have
received one of my other packages from github (ansiterm, bits, etc ...) you may have
gotten __mdr__ as a dependency.  It pulls in a fairly large range of standard lib
packages so if you only need a smallish set of things it might make sense to just
copy them individually or possibly make a package subset from __mdr__.  If you find it
useful - or find a bug - please send an email.  

## Installation

If you have a working go installation on a Unix-like OS:

> ```go get github.com/hotei/mdr```

### Features

Functions/Methods are grouped (loosely) in the following categories

* Hash Helpers for CRC64, MD5, SHA256 etc
  * return string digest of a buffer
  * return string digest of a file
* Random number generators
  * Uniform
  * Uniform numbers between two values (float64 and int)
  * Poisson
  * FlipCoin
  * Normal and Normal Z (unit normal)
  * Crypto-Quality random []byte
* archives
  * whether the file is a "collection" of files (like tar or zip) or just "itself"
  * determining compression types given file name
* Formatting helpers
  * add commas to integer numbers
  * print dates as human friendly times (1 day instead of 86400 seconds)
  * test dates for "legality" and leapyears
  * progress displays - spinner or progressbar (if endpoint known)
  * StarDate
  * validating decimal or hex strings
* Getting obscure properties of files from the os 
  * User and Group numbers
  * Number of links 
* Integer and float64 helpers
  * test if value is in range between two values
  * convert an int to and from []byte of various sizes
  * compute max and min of a []int
  * compute max and min points given []Point  (Bounds)
  * creating Bezier interpolations (quadratic)
* Polar <--> Rectangular coordinate conversion
* Parallel computation helpers
  * JobSplit figures out how to split a job into pieces for each CPU to work on
* Directory operations
  * maps of paths to SHA256 of contents

If you want details you can see them at this link:
<a href="http://godoc.org/github.com/hotei/mdr">
<img src="https://godoc.org/github.com/hotei/mdr?status.png" alt="mdr" />
</a>

### Change Log
* 2020-02-01 built with go 1.15.7 ok
* 2020-02-24 built with go 1.13.7 ok
* 2015-08-18 builds with 1.5
  * added Travis (but not updating it of late)
* 2015-05-01 updated progress bar functions
  * builds with 1.4.2
* 2013-04-10 updated docs, posted at github.com/hotei/MDR.git
* additions
* 2010-04-20 started, working

### Resources

* [go language reference] [1] 
* [go standard library package docs] [2]
* [Source for mdr package on github] [3]

[1]: http://golang.org/ref/spec/ "go reference spec"
[2]: http://golang.org/pkg/ "go package docs"
[3]: http://github.com/hotei/mdr "github.com/hotei/mdr"


License
-------
The 'mdr' go package is distributed under the Simplified BSD License:

> Copyright (c) 2010-2021 David Rook. All rights reserved.
> 
> Redistribution and use in source and binary forms, with or without modification, are
> permitted provided that the following conditions are met:
> 
>    1. Redistributions of source code must retain the above copyright notice, this list of
>       conditions and the following disclaimer.
> 
>    2. Redistributions in binary form must reproduce the above copyright notice, this list
>       of conditions and the following disclaimer in the documentation and/or other materials
>       provided with the distribution.
> 
> THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDER ``AS IS'' AND ANY EXPRESS OR IMPLIED
> WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
> FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> OR
> CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
> CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
> SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
> ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
> NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
> ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// EOF README-mdr-pkg.md  (this is a markdown document and tested OK with blackfriday)
