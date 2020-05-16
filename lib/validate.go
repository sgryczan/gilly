package lib

import (
	"log"
)

// GetBrokenContainers checks a PodList for pods that cannot be scheduled
func GetBrokenContainers(plist *PodList) (*[]BrokenContainer, error) {
	brokenContainers := []BrokenContainer{}
	for _, pod := range plist.Items {
		if pod.Status.Phase == "Pending" {
			log.Printf("[Validate]  Detected pod %s in Pending state.", pod.Metadata.Name)
			for _, status := range pod.Status.ContainerStatuses {
				if val, exists := status.State["waiting"]; exists {
					log.Printf("[Validate]  Detected container %s in pod %s status is %s", status.Name, pod.Metadata.Name, val.Reason)
					if val.Reason == "ImagePullBackOff" || val.Reason == "ErrImagePull" {
						log.Printf("[Validate]  Detected container %s in pod %s reason: %s", status.Name, pod.Metadata.Name, val.Reason)
						c := BrokenContainer{
							Name:  status.Name,
							Image: status.Image,
							Pod:   pod.Metadata.Name,
						}
						brokenContainers = append(brokenContainers, c)
					}
				}
			}
		} else {
			continue
		}
	}
	log.Printf("[Validate]  Checked %d Pods, found %d waiting", len(plist.Items), len(brokenContainers))
	return &brokenContainers, nil
}
