VERSION=0.1.0

build: ssl build-image

build-image:
	docker build -t sgryczan/gilly:$(VERSION) .
build-bigtop:
	docker build -t sgryczan/gilly:$(VERSION) -f Dockerfile.bigtop .

ssl:
	make -C ssl cert
	
push:
	docker push sgryczan/gilly:$(VERSION) 
build-bin:
	go build -o gilly gilly.go

kube-get-ca-bundle:
	kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}'

k3d-build:
	make ssl && \
	make build-bigtop && \
	k3d import-images sgryczan/gilly:$(VERSION) && \
	docker pull docker.repo.eng.netapp.com/rancher/pause:3.1 && \
	docker tag docker.repo.eng.netapp.com/rancher/pause:3.1 docker.io/rancher/pause:3.1 && \
	k3d import-images docker.io/rancher/pause:3.1

.PHONY: build build-bigtop ssl push build-bin kube-get-ca-bundle