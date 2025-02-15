package misc

import (
	"fmt"
	"io/ioutil"

	markdown "github.com/n1rna/go-term-markdown"
)

func main() {
	path := "Readme.md"
	source, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	result := markdown.Render(string(source), 80, 1, 6)

	fmt.Println(result)
}
