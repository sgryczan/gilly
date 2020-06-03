VERSION=0.1.6
ARTIFACTORY_REGISTRY=sf-artifactory.solidfire.net:9004
BIGTOP_REGISTRY=docker.repo.eng.netapp.com\/sgryczan
IMAGE_NAME=gilly

build: ssl build-image template

build-bigtop: ssl build-image-bigtop template-bigtop

install: ssl template deploy
uninstall: deploy-cleanup

deploy:
	kubectl apply -f deploy/stack.yaml

deploy-cleanup:
	kubectl delete -f deploy/stack.yaml

build-image:
	docker build -t sgryczan/$(IMAGE_NAME):$(VERSION) .
	sed "s/KUBE_CA_BUNDLE/$(make -s kube-get-ca-bundle)/" deploy/template.yaml > deploy/stack.yaml

build-image-bigtop:
	docker build -t $(ARTIFACTORY_REGISTRY)/$(IMAGE_NAME):$(VERSION) -f Dockerfile.bigtop .

push-bigtop:
	docker push $(ARTIFACTORY_REGISTRY)/$(IMAGE_NAME):$(VERSION)

ssl:
	make -C ssl cert

template:
	sed "s/KUBE_CA_BUNDLE/$(shell make -s kube-get-ca-bundle)/; s/IMAGE_VERSION/$(VERSION)/; s/IMAGE_REGISTRY/$(BIGTOP_REGISTRY)/; s/IMAGE_NAME/$(IMAGE_NAME)/" deploy/template.yaml > deploy/stack.yaml
	kubectl create secret generic gilly-certs \
		--from-file=./ssl/gilly.pem --from-file ./ssl/gilly.key -n gilly --dry-run -o yaml >> deploy/stack.yaml

template-bigtop:
	sed "s/KUBE_CA_BUNDLE/$(shell make -s kube-get-ca-bundle)/; s/IMAGE_VERSION/$(VERSION)/; s/IMAGE_REGISTRY/$(ARTIFACTORY_REGISTRY)/; s/IMAGE_NAME/$(IMAGE_NAME)/" deploy/template.yaml > deploy/stack.yaml	
	kubectl create secret generic gilly-certs \
		--from-file=./ssl/gilly.pem --from-file ./ssl/gilly.key -n gilly --dry-run -o yaml >> deploy/stack.yaml
	

push:
	docker push sgryczan/$(IMAGE_NAME):$(VERSION) 

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
	k3d import-images $(ARTIFACTORY_REGISTRY)/$(IMAGE_NAME):$(VERSION) && \
	make template

.PHONY: build build-bigtop build-image build-image-bigtop ssl push build-bin kube-get-ca-bundle k3d-deps k3d-build template deploy