VERSION=0.1.3

build: ssl build-image template

build-bigtop: ssl build-image-bigtop template

deploy:
	kubectl apply -f deploy/stack.yaml

deploy-cleanup:
	kubectl delete -f deploy/stack.yaml

build-image:
	docker build -t sgryczan/gilly:$(VERSION) .
	sed "s/KUBE_CA_BUNDLE/$(make -s kube-get-ca-bundle)/" deploy/template.yaml > deploy/stack.yaml

build-image-bigtop:
	docker build -t sf-artifactory.solidfire.net:9004/gilly:$(VERSION) -f Dockerfile.bigtop .

push-image-bigtop:
	docker push sf-artifactory.solidfire.net:9004/gilly:$(VERSION)

ssl:
	make -C ssl cert

template:
	sed "s/KUBE_CA_BUNDLE/$(shell make -s kube-get-ca-bundle)/; s/IMAGE_VERSION/$(VERSION)/" deploy/template.yaml > deploy/stack.yaml

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
	k3d import-images sf-artifactory.solidfire.net:9004/gilly:$(VERSION) && \
	make template

.PHONY: build build-bigtop build-image build-image-bigtop ssl push build-bin kube-get-ca-bundle k3d-deps k3d-build template deploy