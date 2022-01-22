package logs

import "github.com/flanksource/kommons"

var GlobalBackends []SearchBackend

type SearchConfig struct {
	Backends []SearchBackend `yaml:"backends,omitempty"`
}

type SearchBackend struct {
	Routes     []SearchRoute            `json:"routes,omitempty"`
	Backend    SearchAPI                `json:"-"`
	Kubernetes *KubernetesSearchBackend `json:"kubernetes,omitempty"`
}

type SearchRoute struct {
	Type     string            `json:"type,omitempty"`
	IdPrefix string            `json:"idPrefix,omitempty"`
	Labels   map[string]string `json:"labels,omitempty"`
}

type KubernetesSearchBackend struct {
	// empty kubeconfig indicates to use the current kubeconfig for connection
	Kubeconfig *kommons.EnvVar `json:"kubeconfig,omitempty"`
	//namespace to search the kommons.EnvVar in
	Namespace string `json:"namespace,omitempty"`
}

type SearchParams struct {
	// Limit is the maximum number of results to return.
	Limit int `json:"limit,omitempty"`
	// The page token, returned by a previous call, to request the next page of results.
	Page string `json:"page,omitempty"`
	// comma separated list of labels to filter the results. key1=value1,key2=value2
	Labels map[string]string `json:"labels,omitempty"`
	// A generic query string, that is rewritten to the underlying system,
	// If the underlying system does not support queries, than this query is applied on the returned results
	Query string `json:"query,omitempty"`
	// A RFC3339 timestamp or an age string (e.g. "1h", "2d", "1w")
	Start string `json:"start,omitempty"`
	// A RFC3339 timestamp or an age string (e.g. "1h", "2d", "1w")
	End string `json:"end,omitempty"`
	// The type of logs to find, e.g. KubernetesNode, KubernetesService, KubernetesPod, VM, etc. Type and ID are used to route search requests
	Type string `json:"type,omitempty"`
	// The identifier of the type of logs to find, e.g. k8s-node-1, k8s-service-1, k8s-pod-1, vm-1, etc.
	// The ID should include include any cluster/namespace/account information required for routing
	Id string `json:"id,omitempty"`
}

type SearchResults struct {
	Total    int      `json:"total,omitempty"`
	Results  []Result `json:"results,omitempty"`
	NextPage string   `json:"nextPage,omitempty"`
}

type Result struct {
	// Id is the unique identifier provided by the underlying system, use to link to a point in time of a log stream
	Id string `json:"id,omitempty"`
	// RFC3339 timestamp
	Time    string            `json:"timestamp,omitempty"`
	Message string            `json:"message,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type SearchAPI interface {
	Search(q *SearchParams) (r SearchResults, err error)
}

type SearchMapper interface {
	MapSearchParams(p *SearchParams) ([]SearchParams, error)
}
