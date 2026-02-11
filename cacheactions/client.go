package cacheactions

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type client struct {
	baseURL string
	token   string
	http    *http.Client
}

func newClient() (*client, error) {
	baseURL := os.Getenv("ACTIONS_CACHE_URL")
	token := os.Getenv("ACTIONS_RUNTIME_TOKEN")
	if baseURL == "" || token == "" {
		return nil, fmt.Errorf("missing ACTIONS_CACHE_URL or TOKEN")
	}
	return &client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		token:   token,
		http:    &http.Client{},
	}, nil
}

func (c *client) version(paths []string) string {
	h := sha256.New()
	for _, p := range paths {
		h.Write([]byte(p))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (c *client) request(
	method, url string, body io.Reader,
) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json;api-version=6.0-preview.1")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *client) exists(key string, paths []string) (bool, error) {
	version := c.version(paths)
	url := fmt.Sprintf(
		"%s_apis/artifactcache/cache?keys=%s&version=%s",
		c.baseURL, key, version,
	)
	req, err := c.request("GET", url, nil)
	if err != nil {
		return false, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}

func (c *client) restore(key string, paths []string) error {
	version := c.version(paths)
	url := fmt.Sprintf(
		"%s_apis/artifactcache/cache?keys=%s&version=%s",
		c.baseURL, key, version,
	)
	req, err := c.request("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return fmt.Errorf("cache not found for key: %s", key)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get cache: %s", resp.Status)
	}
	var result struct {
		ArchiveLocation string `json:"archiveLocation"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	return c.downloadAndExtract(result.ArchiveLocation)
}

func (c *client) downloadAndExtract(url string) error {
	resp, err := c.http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			mode := os.FileMode(header.Mode)
			if err := os.MkdirAll(target, mode); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			mode := os.FileMode(header.Mode)
			flags := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
			f, err := os.OpenFile(target, flags, mode)
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
	return nil
}

func (c *client) save(key string, paths []string) error {
	version := c.version(paths)
	cacheID, err := c.reserveCache(key, version)
	if err != nil {
		return err
	}
	archive, err := c.createArchive(paths)
	if err != nil {
		return err
	}
	if err := c.uploadCache(cacheID, archive); err != nil {
		return err
	}
	return c.commitCache(cacheID, int64(len(archive)))
}

func (c *client) reserveCache(key, version string) (int64, error) {
	payload := map[string]string{"key": key, "version": version}
	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s_apis/artifactcache/caches", c.baseURL)
	req, err := c.request("POST", url, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("failed to reserve cache: %s", resp.Status)
	}
	var result struct {
		CacheID int64 `json:"cacheId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	return result.CacheID, nil
}

func (c *client) createArchive(paths []string) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	for _, p := range paths {
		p = expandPath(p)
		walkFn := func(file string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(fi, "")
			if err != nil {
				return err
			}
			header.Name = file
			if err := tw.WriteHeader(header); err != nil {
				return err
			}
			if !fi.Mode().IsRegular() {
				return nil
			}
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(tw, f)
			return err
		}
		if err := filepath.Walk(p, walkFn); err != nil {
			return nil, err
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *client) uploadCache(cacheID int64, data []byte) error {
	url := fmt.Sprintf("%s_apis/artifactcache/caches/%d", c.baseURL, cacheID)
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Range", fmt.Sprintf("bytes 0-%d/*", len(data)-1))
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to upload cache: %s", resp.Status)
	}
	return nil
}

func (c *client) commitCache(cacheID int64, size int64) error {
	payload := map[string]int64{"size": size}
	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s_apis/artifactcache/caches/%d", c.baseURL, cacheID)
	req, err := c.request("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to commit cache: %s", resp.Status)
	}
	return nil
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
