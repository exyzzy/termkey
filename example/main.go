package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"

	"github.com/exyzzy/termkey"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	fmt.Println("Type 'q' to quit.")

	//put terminal in Raw mode
	in := os.Stdin
	oldState, err := terminal.MakeRaw(int(in.Fd()))
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(int(in.Fd()), oldState)

	//make a *bufio.Reader
	reader := bufio.NewReader(in)

	//make a *termkey.TermKey
	tk := termkey.NewTermKey(reader)

	for {

		// read a key rune from the reader
		key, err := tk.ReadKey()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading input:", err)
			return
		}
		if string(key) == "q" || key == 3 {
			break
		}

		// print the name of a remapped key
		if txt, ok := termkey.IsRemapped(key); ok {
			fmt.Printf("%d (%s)\r\n", key, txt)
		} else if unicode.IsControl(key) {
			fmt.Printf("%d\r\n", key)
		} else {
			fmt.Printf("%d ('%c')\r\n", key, key)
		}

		// check specific keys
		switch key {
		case termkey.KeyUp:
			fmt.Print(">>KeyUp\r\n")
		case termkey.KeyDel:
			fmt.Print(">>KeyDel\r\n")
		case 13:
			fmt.Print(">>Enter\r\n")
		}

	}

	fmt.Print("Quitting...\r\n")
}
