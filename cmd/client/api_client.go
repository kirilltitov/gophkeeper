package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/kirilltitov/gophkeeper/pkg/api"
)

type client struct {
	baseURL    string
	authCookie string
	httpClient *http.Client
}

func newClient(url string, authCookie string) *client {
	return &client{
		baseURL:    strings.TrimRight(url, "/"),
		authCookie: authCookie,
		httpClient: &http.Client{
			Timeout: time.Second,
		},
	}
}

// SendRawRequest Sends an API request and returns HTTP response or an error.
func (c *client) SendRawRequest(
	ctx context.Context,
	url string,
	method string,
	request any,
) (*http.Response, error) {
	fullURL := c.baseURL + url
	logger.Debugf("About to send API request to %s '%s'", method, fullURL)
	if request != nil {
		logger.Tracef("API request body: %+v", request)
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rawRequest, err := http.NewRequestWithContext(
		ctx,
		method,
		fullURL,
		bytes.NewReader(requestBody),
	)

	if err != nil {
		return nil, err
	}

	if c.authCookie != "" {
		rawRequest.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: c.authCookie,
		})
	}

	result, err := c.httpClient.Do(rawRequest)

	if err != nil {
		logger.Debugf("Failed to do API request: %s", err.Error())
		return nil, err
	} else {
		logger.Debugf("Received API response. Status code: %d", result.StatusCode)
	}

	switch result.StatusCode {
	case http.StatusInternalServerError:
		return nil, errInternalServerError
	case http.StatusBadRequest:
		return nil, errBadRequest
	case http.StatusNotFound:
		return nil, errAPIEndpointNotFound
	case http.StatusUnauthorized:
		return nil, errUnauthorized
	}

	return result, err
}

// SendRequest Sends an API request using given client, assigning unmarshalled response to a given pointer,
// and returns status code or an error.
func SendRequest[R any](
	c *client,
	ctx context.Context,
	url string,
	method string,
	request any,
	response *R,
) (code int, err error) {
	rawResponse, err := c.SendRawRequest(ctx, url, method, request)
	if err != nil {
		return 0, err
	}

	defer rawResponse.Body.Close()
	responseBytes, err := io.ReadAll(rawResponse.Body)
	if len(responseBytes) == 0 {
		if response != nil {
			return 0, errNoAPIResponse
		}
		return 0, nil
	}
	logger.Tracef("Raw response body: %s", string(responseBytes))

	var envelopedResponse api.BaseResponse[R]
	if err := json.Unmarshal(responseBytes, &envelopedResponse); err != nil {
		return 0, err
	}

	logger.Tracef("Unmarshalled response envelope: %+v", envelopedResponse)
	if envelopedResponse.Result != nil {
		*response = *envelopedResponse.Result
		logger.Tracef("Unmarshalled response value: %+v", *envelopedResponse.Result)
	}

	return rawResponse.StatusCode, nil
}

func isOffline(err error) bool {
	var opError *net.OpError
	return errors.As(err, &opError)
}
