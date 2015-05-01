<center>
MDR
===
</center>

License details are at the end of this document. 
This document is (c) 2012-2013 David Rook.

Comments can be sent to <hotei1352@gmail.com> .

This is a package of "utility" code I wrote.  I use it frequently.  If you have
received one of my other packages from github (ansiterm, bits, etc ...) you may have
gotten MDR as a dependency.  It pulls in a fairly large range of standard lib
packages so if you only need a smallish set of things it might make sense to just
copy them individually or possibly make a package subset from MDR.  If you find it
useful - or find a bug - please send an email.  

```
A note on style.  I prefer not to use the single line form :

	if x,err := foo(); err != nil {
	// stuff
	}

instead you'll see:
	x,err := foo()
	if err != nil {
	// stuff
	}
```

* Note that the CRC64 table is Public if you need to replace it without changing
the library.
* I use fatal errors rather than panics in most places.  This is a habit of mine and may
not always lead to the fastest debugging.  I just prefer the end user see a smaller
understandable (I hope) message and not the gut-spilling verbosity of panic.  I'm
considering a user-setable switch for this behavior perhaps tied to the Verbose var.

Journal
-------
* 2013-04-10 updated docs, posted at github.com/hotei/MDR.git
* additions
* 2010-04-20 started, working

License
-------
The 'MDR' go package is distributed under the Simplified BSD License:

> Copyright (c) 2010-2013 David Rook. All rights reserved.
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

// EOF README-MDR.md  (this is a markdown document and tested OK with blackfriday)