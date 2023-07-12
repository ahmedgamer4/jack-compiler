package jacktokenizer

var (
	Keywords = map[string]int{
		"class":       0,
		"function":    0,
		"constructor": 0,
		"int":         0,
		"var":         0,
		"boolean":     0,
		"char":        0,
		"void":        0,
		"static":      0,
		"field":       0,
		"let":         0,
		"if":          0,
		"else":        0,
		"while":       0,
		"return":      0,
		"true":        0,
		"false":       0,
		"null":        0,
		"this":        0,
		"do":          0,
		"method":      0,
	}

	Symbols = map[string]int{
		"{": 0,
		"}": 0,
		"(": 0,
		")": 0,
		"[": 0,
		"]": 0,
		".": 0,
		",": 0,
		";": 0,
		"+": 0,
		"-": 0,
		"*": 0,
		"/": 0,
		"&": 0,
		"|": 0,
		"<": 0,
		">": 0,
		"=": 0,
		"~": 0,
	}
)
