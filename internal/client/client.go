package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MaximBayurov/rate-limiter/internal/configuration"
)

type Client interface {
	AddIP(ctx context.Context, ip, listType string, overwrite bool) (Response, error)
	DeleteIP(ctx context.Context, ip, listType string) (Response, error)
	ClearBucket(ctx context.Context, ip, login string) (Response, error)
	TryAuth(ctx context.Context, login, password, ip string) (Response, error)
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func New(configs configuration.ClientConf) Client {
	urlParts := []string{
		configs.Host,
		configs.Port,
	}
	return &AppClient{
		baseURL: strings.Join(urlParts, ":"),
		client: http.Client{
			Timeout: configs.Timeout * time.Second, //nolint:durationcheck
		},
	}
}

type AppClient struct {
	baseURL string
	client  http.Client
}

func (c *AppClient) AddIP(ctx context.Context, ip, listType string, overwrite bool) (Response, error) {
	payload := struct {
		IP        string `json:"ip"`
		ListType  string `json:"type"`
		Overwrite bool   `json:"overwrite"`
	}{
		IP:        ip,
		ListType:  listType,
		Overwrite: overwrite,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("client payload marshaling: %w", err)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		"PUT",
		c.baseURL+"/ip/list",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return Response{}, fmt.Errorf("client request create: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("client request execution: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Response{}, fmt.Errorf("client response decode: %w", err)
	}

	return result, nil
}

func (c *AppClient) DeleteIP(ctx context.Context, ip, listType string) (Response, error) { //nolint:dupl
	payload := struct {
		IP       string `json:"ip"`
		ListType string `json:"type"`
	}{
		IP:       ip,
		ListType: listType,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("client payload marshaling: %w", err)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		"DELETE",
		c.baseURL+"/ip/list",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return Response{}, fmt.Errorf("client request create: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("client request execution: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Response{}, fmt.Errorf("client response decode: %w", err)
	}

	return result, nil
}

func (c *AppClient) ClearBucket(ctx context.Context, ip, login string) (Response, error) { //nolint:dupl
	payload := struct {
		IP    string `json:"ip"`
		Login string `json:"login"`
	}{
		IP:    ip,
		Login: login,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("client payload marshaling: %w", err)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		"DELETE",
		c.baseURL+"/auth",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return Response{}, fmt.Errorf("client request create: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("client request execution: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Response{}, fmt.Errorf("client response decode: %w", err)
	}

	return result, nil
}

func (c *AppClient) TryAuth(ctx context.Context, login, password, ip string) (Response, error) { //nolint:dupl
	payload := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		IP       string `json:"ip"`
	}{
		Login:    login,
		Password: password,
		IP:       ip,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("client payload marshaling: %w", err)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		"PUT",
		c.baseURL+"/auth",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return Response{}, fmt.Errorf("client request create: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("client request execution: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Response{}, fmt.Errorf("client response decode: %w", err)
	}

	return result, nil
}
