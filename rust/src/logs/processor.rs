use std::sync::atomic::{AtomicBool, Ordering};

use opentelemetry_sdk::error::OTelSdkResult;

use crate::wit::wasi::otel::logs::emit;

#[derive(Debug)]
pub struct WasiLogProcessor {
    is_shutdown: AtomicBool,
}

impl WasiLogProcessor {
    pub fn new() -> Self {
        Self {
            is_shutdown: AtomicBool::new(false),
        }
    }
}

impl Default for WasiLogProcessor {
    fn default() -> Self {
        Self::new()
    }
}

impl opentelemetry_sdk::logs::LogProcessor for WasiLogProcessor {
    fn emit(
        &self,
        data: &mut opentelemetry_sdk::logs::SdkLogRecord,
        _: &opentelemetry::InstrumentationScope,
    ) {
        match emit(&data.into()) {
            Ok(_) => (),
            Err(v) => println!("ERROR: opentelemetry_wasi failed to emit log: {}", v),
        }
    }

    fn force_flush(&self) -> opentelemetry_sdk::error::OTelSdkResult {
        if self.is_shutdown.load(Ordering::Relaxed) {
            return OTelSdkResult::Err(opentelemetry_sdk::error::OTelSdkError::AlreadyShutdown);
        }
        Ok(())
    }

    fn shutdown(&self) -> opentelemetry_sdk::error::OTelSdkResult {
        let result = self.force_flush();
        if self.is_shutdown.swap(true, Ordering::Relaxed) {
            return OTelSdkResult::Err(opentelemetry_sdk::error::OTelSdkError::AlreadyShutdown);
        }
        result
    }
}
