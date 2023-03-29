package termkey

// termkey is a lightweight and simple package to read keyboard keys. It uses and extends bytesToKey() from the golang term package
// unlike ReadRune() it also decodes arrow keys, delete, and untab keys.
// See the main.go example
// Only tested on mac keyboard

// borrowed and embellished, from https://github.com/golang/term/blob/master/terminal.go
// Parts below are Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bufio"
	"bytes"
	"unicode/utf8"
)

// TermKey contains the state for a keypress input session
type TermKey struct {
	inBuf  []byte
	reader *bufio.Reader
}

// NewTermKey runs an input session on the given *bufio.Reader.
// If the Reader is a local terminal, that terminal should first
// be put in raw mode.
func NewTermKey(reader *bufio.Reader) *TermKey {
	var inBuf []byte = make([]byte, 0, 16)
	return &TermKey{
		inBuf:  inBuf,
		reader: reader,
	}
}

// ReadKey returns a keypress rune from the input session.
func (t *TermKey) ReadKey() (rune, error) {
	key := utf8.RuneError
	for key == utf8.RuneError {
		b, err := t.reader.ReadByte()
		if err != nil {
			return KeyUnknown, err
		}
		t.inBuf = append(t.inBuf, b)
		key, t.inBuf = bytesToKey(t.inBuf, false)
		if len(t.inBuf) > 6 {
			t.inBuf = t.inBuf[6:]
		}
	}
	return key, nil
}

// A subset of non-visible public constants that are returned from ReadKey
const (
	KeyTab       = 9
	KeyEnter     = '\r'
	KeyEscape    = 27
	KeyBackspace = 127
	KeyUnknown   = 0xd800 /* UTF-16 surrogate area */ + iota
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyAltLeft
	KeyAltRight
	KeyHome
	KeyEnd
	KeyDeleteWord
	KeyDeleteLine
	KeyClearScreen
	KeyPasteStart
	KeyPasteEnd
	KeyTabOut
	KeyDel
)

var (
	crlf       = []byte{'\r', '\n'}
	pasteStart = []byte{KeyEscape, '[', '2', '0', '0', '~'}
	pasteEnd   = []byte{KeyEscape, '[', '2', '0', '1', '~'}
)

// Standard ascii Keys with names
var standardMap = map[rune]string{
	KeyTab:       "Tab",
	KeyEnter:     "Enter",
	KeyEscape:    "Escape",
	KeyBackspace: "Backspace",
}

// Remapped Keys
var keyMap = map[rune]string{
	KeyUnknown:     "Unknown",
	KeyUp:          "Up",
	KeyDown:        "Down",
	KeyLeft:        "Left",
	KeyRight:       "Right",
	KeyAltLeft:     "AltLeft",
	KeyAltRight:    "AltRight",
	KeyHome:        "Home",
	KeyEnd:         "End",
	KeyDeleteWord:  "DeleteWord",
	KeyDeleteLine:  "DeleteLine",
	KeyClearScreen: "ClearScreen",
	KeyPasteStart:  "PasteStart",
	KeyPasteEnd:    "PasteEnd",
	KeyTabOut:      "TabOut",
	KeyDel:         "Delete",
}

// IsRemapped returns the textual name and true for any remapped keys
func IsRemapped(r rune) (string, bool) {
	s, b := keyMap[r]
	return s, b
}

// bytesToKey tries to parse a key sequence from b. If successful, it returns
// the key and the remainder of the input. Otherwise it returns utf8.RuneError.
func bytesToKey(b []byte, pasteActive bool) (rune, []byte) {

	if len(b) == 0 {
		return utf8.RuneError, nil
	}

	if !pasteActive {
		switch b[0] {
		case 1: // ^A
			return KeyHome, b[1:]
		case 2: // ^B
			return KeyLeft, b[1:]
		case 5: // ^E
			return KeyEnd, b[1:]
		case 6: // ^F
			return KeyRight, b[1:]
		case 8: // ^H
			return KeyBackspace, b[1:]
		case 11: // ^K
			return KeyDeleteLine, b[1:]
		case 12: // ^L
			return KeyClearScreen, b[1:]
		case 23: // ^W
			return KeyDeleteWord, b[1:]
		case 14: // ^N
			return KeyDown, b[1:]
		case 16: // ^P
			return KeyUp, b[1:]
		}
	}

	if b[0] != KeyEscape {
		if !utf8.FullRune(b) {
			return utf8.RuneError, b
		}
		r, l := utf8.DecodeRune(b)
		return r, b[l:]
	}

	if !pasteActive && len(b) >= 3 && b[0] == KeyEscape && b[1] == '[' {
		switch b[2] {
		case 'A':
			return KeyUp, b[3:]
		case 'B':
			return KeyDown, b[3:]
		case 'C':
			return KeyRight, b[3:]
		case 'D':
			return KeyLeft, b[3:]
		case 'H':
			return KeyHome, b[3:]
		case 'F':
			return KeyEnd, b[3:]
		case 'Z':
			return KeyTabOut, b[3:] //added
		}
	}
	if !pasteActive && len(b) >= 4 && b[0] == KeyEscape && b[1] == '[' {
		if b[2] == 0x33 && b[3] == 0x7e {
			return KeyDel, b[4:] //added
		}
	}

	if !pasteActive && len(b) >= 6 && b[0] == KeyEscape && b[1] == '[' && b[2] == '1' && b[3] == ';' && b[4] == '3' {
		switch b[5] {
		case 'C':
			return KeyAltRight, b[6:]
		case 'D':
			return KeyAltLeft, b[6:]
		}
	}

	if !pasteActive && len(b) >= 6 && bytes.Equal(b[:6], pasteStart) {
		return KeyPasteStart, b[6:]
	}

	if pasteActive && len(b) >= 6 && bytes.Equal(b[:6], pasteEnd) {
		return KeyPasteEnd, b[6:]
	}

	// If we get here then we have a Key that we don't recognise, or a
	// partial sequence. It's not clear how one should find the end of a
	// sequence without knowing them all, but it seems that [a-zA-Z~] only
	// appears at the end of a sequence.
	for i, c := range b[0:] {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '~' {
			return KeyUnknown, b[i+1:]
		}
	}

	return utf8.RuneError, b
}
