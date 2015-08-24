# Makefile for mdr package

PROG = mdr
VERSION = 0.1.3
TARDIR = 	$(HOME)/Desktop/TarPit/
DATE = 	`date "+%Y-%m-%d.%H_%M_%S"`
DOCOUT = README-$(PROG)-godoc.md

all:
	go build -v
	
install:
	go build
	go tool vet .
	go tool vet -shadow .
	gofmt -w *.go
	go install
#	cp $(PROG) $(HOME)/bin
	

# note that godepgraph can be used to derive .travis.yml install: section
docs:
	godoc2md . > $(DOCOUT)
	godepgraph -md -p . >> $(DOCOUT)
	deadcode -md >> $(DOCOUT)
	echo "\`\`\`" >> $(DOCOUT)
	echo built with go version = $(GOVERSION) >> $(DOCOUT)
	echo "\`\`\`" >> $(DOCOUT)
	cp README-$(PROG).md README.md
	cat $(DOCOUT) >> README.md
	cp README.md README2.md
	
neat:
	go fmt ./...

dead:
	deadcode > problems.dead

index:
	cindex .

clean:
	go clean ./...
	rm -f *~ problems.dead count.out README2.md $(DOCOUT)

tar:
	echo $(TARDIR)$(PROG)_$(VERSION)_$(DATE).tar
	tar -ncvf $(TARDIR)$(PROG)_$(VERSION)_$(DATE).tar .

# Coverage test maker

cover:
	go test -run=01_test -covermode=count -coverprofile=count.out
	cover -html=count.out
