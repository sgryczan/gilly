# Gilly
![](./img/logo2.png)

*noun; a truck or wagon used to transport the equipment of a circus or carnival.*

Gilly is a tool meant for fixing Kubernetes workloads for hosts running under BigTop, that reference images from gcr.io or quay.io.

## Installation
**Requirements**:
* A running k8s cluster
* `kubectl` installed and configured to access above cluster

## Quick start
`make install`

## Detail

**Provision certificates**

All requests from the cluster control plane to the application must be made over HTTPS. In order to achieve this, we'll need to provision a certificate from our cluster, which gilly will be configured to use.

To do this, run `make ssl`:
```
make -C ssl cert
... creating gilly.key
Generating RSA private key, 2048 bit long modulus
............+++
......................................................................................................+++
e is 65537 (0x10001)
... creating gilly.csr
openssl req -new -key gilly.key -subj "/CN=gilly.gilly.svc" -out gilly.csr -config csr.conf
... deleting existing csr, if any
kubectl delete csr gilly.gilly.svc || :
Error from server (NotFound): certificatesigningrequests.certificates.k8s.io "gilly.gilly.svc" not found
... creating kubernetes CSR object
kubectl create -f -
certificatesigningrequest.certificates.k8s.io/gilly.gilly.svc created
... waiting for csr to be present in kubernetes
kubectl get csr gilly.gilly.svc
certificatesigningrequest.certificates.k8s.io/gilly.gilly.svc approved
... waiting for serverCert to be present in kubernetes
kubectl get csr gilly.gilly.svc -o jsonpath='{.status.certificate}'
... creating gilly.pem cert file
$serverCert | openssl base64 -d -A -out gilly.pem
```

**Build the Template**

Next, render the deployment template with `make template`:
```
sed "s/KUBE_CA_BUNDLE/LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJWakNCL3FBREFnRUNBZ0VBTUFvR0NDcUdTTTQ5QkFNQ01DTXhJVEFmQmdOVkJBTU1HR3N6Y3kxelpYSjIKWlhJdFkyRkFNVFU0T1RreU56RTRNakFlRncweU1EQTFNVGt5TWpJMk1qSmFGdzB6TURBMU1UY3lNakkyTWpKYQpNQ014SVRBZkJnTlZCQU1NR0dzemN5MXpaWEoyWlhJdFkyRkFNVFU0T1RreU56RTRNakJaTUJNR0J5cUdTTTQ5CkFnRUdDQ3FHU000OUF3RUhBMElBQkc0bFZ4Q2N6SCs1bVB1Und6NVhXcEJxTXVYNEN6cTFKdVhYeWNFODBKVFAKdlNtbk96VWtEbDlLdzRyUWhFVm01MkRSY1RKRGliTjhzdGdzUmtaWHRPdWpJekFoTUE0R0ExVWREd0VCL3dRRQpBd0lDcERBUEJnTlZIUk1CQWY4RUJUQURBUUgvTUFvR0NDcUdTTTQ5QkFNQ0EwY0FNRVFDSUVTRXNadnc5akcxCmwzZFBkWkV5b3Z4S2ZYMkJtT1c4akpzcXZQZC8wYTd5QWlBRWY0YVVLaUNMNnJtVU9SWnB4UkZCQk5HSkZFSFAKb3VuV1VUT0lnQkh4QkE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==/; s/IMAGE_VERSION/0.1.6/; s/IMAGE_REGISTRY/docker.repo.eng.netapp.com\/sgryczan/; s/IMAGE_NAME/gilly/" deploy/template.yaml > deploy/stack.yaml
kubectl create secret generic gilly-certs \
                --from-file=./ssl/gilly.pem --from-file ./ssl/gilly.key -n gilly --dry-run -o yaml >> deploy/stack.yaml
```

**Deploy the Stack**

Finally, deploy the stack to your cluster with `make deploy`:
```
gilly git:(master)$ make deploy
kubectl apply -f deploy/stack.yaml
namespace/gilly created
service/gilly created
deployment.apps/gilly created
mutatingwebhookconfiguration.admissionregistration.k8s.io/gilly created
secret/gilly-certs created
```

### Why does this exist?
Under the BigTop, images cannot be pulled from popular Docker registries such as gcr.io or quay.io

So, if we have a pod with the following spec:
```
apiVersion: v1
kind: Pod
metadata:
  name: myapp
spec:
  - name: myapp
    image: gcr.io/sgryczan/rocket:latest
```

This will fail with an error like this:
```
NAME    READY   STATUS              RESTARTS   AGE
myapp   0/1     ImagePullBackOff    2          3m30s
```

Ordinarily, we can fix this issue by updating the `image:` field of our manifest to reference an interally-reachable repository:

```
spec:
  - name: myapp
    image: docker.repo.eng.netapp.com/sgryczan/rocket:latest
```

However, in an environment with hundreds of workloads, this represents a significant effort. In addition, using upstream templating tools like `helm` or `kustomize` may not always offer a parameter for overriding every image included in the template or chart.

Gilly removes the need to do this, by automatically rewriting the `image:` field on pods in real-time, as they are scheduled. Using our example above, Gilly will automatically rewrite `gcr.io/sgryczan/rocket:latest` to `docker.repo.eng.netapp.com/sgryczan/rocket:latest`, allowing the image to be pulled successfully:
```
2020/06/03 16:40:46 [Gilly]  started v0.1.6
2020/06/03 16:41:06 [Mutate]  Received POD create event. Name: myapp2, Namespace: flux
2020/06/03 16:41:06 [ProcessPod]  Found registry => gcr.io
2020/06/03 16:41:06 [ProcessPod] image registry for container myapp is gcr.io - updating
2020/06/03 16:41:06 [ReplaceImageRegistry]  Replacing registry for image gcr.io/sgryczan/rocket:latest with docker.repo.eng.netapp.com
2020/06/03 16:41:06 [ReplaceImageRegistry]  Image gcr.io/sgryczan/rocket:latest => docker.repo.eng.netapp.com/sgryczan/rocket:latest
2020/06/03 16:41:06 [ProcessPod] updated registry for container myapp to docker.repo.eng.netapp.com/sgryczan/rocket:latest
```

It does this without the need to modify the `image:` tag in the original manifest.

### How Gilly Works
Gilly runs registers itself as an API against it's cluster, using a Mutating Admission Webhook. Essentially, it instructs the Kubernetes API to notify the api every time a pod is created in the cluster. The API will send an AdmissionReview, including the spec of the Pod to be created. Gilly will modify the image field of these new pods automatically, and will respond to the API server with an AdmissionResponse, specifying how the image fields in the Pod definition should be modified.


### Example Output

If we create the following pod:

```
apiVersion: v1
kind: Pod
metadata:
  name: myapp2
spec:
  containers:
  - name: myapp
    image: gcr.io/sgryczan/rocket:latest
```

Gilly will modify the image field to reference the mirror:

```
2020/05/20 14:37:22 [Gilly]  started
2020/05/20 14:37:33 [Mutate]  Received POD create event. Name: myapp2, Namespace: default
2020/05/20 14:37:33 [Mutate]  Found registry => gcr.io
2020/05/20 14:37:33 [Mutate] image registry for container myapp is gcr.io - updating
2020/05/20 14:37:33 [ReplaceImageRegistry]  Replacing registry for image gcr.io/sgryczan/rocket:latest with docker.repo.eng.netapp.com
2020/05/20 14:37:33 [ReplaceImageRegistry]  Image gcr.io/sgryczan/rocket:latest => docker.repo.eng.netapp.com/sgryczan/rocket:latest
2020/05/20 14:37:33 [Mutate] updated registry for container myapp to docker.repo.eng.netapp.com/sgryczan/rocket:latest
```

Now our pod can run:
```
NAME    READY   STATUS     RESTARTS   AGE
myapp   1/1     Running    2          3m30s
```


### Why not just rename the image manually?
We could manually update the `image:` field to reference the image through the mirror, e.g.:
```
gcr.io/sgryczan/ansible-runner:0.0.0 ->
docker.repo.eng.netapp.com/sgryczan/ansible-runner:0.0.0
```
There are a few reasons we **dont want to do this**:

* **Every workload will need to be modified in this way**- this will include control plane components as well, as those are referenced from k8s.gcr.io. Such a change requires rewriting a ton of workloads, which presents a significant effort. Furthermore, if the manifests are rendered by a tool such as `kustomize` or `helm`, the image tags may not be exposed in such a way to allow this to be done without manual modification.

* **The registry mirror clobbers the upstream registry**. Essentially, each mirrored registry is checked sequentially when pulling images. This means that effectively, each image pulled through the mirror is tagged to either docker.io, or the mirror. This obscures the original registry of the image. Since this is part of the FQDN used to identify the image, it becomes impossible to determine which upstream registry the image is actually being pulled from.

For example, if I want to pull this image: `docker pull gcr.io/linkerd-io/proxy:stable-2.7.1`:

* Pulling directly fails:
```
$ docker pull gcr.io/linkerd-io/proxy:stable-2.7.1
Error response from daemon: received unexpected HTTP status: 503 Service Unavailable
```
* But pulling via the mirror works by replacing the registry domain with docker.repo.eng.netapp.com:
```
$ docker pull docker.repo.eng.netapp.com/linkerd-io/proxy:stable-2.7.1
stable-2.7.1: Pulling from linkerd-io/proxy
Digest: sha256:22af88a12d252f71a3aee148d32d727fe7161825cb4164776c04aa333a1dcd28
Status: Image is up to date for linkerd-io/proxy:stable-2.7.1
docker.repo.eng.netapp.com/linkerd-io/proxy:stable-2.7.1
```
* Even though an image with this FQDN doesn't actually exist!


