// Copyright (c) The RPSL Go Authors.
// SPDX-License-Identifier: Apache-2.0

package rpsl

import (
	"fmt"
	"strings"
)

type Attribute struct {
	Name  string
	Value string
}

func newAttribute(name, value string) Attribute {
	key := strings.ToLower(name)

	var cleanedLines []string
	lines := strings.Split(value, "\n")
	for i, line := range lines {
		if i > 0 && strings.HasPrefix(line, "+") {
			line = line[1:]
		}

		if strings.Contains(line, "#") {
			line = strings.Split(line, "#")[0]
		}

		line = strings.TrimSpace(line)
		if line != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	value = strings.Join(cleanedLines, " ")
	return Attribute{Name: key, Value: value}
}

func parseAttributes(buf string) ([]Attribute, error) {
	if buf == "" {
		return nil, fmt.Errorf("object cannot be null")
	}

	var attributes []Attribute
	var pos int

	for pos < len(buf) {
		key, newPos, err := parseKey(buf, pos)
		pos = newPos
		if err != nil {
			return nil, err
		}

		value, newPos := parseValue(buf, pos)
		pos = newPos

		attributes = append(attributes, newAttribute(key, value))
	}

	return attributes, nil
}

func isValidKeyChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '*'
}

func isLineContinuationChar(c byte) bool {
	return c == ' ' || c == '\t' || c == '+'
}

func parseKey(buf string, pos int) (string, int, error) {
	start := pos
	for pos < len(buf) {
		c := buf[pos]
		if c == ':' {
			if pos == start {
				return "", 0, fmt.Errorf("read zero sized key")
			}

			return buf[start:pos], pos + 1, nil
		}

		if !isValidKeyChar(c) {
			return "", 0, fmt.Errorf("read illegal character in key: '%c'", c)
		}

		pos++
	}

	return "", 0, fmt.Errorf("no key found")
}

func parseValue(buf string, pos int) (string, int) {
	start := pos
	stop := pos

	for pos < len(buf) {
		c := buf[pos]
		pos++

		if c == '\r' {
			continue
		}

		if c == '\n' && pos < len(buf) {
			next := buf[pos]
			if isLineContinuationChar(next) {
				continue
			}
			break
		}

		stop = pos
	}

	return buf[start:stop], pos
}

func (a *Attribute) String() string {
	var str strings.Builder
	str.WriteString(a.Name)
	str.WriteString(":")
	str.WriteString(a.Value)
	return str.String()
}
