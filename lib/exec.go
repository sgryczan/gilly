package lib

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

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
	fmt.Printf("[ReplaceImageRegistry]  Replacing registry for image %s with %s", image, registry)
	// gcr.io/linkerd-io/image:tag
	chunk := strings.Split(image, "/")
	log.Printf("[ReplaceImageRegistry]  got chunks: %v", chunk)
	if !(strings.Contains(chunk[0], ".")) {
		fmt.Printf("[ReplaceImageRegistry]  Image %s uses inferred registry", image)
		// Image was passed with no registry e.g. sgryczan/hello-world:0.0.0
		return registry + "/" + image, nil
	}
	// Image was passed with a registry e.g. gcr.io/sgryczan/hello-world:0.0.0
	result := registry + "/" + chunk[1] + "/" + chunk[2]
	fmt.Printf("[ReplaceImageRegistry]  Image %s transformed to %s", image, result)
	return result, nil
}

// PullandRetag runs the docker commands to pull images through the mirror, and retag
// back to the original name
func PullandRetag(containers *[]BrokenContainer) {
	fmt.Printf("[PullandRetag]  processing %d containers.", len(*containers))
	for _, container := range *containers {
		if registry := GetImageRegistry(container.Image); registry != "docker.io" {
			tempName, err := ReplaceImageRegistry(container.Image, "docker.repo.eng.netapp.com")
			if err != nil {
				log.Fatalf("Error replacing image name %s: %v\n", registry, err)
			}

			err = PullImage(tempName)
			if err != nil {
				log.Fatalf("Error pulling image name %s: %v\n", registry, err)
			}

			err = TagImage(tempName, container.Image)
			if err != nil {
				log.Fatalf("Error pulling image name %s: %v\n", registry, err)
			}
		}
		continue
	}
	return
}

// PullImage pulls the specified image with the docker daemon
func PullImage(image string) error {
	log.Printf("[PullImage]  pulling image: %s", image)
	pullCmd := exec.Command("/usr/bin/docker", "pull", image)
	var out bytes.Buffer
	var eout bytes.Buffer
	pullCmd.Stdout = &out
	pullCmd.Stdout = &eout
	err := pullCmd.Run()
	if err != nil {
		log.Printf("Error running command: %v", err)
		log.Printf("stdout: %v", out)
		log.Printf("stderr: %v", eout)
		return err
	}
	return nil
}

// TagImage tags image oldname to newname
func TagImage(oldname string, newname string) error {

	tagCmd := exec.Command("/usr/bin/docker", "tag", oldname, newname)
	var out bytes.Buffer
	var eout bytes.Buffer
	tagCmd.Stdout = &out
	tagCmd.Stdout = &eout
	err := tagCmd.Run()
	if err != nil {
		log.Printf("Error running command: %v", err)
		log.Printf("stdout: %v", out)
		log.Printf("stderr: %v", eout)
		return err
	}
	log.Printf("[TagImage] tagged image: %s to %s", oldname, newname)
	return nil
}
