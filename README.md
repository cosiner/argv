# Argv
Argv is a  library for [Go](https://golang.org) to parse command line string into arguments array. 

# Documentation
Documentation can be found at [Godoc](https://godoc.org/github.com/cosiner/argv)

# Example
```Go
func TestArgv(t *testing.T) {
	cases := map[string][]string{
		`     ./argv      `:  []string{`./argv`},
		`./argv ""`:          []string{`./argv`},
		`./argv " \ a \ a "`: []string{`./argv`, ` \ a \ a `},
		`./argv "`:           []string{`./argv`, `"`},
		`./argv " '`:         []string{`./argv`, `"`, `'`},
		`./argv " "`:         []string{`./argv`, ` `},
		`./argv "'"`:         []string{`./argv`, `'`},
		`./argv ''`:          []string{`./argv`},
		`./argv ' '`:         []string{`./argv`, ` `},
		`./argv 'a'`:         []string{`./argv`, `a`},
		`./argv 'a'aa`:       []string{`./argv`, `aaa`},
		`./argv 'a' aa`:      []string{`./argv`, `a`, `aa`},
		`./argv   "a"  `:     []string{`./argv`, `a`},
		`./argv   "'a"  `:    []string{`./argv`, `'a`},
		`./argv   "'a'"  `:   []string{`./argv`, `'a'`},
		`./argv  \" "'a'"  `: []string{`./argv`, `"`, `'a'`},
		`./argv  \" "'a'"\ `: []string{`./argv`, `"`, `'a' `},
		`./argv \  aaa`:      []string{`./argv`, ` `, `aaa`},
		`./argv \ aaa`:       []string{`./argv`, ` aaa`},
		`./argv \\`:          []string{`./argv`, `\`},
		`./argv  \中 `:        []string{`./argv`, `中`},
		`./argv  \ 中 `:       []string{`./argv`, ` 中`},
		`./argv  " 中 " `:     []string{`./argv`, ` 中 `},
	}

	for argv, expect := range cases {
		got := Argv(argv)
		if !reflect.DeepEqual(expect, got) {
			t.Errorf("parse argv '%s' failed, expect: %+v, got: %+v", argv, expect, got)
		}
	}
}
```

# LICENSE
MIT.
