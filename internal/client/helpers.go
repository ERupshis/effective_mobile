package client

import (
	"fmt"
	"net/url"
)

func addQueryInURI(base string, params map[string]string) (string, error) {
	if len(params) == 0 {
		return base, nil
	}

	u, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("Error parsing URL: %w", err)
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
