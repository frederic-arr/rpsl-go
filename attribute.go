// Copyright (c) The RPSL Go Authors.
// SPDX-License-Identifier: Apache-2.0

package rpsl

import (
	"errors"
	"fmt"
	"strings"
)

// Attribute represents a parsed attribute with a normalized key.
type Attribute struct {
	Name  string
	Value string
}

// newAttribute creates an Attribute from a name and value string,
// normalizing the key to lowercase and cleaning the value.
// This version avoids allocating a slice for lines by iterating over the value once.
func newAttribute(name, value string) Attribute {
	var builder strings.Builder
	key := strings.ToLower(name)
	firstLine := true
	start := 0
	n := len(value)

	for i := 0; i <= n; i++ {
		// Look for newline or end-of-string.
		if i == n || value[i] == '\n' {
			line := value[start:i]
			start = i + 1

			// For continuation lines (all but the first), remove a leading '+' if present.
			if !firstLine && len(line) > 0 && line[0] == '+' {
				line = line[1:]
			}

			// Remove any inline comment.
			if idx := strings.IndexByte(line, '#'); idx >= 0 {
				line = line[:idx]
			}

			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				// Separate multiple lines with a space.
				if builder.Len() > 0 {
					builder.WriteByte(' ')
				}
				builder.WriteString(trimmed)
			}

			firstLine = false
		}
	}

	return Attribute{Name: key, Value: builder.String()}
}

// parseAttributes parses the given buffer into a slice of Attributes.
func parseAttributes(buf string) ([]Attribute, error) {
	if buf == "" {
		return nil, errors.New("parseAttributes: object cannot be null")
	}

	var attributes []Attribute
	var pos int

	for pos < len(buf) {
		key, newPos, err := parseKey(buf, pos)
		if err != nil {
			return nil, fmt.Errorf("parseAttributes: %w", err)
		}

		pos = newPos
		value, newPos := parseValue(buf, pos)
		pos = newPos

		attributes = append(attributes, newAttribute(key, value))
	}

	return attributes, nil
}

// isValidKeyChar returns true if c is allowed in a key.
func isValidKeyChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '-' || c == '*'
}

// isLineContinuationChar returns true if c signifies a continuation of a line.
func isLineContinuationChar(c byte) bool {
	return c == ' ' || c == '\t' || c == '+'
}

// parseKey extracts a key ending at the first ':' and returns the key,
// the position after the colon, and an error if any.
func parseKey(buf string, pos int) (string, int, error) {
	start := pos

	for pos < len(buf) {
		c := buf[pos]
		if c == ':' {
			if pos == start {
				return "", 0, fmt.Errorf("parseKey: zero-sized key at pos %d", pos)
			}
			return buf[start:pos], pos + 1, nil
		}

		if !isValidKeyChar(c) {
			return "", 0, fmt.Errorf("parseKey: illegal character '%c' at pos %d", c, pos)
		}

		pos++
	}

	return "", 0, fmt.Errorf("parseKey: no key found starting at pos %d", start)
}

// parseValue extracts a value until a newline that is not followed by a continuation char.
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

// String returns a string representation of the Attribute.
func (a *Attribute) String() string {
	var str strings.Builder
	str.WriteString(a.Name)
	str.WriteString(":")
	str.WriteString(a.Value)
	return str.String()
}
