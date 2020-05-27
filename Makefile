VERSION=0.1.0

build: ssl build-image template

build-bigtop: ssl build-image-bigtop template

build-image:
	docker build -t sgryczan/gilly:$(VERSION) .

build-image-bigtop:
	docker build -t sgryczan/gilly:$(VERSION) -f Dockerfile.bigtop .

ssl:
	make -C ssl cert

template:
	sed "s/KUBE_CA_BUNDLE/$(make -s kube-get-ca-bundle)/" deploy/template.yaml > deploy/stack.yaml

push:
	docker push sgryczan/gilly:$(VERSION) 

build-bin:
	go build -o gilly gilly.go

kube-get-ca-bundle:
	kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}'

k3d-deps:
	docker pull docker.repo.eng.netapp.com/rancher/pause:3.1 && \
	docker tag docker.repo.eng.netapp.com/rancher/pause:3.1 docker.io/rancher/pause:3.1 && \
	k3d import-images docker.io/rancher/pause:3.1
	
k3d-build:
	make ssl && \
	make build-image-bigtop && \
	k3d import-images sgryczan/gilly:$(VERSION) && \
	make k3d-deps

.PHONY: build build-bigtop build-image build-image-bigtop ssl push build-bin kube-get-ca-bundle k3d-deps k3d-build