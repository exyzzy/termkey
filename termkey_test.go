package termkey

// borrowed and embellished, from https://github.com/golang/term/blob/master/terminal.go

// Parts below are Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bufio"
	"io"
	"testing"
)

type MockTerminal struct {
	toSend       []byte
	bytesPerRead int
	received     []byte
}

func (c *MockTerminal) Read(data []byte) (n int, err error) {
	n = len(data)
	if n == 0 {
		return
	}
	if n > len(c.toSend) {
		n = len(c.toSend)
	}
	if n == 0 {
		return 0, io.EOF
	}
	if c.bytesPerRead > 0 && n > c.bytesPerRead {
		n = c.bytesPerRead
	}
	copy(data, c.toSend[:n])
	c.toSend = c.toSend[n:]
	return
}

func (c *MockTerminal) Write(data []byte) (n int, err error) {
	c.received = append(c.received, data...)
	return len(data), nil
}

var keyPressTests = []struct {
	in   string
	keys []rune
}{
	{ //0
	},
	{ //1
		in:   "\r",
		keys: []rune{13},
	},
	{ //2
		in:   "foo\r",
		keys: []rune{102, 111, 111, 13},
	},
	{ //3
		in:   "a\x1b[Cb\r", // right
		keys: []rune{97, 55304, 98, 13},
	},
	{ //4
		in:   "a\x1b[Db\r", // left
		keys: []rune{97, 55303, 98, 13},
	},
	{ //5
		in:   "a\006b\r", // ^F
		keys: []rune{97, 55304, 98, 13},
	},
	{ //6
		in:   "a\002b\r", // ^B
		keys: []rune{97, 55303, 98, 13},
	},
	{ //7
		in:   "a\177b\r", // backspace
		keys: []rune{97, 127, 98, 13},
	},
	{ //8
		in:   "\x1b[A\r", // up
		keys: []rune{55301, 13},
	},
	{ //9
		in:   "\x1b[B\r", // down
		keys: []rune{55302, 13},
	},
	{ //10
		in:   "\016\r", // ^P
		keys: []rune{55302, 13},
	},
	{ //11
		in:   "\014\r", // ^N
		keys: []rune{55311, 13},
	},
	{ //12
		in:   "line\x1b[A\x1b[B\r", // up then down
		keys: []rune{108, 105, 110, 101, 55301, 55302, 13},
	},
	{ //13
		in:   "line1\rline2\x1b[A\r",
		keys: []rune{108, 105, 110, 101, 49, 13, 108, 105, 110, 101, 50, 55301, 13},
	},
	{ //14
		// line.
		in:   "a b \001\013\r",
		keys: []rune{97, 32, 98, 32, 55307, 55310, 13},
	},
	{ //15
		in:   "a b \001\005\013\r",
		keys: []rune{97, 32, 98, 32, 55307, 55308, 55310, 13},
	},
	{ //16
		in:   "\027\r",
		keys: []rune{55309, 13},
	},
	{ //17
		in:   "Ξεσκεπάζω\r",
		keys: []rune{926, 949, 963, 954, 949, 960, 940, 950, 969, 13},
	},
	{ //18
		in:   "£\r\x1b[A\177\r", // non-ASCII char, enter, up, backspace.
		keys: []rune{163, 13, 55301, 127, 13},
	},

	{ //19
		// Bracketed paste mode: control sequences should be returned
		// verbatim in paste mode.
		in:   "abc\x1b[200~de\177f\x1b[201~\177\r",
		keys: []rune{97, 98, 99, 55312, 100, 101, 127, 102, 55300, 127, 13},
	},
	{ //20
		// Enter in bracketed paste mode should still work.
		in:   "abc\x1b[200~d\refg\x1b[201~h\r",
		keys: []rune{97, 98, 99, 55312, 100, 13, 101, 102, 103, 55300, 104, 13},
	},
	{ //21
		// Ctrl-C at the end of line
		in:   "a\003\r",
		keys: []rune{97, 3, 13},
	},
	{ //22
		in:   "\x1b[Z\r", // untab
		keys: []rune{55314, 13},
	},
	{ //23
		in:   "\x1b[\x33\x7e\r", // del
		keys: []rune{55315, 13},
	},
	{ //24
		in:   "\x1b\x1b\x99\x1b[\x33\x7e\rA", // confuse, recover
		keys: []rune{55300, 13, 65},
	},
}

func TestKeyPresses(t *testing.T) {
	for i, test := range keyPressTests {
		c := &MockTerminal{
			toSend:       []byte(test.in),
			bytesPerRead: 1,
		}
		reader := bufio.NewReader(c)
		ss := NewTermKey(reader)
		var keysRead []rune
		var err error
		for {
			var k rune
			k, err = ss.ReadKey()
			if err != nil {
				break
			}
			keysRead = append(keysRead, k)
		}

		if len(keysRead) != len(test.keys) {
			t.Errorf("Number of keys read do not match test %d were %v, expected %v", i, keysRead, test.keys)
			break
		}
		for x := 0; x < len(test.keys); x++ {
			if keysRead[x] != test.keys[x] {
				t.Errorf("Keys read in test %d were %v, expected %v", i, keysRead, test.keys)
				break
			}
		}
	}
}
