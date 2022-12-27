mondu: mondu.go
	go build mondu.go

.PHONY: clean
clean:
	-go clean mondu.go
