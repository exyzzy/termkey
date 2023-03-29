# Package TermKey

termkey is a lightweight and simple package to read keyboard keys. It uses and extends bytesToKey() from the golang term package (https://github.com/golang/term/blob/master/terminal.go) 

Unlike ReadRune() it also decodes arrow keys, delete, and untab keys.
So far, only tested on a mac pc.

See /example/main.go

To use:

1. add the import

```
import 	"github.com/exyzzy/termkey"
```

2. put the terminal in raw mode

```
	in := os.Stdin
	oldState, err := terminal.MakeRaw(int(in.Fd()))
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(int(in.Fd()), oldState)
```

3. make a *bufio.Reader
```	
reader := bufio.NewReader(in)
```

4. make a *termkey.TermKey

```
tk := termkey.NewTermKey(reader)
```

5. read keypresses as individual runes from the reader

```
    key, err := tk.ReadKey()
    if err != nil {
        if err == io.EOF {
            break
        }
        fmt.Println("Error reading input:", err)
        return
    }
```

6. check IsRemapped to see if the key is remapped and to the textual name

```
    if txt, ok := termkey.IsRemapped(key); ok {
        fmt.Printf("%d (%s)\r\n", key, txt)
    } else if unicode.IsControl(key) {
        fmt.Printf("%d\r\n", key)
    } else {
        fmt.Printf("%d ('%c')\r\n", key, key)
    }
```

7. check the key against any termkey constant 

```
    switch key {
    case termkey.KeyUp:
        fmt.Print(">>KeyUp\r\n")
    case termkey.KeyDel:
        fmt.Print(">>KeyDel\r\n")
    case 13:
        fmt.Print(">>Enter\r\n")
    }
```