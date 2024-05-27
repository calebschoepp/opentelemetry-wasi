pub mod exporter;
pub mod processor;
pub mod propagation;

#[doc(hidden)]
/// Module containing wit bindgen generated code.
///
/// This is only meant for internal consumption.
pub mod wit {
    #![allow(missing_docs)]
    #![allow(clippy::missing_safety_doc)]
    wit_bindgen::generate!({
        world: "platform",
        path: "./wit",
    });
    pub use fermyon::spin2_0_0 as v2;
}
