module github.com/spin_basic

go 1.25

require github.com/spinframework/spin-go-sdk/v3 v3.0.0-00010101000000-000000000000

require (
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/spinframework/spin-go-sdk/v3/wit_component v0.1.0 // indirect
)

replace github.com/spinframework/spin-go-sdk/v3 => github.com/asteurer/spin-go-sdk/v3 v3.0.0-20251215051036-e1273e70063a

replace github.com/spinframework/spin-go-sdk/v3/wit_component => github.com/asteurer/spin-go-sdk/v3/wit_component v0.0.0-20251215051036-e1273e70063a

replace github.com/calebschoepp/opentelemetry_wasi => ../../

replace github.com/calebschoepp/opentelemetry_wasi/wit_component => ../../wit_component
