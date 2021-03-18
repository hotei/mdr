# Makefile for mdr
# note that merc (primary) target has extra checks that need not be duplicated for other targets
# status is working

PROJECT = mdr
CMDNAME = mdr
VERSION = 0.1.5

TARDIR = 	$(HOME)/Desktop/TarPit/
DATETIME = 	`date "+%Y-%m-%d.%H_%M_%S"`
DOCOUT = README-$(CMDNAME)-godoc.md
TARGETS = yoda # thor loki
SRC = *.go	# individual files or *.go if that works ok

all: $(PROJECT)_linux

#	GOOS=linux   GOARCH=amd64 go build -v -o $(PROG)_linux *.go
#	GOOS=freebsd GOARCH=amd64 go build -v -o $(PROG)_freebsd *.go

$(PROJECT)_linux: $(SRC)
	godatetime -pkg="$(PROJECT)" > compileDate.go
	echo "var CompilerVersion = \"" `go version`\" >> compileDate.go
	GOOS=linux GOARCH=amd64 go build -v -o $(PROJECT)_linux $(SRC)
	gofmt -w *.go
	go vet . 2> vetProbs.txt
# more vet could be useful here - shadow syntax changed to ...

$(PROJECT)_freebsd: $(SRC)
	godatetime -pkg="$(PROJECT)" > compileDate.go
	echo "var CompilerVersion = \"" `go version`\" >> compileDate.go
	GOOS=freebsd GOARCH=amd64 go build -v -o $(PROJECT)_freebsd *.go
	gofmt -w *.go
	go vet . 2> vetProbs.txt

#install: not relevant for pkg but use this for commands
#	cp $(PROJECT)_linux $(HOME)/bin/$(CMDNAME)

yoda: $(PROJECT)_linux
#	cp $(PROJECT)_linux $(HOME)/bin/$(CMDNAME)
	touch yoda

loki: $(PROJECT)_linux
#	scp -p ./$(PROJECT)_linux mdr@loki:bin/$(CMDNAME)
	touch loki

thor: $(PROJECT)_freebsd
#	scp -p ./$(PROJECT)_freebsd mdr@thor:bin/$(CMDNAME)
	touch thor

wolf: $(PROJECT)_freebsd
#	scp -p ./$(PROJECT)_freebsd mdr@wolf:bin/$(CMDNAME)
	touch wolf

dead:
	deadcode > problems.dead

# note that godepgraph can be used to derive .travis.yml install: section
docs:
	gofmt -w *.go
	echo "\n" > $(DOCOUT)
	deadcode -md >> $(DOCOUT)
	echo "\n" >> $(DOCOUT)
#	godoc2md . >> $(DOCOUT)
#	godepgraph -md . >> $(DOCOUT) // buggy if source not in ~/src/~
# 	godepgraph -md -v . >> $(DOCOUT)
	sloc -md >> $(DOCOUT)
#	findTags -md *.go >> $(DOCOUT)
	echo "\n" >> $(DOCOUT)
	echo "\`\`\`" >> $(DOCOUT)
	echo built with `go version` >> $(DOCOUT)
	echo "\`\`\`" >> $(DOCOUT)
	echo "\n" >> $(DOCOUT)
	cat compileDate.go >> $(DOCOUT)
	cp README-$(CMDNAME).md README.md
	cat $(DOCOUT) >> README.md
	rm $(DOCOUT)

paranoid:
	deadcode

# I use github.com/junkblocker/codesearch here
index:
	cindex .

neat:
	go fmt ./...

clean:
	go clean ./...
	rm -f *~ $(PROJECT)_linux $(PROJECT)_freebsd

test:
	echo $(DATETIME)
#	go test # -test.run="Test_44"
	go test -test.run="Test_30"

tar:
	echo $(DATETIME)
	go fmt ./...
	go clean ./...
	find . -type f -name "*~" -delete
	find . -type f -name "*linux" -delete
	find . -type f -name "*freebsd" -delete
	tar -ncvf $(TARDIR)$(PROJECT)_$(VERSION)_$(DATETIME).tar .


# Coverage test maker may be more difficult for programs than for packages
# suggestions may be in The Go Programming Language by Donovan and Kernighan
#cover:
#	go test -covermode=count -coverprofile=count.out
#	cover -html=count.out

cover:
	go test -run=01_test -covermode=count -coverprofile=count.out
	cover -html=count.out

bench:
	go test '-bench=.'
