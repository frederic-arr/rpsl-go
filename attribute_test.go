// Copyright (c) The RPSL Go Authors.
// SPDX-License-Identifier: Apache-2.0

package rpsl

import (
	"bytes"
	"testing"
)

func TestSingleObject(t *testing.T) {
	data := []byte("" +
		"mntner:          DEV-MNT  # Comment \n" +
		"descr:           DEV maintainer\n" +
		"admin-c:         VM1-DEV\n" +
		"tech-c:          VM1-DEV\n" +
		"upd-to:          v.m@example.net\n" +
		"mnt-nfy:         auto@example.net\n" +
		"auth:            MD5-PW $1$q8Su3Hq/$rJt5M3TNLeRE4UoCh5bSH/\n" +
		"remarks:         password: secret\n" +
		"mnt-by:          DEV-MNT\n" +
		"source:          DEV\n")

	attr, err := parseAttributes(data)
	if err != nil {
		t.Fatalf(`(error): %v`, err)
	}

	if len(attr) != 10 {
		t.Fatalf(`(length): got %v, want %v`, len(attr), 10)
	}

	if attr[0].Name != "mntner" {
		t.Fatalf(`mntner: got %v, want %v`, attr[0].Name, "mntner")
	}

	if attr[0].Value != "DEV-MNT" {
		t.Fatalf(`mntner: got %v, want %v`, attr[0].Value, "DEV-MNT")
	}

	if attr[9].Name != "source" {
		t.Fatalf(`source: got %v, want %v`, attr[9].Name, "source")
	}

	if attr[9].Value != "DEV" {
		t.Fatalf(`source: got %v, want %v`, attr[9].Value, "DEV")
	}
}

func TestAllowedCharactersInKey(t *testing.T) {
	data := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-:")
	attr, err := parseAttributes(data)
	if err != nil {
		t.Fatalf(`(error): %v`, err)
	}

	if len(attr) != 1 {
		t.Fatalf(`(length): got %v, want %v`, len(attr), 1)
	}
}

func TestSpecialCharacters(t *testing.T) {
	data := []byte("" +
		"person:  New Test Person\n" +
		"address: Flughafenstraße 120\n" +
		"address: D - 40474 Düsseldorf\n" +
		"nic-hdl: ABC-RIPE\n")

	attr, err := parseAttributes(data)
	if err != nil {
		t.Fatalf(`(error): %v`, err)
	}

	if len(attr) != 4 {
		t.Fatalf(`(length): got %v, want %v`, len(attr), 4)
	}

	if attr[1].Value != "Flughafenstraße 120" {
		t.Fatalf(`address: got %v, want %v`, attr[1].Value, "Flughafenstraße 120")
	}

	if attr[2].Value != "D - 40474 Düsseldorf" {
		t.Fatalf(`address: got %v, want %v`, attr[2].Value, "D - 40474 Düsseldorf")
	}
}

func TestGarbageValue(t *testing.T) {
	data := []byte("mntner::!@$%^&*()_+~![]{};':<>,./?\\/")

	attr, err := parseAttributes(data)
	if err != nil {
		t.Fatalf(`(error): %v`, err)
	}

	if len(attr) != 1 {
		t.Fatalf(`(length): got %v, want %v`, len(attr), 1)
	}

	if attr[0].Name != "mntner" {
		t.Fatalf(`(0.name): got %v, want %v`, attr[0].Name, "mntner")
	}

	actualBytes := []byte(attr[0].Value)
	expectedBytes := []byte(":!@$%^&*()_+~![]{};':<>,./?\\/")

	if !bytes.Equal(actualBytes, expectedBytes) {
		t.Fatalf("Value bytes don't match")
	}
}

func TestContinuationLines(t *testing.T) {
	data := []byte("" +
		"mntner:   DEV-MNT\n" +
		"descr:    \n" +
		"+1\n" +
		" 2\n" +
		"\t3")

	attr, err := parseAttributes(data)
	if err != nil {
		t.Fatalf(`(error): %v`, err)
	}

	if len(attr) != 2 {
		t.Fatalf(`(length): got %v, want %v`, len(attr), 2)
	}

	if attr[1].Value != "1 2 3" {
		t.Fatalf(`descr: got %v, want %v`, attr[1].Value, "1 2 3")
	}
}

func TestNumberInKey(t *testing.T) {
	data := []byte("" +
		"route6:         2001:0000::/32\n" +
		"origin:AS10")

	attr, err := parseAttributes(data)
	if err != nil {
		t.Fatalf(`(error): %v`, err)
	}

	if len(attr) != 2 {
		t.Fatalf(`(length): got %v, want %v`, len(attr), 2)
	}

	if attr[0].Name != "route6" {
		t.Fatalf(`route6: got %v, want %v`, attr[0].Name, "route6")
	}
}

func TestCasingAndNumberInKey(t *testing.T) {
	data := []byte("" +
		"roUte6:         2001:0000::/32\n" +
		"origin:AS10")

	attr, err := parseAttributes(data)
	if err != nil {
		t.Fatalf(`(error): %v`, err)
	}

	if len(attr) != 2 {
		t.Fatalf(`(length): got %v, want %v`, len(attr), 2)
	}

	if attr[0].Name != "route6" {
		t.Fatalf(`route6: got %v, want %v`, attr[0].Name, "route6")
	}
}

func TestComment(t *testing.T) {
	data := []byte("" +
		"person: foo\n" +
		"nic-hdl: VM1-DEV  # SOME COMMENT \n")

	attr, err := parseAttributes(data)
	if err != nil {
		t.Fatalf(`(error): %v`, err)
	}

	if len(attr) != 2 {
		t.Fatalf(`(length): got %v, want %v`, len(attr), 2)
	}

	if attr[1].Name != "nic-hdl" {
		t.Fatalf(`nic-hdl: got %v, want %v`, attr[1].Name, "nic-hdl")
	}

	if attr[1].Value != "VM1-DEV" {
		t.Fatalf(`nic-hdl: got %v, want %v`, attr[1].Value, "VM1-DEV")
	}
}
