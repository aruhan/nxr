package nxr

import (
	"encoding/xml"
	"io"
	"log"

	"golang.org/x/net/html/charset"
)

// FormatXML is format XML
func FormatXML(r io.Reader, w io.Writer) error {

	decoder := xml.NewDecoder(r)
	charsetReader := func(label string, input io.Reader) (io.Reader, error) {
		log.Printf("Detect charset: %v", label)
		return charset.NewReaderLabel(label, input)
	}
	decoder.CharsetReader = charsetReader

	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ") // 2 space

	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		switch t.(type) {
		case xml.ProcInst:
			// 出力は utf-8 のみ
			encoder.EncodeToken(xml.ProcInst{
				Target: "xml",
				Inst:   []byte(`version="1.0" charset="utf-8"`),
			})
		default:
			encoder.EncodeToken(t)
		}
	}

	encoder.Flush()
	return nil
}
