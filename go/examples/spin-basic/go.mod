module github.com/calebschoepp/opentelemetry-wasi/examples/spin-basic

go 1.25

require (
	github.com/calebschoepp/opentelemetry-wasi v0.0.0-00010101000000-000000000000
	github.com/spinframework/spin-go-sdk/v3 v3.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.39.0
	go.opentelemetry.io/otel/log v0.14.0
	go.opentelemetry.io/otel/metric v1.39.0
	go.opentelemetry.io/otel/sdk v1.39.0
	go.opentelemetry.io/otel/sdk/log v0.15.0
	go.opentelemetry.io/otel/sdk/metric v1.39.0
)

require (
	github.com/calebschoepp/opentelemetry-wasi/wit_component v0.0.0-00010101000000-000000000000 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/spinframework/spin-go-sdk/v3/wit_component v0.1.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
)

replace github.com/spinframework/spin-go-sdk/v3 => github.com/asteurer/spin-go-sdk/v3 v3.0.0-20251215051036-e1273e70063a

replace github.com/spinframework/spin-go-sdk/v3/wit_component => github.com/asteurer/spin-go-sdk/v3/wit_component v0.0.0-20251215051036-e1273e70063a

replace github.com/calebschoepp/opentelemetry-wasi => ../../

replace github.com/calebschoepp/opentelemetry-wasi/wit_component => ../../wit_component
