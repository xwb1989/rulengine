package parser

import (
	"testing"
)

func BenchmarkReUseParser(b *testing.B) {
	parser := yyNewParser()
	for i := 0; i < b.N; i++ {
		expr := "1+2+3\n"
		parser.Parse(&CalcLex{s: expr})
	}
}

func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser := yyNewParser()
		expr := "1+2+3\n"
		parser.Parse(&CalcLex{s: expr})
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		expr := "1+2+3\n"
		yyParse(&CalcLex{s: expr})
	}
}
