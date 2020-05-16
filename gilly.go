package main

import (
	"log"
	"time"

	"bitbucket.ngage.netapp.com/scm/hcit/gilly/lib"
)

var kubeletEndpoint string

func main() {
	client, err := lib.NewConnection()
	if err != nil {
		log.Fatalf("Failed to establish connection to Kubelet: %v", err)
	}
	log.Printf("[Main]  host: %s", client.EndPoint)
	//log.Printf("token: %s", client.Token)
	for {
		//Get current running pods from kubelet (http://localhost:10250/pods)
		pods, err := client.GetPods()
		if err != nil {
			log.Fatalf("Failed to get running pods: %v", err)
		}

		// Determine if any pods are in targeted states
		brokenContainers, err := lib.GetBrokenContainers(pods)
		if err != nil {
			log.Fatalf("Failed to get running pods: %v", err)
		}

		// For affected images, run docker commands to fix
		go lib.PullandRetag(brokenContainers)
		log.Printf("[Main]  Checks done. Sleeping for awhile")
		time.Sleep(time.Second * 300)
	}
}
