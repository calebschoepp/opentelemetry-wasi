use std::borrow::Cow;

use opentelemetry::trace::SpanContext;

pub struct Span {
    inner: crate::wit::fermyon::spin2_0_0::observe::Span,
    span_context: SpanContext,
}

impl Span {
    pub(crate) fn new(name: Cow<'static, str>, span_context: SpanContext) -> Self {
        Self {
            inner: crate::wit::fermyon::spin2_0_0::observe::Span::enter(&name),
            span_context,
        }
    }
}

impl opentelemetry::trace::Span for Span {
    fn add_event_with_timestamp<T>(
        &mut self,
        _name: T,
        _timestamp: std::time::SystemTime,
        _attributes: Vec<opentelemetry::KeyValue>,
    ) where
        T: Into<std::borrow::Cow<'static, str>>,
    {
        todo!()
    }

    fn span_context(&self) -> &opentelemetry::trace::SpanContext {
        &self.span_context
    }

    fn is_recording(&self) -> bool {
        todo!()
    }

    fn set_attribute(&mut self, attribute: opentelemetry::KeyValue) {
        self.inner
            .set_attribute(attribute.key.as_str(), &attribute.value.as_str());
    }

    fn set_status(&mut self, _status: opentelemetry::trace::Status) {
        todo!()
    }

    fn update_name<T>(&mut self, _new_name: T)
    where
        T: Into<std::borrow::Cow<'static, str>>,
    {
        todo!()
    }

    fn end_with_timestamp(&mut self, _timestamp: std::time::SystemTime) {
        // Note: This does not respect the timestamp
        self.inner.close();
    }

    fn add_link(
        &mut self,
        span_context: opentelemetry::trace::SpanContext,
        attributes: Vec<opentelemetry::KeyValue>,
    ) {
        todo!()
    }
}

impl Drop for Span {
    fn drop(&mut self) {
        self.inner.close();
    }
}
