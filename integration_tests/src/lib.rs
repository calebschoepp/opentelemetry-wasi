#![cfg(unix)]
#[cfg(test)]
mod tests {
    use anyhow::{Error, anyhow};
    use fake_opentelemetry_collector::{
        ExportedLog, ExportedMetric, ExportedSpan, FakeCollectorServer,
    };
    use nix::{
        sys::signal::{Signal, kill},
        unistd::Pid,
    };
    use serial_test::serial;
    use std::{
        process::{Child, Command},
        time::Duration,
    };

    #[tokio::test]
    #[serial]
    async fn rust_spin_basic() {
        // Retrieve telemetry.
        let (spans, metrics, logs) = get_telemetry_from_spin_app("../rust/examples/spin-basic")
            .await
            .expect("Failed to retrieve telemetry from Spin app");

        // Run tests.
        basic_signal_validation("rust_spin_basic", Some(&spans), Some(&metrics), Some(&logs));
    }

    #[tokio::test]
    #[serial]
    async fn rust_spin_tracing() {
        // Retrieve telemetry.
        let (spans, _metrics, logs) = get_telemetry_from_spin_app("../rust/examples/spin-tracing")
            .await
            .expect("Failed to retrieve telemetry from Spin app");

        // Run tests.
        basic_signal_validation("rust_spin_tracing", Some(&spans), None, Some(&logs));
    }

    #[tokio::test]
    #[serial]
    async fn go_spin_basic() {
        // Retrieve telemetry.
        let (spans, metrics, logs) = get_telemetry_from_spin_app("../go/examples/spin-basic")
            .await
            .expect("Failed to retrieve telemetry from Spin app");

        // Run tests.
        basic_signal_validation("go_spin_basic", Some(&spans), Some(&metrics), Some(&logs));
    }

    /// Performs a basic validation on each telemetry signal's struct field.
    fn basic_signal_validation(
        prefix: &str,
        spans: Option<&[ExportedSpan]>,
        metrics: Option<&[ExportedMetric]>,
        logs: Option<&[ExportedLog]>,
    ) {
        if let Some(span_data) = spans {
            insta::assert_yaml_snapshot!(format!("{}_tracing", prefix), span_data, {
                "[].start_time_unix_nano" => "[timestamp]",
                "[].end_time_unix_nano" => "[timestamp]",
                "[].attributes" => "[attributes]",
                "[].events[].time_unix_nano" => "[timestamp]",

                // These correspond to spans emitted directly from Spin,
                // and may break the tests if they are changed.
                "[].events[].attributes[\"code.filepath\"]" => "[value]",
                "[].events[].attributes[\"code.lineno\"]" => "[value]",
                "[].events[].attributes[\"code.namespace\"]" => "[value]",
                "[].events[].attributes.level" => "[value]",
                "[].events[].attributes.target" => "[value]",


                "[].trace_id" => insta::dynamic_redaction(|value, _path| {
                    assert2::let_assert!(Some(trace_id) = value.as_str());
                    format!("[trace_id:len({})]", trace_id.len())
                }),
                "[].span_id" => insta::dynamic_redaction(|value, _path| {
                    assert2::let_assert!(Some(span_id) = value.as_str());
                    format!("[span_id:len({})]", span_id.len())
                }),
                "[].parent_span_id" => insta::dynamic_redaction(|value, _path| {
                    assert2::let_assert!(Some(span_id) = value.as_str());
                    format!("[parent_span_id:len({})]", span_id.len())
                }),
            });
        }

        if let Some(metric_data) = metrics {
            insta::assert_yaml_snapshot!(format!("{}_metrics", prefix), metric_data, {
                "[].**.start_time_unix_nano" => "[timestamp]",
                "[].**.time_unix_nano" => "[timestamp]",
            });
        }

        if let Some(log_data) = logs {
            insta::assert_yaml_snapshot!(format!("{}_logs", prefix), log_data, {
                "[].observed_time_unix_nano" => "[timestamp]",
                "[].trace_id" => insta::dynamic_redaction(|value, _path| {
                    assert2::let_assert!(Some(trace_id) = value.as_str());
                    format!("[trace_id:len({})]", trace_id.len())
                }),
                "[].span_id" => insta::dynamic_redaction(|value, _path| {
                    assert2::let_assert!(Some(span_id) = value.as_str());
                    format!("[span_id:len({})]", span_id.len())
                }),
            });
        }
    }

    async fn get_telemetry_from_spin_app(
        path: &str,
    ) -> Result<(Vec<ExportedSpan>, Vec<ExportedMetric>, Vec<ExportedLog>), Error> {
        // Start the collector.
        let (mut collector, collector_endpoint) = start_collector().await;

        // Build, instantiate, and invoke the Spin app.
        let mut app = SpinApp::new(path, &collector_endpoint);
        app.build()?;
        app.instantiate()?;
        app.invoke().await?;

        // Retrieve telemetry data.
        let timeout = Duration::from_secs(10);
        let collector_min = 6;
        let spans = collector.exported_spans(collector_min, timeout).await;
        let metrics = collector.exported_metrics(collector_min, timeout).await;
        let logs = collector.exported_logs(collector_min, timeout).await;

        Ok((spans, metrics, logs))
    }

    /// Starts a fake collector server.
    async fn start_collector() -> (FakeCollectorServer, String) {
        let collector = FakeCollectorServer::start()
            .await
            .expect("fake collector should start");
        let endpoint = collector.endpoint().clone();
        (collector, endpoint)
    }

    struct SpinApp<'a> {
        path: &'a str,
        collector_endpoint: &'a str,
        process: Option<Child>,
    }

    impl<'a> Drop for SpinApp<'a> {
        fn drop(&mut self) {
            if let Some(app_process) = &mut self.process {
                // Gracefully shut down the process using 'ctrl + c'
                let _ = kill(Pid::from_raw(app_process.id() as i32), Signal::SIGINT);
                let _ = app_process.wait();
            }
        }
    }

    impl<'a> SpinApp<'a> {
        fn new(path: &'a str, collector_endpoint: &'a str) -> Self {
            Self {
                path,
                collector_endpoint,
                process: None,
            }
        }

        /// Compiles a Spin app.
        fn build(&self) -> Result<(), Error> {
            let spin_build_output = Command::new("spin")
                .args(["build", "-f", self.path])
                .output()
                .expect("Failed to execute 'spin build'");
            if !spin_build_output.status.success() {
                return Err(anyhow!(
                    "{} -> build: {}",
                    self.path,
                    String::from_utf8_lossy(&spin_build_output.stderr)
                ));
            }
            Ok(())
        }

        /// Instantiates a Spin app.
        fn instantiate(&mut self) -> Result<(), Error> {
            let child = Command::new("spin")
                .env("OTEL_EXPORTER_OTLP_ENDPOINT", self.collector_endpoint)
                .env("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc")
                .args(["up", "-f", self.path])
                .spawn()?;

            self.process = Some(child);
            Ok(())
        }

        /// Sends an HTTP request to a Spin app.
        async fn invoke(&self) -> Result<(), Error> {
            let start = std::time::Instant::now();
            loop {
                match reqwest::get("http://localhost:3000").await {
                    Ok(_) => return Ok(()),
                    Err(e) => {
                        if start.elapsed() > Duration::from_secs(5) {
                            return Err(anyhow!("Unable to reach the Spin app: {e}"));
                        }
                    }
                }
            }
        }
    }
}
