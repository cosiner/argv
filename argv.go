// Package argv parse command line string into arguments array.
//
// Support backslash '\' and quotes ''' and '"'.
package argv

import "strings"

// Argv parse command line string into arguments array.
func Argv(str string) []string {
	var (
		argv []string
		curr []rune

		skip        = -1
		prev        rune
		prevSpecial bool
	)
	for i, r := range str {
		if i <= skip {
			continue
		}

		var isSpecial bool
		if prev != '\\' || !prevSpecial {
			switch r {
			case '"', '\'':
				nextIndex := strings.IndexRune(str[i+1:], r)
				if nextIndex >= 0 {
					isSpecial = true

					skip = i + nextIndex + 1
					curr = append(curr, []rune(str[i+1:skip])...)
				}
			case '\\':
				isSpecial = true
			case ' ':
				isSpecial = true
				if prev != ' ' || !prevSpecial {
					isSpecial = true
					if len(curr) > 0 {
						argv = append(argv, string(curr))
						curr = curr[:0]
					}
				}
			}
		}
		if !isSpecial {
			curr = append(curr, r)
		}
		prev = r
		prevSpecial = isSpecial
	}

	if len(curr) > 0 {
		argv = append(argv, string(curr))
	}
	return argv
}
