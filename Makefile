build: ssl build-image

build-image:
	docker build -t sgryczan/gilly:0.1.0 .
build-bigtop:
	docker build -t sgryczan/gilly:0.1.0 -f Dockerfile.bigtop .

ssl:
	make -C ssl cert
	
push:
	docker push sgryczan/gilly:0.1.0 
build-bin:
	go build -o gilly gilly.go

kube-get-ca-bundle:
	kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}'

.PHONY: build build-bigtop ssl push build-bin kube-get-ca-bundle