use crate::tracer::Tracer;

pub struct TracerProvider {}

impl Default for TracerProvider {
    fn default() -> Self {
        TracerProvider::builder().build()
    }
}

impl TracerProvider {
    pub fn builder() -> Builder {
        Builder::default()
    }
}

impl opentelemetry::trace::TracerProvider for TracerProvider {
    type Tracer = crate::tracer::Tracer;

    fn library_tracer(
        &self,
        library: std::sync::Arc<opentelemetry::InstrumentationLibrary>,
    ) -> Self::Tracer {
        Tracer::new(library)
    }
}

#[derive(Debug, Default)]
pub struct Builder {}

impl Builder {
    pub fn build(self) -> TracerProvider {
        TracerProvider {}
    }
}
