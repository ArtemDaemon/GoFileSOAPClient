package main

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"go-file-soap-client/internal/soap"
	"io"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: client <json-file>")
		return
	}
	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Failed to read file:", err)
		return
	}

	encoded := base64.StdEncoding.EncodeToString(data)

	envelope := soap.SOAPEnvelope{
		Body: soap.SOAPBody{
			UploadJSONRequest: soap.UploadJSONRequest{
				Filename: filename,
				Content:  encoded,
			},
		},
	}

	var buf bytes.Buffer
	if err := xml.NewEncoder(&buf).Encode(envelope); err != nil {
		fmt.Println("Failed to encode SOAP enveloper:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/soap", "text/xml; charset=utf-8", &buf)
	if err != nil {
		fmt.Println("Failed to send: request:", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Server response:")
	fmt.Println(string(respBody))
}
