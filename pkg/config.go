package pkg

import (
	"os"
	"path/filepath"

	"github.com/flanksource/flanksource-ui/apm-hub/api/logs"
	"github.com/flanksource/flanksource-ui/apm-hub/pkg/files"
	k8s "github.com/flanksource/flanksource-ui/apm-hub/pkg/kubernetes"
	"github.com/flanksource/kommons"
	"gopkg.in/yaml.v3"
)

// Initilize all the backends mentioned in the config
// GlobalBackends, error
func ParseConfig(kommonsClient *kommons.Client, configFile string) ([]logs.SearchBackend, error) {
	searchConfig := &logs.SearchConfig{}
	backends := []logs.SearchBackend{}
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, searchConfig); err != nil {
		return nil, err
	}

	for _, backend := range searchConfig.Backends {
		if backend.Kubernetes != nil {
			client, err := k8s.GetKubeClient(kommonsClient, backend.Kubernetes)
			if err != nil {
				return nil, err
			}
			backend.Backend = &k8s.KubernetesSearch{
				Client: client,
			}
			backends = append(backends, backend)
		}

		if len(backend.Files) != 0 {
			// If the paths are not absolute,
			// They should be parsed with respect to the provided config file
			for i, f := range backend.Files {
				for j, p := range f.Paths {
					if !filepath.IsAbs(p) {
						backend.Files[i].Paths[j] = filepath.Join(filepath.Dir(configFile), p)
					}
				}
			}

			backend.Backend = &files.FileSearch{
				FilesBackend: backend.Files,
			}
			backends = append(backends, backend)
		}
	}

	return backends, nil
}
