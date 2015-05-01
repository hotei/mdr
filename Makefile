# Coverage test maker

cover: mdr.go
	go test -covermode=count -coverprofile=count.out
	cover -html=count.out
