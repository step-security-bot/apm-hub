package kubernetes

import "strings"

func GetLabelString(labels map[string]string) string {
	labelsString := ""
	for key, value := range labels {
		labelsString = labelsString + key + "=" + value + ","
	}
	return strings.TrimSuffix(labelsString, ",")
}
