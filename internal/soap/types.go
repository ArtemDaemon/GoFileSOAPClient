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

type SOAPResponseEnvelope struct {
	XMLName xml.Name         `xml:"Envelope"`
	Xmlns   string           `xml:"xmlns:soap,attr"`
	Body    SOAPResponseBody `xml:"Body"`
}

type SOAPResponseBody struct {
	UploadJSONResponse UploadJSONResponse `xml:"UploadJSONResponse"`
}

type UploadJSONResponse struct {
	Status  string `xml:"Status"`
	Message string `xml:"Message"`
}

func UnmarshalEnvelope(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}
