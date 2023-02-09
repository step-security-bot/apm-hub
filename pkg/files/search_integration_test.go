//go:build integration

package files_test

import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/flanksource/flanksource-ui/apm-hub/api/logs"
	"github.com/flanksource/flanksource-ui/apm-hub/cmd"
	"github.com/flanksource/flanksource-ui/apm-hub/pkg"
	"github.com/stretchr/testify/assert"
)

func TestFileSearch(t *testing.T) {
	confPath := "../../samples/config-file.yaml"
	backend, err := pkg.ParseConfig(nil, confPath)
	if err != nil {
		t.Fatal("Fail to parse the config file", err)
	}
	logs.GlobalBackends = append(logs.GlobalBackends, backend...)

	sp := logs.SearchParams{
		Labels: map[string]string{
			"name": "acmehost",
			"type": "Nginx",
		},
	}
	b, err := json.Marshal(sp)
	if err != nil {
		t.Fatal("Fail to marshal search param")
	}

	req := httptest.NewRequest(http.MethodPost, "/search", strings.NewReader(string(b)))
	req.Header.Add("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e := cmd.SetupServer(nil)
	e.ServeHTTP(rec, req)

	var res logs.SearchResults
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatal("Failed to decode the search result")
	}

	filePath := "../../samples/nginx-access.log"
	nginxLogFile, err := os.Open(filePath)
	if err != nil {
		t.Fatal("Fail to read nginx log", err)
	}

	scanner := bufio.NewScanner(nginxLogFile)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(res.Results) != len(lines) {
		t.Fatalf("Expected [%d] lines but got [%d]", len(lines), len(res.Results))
	}

	for i, r := range res.Results {
		if r.Message != lines[i] {
			t.Fatalf("Incorrect line [%d]. Expected %s got %s", i+1, lines[i], r.Message)
		}

		assert.Equal(t, r.Labels, map[string]string{
			"filepath": filePath,
			"name":     "acmehost",
			"type":     "Nginx",
		})
	}
}
