package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/html/charset"
)

func main() {
	fi, _ := os.Open("housemodel.nxr")
	defer fi.Close()

	fo, _ := os.Create("convert.xml")
	defer fo.Close()

	decoder := xml.NewDecoder(fi)
	charsetReader := func(label string, input io.Reader) (io.Reader, error) {
		fmt.Printf("Input charset = %s\n", label)
		return charset.NewReaderLabel(label, input)
	}
	decoder.CharsetReader = charsetReader

	encoder := xml.NewEncoder(fo)
	encoder.Indent("", "  ")

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
			return
		}

		switch tok.(type) {
		case xml.ProcInst:
			encoder.EncodeToken(xml.ProcInst{Target: "xml", Inst: []byte(`version="1.0" charset="utf-8"`)})
		default:
			encoder.EncodeToken(tok)
		}
	}

	encoder.Flush()
}
