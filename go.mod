module github.com/mellowplace/go-pact-provider

go 1.12

require (
	github.com/hashicorp/go-version v1.2.0 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/pact-foundation/pact-go v1.0.0-beta.7
)

replace github.com/pact-foundation/pact-go => github.com/karhoo/pact-go v1.0.0-karhoo.0.20190802160034-5e030d39530c
