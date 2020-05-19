package lib

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	v1beta1 "k8s.io/api/admission/v1beta1"
)

func TestMutatesValidRequest(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1beta1",
		"request": {
			"uid": "7f0b2891-916f-4ed6-b7cd-27bff1815a8c",
			"kind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"resource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"requestKind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"requestResource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"namespace": "yolo",
			"operation": "CREATE",
			"userInfo": {
				"username": "kubernetes-admin",
				"groups": [
					"system:masters",
					"system:authenticated"
				]
			},
			"object": {
				"kind": "Pod",
				"apiVersion": "v1",
				"metadata": {
					"name": "c7m",
					"namespace": "yolo",
					"creationTimestamp": null,
					"labels": {
						"name": "c7m"
					},
					"annotations": {
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"labels\":{\"name\":\"c7m\"},\"name\":\"c7m\",\"namespace\":\"yolo\"},\"spec\":{\"containers\":[{\"args\":[\"-c\",\"trap \\\"killall sleep\\\" TERM; trap \\\"kill -9 sleep\\\" KILL; sleep infinity\"],\"command\":[\"/bin/bash\"],\"image\":\"centos:7\",\"name\":\"c7m\"}]}}\n"
					}
				},
				"spec": {
					"volumes": [
						{
							"name": "default-token-5z7xl",
							"secret": {
								"secretName": "default-token-5z7xl"
							}
						}
					],
					"containers": [
						{
							"name": "c7m",
							"image": "centos:7",
							"command": [
								"/bin/bash"
							],
							"args": [
								"-c",
								"trap \"killall sleep\" TERM; trap \"kill -9 sleep\" KILL; sleep infinity"
							],
							"resources": {},
							"volumeMounts": [
								{
									"name": "default-token-5z7xl",
									"readOnly": true,
									"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
								}
							],
							"terminationMessagePath": "/dev/termination-log",
							"terminationMessagePolicy": "File",
							"imagePullPolicy": "IfNotPresent"
						}
					],
					"restartPolicy": "Always",
					"terminationGracePeriodSeconds": 30,
					"dnsPolicy": "ClusterFirst",
					"serviceAccountName": "default",
					"serviceAccount": "default",
					"securityContext": {},
					"schedulerName": "default-scheduler",
					"tolerations": [
						{
							"key": "node.kubernetes.io/not-ready",
							"operator": "Exists",
							"effect": "NoExecute",
							"tolerationSeconds": 300
						},
						{
							"key": "node.kubernetes.io/unreachable",
							"operator": "Exists",
							"effect": "NoExecute",
							"tolerationSeconds": 300
						}
					],
					"priority": 0,
					"enableServiceLinks": true
				},
				"status": {}
			},
			"oldObject": null,
			"dryRun": false,
			"options": {
				"kind": "CreateOptions",
				"apiVersion": "meta.k8s.io/v1"
			}
		}
	}`
	response, err := Mutate([]byte(rawJSON), false)
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(response), err)
	}

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	//assert.Equal(t, `[{"op":"replace","path":"/spec/containers/0/image","value":"debian"}]`, string(rr.Patch))
	assert.Contains(t, rr.AuditAnnotations, "gilly")

}

func TestMutatesRequestInternetRegistry(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1beta1",
		"request": {
			"uid": "7c6d0579-80c2-4a9c-bc85-a43460300906",
			"kind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"resource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"requestKind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"requestResource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"name": "myapp2",
			"namespace": "gilly",
			"operation": "CREATE",
			"userInfo": {
				"username": "admin",
				"uid": "admin",
				"groups": ["system:masters", "system:authenticated"]
			},
			"object": {
				"kind": "Pod",
				"apiVersion": "v1",
				"metadata": {
					"name": "myapp2",
					"namespace": "gilly",
					"creationTimestamp": null,
					"labels": {
						"name": "myapp"
					},
					"annotations": {
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"labels\":{\"name\":\"myapp\"},\"name\":\"myapp2\",\"namespace\":\"gilly\"},\"spec\":{\"containers\":[{\"image\":\"gcr.io/sgryczan/rocket:latest\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"myapp\"}]}}\n"
					}
				},
				"spec": {
					"volumes": [{
						"name": "default-token-4rjmb",
						"secret": {
							"secretName": "default-token-4rjmb"
						}
					}],
					"containers": [{
						"name": "myapp",
						"image": "gcr.io/sgryczan/rocket:latest",
						"resources": {},
						"volumeMounts": [{
							"name": "default-token-4rjmb",
							"readOnly": true,
							"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
						}],
						"terminationMessagePath": "/dev/termination-log",
						"terminationMessagePolicy": "File",
						"imagePullPolicy": "IfNotPresent"
					}],
					"restartPolicy": "Always",
					"terminationGracePeriodSeconds": 30,
					"dnsPolicy": "ClusterFirst",
					"serviceAccountName": "default",
					"serviceAccount": "default",
					"securityContext": {},
					"schedulerName": "default-scheduler",
					"tolerations": [{
						"key": "node.kubernetes.io/not-ready",
						"operator": "Exists",
						"effect": "NoExecute",
						"tolerationSeconds": 300
					}, {
						"key": "node.kubernetes.io/unreachable",
						"operator": "Exists",
						"effect": "NoExecute",
						"tolerationSeconds": 300
					}],
					"priority": 0,
					"enableServiceLinks": true
				},
				"status": {}
			},
			"oldObject": null,
			"dryRun": false,
			"options": {
				"kind": "CreateOptions",
				"apiVersion": "meta.k8s.io/v1"
			}
		}
	}`
	response, err := Mutate([]byte(rawJSON), false)
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(response), err)
	}

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.Equal(t, `[{"op":"replace","path":"/spec/containers/0/image","value":"docker.repo.eng.netapp.com/sgryczan/rocket:latest"},{"op":"add","path":"/metadata/annotations/gilly-original-image","value":"gcr.io/sgryczan/rocket:latest"}]`, string(rr.Patch))
	assert.Contains(t, rr.AuditAnnotations, "gilly")
}

func TestMutatesRequestNoRegistry(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1beta1",
		"request": {
			"uid": "04026241-9672-4696-93d6-e290cc783e2d",
			"kind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"resource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"requestKind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"requestResource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"name": "test-dockerhub",
			"namespace": "gilly",
			"operation": "CREATE",
			"userInfo": {
				"username": "admin",
				"uid": "admin",
				"groups": ["system:masters", "system:authenticated"]
			},
			"object": {
				"kind": "Pod",
				"apiVersion": "v1",
				"metadata": {
					"name": "test-dockerhub",
					"namespace": "gilly",
					"creationTimestamp": null,
					"annotations": {
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"name\":\"test-dockerhub\",\"namespace\":\"gilly\"},\"spec\":{\"containers\":[{\"image\":\"sgryczan/rocket:latest\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"test-dh\"}]}}\n"
					}
				},
				"spec": {
					"volumes": [{
						"name": "default-token-4rjmb",
						"secret": {
							"secretName": "default-token-4rjmb"
						}
					}],
					"containers": [{
						"name": "test-dh",
						"image": "sgryczan/rocket:latest",
						"resources": {},
						"volumeMounts": [{
							"name": "default-token-4rjmb",
							"readOnly": true,
							"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
						}],
						"terminationMessagePath": "/dev/termination-log",
						"terminationMessagePolicy": "File",
						"imagePullPolicy": "IfNotPresent"
					}],
					"restartPolicy": "Always",
					"terminationGracePeriodSeconds": 30,
					"dnsPolicy": "ClusterFirst",
					"serviceAccountName": "default",
					"serviceAccount": "default",
					"securityContext": {},
					"schedulerName": "default-scheduler",
					"tolerations": [{
						"key": "node.kubernetes.io/not-ready",
						"operator": "Exists",
						"effect": "NoExecute",
						"tolerationSeconds": 300
					}, {
						"key": "node.kubernetes.io/unreachable",
						"operator": "Exists",
						"effect": "NoExecute",
						"tolerationSeconds": 300
					}],
					"priority": 0,
					"enableServiceLinks": true
				},
				"status": {}
			},
			"oldObject": null,
			"dryRun": false,
			"options": {
				"kind": "CreateOptions",
				"apiVersion": "meta.k8s.io/v1"
			}
		}
	}`
	response, err := Mutate([]byte(rawJSON), false)
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(response), err)
	}

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.Equal(t, "[]", string(rr.Patch))
	assert.Contains(t, rr.AuditAnnotations, "gilly")
}

func TestMutatesRequestInternalRegistry(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1beta1",
		"request": {
			"uid": "2066f94b-9d89-4e2d-8795-d56429ceddc5",
			"kind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"resource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"requestKind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"requestResource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"name": "test-artifactory",
			"namespace": "gilly",
			"operation": "CREATE",
			"userInfo": {
				"username": "admin",
				"uid": "admin",
				"groups": ["system:masters", "system:authenticated"]
			},
			"object": {
				"kind": "Pod",
				"apiVersion": "v1",
				"metadata": {
					"name": "test-artifactory",
					"namespace": "gilly",
					"creationTimestamp": null,
					"annotations": {
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"name\":\"test-artifactory\",\"namespace\":\"gilly\"},\"spec\":{\"containers\":[{\"image\":\"sf-artifactory.solidfire.net:9004/pixiecore-dynamic-rom:sidecar-v0.0.7\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"test-artifactory\"}]}}\n"
					}
				},
				"spec": {
					"volumes": [{
						"name": "default-token-4rjmb",
						"secret": {
							"secretName": "default-token-4rjmb"
						}
					}],
					"containers": [{
						"name": "test-artifactory",
						"image": "sf-artifactory.solidfire.net:9004/pixiecore-dynamic-rom:sidecar-v0.0.7",
						"resources": {},
						"volumeMounts": [{
							"name": "default-token-4rjmb",
							"readOnly": true,
							"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
						}],
						"terminationMessagePath": "/dev/termination-log",
						"terminationMessagePolicy": "File",
						"imagePullPolicy": "IfNotPresent"
					}],
					"restartPolicy": "Always",
					"terminationGracePeriodSeconds": 30,
					"dnsPolicy": "ClusterFirst",
					"serviceAccountName": "default",
					"serviceAccount": "default",
					"securityContext": {},
					"schedulerName": "default-scheduler",
					"tolerations": [{
						"key": "node.kubernetes.io/not-ready",
						"operator": "Exists",
						"effect": "NoExecute",
						"tolerationSeconds": 300
					}, {
						"key": "node.kubernetes.io/unreachable",
						"operator": "Exists",
						"effect": "NoExecute",
						"tolerationSeconds": 300
					}],
					"priority": 0,
					"enableServiceLinks": true
				},
				"status": {}
			},
			"oldObject": null,
			"dryRun": false,
			"options": {
				"kind": "CreateOptions",
				"apiVersion": "meta.k8s.io/v1"
			}
		}
	}`
	response, err := Mutate([]byte(rawJSON), false)
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(response), err)
	}

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.Equal(t, "[]", string(rr.Patch))
	assert.Contains(t, rr.AuditAnnotations, "gilly")
}

func TestErrorsOnInvalidJson(t *testing.T) {
	rawJSON := `Wut ?`
	_, err := Mutate([]byte(rawJSON), false)
	if err == nil {
		t.Error("did not fail when sending invalid json")
	}
}

func TestErrorsOnInvalidPod(t *testing.T) {
	rawJSON := `{
		"request": {
			"object": 111
		}
	}`
	_, err := Mutate([]byte(rawJSON), false)
	if err == nil {
		t.Error("did not fail when sending invalid pod")
	}
}

func TestGetImageRegistry(t *testing.T) {
	cases := []struct{ input, expected string }{
		{"gcr.io/sgryczan/hello-world:0.0.0", "gcr.io"},
		{"quay.io/some-owner/image-name", "quay.io"},
		{"docker.io/iwilltry42/k3d-tools", "docker.io"},
		{"sgryczan/rocket", "docker.io"},
		{"busybox", "docker.io"},
	}
	for _, test := range cases {
		if r := GetImageRegistry(test.input); r != test.expected {
			t.Fatalf("GetImageRegistry(%s) = %s, want %s", test.input, r, test.expected)
		}
	}
}

func TestReplaceImageRegistry(t *testing.T) {
	cases := []struct{ inputImage, inputRegistry, expected string }{
		{"gcr.io/sgryczan/hello-world:0.0.0", "docker.repo.eng.netapp.com", "docker.repo.eng.netapp.com/sgryczan/hello-world:0.0.0"},
		{"sf-artifactory.solidfire.net:9004/pixiecore-dynamic-rom:0.0.10", "gcr.io", "gcr.io/pixiecore-dynamic-rom:0.0.10"},
	}
	for _, test := range cases {
		if r, err := ReplaceImageRegistry(test.inputImage, test.inputRegistry); r != test.expected || err != nil {
			if err != nil {
				t.Fatalf("ReplaceImageRegistry(%s, %s) = %s", test.inputImage, test.inputRegistry, err)
			}
			t.Fatalf("ReplaceImageRegistry(%s, %s) = %s, want %s", test.inputImage, test.inputRegistry, r, test.expected)
		}
	}
}
