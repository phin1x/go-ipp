package main

import (
	"bytes"
	"fmt"
	"github.com/phin1x/go-ipp"
)

func main() {
	doc := new(bytes.Buffer)
	doc.WriteString("asdfasdf")

	ippclt := ipp.NewIPPClient("", 0, "", "", false)
	fmt.Println(ippclt.Print(doc, doc.Len(), "blub", "testprint", 1, 50))
}
