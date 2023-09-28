package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/erupshis/effective_mobile/internal/helpers"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/retryer"
)

type DefaultClient struct {
	client *http.Client
	log    logger.BaseLogger
}

func CreateDefault(log logger.BaseLogger) BaseClient {
	return &DefaultClient{client: &http.Client{}, log: log}
}

func (c *DefaultClient) DoGetURIWithQuery(ctx context.Context, url string, params map[string]string) (int64, []byte, error) {
	uri, err := addQueryInURI(url, params)
	if err != nil {
		return http.StatusInternalServerError, []byte{}, fmt.Errorf("couldn't send post request")
	}

	request := func(context context.Context) (int64, []byte, error) {
		return c.makeEmptyBodyRequest(context, http.MethodGet, uri)
	}

	c.log.Info("[DefaultClient:DoGetURIWithQuery] make request to: %s", uri)
	status, respBody, err := retryer.RetryCallWithTimeout(ctx, c.log, nil, nil, request)
	if err != nil {
		err = fmt.Errorf("couldn't send post request")
	}
	return status, respBody, err
}

func (c *DefaultClient) makeEmptyBodyRequest(ctx context.Context, method string, url string) (int64, []byte, error) {
	bufReq := bytes.NewBuffer([]byte{})
	req, err := http.NewRequestWithContext(ctx, method, url, bufReq)
	if err != nil {
		return http.StatusInternalServerError, []byte{}, fmt.Errorf("request creation: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return http.StatusBadRequest, []byte{}, fmt.Errorf("client request error: %w", err)
	}
	defer helpers.ExecuteWithLogError(resp.Body.Close, c.log)

	bufResp := bytes.Buffer{}
	_, err = bufResp.ReadFrom(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, []byte{}, fmt.Errorf("read response body: %w", err)
	}

	return int64(resp.StatusCode), bufResp.Bytes(), err
}
