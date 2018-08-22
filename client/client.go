// Copyright 2017, 2018 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/AtlantPlatform/atlant-go/proto"
)

type Client interface {
	Ping(ctx context.Context) (id string, err error)
	Version(ctx context.Context) (ver string, err error)
	PutObject(ctx context.Context, path string, obj *PutObjectInput) (*ObjectMeta, error)
	GetContents(ctx context.Context, path, version string) ([]byte, error)
	GetMeta(ctx context.Context, path, version string) (*ObjectMeta, error)
	DeleteObject(ctx context.Context, id string) (*ObjectMeta, error)
	ListVersions(ctx context.Context, path string) (id string, versions []*ObjectMeta, err error)
	ListObjects(ctx context.Context, prefix string) (dirs []string, files []*ObjectMeta, err error)
}

func NewID() string {
	return proto.NewID()
}

type ObjectMeta struct {
	ID              string `json:"id"`
	Path            string `json:"path,omitempty"`
	CreatedAt       int64  `json:"createdAt,omitempty"`
	Version         string `json:"version,omitempty"`
	VersionPrevious string `json:"versionPrevious,omitempty"`
	IsDeleted       bool   `json:"isDeleted,omitempty"`
	Size            int64  `json:"size,omitempty"`
	UserMeta        string `json:"userMeta,omitempty"`
}

type rpcClient struct {
	apiURL string
	cli    *http.Client
}

func New(apiURL string) Client {
	return &rpcClient{
		apiURL: apiURL,
		cli: &http.Client{
			Transport: &http.Transport{
				ResponseHeaderTimeout: 30 * time.Second,
			},
		},
	}
}

func (client *rpcClient) Ping(ctx context.Context) (id string, err error) {
	data, err := client.get(ctx, "/api/v1/ping", nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (client *rpcClient) Version(ctx context.Context) (ver string, err error) {
	data, err := client.get(ctx, "/api/v1/version", nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type PutObjectInput struct {
	Body     io.ReadCloser
	Size     int64
	UserMeta string
}

func (client *rpcClient) PutObject(ctx context.Context, path string, obj *PutObjectInput) (*ObjectMeta, error) {
	contentType := mime.TypeByExtension(filepath.Base(path))
	if len(contentType) == 0 {
		contentType = "application/binary"
	}
	headers := map[string]string{
		"X-Meta-UserMeta": obj.UserMeta,
	}
	respData, err := client.post(ctx, filepath.Join("/api/v1/put", path), contentType, obj.Body, obj.Size, headers)
	if err != nil {
		return nil, err
	}
	var meta *ObjectMeta
	if err := json.Unmarshal(respData, &meta); err != nil {
		err = fmt.Errorf("response unmarshal failed: %v", err)
		return nil, err
	}
	return meta, err
}

func (client *rpcClient) GetContents(ctx context.Context, path, version string) ([]byte, error) {
	return client.get(ctx, filepath.Join("/api/v1/content", path)+"?ver="+version, nil)
}

func (client *rpcClient) GetMeta(ctx context.Context, path, version string) (*ObjectMeta, error) {
	respData, err := client.get(ctx, filepath.Join("/api/v1/meta", path)+"?ver="+version, nil)
	if err != nil {
		return nil, err
	}
	var meta *ObjectMeta
	if err := json.Unmarshal(respData, &meta); err != nil {
		err = fmt.Errorf("response unmarshal failed: %v", err)
		return nil, err
	}
	return meta, nil
}

func (client *rpcClient) DeleteObject(ctx context.Context, id string) (*ObjectMeta, error) {
	respData, err := client.post(ctx, "/api/v1/delete/"+id, "", nil, 0, nil)
	if err != nil {
		return nil, err
	}
	var meta *ObjectMeta
	if err := json.Unmarshal(respData, &meta); err != nil {
		err = fmt.Errorf("response unmarshal failed: %v", err)
		return nil, err
	}
	return meta, nil
}

type listVersionsResponse struct {
	ID       string        `json:"id"`
	Versions []*ObjectMeta `json:"versions"`
}

func (client *rpcClient) ListVersions(ctx context.Context, path string) (id string, versions []*ObjectMeta, err error) {
	respData, err := client.get(ctx, filepath.Join("/api/v1/listVersions", path), nil)
	if err != nil {
		return "", nil, err
	}
	var resp listVersionsResponse
	if err := json.Unmarshal(respData, &resp); err != nil {
		err = fmt.Errorf("response unmarshal failed: %v", err)
		return "", nil, err
	}
	return resp.ID, resp.Versions, nil

}

type listResponse struct {
	Dirs  []string
	Files []*ObjectMeta
}

func (client *rpcClient) ListObjects(ctx context.Context, prefix string) (dirs []string, files []*ObjectMeta, err error) {
	respData, err := client.get(ctx, filepath.Join("/api/v1/listAll", prefix), nil)
	if err != nil {
		return nil, nil, err
	}
	var resp listResponse
	if err := json.Unmarshal(respData, &resp); err != nil {
		err = fmt.Errorf("response unmarshal failed: %v", err)
		return nil, nil, err
	}
	return resp.Dirs, resp.Files, nil
}

func (client *rpcClient) post(ctx context.Context,
	endpoint, contentType string, r io.ReadCloser, length int64,
	headers map[string]string) ([]byte, error) {
	u, err := url.Parse(client.apiURL + endpoint)
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		Method:        "POST",
		URL:           u,
		Body:          r,
		ContentLength: length,
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	req = req.WithContext(ctx)
	resp, err := client.cli.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if len(respBody) > 0 {
			err := fmt.Errorf("error %d: %s", resp.StatusCode, respBody)
			return nil, err
		}
		err := errors.New(resp.Status)
		return nil, err
	}
	return respBody, nil
}

func (client *rpcClient) get(ctx context.Context, endpoint string,
	headers map[string]string) ([]byte, error) {
	u, err := url.Parse(client.apiURL + endpoint)
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		Method: "GET",
		URL:    u,
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	req = req.WithContext(ctx)
	resp, err := client.cli.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if len(respBody) > 0 {
			err := fmt.Errorf("error %d: %s", resp.StatusCode, respBody)
			return nil, err
		}
		err := errors.New(resp.Status)
		return nil, err
	}
	return respBody, nil
}
