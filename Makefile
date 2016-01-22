bench:
	go test -bench Memory -benchmem -run NoTests -memprofile pp/mem.out -cpuprofile pp/cpu.out
	go tool pprof -svg --alloc_space delay.test pp/mem.out > pp/mem.svg
	open pp/mem.svg

cpu:
	go test -bench Memory -run NoTests -cpuprofile pp/cpu.out
	go tool pprof -svg delay.test pp/cpu.out > pp/cpu.svg
	open pp/cpu.svg

.PHONY: bench cpu
