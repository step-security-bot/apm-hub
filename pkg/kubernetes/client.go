package kubernetes

import (
	"fmt"

	"github.com/flanksource/flanksource-ui/apm-hub/api/logs"
	"github.com/flanksource/kommons"
)

func GetClient(kommonsClient *kommons.Client, kubernetesSeachBackend *logs.KubernetesSearchBackend) (*kommons.Client, error) {
	if kubernetesSeachBackend.Kubeconfig != nil {
		if kommonsClient != nil {
			_, value, err := kommonsClient.GetEnvValue(*kubernetesSeachBackend.Kubeconfig, kubernetesSeachBackend.Namespace)
			if err != nil {
				return nil, err
			}
			kommonsClient, err = kommons.NewClientFromBytes([]byte(value))
			return kommonsClient, err
		}
		return nil, fmt.Errorf("default client is nil and kubeconfig is not set")
	}
	return kommonsClient, nil
}
