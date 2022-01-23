package kubernetes

import (
	"fmt"
	"strings"

	"github.com/flanksource/commons/logger"
	"github.com/flanksource/flanksource-ui/apm-hub/api/logs"
	v1 "k8s.io/api/core/v1"
)

type KubernetesSearch struct {
	Client *Client
}

func (s *KubernetesSearch) Search(q *logs.SearchParams) (r logs.SearchResults, err error) {
	var pods *v1.PodList
	var resultLabels map[string]string
	namespace, name := s.GetNameNamespace(q)
	fmt.Println(namespace)
	fmt.Println(name)
	switch {
	case strings.Contains(strings.ToLower(q.Type), "kubernetespod"):
		pods, err = s.Client.GetPodsWithNameAndLabels(name, namespace, q.Labels)

	case strings.Contains(strings.ToLower(q.Type), "kubernetesnode"):
		pods, err = s.Client.GetAllPodsForNode(q.Id, q.Labels)

	case strings.Contains(strings.ToLower(q.Type), "kubernetesdeployment"):
		pods, err = s.Client.GetPodsForDeployment(name, namespace, q.Labels)
		resultLabels = map[string]string{
			"deployment": q.Id,
		}
	case strings.Contains(strings.ToLower(q.Type), "kubernetesservice"):
		pods, err = s.Client.GetPodsForService(name, namespace, q.Labels)
		resultLabels = map[string]string{
			"service": q.Id,
		}
	}
	if err != nil {
		return r, fmt.Errorf("error fetching the pods for node %v: %v", q.Id, err)
	}
	r.Results = s.getLogResultsForPods(pods, resultLabels)
	r.Total = len(r.Results)
	return
}

func (s *KubernetesSearch) getLogResultsForPods(pods *v1.PodList, resultLabels map[string]string) []logs.Result {
	var results []logs.Result
	for _, pod := range pods.Items {
		podLogs, err := s.Client.GetLogsForPod(pod)
		if err != nil {
			logger.Errorf("error fetching logs for pod: %v in namespace: %v, err: ", pod.Name, pod.Namespace, err)
			continue
		}
		for containerName, containerLogs := range podLogs {
			var labels = map[string]string{
				"pod":           pod.Name,
				"containerName": containerName,
				"nodeName":      pod.Spec.NodeName,
				"namespace":     pod.Namespace,
			}
			for k, v := range resultLabels {
				labels[k] = v
			}
			results = append(results, logs.Result{
				Id:      pod.Name,
				Message: containerLogs,
				Labels:  labels,
			})
		}
	}
	return results
}


func (s *KubernetesSearch) GetNameNamespace(q *logs.SearchParams) (namespace, name string) {
	if strings.Contains(q.Id, "/"){
		// namespace is provided as a prefix in the ID
		namespaceName := strings.Split(q.Id, "/")
		if len(namespaceName) < 2 {
			logger.Errorf("expected id in format <namespace>/<name>")
			return "", ""
		}
		return namespaceName[0], namespaceName[1]
	}
	// namespace is provided in the labels. if no label is there we just return the empty string which extends the search to all namespaces
	namespace = q.Labels["namespace"]
	// deleting namespace label from the map so it doesn't filter out the result based on the namespace label
	delete(q.Labels, "namespace")
	return q.Id, namespace
}
