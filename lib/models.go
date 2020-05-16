package lib

// KubeletConn represents a connection to a Kubelet
type KubeletConn struct {
	EndPoint string
	Token    string
}

// PodList is the response to the kubelet /pods endpoint
type PodList struct {
	Kind       string            `json:"kind"`
	APIVersion string            `json:"apiVersion"`
	MetaData   map[string]string `json:"metadata"`
	Items      []ScheduledPod    `json:"items"`
}

// ScheduledPod represents .. a scheduled pod
type ScheduledPod struct {
	Metadata struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`

	Spec struct {
		Containers []struct {
			Name            string `json:"name"`
			Image           string `json:"image"`
			ImagePullPolicy string `json:"imagePullPolicy"`
		} `json:"containers"`
	}

	Status struct {
		Phase             string `json:"phase"` // e.g. "Pending"
		ContainerStatuses []struct {
			Name  string `json:"name"`
			State map[string]struct {
				Reason  string `json:"reason"`  // e.g. "ImagePullBackOff" or "ErrImagePull"
				Message string `json:"message"` // e.g. "Back-off pulling image \"docker:IDONTEXIST\""
			} `json:"state"`
			Image string `json:"image"`
		} `json:"containerStatuses"`
	} `json:"status"`
}

// BrokenPod represents an unschedulable pod
type BrokenPod struct {
	Name       string
	Status     string
	Containers []string
}

// BrokenContainer represents an unschedulable container
type BrokenContainer struct {
	Name  string
	Image string
	Pod   string
}
