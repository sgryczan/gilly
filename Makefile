build:
	docker build -t sgryczan/gilly:0.1.0 .

push:
	docker push sgryczan/gilly:0.1.0 
build-bin:
	go build -o gilly gilly.go

kube-get-ca-bundle:
	kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}'