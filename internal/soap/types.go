package soap

import (
	"encoding/xml"
)

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    SOAPBody `xml:"Body"`
}

type SOAPBody struct {
	UploadJSONRequest UploadJSONRequest `xml:"UploadJSONRequest"`
}

type UploadJSONRequest struct {
	Filename string `xml:"Filename"`
	Content  string `xml:"Content"`
}
