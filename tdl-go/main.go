package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var file File
	if _, err := file.Parser().Parse(Input{Source: data}); err != nil {
		log.Fatal(err)
	}
	fmt.Println(file.ToGo().Render())
}
