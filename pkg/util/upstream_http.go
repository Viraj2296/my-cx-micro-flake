package util

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func UpstreamGet(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	var response []byte
	if err != nil {
		return response, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	//client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return response, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

// Creates a new file upload http request with optional extra params
func fileUploadFromDisk(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

// Creates a new file upload http request with optional extra params
func fileUploadFromRawData(uri string, params map[string]string, paramName string, filePath string, rawData []byte) (*http.Request, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, bytes.NewReader(rawData))

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
func FileUploadFromDisk(upstreamFileServerURL string, params map[string]string, paramName string, path string) (error, *http.Response, []byte) {
	request, err := fileUploadFromDisk(upstreamFileServerURL, params, paramName, path)
	var responseBody []byte
	if err != nil {
		return err, nil, responseBody
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err, nil, responseBody
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			return err, nil, responseBody
		}
		responseBody = body.Bytes()
		resp.Body.Close()
		return nil, resp, responseBody
	}
}

func FileUploadFromRawData(upstreamFileServerURL string, params map[string]string, paramName string, filePath string, rawData []byte) (error, *http.Response, []byte) {
	request, err := fileUploadFromRawData(upstreamFileServerURL, params, paramName, filePath, rawData)
	var responseBody []byte
	if err != nil {
		return err, nil, responseBody
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err, nil, responseBody
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			return err, nil, responseBody
		}
		responseBody = body.Bytes()
		resp.Body.Close()
		return nil, resp, responseBody
	}
}
