package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli"
	"golang.org/x/net/html/charset"
)

func main() {
	app := cli.NewApp()
	app.Name = "nxrfmt"
	app.Usage = "nxr formatter"
	app.Author = "HARADA Takashi"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input, i",
			Usage: "input file",
			Value: "",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "output file",
		},
	}
	app.Commands = []cli.Command{
		formatCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var formatCommand = cli.Command{
	Name:  "format",
	Usage: "format nxr/mpb/lit/(xml) file",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Usage: "output file",
		},
	},
	Action: func(c *cli.Context) error {
		input := c.GlobalString("input")
		output := c.String("output")
		if output == "" {
			output = input + ".out"
		}
		doFormatCommand(input, output)
		return nil
	},
}

func doFormatCommand(input string, output string) {
	fi, _ := os.Open(input)
	defer fi.Close()

	fo, _ := os.Create(output)
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
