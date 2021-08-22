package fuzzing

import markdown "github.com/n1rna/go-term-markdown"

func Fuzz(data []byte) int {
	markdown.Render(string(data), 50, 1, 4)
	return 1
}
