# Gilly
*noun; a truck or wagon used to transport the equipment of a circus or carnival.*

Gilly is a tool meant for fixing Kubernetes workloads for hosts running under BigTop, that reference images from gcr.io or quay.io.

### Why?
Under the BigTop, images cannot be pulled from popular Docker registries such as gcr.io or quay.io

So, if we have a pod with the following spec:
```
apiVersion: v1
kind: Pod
metadata:
  name: myapp
spec:
  - name: myapp
    image: gcr.io/sgryczan/ansible-runner:0.0.0
```

This will fail with an error like this:
```
NAME    READY   STATUS              RESTARTS   AGE
myapp   0/1     ImagePullBackOff    2          3m30s
```

### How Gilly Works
Gilly runs a daemon on each worker in the cluster. Each daemon will [connect to the Kubelet service](https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/server/server.go#L300-L314) on the host under which it is running, and will look for pods that are failing with `ImagePullbackOff` or `ErrImagePull` errors. 

It will then attempt to fix these images by pulling the image through the mirror, and re-tagging the mirrored image back to the original name. Once this is done, the Kubelet can start the pod, since the referenced image is now accessible by the daemon.

### Example Output

```
2020/05/16 01:11:53 [NewConnection]  Established connection to kubelet
2020/05/16 01:11:53 [Main]  host: https://hg-sre-devel-worker10:10250
2020/05/16 01:11:53 [GetPods]  Pulled 12 running pods.
2020/05/16 01:11:53 [Validate]  Detected pod myapp in Pending state.
2020/05/16 01:11:53 [Validate]  Detected container myapp in pod myapp status is ImagePullBackOff
2020/05/16 01:11:53 [Validate]  Detected container myapp in pod myapp reason: ImagePullBackOff
2020/05/16 01:11:53 [Validate]  Checked 12 Pods, found 1 waiting
2020/05/16 01:11:53 [Main]  Checks done. Sleeping for awhile
2020/05/16 01:11:53 [ReplaceImageRegistry]  got chunks: [gcr.io sgryczan ansible-runner:0.0.0]
2020/05/16 01:11:53 [PullImage]  pulling image: docker.repo.eng.netapp.com/sgryczan/ansible-runner:0.0.0
2020/05/16 01:11:54 [TagImage] tagged image: docker.repo.eng.netapp.com/sgryczan/ansible-runner:0.0.0 to gcr.io/sgryczan/ansible-runner:0.0.0
```

Now our pod is running:
```
NAME    READY   STATUS     RESTARTS   AGE
myapp   1/1     Running    2          3m30s
```


### Why not just rename the image?
We could update the `image:` field to reference the image through the mirror, e.g.:
```
gcr.io/sgryczan/ansible-runner:0.0.0 ->
docker.repo.eng.netapp.com/sgryczan/ansible-runner:0.0.0
```
There are a few reasons we **dont want to do this**:

* **Every workload will need to be modified in this way**- this will include control plane components as well, as those are referenced from k8s.gcr.io. Such a change requires rewriting a ton of workloads, which presents a significant effort. Furthermore, if the manifests are rendered by a tool such as `kustomize` or `helm`, the image tags may not be exposed in such a way to allow this to be done easily.

* **The registry mirror clobbers the upstream registry**. Essentially, each mirrored registry is checked sequentially when pulling images. This means that effectively, each image pulled through the mirror is tagged to either docker.io, or the mirror. This obscures the original registry of the image. Since this is part of the FQDN used to identify the image, it becomes impossible to determine which upstream registry the image is actually being pulled from.

For example, if I want to pull this image: `docker pull gcr.io/linkerd-io/proxy:stable-2.7.1`:

* Pulling directly fails:
```
$ docker pull gcr.io/linkerd-io/proxy:stable-2.7.1
Error response from daemon: received unexpected HTTP status: 503 Service Unavailable
```
* But pulling via the mirror works by replacing the registry domain with docker.io or docker.repo.eng.netapp.com:
```
$ docker pull docker.io/linkerd-io/proxy:stable-2.7.1
stable-2.7.1: Pulling from linkerd-io/proxy
Digest: sha256:22af88a12d252f71a3aee148d32d727fe7161825cb4164776c04aa333a1dcd28
Status: Image is up to date for linkerd-io/proxy:stable-2.7.1
docker.io/linkerd-io/proxy:stable-2.7.1
```
* Even though `docker.io/linkerd-io/proxy:stable-2.7.1` doesn't actually exist!
