package wordlist

import (
	_ "embed"
)

//go:embed words.txt
var Source string
