mondu: main.go
	go build -o mondu main.go

.PHONY: clean
clean:
	-rm -rf mondu
