# Argv
Argv is a  library for [Go](https://golang.org) to split command line string into arguments array. 

# Documentation
Documentation can be found at [Godoc](https://godoc.org/github.com/cosiner/argv)

# Example
```Go
func TestArgv(t *testing.T) {
	args, err := argv.Argv([]rune(" ls   `echo /`   |  wc  -l "), os.Environ(), argv.Run)
	if err != nil {
	    t.Fatal(err)
	}
	expects := [][]string{
	    []string{"ls", "/"},
	    []string{"wc", "-l"},
	}
	if !reflect.DeepDqual(args, expects) {
	    t.Fatal(args)
	}
}
```

# LICENSE
MIT.
