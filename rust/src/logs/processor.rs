use crate::{logs::conversion::to_wasi_log_record, wit::wasi};
use opentelemetry_sdk::error::OTelSdkResult;
use std::sync::atomic::{AtomicBool, Ordering};

#[derive(Debug)]
pub struct WasiLogProcessor {
    is_shutdown: AtomicBool,
    resource: Option<opentelemetry_sdk::Resource>,
}

impl WasiLogProcessor {
    pub fn new(resource: Option<opentelemetry_sdk::Resource>) -> Self {
        Self {
            is_shutdown: AtomicBool::new(false),
            resource,
        }
    }
}

impl opentelemetry_sdk::logs::LogProcessor for WasiLogProcessor {
    fn emit(
        &self,
        data: &mut opentelemetry_sdk::logs::SdkLogRecord,
        scope: &opentelemetry::InstrumentationScope,
    ) {
        wasi::otel::logs::on_emit(&to_wasi_log_record(data, scope, self.resource.as_ref()))
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
