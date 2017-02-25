package argv

import (
	"math"
	"testing"
)

var (
	parseText = ` a aa a'aa' a"aa"a
		 a$PATH a"$PATH" a'$PATH'
		 a"$*" a"$0" a"$\"
		 a| a|a
		 a"\A" a"\a\b\f\n\r\t\v\\\$" \t a'\A' a'\t'` +
		" a`ls /` `ls ~`"
	env = ParseEnv([]string{
		"PATH=/bin",
		"*=a",
	})
)

func TestScanner(t *testing.T) {
	gots, err := Scan(
		[]rune(parseText),
		env,
	)
	if err != nil {
		t.Fatal(err)
	}
	expects := []Token{
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("a")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("aa")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("aaa")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("aaaa")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("a/bin")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("a/bin")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("a$PATH")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("aa")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("a")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("a$\\")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("a")},
		{Type: TOK_PIPE},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("a")},
		{Type: TOK_PIPE},
		{Type: TOK_STRING, Value: []rune("a")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("a\\A")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("a\a\b\f\n\r\t\v\\$")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("t")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("a\\A")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("a\t")},
		{Type: TOK_SPACE},
		{Type: TOK_STRING, Value: []rune("a")},
		{Type: TOK_REVERSEQUOTE, Value: []rune("ls /")},
		{Type: TOK_SPACE},
		{Type: TOK_REVERSEQUOTE, Value: []rune("ls ~")},
		{Type: TOK_EOF},
	}
	if len(gots) != len(expects) {
		t.Errorf("token count is not equal: expect %d, got %d", len(expects), len(gots))
	}
	l := int(math.Min(float64(len(gots)), float64(len(expects))))
	for i := 0; i < l; i++ {
		got := gots[i]
		expect := expects[i]
		if got.Type != expect.Type {
			t.Errorf("token type is not equal: %d: expect %d, got %d", i, expect.Type, got.Type)
		}

		if expect.Type != TOK_SPACE && string(got.Value) != string(expect.Value) {
			t.Errorf("token value is not equal: %d: expect %s, got %s", i, string(expect.Value), string(got.Value))
		}
	}

	for _, text := range []string{
		`a"`, `a'`, `a"\`, "`ls ~", `a\`,
	} {
		_, err := Scan([]rune(text), nil)
		if err != ErrInvalidSyntax {
			t.Errorf("expect unexpected eof error, but got: %v", err)
		}
	}
}
