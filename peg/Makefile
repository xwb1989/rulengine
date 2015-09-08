default:
	${GOPATH}/bin/peg -inline -switch parser.peg 

.PHONY: clean
clean:
	rm -rf parser.peg.go

.PHONY: bench
bench:
	go test -bench=. -run=XXX
