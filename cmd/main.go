package main

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"go-file-soap-client/internal/soap"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Читаем .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file:", err)
	}

	// Читаем значение API_TOKEN
	token := os.Getenv("API_TOKEN")
	if token == "" {
		log.Fatal("API_TOKEN parameter not found in environment variables file:", err)
	}

	// Читаем значения аргументов командной строки
	if len(os.Args) < 2 {
		log.Fatal("Usage: client <json-file>")
	}
	filename := os.Args[1]

	// Читаем JSON-файл
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Failed to read JSON file:", err)
	}

	// Кодируем JSON-файл
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
		log.Fatal("Failed to encode SOAP enveloper:", err)
	}

	// Формируем зарпос
	req, err := http.NewRequest("POST", "http://localhost:8080/api/soap", &buf)
	if err != nil {
		log.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+token)

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to send: request:", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var soapResp soap.SOAPResponseEnvelope
	if err := soap.UnmarshalEnvelope(respBody, &soapResp); err != nil {
		fmt.Println("Failed to parse server response:", err)
		fmt.Println("Raw response")
		fmt.Println(string(respBody))
		return
	}

	fmt.Println("Server response:")
	fmt.Printf("Status: %s\n", soapResp.Body.UploadJSONResponse.Status)
	fmt.Printf("Message: %s\n", soapResp.Body.UploadJSONResponse.Message)
}
