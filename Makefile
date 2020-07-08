VERSION=0.1.6
IMAGE_NAME=gilly
IMAGE_REGISTRY=docker.io

build: ssl build-image template

install: ssl template deploy
uninstall: deploy-cleanup

deploy:
	kubectl apply -f deploy/stack.yaml

deploy-cleanup:
	kubectl delete -f deploy/stack.yaml

build-image:
	docker build -t sgryczan/$(IMAGE_NAME):$(VERSION) .
	sed "s/KUBE_CA_BUNDLE/$(make -s kube-get-ca-bundle)/" deploy/template.yaml > deploy/stack.yaml

ssl:
	make -C ssl cert

template:
	sed "s/KUBE_CA_BUNDLE/$(shell make -s kube-get-ca-bundle)/; s/IMAGE_VERSION/$(VERSION)/; s/IMAGE_REGISTRY/$(IMAGE_REGISTRY)/; s/IMAGE_NAME/$(IMAGE_NAME)/" deploy/template.yaml > deploy/stack.yaml
	kubectl create secret generic gilly-certs \
		--from-file=./ssl/gilly.pem --from-file ./ssl/gilly.key -n gilly --dry-run -o yaml >> deploy/stack.yaml

push:
	docker push sgryczan/$(IMAGE_NAME):$(VERSION) 

build-bin:
	go build -o gilly gilly.go

kube-get-ca-bundle:
	kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}'

.PHONY: build build-image build-image-bigtop ssl push build-bin kube-get-ca-bundle template deploy