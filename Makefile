test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

bench:
	go test -bench=. -cpuprofile=cpu.out -memprofile=mem.out
