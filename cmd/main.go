package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"go-file-soap-client/internal/soap"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

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

	// Создаём multipart/related запрос для MTOM
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Генерируем Content-ID для файла
	fileCID := "file1@mtom"

	// 1. SOAP Enveloper part
	soapEnvelope := soap.SOAPEnvelope{
		Body: soap.SOAPBody{
			UploadFileRequest: soap.UploadFileRequest{
				Filename: filepath.Base(filename),
				File: soap.XOPInclude{
					Href: "cid:" + fileCID,
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := xml.NewEncoder(&buf).Encode(soapEnvelope); err != nil {
		log.Fatal("Failed to encode SOAP envelope:", err)
	}

	soapPartHeaders := textproto.MIMEHeader{}
	soapPartHeaders.Set("Content-Type", `application/xop+xml; charset=UTF-8; type="text/xml"`)
	soapPartHeaders.Set("Content-Transfer-Encoding", "8bit")
	soapPartHeaders.Set("Content-ID", "<rootpart@mtom>")
	soapPart, err := writer.CreatePart(soapPartHeaders)
	if err != nil {
		log.Fatal("Failed to create SOAP part:", err)
	}
	_, err = soapPart.Write(buf.Bytes())
	if err != nil {
		log.Fatal("Failed to write SOAP part:", err)
	}

	// 2. File part
	filePartHeaders := textproto.MIMEHeader{}
	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	filePartHeaders.Set("Content-Type", mimeType)
	filePartHeaders.Set("Content-Transfer-Encoding", "binary")
	filePartHeaders.Set("Content-ID", "<"+fileCID+">")
	filePart, err := writer.CreatePart(filePartHeaders)
	if err != nil {
		log.Fatal("Failed to create flie part:", err)
	}
	_, err = filePart.Write(data)
	if err != nil {
		log.Fatal("Failed to write file part:", err)
	}
	writer.Close()

	// Формируем зарпос
	req, err := http.NewRequest("POST", "http://localhost:8080/api/soap", &body)
	if err != nil {
		log.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Content-Type", `multipart/related; type="application/xop+xml"; start="<rootpart@mtom>"; boundary=`+writer.Boundary())
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
	fmt.Printf("Status: %s\n", soapResp.Body.UploadFileResponse.Status)
	fmt.Printf("Message: %s\n", soapResp.Body.UploadFileResponse.Message)
}
