// Package main contains a runnable Provider Pact test example.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
)

var dir, _ = os.Getwd()

// Example Provider Pact: How to run me!
// 1. Start the daemon with `./pact-go daemon`
// 2. cd <pact-go>/examples
// 3. go test -v -run TestProvider
func TestProvider(t *testing.T) {

	// Create Pact connecting to local Daemon
	pact := &dsl.Pact{
		Provider: "MyProvider",
		Host:     "localhost",
		LogLevel: "DEBUG",
	}

	// Start provider API in the background
	go startServer()

	publishPact := false
	if os.Getenv("CI") != "" {
		publishPact = true
	}

	request := types.VerifyRequest{
		ProviderBaseURL:            "http://localhost:8000",
		ProviderVersion:            "0.0.1-" + os.Getenv("TRAVIS_COMMIT"), // NEED TO USE GIT SHA TO GET A REPEATABLE VERSION BETWEEN JOB RUNS
		PublishVerificationResults: publishPact,
		ProviderStatesSetupURL:     "http://localhost:8000/setup",
		CustomProviderHeaders:      []string{"Authorization: basic e5e5e5e5e5e5e5"},
		Verbose:                    true,
	}

	if os.Getenv("TRAVIS_PACT_URL") != "" {
		// we're supposed to validate a specific pact
		request.PactURLs = []string{os.Getenv("TRAVIS_PACT_URL")}
	} else {
		// validate against all consumers of this service who have
		// "production" pacts published
		request.BrokerURL = "http://localhost:8080"

		// Provider and Tags seems to just be completely ignored?!!
		// https://github.com/pact-foundation/pact-go/issues/116
		request.Provider = "MyProvider"
		request.Tags = []string{"production"}
	}

	details, err := pact.VerifyProvider(t, request)

	if err != nil {
		// TODO insane they don't provide a nice test output formatter
		log.Fatalf("Error on Verify: %v", details)
	}
}

func startServer() {
	mux := http.NewServeMux()
	lastName := "billy"

	mux.HandleFunc("/foobar", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, fmt.Sprintf(`{"lastName":"%s", "name":"%s", "age": 3}`, lastName, "jean"))

		// Break the API by replacing the above and uncommenting one of these
		// w.WriteHeader(http.StatusUnauthorized)
		// fmt.Fprintf(w, fmt.Sprintf(`{"NoName":"%s"}`, lastName))
	})

	// This function handles state requests for a particular test
	// In this case, we ensure that the user being requested is available
	// before the Verification process invokes the API.
	mux.HandleFunc("/setup", func(w http.ResponseWriter, req *http.Request) {
		var s *types.ProviderState
		decoder := json.NewDecoder(req.Body)
		decoder.Decode(&s)
		if s.State == "User foo exists" {
			lastName = "bar"
		}

		w.Header().Add("Content-Type", "application/json")
	})
	log.Fatal(http.ListenAndServe(":8000", mux))
}
