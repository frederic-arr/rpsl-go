// Copyright (c) The RPSL Go Authors.
// SPDX-License-Identifier: Apache-2.0

package rpsl

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func BenchmarkParse(b *testing.B) {
	raw := "person:	John Doe\n" +
		"address:	1234 Elm Street\n" +
		"phone:		+1 555 123456\n" +
		"nic-hdl:	JD1234-RIPE\n" +
		"mnt-by:	EXAMPLE-MNT\n" +
		"source:	RIPE\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Parse(raw)
	}
}

func BenchmarkParseMany(b *testing.B) {
	raw := "person:	John Doe\n" +
		"address:	1234 Elm Street\n" +
		"phone:		+1 555 123456\n" +
		"nic-hdl:	JD1234-RIPE\n" +
		"mnt-by:	EXAMPLE-MNT\n" +
		"source:	RIPE\n" +
		"\n" +
		"person:	Jane Smith\n" +
		"address:	5678 Oak Street\n" +
		"phone:		+1 555 654321\n" +
		"nic-hdl:	JS5678-RIPE\n" +
		"mnt-by:	EXAMPLE-MNT\n" +
		"source:	RIPE"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseMany(raw)
	}
}

func BenchmarkParseLargeFile(b *testing.B) {
	data, err := os.ReadFile("tests/data/ripe.db.as-set")
	if err != nil {
		b.Skip("skipping large file benchmark; file not found")
	}
	s := string(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseMany(s)
	}
}

func BenchmarkReaderNext(b *testing.B) {
	data, err := os.ReadFile("tests/data/ripe.db.as-set")
	if err != nil {
		b.Skip("skipping large file benchmark; file not found")
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := NewReader(bytes.NewReader(data))
		for {
			_, err = r.Next()
			if err != nil {
				if err == io.EOF {
					break
				}
				b.Fatal(err)
			}
		}
	}
}
