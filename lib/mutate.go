// Package lib deals with AdmissionReview requests and responses, it takes in the request body and returns a readily converted JSON []byte that can be
// returned from a http Handler w/o needing to further convert or modify it, it also makes testing Mutate() kind of easy w/o need for a fake http server, etc.
package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	v1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Mutate mutates
func Mutate(body []byte, verbose bool) ([]byte, error) {
	if verbose {
		log.Printf("recv: %s\n", string(body)) // untested section
	}

	// unmarshal request into AdmissionReview struct
	admReview := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}

	var err error
	var pod *corev1.Pod

	responseBody := []byte{}
	ar := admReview.Request
	resp := v1beta1.AdmissionResponse{}

	if ar != nil {

		if err := json.Unmarshal(ar.Object.Raw, &pod); err != nil {
			return nil, fmt.Errorf("unable unmarshal pod json object %v", err)
		}

		log.Printf("[Mutate]  Received POD create event. Name: %s, Namespace: %s", pod.ObjectMeta.Name, pod.ObjectMeta.Namespace)
		resp.Allowed = true
		resp.UID = ar.UID
		pT := v1beta1.PatchTypeJSONPatch
		resp.PatchType = &pT

		resp.AuditAnnotations = map[string]string{
			"gilly": "review complete",
		}

		p := []map[string]string{}

		for i, c := range pod.Spec.Containers {

			registry := GetImageRegistry(c.Image)
			log.Printf("[Mutate]  Found registry => %s", registry)
			if !(strings.Contains(registry, "sf-artifactory.solidfire.net")) {
				log.Printf("[Mutate] image registry for container %s is %s - updating", c.Name, registry)
				patchedRegistry, _ := ReplaceImageRegistry(c.Image, "docker.repo.eng.netapp.com")
				imagePatch := map[string]string{
					"op":    "replace",
					"path":  fmt.Sprintf("/spec/containers/%d/image", i),
					"value": patchedRegistry,
				}
				p = append(p, imagePatch)

				annotationPatch := map[string]string{
					"op":    "add",
					"path":  "/metadata/annotations/gilly-original-image",
					"value": c.Image,
				}
				p = append(p, annotationPatch)
				log.Printf("[Mutate] updated registry for container %s to %s - updating", c.Name, patchedRegistry)
			}
		}

		for i, c := range pod.Spec.InitContainers {

			registry := GetImageRegistry(c.Image)
			log.Printf("[Mutate]  Found registry => %s", registry)
			if !(strings.Contains(registry, "sf-artifactory.solidfire.net")) {
				log.Printf("[Mutate] image registry for container %s is %s - updating", c.Name, registry)
				patchedRegistry, _ := ReplaceImageRegistry(c.Image, "docker.repo.eng.netapp.com")
				imagePatch := map[string]string{
					"op":    "replace",
					"path":  fmt.Sprintf("/spec/containers/%d/image", i),
					"value": patchedRegistry,
				}
				p = append(p, imagePatch)

				annotationPatch := map[string]string{
					"op":    "add",
					"path":  "/metadata/annotations/gilly-init-original-image",
					"value": c.Image,
				}
				p = append(p, annotationPatch)
				log.Printf("[Mutate] updated registry for container %s to %s - updating", c.Name, patchedRegistry)
			}
		}

		resp.Patch, err = json.Marshal(p)

		resp.Result = &metav1.Status{
			Status: "Success",
		}

		admReview.Response = &resp
		responseBody, err = json.Marshal(admReview)
		if err != nil {
			return nil, err
		}
	}

	if verbose {
		log.Printf("resp: %s\n", string(responseBody))
	}

	return responseBody, nil
}

func checkContainer(c *corev1.Container) {
	return
}

// GetImageRegistry accepts an image name and returns it's registry
func GetImageRegistry(image string) string {
	// gcr.io/linkerd-io/image:tag
	chunk := strings.Split(image, "/")
	if strings.Contains(chunk[0], ".") {
		return chunk[0]
	}
	return "docker.io"
}

// ReplaceImageRegistry accepts an image name and returns it's registry
func ReplaceImageRegistry(image string, registry string) (string, error) {
	log.Printf("[ReplaceImageRegistry]  Replacing registry for image %s with %s", image, registry)
	// gcr.io/linkerd-io/image:tag
	chunk := strings.Split(image, "/")
	//log.Printf("[ReplaceImageRegistry]  got chunks: %v", chunk)
	if !(strings.Contains(chunk[0], ".")) {
		fmt.Printf("[ReplaceImageRegistry]  Image %s uses inferred registry", image)
		// Image was passed with no registry e.g. sgryczan/hello-world:0.0.0
		return registry + "/" + image, nil
	}
	if len(chunk) > 1 {
		// Image was passed with a registry e.g. gcr.io/sgryczan/hello-world:0.0.0
		result := registry + "/" + strings.Join(chunk[1:], "/")
		log.Printf("[ReplaceImageRegistry]  Image %s => %s", image, result)
		return result, nil
	}
	result := registry + "/" + chunk[0]
	log.Printf("[ReplaceImageRegistry]  Image %s => %s", image, result)
	return result, nil
}
