# Coverage test maker

cover: mdr.go
	go test -covermode=count -coverprofile=count.out
	go tool cover -html=count.out