package ocrclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/couchbaselabs/logg"
)

type HttpClient struct {
	ApiEndpointUrl string // the url of the server, eg, http://api.openocr.net
}

func NewHttpClient(apiEndpointUrl string) *HttpClient {
	return &HttpClient{
		ApiEndpointUrl: apiEndpointUrl,
	}
}

func (c HttpClient) DecodeImageReader(imageReader io.Reader, e OcrEngineType) (string, error) {

	ocrRequest := OcrRequest{
		EngineType: e,
	}

	// create JSON for POST reqeust
	jsonBytes, err := json.Marshal(ocrRequest)
	if err != nil {
		return "", err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	mimeHeader := textproto.MIMEHeader{}
	mimeHeader.Set("Content-Type", "application/json")

	part, err := writer.CreatePart(mimeHeader)
	if err != nil {
		logg.LogTo("OCRCLIENT", "Unable to create json multipart part")
		return "", err
	}

	_, err = part.Write(jsonBytes)
	if err != nil {
		logg.LogTo("OCRCLIENT", "Unable to write json multipart part")
		return "", err
	}

	partHeaders := textproto.MIMEHeader{}

	// TODO: pass these vals in instead of hardcoding
	partHeaders.Set("Content-Type", "image/png")
	partHeaders.Set("Content-Disposition", "attachment; filename=\"attachment.txt\".")

	partAttach, err := writer.CreatePart(partHeaders)
	if err != nil {
		logg.LogTo("OCRCLIENT", "Unable to create image multipart part")
		return "", err
	}

	_, err = io.Copy(partAttach, imageReader)
	if err != nil {
		logg.LogTo("OCRCLIENT", "Unable to write image multipart part")
		return "", err
	}

	err = writer.Close()
	if err != nil {
		logg.LogTo("OCRCLIENT", "Error closing writer")
		return "", err
	}

	// create a client
	client := &http.Client{}

	// create POST request
	apiUrl := c.OcrApiFileUploadEndpointUrl()
	req, err := http.NewRequest("POST", apiUrl, bytes.NewReader(body.Bytes()))
	if err != nil {
		logg.LogTo("OCRCLIENT", "Error creating POST request")
		return "", err
	}

	// set content type
	contentType := fmt.Sprintf("multipart/related; boundary=%q", writer.Boundary())
	req.Header.Set("Content-Type", contentType)

	// TODO: code below is all duplicated with DecodeImageUrl()

	// send POST request
	resp, err := client.Do(req)
	if err != nil {
		logg.LogTo("OCRCLIENT", "Error sending POST request")
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("Got error status response: %d", resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logg.LogTo("OCRCLIENT", "Error reading response")
		return "", err
	}

	return string(respBytes), nil

}

func (c HttpClient) DecodeImageUrl(u string, e OcrEngineType) (string, error) {

	ocrRequest := OcrRequest{
		ImgUrl:     u,
		EngineType: e,
	}

	// create JSON for POST reqeust
	jsonBytes, err := json.Marshal(ocrRequest)
	if err != nil {
		return "", err
	}

	// create a client
	client := &http.Client{}

	// create POST request
	apiUrl := c.OcrApiEndpointUrl()
	req, err := http.NewRequest("POST", apiUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		return "", err
	}

	// send POST request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("Got error status response: %d", resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBytes), nil

}

// Get the url of the OCR API endpoint, eg, http://api.openocr.net/ocr
func (c HttpClient) OcrApiEndpointUrl() string {
	return fmt.Sprintf("%v/%v", c.ApiEndpointUrl, "ocr")
}

// Get the url of the OCR API endpoint that supports file upload
func (c HttpClient) OcrApiFileUploadEndpointUrl() string {
	return fmt.Sprintf("%v/%v", c.ApiEndpointUrl, "ocr-file-upload")
}
