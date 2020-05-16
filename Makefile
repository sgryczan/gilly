build:
	docker build -t sgryczan/gilly:0.0.0 .

push:
	docker push sgryczan/gilly:0.0.0 
build-bin:
	go build -o gilly gilly.go