package version

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getServer(version *versionService) *httptest.Server {
	mux := version.registerRoutes()
	return httptest.NewServer(mux)
}

func Test_GetVersion(t *testing.T) {
	service := &versionService{Version: "1.0", Checksum: "checksum"}
	server := getServer(service)
	defer server.Close()

	res, err := http.Get(server.URL + "/version")
	if err != nil {
		assert.Fail(t, "No error expected. Got error: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		assert.Fail(t, "Error reading response body")
	}

	actual := strings.TrimRight(string(body), "\n")
	expected := fmt.Sprintf(`{"version":"%v","checksum":"%v"}`, service.Version, service.Checksum)

	assert.Equal(t, 200, res.StatusCode, "Expected 200 status code in response")
	assert.Equal(t, expected, actual, "Actual response does not match expected response")
}

func Test_GetVersion_Cors(t *testing.T) {
	server := getServer(&versionService{})
	defer server.Close()

	client := &http.Client{}

	req, _ := http.NewRequest("OPTIONS", server.URL+"/version", nil)
	req.Header.Add("Origin", "bar.com")

	res, err := client.Do(req)
	if err != nil {
		assert.Fail(t, "No error expected. Got error: %s", err.Error())
	}

	assertHeaders(t, res.Header, map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET",
	})
}

func assertHeaders(t *testing.T, resHeaders http.Header, reqHeaders map[string]string) {
	for name, value := range reqHeaders {
		if actual := strings.Join(resHeaders[name], ", "); actual != value {
			t.Errorf("Invalid header `%s', wanted `%s', got `%s'", name, value, actual)
		}
	}
}
