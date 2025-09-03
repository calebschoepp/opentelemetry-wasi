use crate::wit::wasi::clocks0_2_0::wall_clock::Datetime;
use std::time::UNIX_EPOCH;

impl From<std::time::SystemTime> for Datetime {
    fn from(value: std::time::SystemTime) -> Self {
        let duration_since_epoch = value
            .duration_since(UNIX_EPOCH)
            .expect("SystemTime should be after UNIX EPOCH");
        Self {
            seconds: duration_since_epoch.as_secs(),
            nanoseconds: duration_since_epoch.subsec_nanos(),
        }
    }
}
