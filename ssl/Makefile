APP = gilly
NAMESPACE = gilly

.PHONY:cert
cert:
	@./ssl.sh $(APP) $(NAMESPACE)

.PHONY:clean
clean:
	rm -vf *.key *.pem *.cert *.crt *.csr