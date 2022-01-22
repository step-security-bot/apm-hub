package kubernetes

import (
	"fmt"
	"strings"

	"github.com/flanksource/commons/logger"
	"github.com/flanksource/flanksource-ui/apm-hub/api/logs"
	"github.com/flanksource/kommons"
	v1 "k8s.io/api/core/v1"
)

type KubernetesSearch struct {
	KommonsClient *kommons.Client
}

func (s *KubernetesSearch) Search(q *logs.SearchParams) (r logs.SearchResults, err error) {
	var pods *v1.PodList
	var resultLabels map[string]string
	switch {
	case strings.Contains(strings.ToLower(q.Type), "kubernetespod"):
		pods, err = s.KommonsClient.GetPodsWithNameAndLabels(q.Id, q.Labels)

	case strings.Contains(strings.ToLower(q.Type), "kubernetesnode"):
		pods, err = s.KommonsClient.GetAllPodsForNode(q.Id, q.Labels)

	case strings.Contains(strings.ToLower(q.Type), "kubernetesdeployment"):
		pods, err = s.KommonsClient.GetPodsForDeployment(q.Id, q.Labels)
		resultLabels = map[string]string{
			"deployment": q.Id,
		}
	case strings.Contains(strings.ToLower(q.Type), "kubernetesservice"):
		pods, err = s.KommonsClient.GetPodsForService(q.Id, q.Labels)
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
		podLogs, err := s.KommonsClient.GetLogsForPod(pod)
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
