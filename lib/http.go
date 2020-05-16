package lib

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// GetPods returns the running pods from the Kubelet
func (c *KubeletConn) GetPods() (*PodList, error) {
	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest("GET", c.EndPoint+"/pods", nil)
	if err != nil {
		log.Printf("Error getting pods: %v", err)
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.Token)

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error getting pods: %v", err)
		return nil, err
	}

	pods := &PodList{}
	err = json.NewDecoder(res.Body).Decode(&pods)
	if err != nil {
		log.Printf("Error decoding pods: %v", err)
		raw, err := ioutil.ReadAll(res.Body)
		fmt.Println(string(raw))
		return nil, err
	}
	log.Printf("[GetPods]  Pulled %d running pods.", len(pods.Items))
	return pods, nil
}

// NewConnection creates a new KubeletConn instance
func NewConnection() (*KubeletConn, error) {
	token, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		log.Printf("Error reading service account token: %v", err)
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Error getting Hostname: %v", err)
		return nil, err
	}

	kc := &KubeletConn{
		EndPoint: fmt.Sprintf("https://%s:10250", hostname),
		Token:    string(token),
	}

	log.Printf("[NewConnection]  Established connection to kubelet")
	return kc, nil
}
