mod conversion;
mod processor;
mod propagation;

pub use processor::WasiProcessor;
pub use propagation::TraceContextPropagator;
pub use propagation::WasiPropagator;

#[doc(hidden)]
/// Module containing wit bindgen generated code.
///
/// This is only meant for internal consumption.
mod wit {
    #![allow(missing_docs)]
    #![allow(clippy::missing_safety_doc)]
    wit_bindgen::generate!({
        world: "wasi:otel/imports@0.2.0-draft",
        path: "../wit",
        generate_all,
    });
}
