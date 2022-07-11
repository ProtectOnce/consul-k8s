package read

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed test_config_dump.json
var fs embed.FS

const testConfigDump string = "test_config_dump.json"

func TestUnmarshaling(t *testing.T) {
	raw, err := fs.ReadFile(testConfigDump)
	require.NoError(t, err)

	var envoyConfig EnvoyConfig
	err = json.Unmarshal(raw, &envoyConfig)
	require.NoError(t, err)

	require.Equal(t, testEnvoyConfig.Clusters, envoyConfig.Clusters)
	require.Equal(t, testEnvoyConfig.Endpoints, envoyConfig.Endpoints)
	require.Equal(t, testEnvoyConfig.Listeners, envoyConfig.Listeners)
	require.Equal(t, testEnvoyConfig.Routes, envoyConfig.Routes)
	require.Equal(t, testEnvoyConfig.Secrets, envoyConfig.Secrets)
}

func TestJSON(t *testing.T) {
	raw, err := fs.ReadFile(testConfigDump)
	require.NoError(t, err)
	expected := bytes.TrimSpace(raw)

	var envoyConfig EnvoyConfig
	err = json.Unmarshal(raw, &envoyConfig)
	require.NoError(t, err)

	actual := envoyConfig.JSON()

	require.Equal(t, expected, actual)
}

func TestFetchConfig(t *testing.T) {
	configResponse, err := fs.ReadFile(testConfigDump)
	require.NoError(t, err)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(configResponse)
	}))
	defer mockServer.Close()

	mpf := &mockPortForwarder{
		openBehavior: func(ctx context.Context) (string, error) {
			return strings.Replace(mockServer.URL, "http://", "", 1), nil
		},
	}

	envoyConfig, err := FetchConfig(context.Background(), mpf)

	require.NoError(t, err)

	require.Equal(t, testEnvoyConfig.Clusters, envoyConfig.Clusters)
	require.Equal(t, testEnvoyConfig.Endpoints, envoyConfig.Endpoints)
	require.Equal(t, testEnvoyConfig.Listeners, envoyConfig.Listeners)
	require.Equal(t, testEnvoyConfig.Routes, envoyConfig.Routes)
	require.Equal(t, testEnvoyConfig.Secrets, envoyConfig.Secrets)
}

func TestClusterFiltering(t *testing.T) {
	cases := map[string]struct {
		fqdnFilter string
		expected   []Cluster
	}{
		"No Filter": {
			fqdnFilter: "",
			expected: []Cluster{
				{
					Name:                     "local_agent",
					FullyQualifiedDomainName: "local_agent",
					Endpoints:                []string{"192.168.79.187:8502"},
					Type:                     "STATIC",
					LastUpdated:              "2022-05-13T04:22:39.553Z",
				},
				{
					Name:                     "client",
					FullyQualifiedDomainName: "client.default.dc1.internal.bc3815c2-1a0f-f3ff-a2e9-20d791f08d00.consul",
					Type:                     "EDS",
					LastUpdated:              "2022-06-09T00:39:12.948Z",
				},
				{
					Name:                     "frontend",
					FullyQualifiedDomainName: "frontend.default.dc1.internal.bc3815c2-1a0f-f3ff-a2e9-20d791f08d00.consul",
					Type:                     "EDS",
					LastUpdated:              "2022-06-09T00:39:12.855Z",
				},
				{
					Name:                     "local_app",
					FullyQualifiedDomainName: "local_app",
					Endpoints:                []string{"127.0.0.1:8080"},
					Type:                     "STATIC",
					LastUpdated:              "2022-05-13T04:22:39.655Z",
				},
				{
					Name:                     "original-destination",
					FullyQualifiedDomainName: "original-destination",
					Type:                     "ORIGINAL_DST",
					LastUpdated:              "2022-05-13T04:22:39.743Z",
				},
				{
					Name:                     "server",
					FullyQualifiedDomainName: "server.default.dc1.internal.bc3815c2-1a0f-f3ff-a2e9-20d791f08d00.consul",
					Type:                     "EDS",
					LastUpdated:              "2022-06-09T00:39:12.754Z",
				},
			},
		},
		"Filter FQDN by default": {
			fqdnFilter: "default",
			expected: []Cluster{
				{
					Name:                     "client",
					FullyQualifiedDomainName: "client.default.dc1.internal.bc3815c2-1a0f-f3ff-a2e9-20d791f08d00.consul",
					Type:                     "EDS",
					LastUpdated:              "2022-06-09T00:39:12.948Z",
				},
				{
					Name:                     "frontend",
					FullyQualifiedDomainName: "frontend.default.dc1.internal.bc3815c2-1a0f-f3ff-a2e9-20d791f08d00.consul",
					Type:                     "EDS",
					LastUpdated:              "2022-06-09T00:39:12.855Z",
				},
				{
					Name:                     "server",
					FullyQualifiedDomainName: "server.default.dc1.internal.bc3815c2-1a0f-f3ff-a2e9-20d791f08d00.consul",
					Type:                     "EDS",
					LastUpdated:              "2022-06-09T00:39:12.754Z",
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual := testEnvoyConfig.Clusters(tc.fqdnFilter)
			require.Equal(t, tc.expected, actual)
		})
	}
}

type mockPortForwarder struct {
	openBehavior func(context.Context) (string, error)
}

func (m *mockPortForwarder) Open(ctx context.Context) (string, error) { return m.openBehavior(ctx) }
func (m *mockPortForwarder) Close()                                   {}

// testEnvoyConfig is what we expect the config at `test_config_dump.json` to be.
var testEnvoyConfig = &EnvoyConfig{
	clusters: []Cluster{
		{
			Name:                     "local_agent",
			FullyQualifiedDomainName: "local_agent",
			Endpoints:                []string{"192.168.79.187:8502"},
			Type:                     "STATIC",
			LastUpdated:              "2022-05-13T04:22:39.553Z",
		},
		{
			Name:                     "client",
			FullyQualifiedDomainName: "client.default.dc1.internal.bc3815c2-1a0f-f3ff-a2e9-20d791f08d00.consul",
			Type:                     "EDS",
			LastUpdated:              "2022-06-09T00:39:12.948Z",
		},
		{
			Name:                     "frontend",
			FullyQualifiedDomainName: "frontend.default.dc1.internal.bc3815c2-1a0f-f3ff-a2e9-20d791f08d00.consul",
			Type:                     "EDS",
			LastUpdated:              "2022-06-09T00:39:12.855Z",
		},
		{
			Name:                     "local_app",
			FullyQualifiedDomainName: "local_app",
			Endpoints:                []string{"127.0.0.1:8080"},
			Type:                     "STATIC",
			LastUpdated:              "2022-05-13T04:22:39.655Z",
		},
		{
			Name:                     "original-destination",
			FullyQualifiedDomainName: "original-destination",
			Type:                     "ORIGINAL_DST",
			LastUpdated:              "2022-05-13T04:22:39.743Z",
		},
		{
			Name:                     "server",
			FullyQualifiedDomainName: "server.default.dc1.internal.bc3815c2-1a0f-f3ff-a2e9-20d791f08d00.consul",
			Type:                     "EDS",
			LastUpdated:              "2022-06-09T00:39:12.754Z",
		},
	},
	Endpoints: []Endpoint{
		{
			Address: "192.168.79.187:8502",
			Cluster: "local_agent",
			Weight:  1,
			Status:  "HEALTHY",
		},
		{
			Address: "127.0.0.1:8080",
			Cluster: "local_app",
			Weight:  1,
			Status:  "HEALTHY",
		},
		{
			Address: "192.168.31.201:20000",
			Weight:  1,
			Status:  "HEALTHY",
		},
		{
			Address: "192.168.47.235:20000",
			Weight:  1,
			Status:  "HEALTHY",
		},
		{
			Address: "192.168.71.254:20000",
			Weight:  1,
			Status:  "HEALTHY",
		},
		{
			Address: "192.168.63.120:20000",
			Weight:  1,
			Status:  "HEALTHY",
		},
		{
			Address: "192.168.18.110:20000",
			Weight:  1,
			Status:  "HEALTHY",
		},
		{
			Address: "192.168.52.101:20000",
			Weight:  1,
			Status:  "HEALTHY",
		},
		{
			Address: "192.168.65.131:20000",
			Weight:  1,
			Status:  "HEALTHY",
		},
	},
	Listeners: []Listener{
		{
			Name:    "public_listener",
			Address: "192.168.69.179:20000",
			FilterChain: []FilterChain{
				{
					FilterChainMatch: "Any",
					Filters:          []string{"* -> local_app/"},
				},
			},
			Direction:   "INBOUND",
			LastUpdated: "2022-06-09T00:39:27.668Z",
		},
		{
			Name:    "outbound_listener",
			Address: "127.0.0.1:15001",
			FilterChain: []FilterChain{
				{
					FilterChainMatch: "10.100.134.173/32, 240.0.0.3/32",
					Filters:          []string{"-> client.default.dc1.internal.bc3815c2-1a0f-f3ff-a2e9-20d791f08d00.consul"},
				},
				{
					FilterChainMatch: "10.100.254.176/32, 240.0.0.4/32",
					Filters:          []string{"* -> server.default.dc1.internal.bc3815c2-1a0f-f3ff-a2e9-20d791f08d00.consul/"},
				},
				{
					FilterChainMatch: "10.100.31.2/32, 240.0.0.2/32",
					Filters: []string{
						"-> frontend.default.dc1.internal.bc3815c2-1a0f-f3ff-a2e9-20d791f08d00.consul",
					},
				},
				{
					FilterChainMatch: "Any",
					Filters:          []string{"-> original-destination"},
				},
			},
			Direction:   "OUTBOUND",
			LastUpdated: "2022-05-24T17:41:59.079Z",
		},
	},
	Routes: []Route{
		{
			Name:               "public_listener",
			DestinationCluster: "local_app/",
			LastUpdated:        "2022-06-09T00:39:27.667Z",
		},
		{
			Name:               "server",
			DestinationCluster: "server.default.dc1.internal.bc3815c2-1a0f-f3ff-a2e9-20d791f08d00.consul/",
			LastUpdated:        "2022-05-24T17:41:59.078Z",
		},
	},
	Secrets: []Secret{
		{
			Name:        "default",
			Type:        "Dynamic Active",
			LastUpdated: "2022-05-24T17:41:59.078Z",
		},
		{
			Name:        "ROOTCA",
			Type:        "Dynamic Warming",
			LastUpdated: "2022-03-15T05:14:22.868Z",
		},
	},
}
