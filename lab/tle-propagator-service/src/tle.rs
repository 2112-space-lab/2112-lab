use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize)]
pub struct TLEUpdate {
    pub id: String,
    pub line_1: String,
    pub line_2: String,
    pub epoch: String,
}

#[derive(Debug, Serialize)]
pub struct SatellitePosition {
    pub id: String,
    pub timestamp: String,
    pub latitude: f64,
    pub longitude: f64,
    pub altitude: f64,
}
