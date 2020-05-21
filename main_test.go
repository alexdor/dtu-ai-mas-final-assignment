package main

import (
	"bytes"
	"strings"
	"testing"
)

func BenchmarkConcat(b *testing.B) {
	b.ReportAllocs()
	stra := []string{}
	for n := 0; n < b.N; n++ {
		stra = append(stra, "x")
	}
	str := strings.Join(stra, "")
	b.StopTimer()

	if s := strings.Repeat("x", b.N); str != s {
		b.Errorf("unexpected result; got=%s, want=%s", str, s)
	}
}
func BenchmarkConcatPrealloc(b *testing.B) {
	b.ReportAllocs()
	stra := make([]string, b.N)
	for n := 0; n < b.N; n++ {
		stra = append(stra, "x")
	}
	str := strings.Join(stra, "")
	b.StopTimer()

	if s := strings.Repeat("x", b.N); str != s {
		b.Errorf("unexpected result; got=%s, want=%s", str, s)
	}
}
func BenchmarkBuffer(b *testing.B) {
	b.ReportAllocs()
	var buffer bytes.Buffer
	for n := 0; n < b.N; n++ {
		buffer.WriteString("x")
	}
	b.StopTimer()

	if s := strings.Repeat("x", b.N); buffer.String() != s {
		b.Errorf("unexpected result; got=%s, want=%s", buffer.String(), s)
	}
}

func BenchmarkCopy(b *testing.B) {
	b.ReportAllocs()
	bs := make([]byte, b.N)
	bl := 0

	for n := 0; n < b.N; n++ {
		bl += copy(bs[bl:], "x")
	}
	b.StopTimer()

	if s := strings.Repeat("x", b.N); string(bs) != s {
		b.Errorf("unexpected result; got=%s, want=%s", string(bs), s)
	}
}

func BenchmarkBytesAppend(b *testing.B) {
	b.ReportAllocs()
	bs := []byte{}

	for n := 0; n < b.N; n++ {
		bs = append(bs, []byte("x")...)
	}
	b.StopTimer()

	if s := strings.Repeat("x", b.N); string(bs) != s {
		b.Errorf("unexpected result; got=%s, want=%s", string(bs), s)
	}
}

func BenchmarkStringBuilder(b *testing.B) {
	b.ReportAllocs()
	var strBuilder strings.Builder

	for n := 0; n < b.N; n++ {
		strBuilder.WriteString("x")
	}
	b.StopTimer()

	if s := strings.Repeat("x", b.N); strBuilder.String() != s {
		b.Errorf("unexpected result; got=%s, want=%s", strBuilder.String(), s)
	}
}

func BenchmarkStringAdd(b *testing.B) {
	b.ReportAllocs()
	var str string
	for n := 0; n < b.N; n++ {
		str += "x"
	}
	b.StopTimer()

	if s := strings.Repeat("x", b.N); str != s {
		b.Errorf("unexpected result; got=%s, want=%s", str, s)
	}
}
