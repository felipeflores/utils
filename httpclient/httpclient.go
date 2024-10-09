package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/felipeflores/utils/httpclient/model"
)

type HttpClient[T any] struct {
	client *http.Client
}

func New[T any]() *HttpClient[T] {
	client := &http.Client{}
	return &HttpClient[T]{
		client: client,
	}
}

func (h *HttpClient[T]) PostFormUrlEncoded(path string, formData map[string]string, response *T) error {
	data := url.Values{}
	for key, value := range formData {
		data.Set(key, value)
	}

	req, err := http.NewRequest(model.Post, path, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set(model.ContentType, model.FormUrlEncoded)

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	return nil
}

func (h *HttpClient[T]) Post(ctx context.Context, path string, headers map[string]string, body T) (int, *http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return 0, nil, err
	}

	buffer := bytes.NewBuffer(b)

	req, err := http.NewRequestWithContext(ctx, model.Post, path, buffer)
	if err != nil {
		return 0, nil, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Dump the request
	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Println("Error dumping request:", err)
		return 0, nil, err
	}
	fmt.Printf("Request:\n%s\n", string(requestDump))

	resp, err := h.client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	// Dump the response
	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println("Error dumping response:", err)
		return 0, nil, err
	}
	fmt.Printf("Response:\n%s\n", string(responseDump))
	return resp.StatusCode, resp, nil
}

func (h *HttpClient[T]) Get(ctx context.Context, path string, headers map[string]string, response *T) error {
	req, err := http.NewRequestWithContext(ctx, model.Get, path, nil)
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	return nil
}
