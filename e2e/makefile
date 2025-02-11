GATEWAY_PKGS = ./test/gateway
.phony: e2e
e2e:
	ginkgo run -r --v --json-report ./gateway-report.json $(GATEWAY_PKGS)
