.PHONY: pact-tests

pact-tests:
	TRAVIS_COMMIT=${USER}-snapshot go test -v -run TestProvider .
