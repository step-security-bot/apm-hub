package kubernetes

import (
	"bufio"
	"container/list"
	"context"
	"fmt"
	"github.com/flanksource/commons/logger"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sync"
	"time"

	"github.com/flanksource/flanksource-ui/apm-hub/api/logs"
	"github.com/flanksource/kommons"
)

type Client struct {
	*kommons.Client
}

func GetKubeClient(kommonsClient *kommons.Client, kubernetesSeachBackend *logs.KubernetesSearchBackend) (*Client, error) {
	if kubernetesSeachBackend.Kubeconfig != nil {
		if kommonsClient != nil {
			_, value, err := kommonsClient.GetEnvValue(*kubernetesSeachBackend.Kubeconfig, kubernetesSeachBackend.Namespace)
			if err != nil {
				return nil, err
			}
			kommonsClient, err = kommons.NewClientFromBytes([]byte(value))
			return &Client{kommonsClient}, err
		}
		return nil, fmt.Errorf("default client is nil and kubeconfig is not set")
	}
	return &Client{kommonsClient}, nil
}


func (c *Client) GetAllPodsForNode(nodeName string, labels map[string]string) (pods *v1.PodList, err error) {
	client, err := c.GetClientset()
	if err != nil {
		return nil, err
	}
	labelsString := GetLabelString(labels)

	if nodeName != "" {
		pods, err = client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
			FieldSelector: "spec.nodeName=" + nodeName,
			LabelSelector: labelsString,
		})
	} else {
		pods, err = client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelsString,
		})
	}
	if err != nil {
		return nil, err
	}
	if len(pods.Items) != 0 {
		return pods, nil
	}
	return nil, nil
}

// empty name will fetch all pods with the specified labels and if labels are nil will fetch the pods with the specified name
func (c *Client) GetPodsWithNameAndLabels(name, namespace string, labels map[string]string) (pods *v1.PodList, err error) {
	client, err := c.GetClientset()
	if err != nil {
		return nil, err
	}
	labelsString := GetLabelString(labels)
	if name != "" {
		pods, err = client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			FieldSelector: "metadata.name=" + name,
			LabelSelector: labelsString,
		})
	} else {
		pods, err = client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelsString,
		})
	}
	if err != nil {
		return nil, err
	}
	if len(pods.Items) != 0 {
		return pods, nil
	}
	return nil, nil
}

func (c *Client) GetPodsForDeployment(name, namespace string, labels map[string]string) (pods *v1.PodList, err error) {
	client, err := c.GetClientset()
	if err != nil {
		return nil, err
	}

	labelsString := GetLabelString(labels)

	var deployments *appsv1.DeploymentList
	if name != "" {
		deployments, err = client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelsString,
			FieldSelector: "metadata.name=" + name,
		})
	} else {
		deployments, err = client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelsString,
		})
	}
	if err != nil {
		return nil, err
	}
	var deploymentPod *v1.PodList
	pods = &v1.PodList{
		Items: []v1.Pod{},
	}
	for _, deployment := range deployments.Items {
		deploymentPod, err = client.CoreV1().Pods(deployment.GetNamespace()).List(context.TODO(), metav1.ListOptions{
			LabelSelector: GetLabelString(deployment.Spec.Template.Labels),
		})
		if err != nil {
			logger.Errorf("error fetching pod for deployment: %v; error: %v", deployment.Name, err)
			continue
		}
		pods.Items = append(pods.Items, deploymentPod.Items...)
	}
	return
}

func (c *Client) GetPodsForService(name, namespace string, labels map[string]string) (pods *v1.PodList, err error) {
	client, err := c.GetClientset()
	if err != nil {
		return nil, err
	}

	labelsString := GetLabelString(labels)

	var services *v1.ServiceList
	if name != "" {
		services, err = client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelsString,
			FieldSelector: "metadata.name=" + name,
		})
	} else {
		services, err = client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelsString,
		})
	}
	if err != nil {
		return nil, err
	}
	pods = &v1.PodList{
		Items: []v1.Pod{},
	}
	for _, service := range services.Items {
		servicePods, err := client.CoreV1().Pods(service.GetNamespace()).List(context.TODO(), metav1.ListOptions{
			LabelSelector: GetLabelString(service.Spec.Selector),
		})
		if err != nil {
			logger.Errorf("error fetching pod for deployment: %v; error: %v", service.Name, err)
			continue
		}
		pods.Items = append(pods.Items, servicePods.Items...)
	}
	return
}

func (c *Client) GetLogsForPod(pod v1.Pod) (map[string]string, error) {
	containerLogs := make(map[string]string)
	c.Tracef("Waiting for %s/%s to be running", pod.Namespace, pod.Name)
	if err := c.WaitForContainerStart(pod.Namespace, pod.Name, 5*time.Second); err != nil {
		return nil, err
	}
	client, err := c.GetClientset()
	if err != nil {
		return nil, err
	}
	pods := client.CoreV1().Pods(pod.Namespace)
	var wg sync.WaitGroup
	var containers = list.New()

	for _, container := range append(pod.Spec.Containers, pod.Spec.InitContainers...) {
		containers.PushBack(container)
	}
	// Loop over container list.
	for element := containers.Front(); element != nil; element = element.Next() {
		var logsPod string
		container := element.Value.(v1.Container)
		logs := pods.GetLogs(pod.Name, &v1.PodLogOptions{
			Container: container.Name,
		})

		podLogs, err := logs.Stream(context.TODO())
		if err != nil {
			containers.PushBack(container)
			logger.Tracef("failed to begin streaming %s/%s: %s", pod.Name, container.Name, err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		wg.Add(1)
		go func() {
			defer podLogs.Close()
			defer wg.Done()

			scanner := bufio.NewScanner(podLogs)
			for scanner.Scan() {
				incoming := scanner.Bytes()
				buffer := make([]byte, len(incoming))
				copy(buffer, incoming)
				logsPod = logsPod + fmt.Sprintf("%s\n", string(buffer))
			}
			containerLogs[container.Name] = logsPod
		}()
		fmt.Println(logsPod)
	}
	wg.Wait()
	return containerLogs, nil
}