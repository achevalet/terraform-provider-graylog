package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shk3bq4d/terraform-provider-graylog/graylog/config"
)

// serverMajor reports the running Graylog's major version, which decides the
// shape of requests that changed between releases.
func serverMajor(cfg config.Config) (int, error) {
	req, err := http.NewRequest(http.MethodGet, strings.TrimSuffix(cfg.Endpoint, "/")+"/system", nil)
	if err != nil {
		return 0, err
	}
	req.SetBasicAuth(cfg.AuthName, cfg.AuthPassword)
	req.Header.Set("X-Requested-By", "terraform-provider-graylog")

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusMultipleChoices {
		return 0, fmt.Errorf("GET /system returned %d", resp.StatusCode)
	}

	var body struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, err
	}
	major, err := strconv.Atoi(strings.SplitN(body.Version, ".", 2)[0])
	if err != nil {
		return 0, fmt.Errorf("unexpected version %q", body.Version)
	}
	return major, nil
}
