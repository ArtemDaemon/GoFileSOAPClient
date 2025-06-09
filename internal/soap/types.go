package soap

import (
	"encoding/xml"
	"net/http"
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

func WriteSOAPResponse(w http.ResponseWriter, status, message string) {
	response := SOAPResponseEnvelope{
		Xmlns: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: SOAPResponseBody{
			UploadJSONResponse: UploadJSONResponse{
				Status:  status,
				Message: message,
			},
		},
	}
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	xml.NewEncoder(w).Encode(response)
}
