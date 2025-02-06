GATEWAY_PKGS = ./test/gateway
.phony: e2e
e2e:
	ginkgo run -r --v $(GATEWAY_PKGS) --ginkgo.json-report=gateway-report.json
