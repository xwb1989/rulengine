MAKEFLAGS = -s

#expr.go: expr.y
	#go tool yacc -o expr.go expr.y
   # gofmt -w expr.go

expr.go: expr.y
	go tool yacc -o expr.go expr.y
	gofmt -w expr.go


clean:
	rm -f y.output expr.go
