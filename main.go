package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/weberc2/gallium/combinator"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	r := ParseFile(combinator.Input(data))
	if r.Err != nil {
		log.Fatal(r.Err)
	}
	fmt.Println(r.Value.(File).ToGo().Render())
}
