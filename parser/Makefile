MAKEFLAGS = -s

#expr.go: expr.y
	#go tool yacc -o expr.go expr.y
   # gofmt -w expr.go

calc.go: calc.y
	go tool yacc -o calc.go calc.y
	gofmt -w calc.go


clean:
	rm -f y.output expr.go
