package lib

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleMutate(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1beta1",
		"request": {
			"uid": "256d5398-3989-44f3-a886-8f2bd791cb0a",
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
	req, err := http.NewRequest("POST", "/mutate", bytes.NewBufferString(rawJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleMutate)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned status code %v want %v", status, http.StatusOK)
	}

	expected := `{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"256d5398-3989-44f3-a886-8f2bd791cb0a","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"requestKind":{"group":"","version":"v1","kind":"Pod"},"requestResource":{"group":"","version":"v1","resource":"pods"},"name":"myapp2","namespace":"gilly","operation":"CREATE","userInfo":{"username":"admin","uid":"admin","groups":["system:masters","system:authenticated"]},"object":{"kind":"Pod","apiVersion":"v1","metadata":{"name":"myapp2","namespace":"gilly","creationTimestamp":null,"labels":{"name":"myapp"},"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"labels\":{\"name\":\"myapp\"},\"name\":\"myapp2\",\"namespace\":\"gilly\"},\"spec\":{\"containers\":[{\"image\":\"gcr.io/sgryczan/rocket:latest\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"myapp\"}]}}\n"}},"spec":{"volumes":[{"name":"default-token-4rjmb","secret":{"secretName":"default-token-4rjmb"}}],"containers":[{"name":"myapp","image":"gcr.io/sgryczan/rocket:latest","resources":{},"volumeMounts":[{"name":"default-token-4rjmb","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}],"priority":0,"enableServiceLinks":true},"status":{}},"oldObject":null,"dryRun":false,"options":{"kind":"CreateOptions","apiVersion":"meta.k8s.io/v1"}},"response":{"uid":"256d5398-3989-44f3-a886-8f2bd791cb0a","allowed":true,"status":{"metadata":{},"status":"Success"},"patch":"W3sib3AiOiJyZXBsYWNlIiwicGF0aCI6Ii9zcGVjL2NvbnRhaW5lcnMvMC9pbWFnZSIsInZhbHVlIjoiZG9ja2VyLnJlcG8uZW5nLm5ldGFwcC5jb20vc2dyeWN6YW4vcm9ja2V0OmxhdGVzdCJ9LHsib3AiOiJhZGQiLCJwYXRoIjoiL21ldGFkYXRhL2Fubm90YXRpb25zL2dpbGx5LW9yaWdpbmFsLWltYWdlIiwidmFsdWUiOiJnY3IuaW8vc2dyeWN6YW4vcm9ja2V0OmxhdGVzdCJ9XQ==","patchType":"JSONPatch","auditAnnotations":{"gilly":"review complete"}}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v\n want \n%v",
			rr.Body.String(), expected)
	}

}

func TestHandleMutateInternalImage(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1beta1",
		"request": {
			"uid": "2efc8337-571a-4b05-a4b0-0ce98c3c2825",
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
	req, err := http.NewRequest("POST", "/mutate", bytes.NewBufferString(rawJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleMutate)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned status code %v want %v", status, http.StatusOK)
	}

	expected := `{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"2efc8337-571a-4b05-a4b0-0ce98c3c2825","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"requestKind":{"group":"","version":"v1","kind":"Pod"},"requestResource":{"group":"","version":"v1","resource":"pods"},"name":"test-artifactory","namespace":"gilly","operation":"CREATE","userInfo":{"username":"admin","uid":"admin","groups":["system:masters","system:authenticated"]},"object":{"kind":"Pod","apiVersion":"v1","metadata":{"name":"test-artifactory","namespace":"gilly","creationTimestamp":null,"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"name\":\"test-artifactory\",\"namespace\":\"gilly\"},\"spec\":{\"containers\":[{\"image\":\"sf-artifactory.solidfire.net:9004/pixiecore-dynamic-rom:sidecar-v0.0.7\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"test-artifactory\"}]}}\n"}},"spec":{"volumes":[{"name":"default-token-4rjmb","secret":{"secretName":"default-token-4rjmb"}}],"containers":[{"name":"test-artifactory","image":"sf-artifactory.solidfire.net:9004/pixiecore-dynamic-rom:sidecar-v0.0.7","resources":{},"volumeMounts":[{"name":"default-token-4rjmb","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}],"priority":0,"enableServiceLinks":true},"status":{}},"oldObject":null,"dryRun":false,"options":{"kind":"CreateOptions","apiVersion":"meta.k8s.io/v1"}},"response":{"uid":"2efc8337-571a-4b05-a4b0-0ce98c3c2825","allowed":true,"status":{"metadata":{},"status":"Success"},"patch":"W10=","patchType":"JSONPatch","auditAnnotations":{"gilly":"review complete"}}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v\n want \n%v",
			rr.Body.String(), expected)
	}

}

func TestHandleMutateDockerHubImage(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1beta1",
		"request": {
			"uid": "c19e8494-4800-4385-abce-6a8a69369b03",
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
	req, err := http.NewRequest("POST", "/mutate", bytes.NewBufferString(rawJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleMutate)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned status code %v want %v", status, http.StatusOK)
	}

	expected := `{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"c19e8494-4800-4385-abce-6a8a69369b03","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"requestKind":{"group":"","version":"v1","kind":"Pod"},"requestResource":{"group":"","version":"v1","resource":"pods"},"name":"test-dockerhub","namespace":"gilly","operation":"CREATE","userInfo":{"username":"admin","uid":"admin","groups":["system:masters","system:authenticated"]},"object":{"kind":"Pod","apiVersion":"v1","metadata":{"name":"test-dockerhub","namespace":"gilly","creationTimestamp":null,"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"name\":\"test-dockerhub\",\"namespace\":\"gilly\"},\"spec\":{\"containers\":[{\"image\":\"sgryczan/rocket:latest\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"test-dh\"}]}}\n"}},"spec":{"volumes":[{"name":"default-token-4rjmb","secret":{"secretName":"default-token-4rjmb"}}],"containers":[{"name":"test-dh","image":"sgryczan/rocket:latest","resources":{},"volumeMounts":[{"name":"default-token-4rjmb","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}],"priority":0,"enableServiceLinks":true},"status":{}},"oldObject":null,"dryRun":false,"options":{"kind":"CreateOptions","apiVersion":"meta.k8s.io/v1"}},"response":{"uid":"c19e8494-4800-4385-abce-6a8a69369b03","allowed":true,"status":{"metadata":{},"status":"Success"},"patch":"W10=","patchType":"JSONPatch","auditAnnotations":{"gilly":"review complete"}}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v\n want \n%v",
			rr.Body.String(), expected)
	}

}
