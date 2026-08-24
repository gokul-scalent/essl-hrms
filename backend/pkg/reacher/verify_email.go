package reacher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/scalent.io/scalent-hrms/pkg/log"
)

func (c *Client) VerifyEmail(ctx context.Context, email string) (*VerifyEmailResponse, error) {
	log.Info("Calling Reacher API for: "+email, "")

	// Create the request payload.
	reqBody := VerifyEmailRequest{
		ToEmail: email,
	}

	// Convert the request payload to JSON.
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// call an http POST request with the provided context.
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/v1/check_email", c.baseURL),
		bytes.NewBuffer(b),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request to the Reacher API.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Info(fmt.Sprintf("Reacher Status Code: %d", resp.StatusCode), "")

	// Return an error if the API responds with a non-200 status.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("reacher returned %d: %s", resp.StatusCode, string(body))
	}
	// Decode the JSON response into the response struct.
	var result VerifyEmailResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	log.Info(fmt.Sprintf("Reacher Response: %+v", result), "")
	return &result, nil
}
