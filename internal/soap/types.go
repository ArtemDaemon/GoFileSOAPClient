package soap

import (
	"encoding/xml"
)

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    SOAPBody `xml:"Body"`
}

type SOAPBody struct {
	UploadFileRequest UploadFileRequest `xml:"UploadFileRequest"`
}

type UploadFileRequest struct {
	Filename string     `xml:"Filename"`
	File     XOPInclude `xml:"File"`
}

type XOPInclude struct {
	Href string `xml:"href,attr"`
}

type SOAPResponseEnvelope struct {
	XMLName xml.Name         `xml:"Envelope"`
	Xmlns   string           `xml:"xmlns:soap,attr"`
	Body    SOAPResponseBody `xml:"Body"`
}

type SOAPResponseBody struct {
	UploadFileResponse UploadFileResponse `xml:"UploadFileResponse"`
}

type UploadFileResponse struct {
	Status  string `xml:"Status"`
	Message string `xml:"Message"`
}

func UnmarshalEnvelope(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}
