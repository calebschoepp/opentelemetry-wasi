# NOTE: To generate the C bindings you need to install a forked version of wit-bindgen.
#
#   cargo install wit-bindgen-cli --git https://github.com/fermyon/wit-bindgen-backport --rev "b89d5079ba5b07b319631a1b191d2139f126c976"
#

.PHONY: generate
generate:
	wit-bindgen c --import wit/logs.wit --out-dir ./logs
	wit-bindgen c --import wit/metrics.wit --out-dir ./metrics
	wit-bindgen c --import wit/traces.wit --out-dir ./traces