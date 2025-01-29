use std::sync::Arc;
use std::time::{SystemTime, Duration};
use redis::Client as RedisClient;
use serde_json::json;
use tokio::sync::mpsc;
use tracing::{info, error, debug};
use chrono::{DateTime, Utc};
use sgp4::{Constants, MinutesSinceEpoch};
use crate::tle::{TLEUpdate, SatellitePosition};
use crate::redis_utils::{publish_to_redis, store_to_redis};

/// Parse eccentricity from the raw TLE string and interpret it correctly.
fn parse_eccentricity(tle_line2: &str) -> f64 {
    tle_line2[26..33]
        .trim()
        .parse::<f64>()
        .unwrap_or(0.0)
        / 10_000_000.0 // Interpret as 0.xxxxxxxx
}

/// Extract inclination from TLE second line as a float.
fn parse_inclination(tle_line2: &str) -> f64 {
    tle_line2[8..16].trim().parse::<f64>().unwrap_or(0.0)
}

/// Convert SystemTime to a readable ISO 8601 timestamp.
fn system_time_to_iso8601(system_time: SystemTime) -> String {
    let datetime: DateTime<Utc> = system_time.into();
    datetime.to_rfc3339()
}

/// Calculate latitude, longitude, and altitude from Cartesian coordinates.
fn calculate_lat_lon_alt(position: [f64; 3]) -> (f64, f64, f64) {
    let [x, y, z] = position;
    let r = (x * x + y * y + z * z).sqrt();
    let latitude = (z / r).asin().to_degrees();
    let longitude = y.atan2(x).to_degrees();
    let altitude = r - 6371.0; // Earth's radius in km
    (latitude, longitude, altitude)
}

/// Propagate satellite positions using raw TLE strings.
pub async fn propagate_satellite_positions(
    satellite_id: &str,
    tle_line1: &str,
    tle_line2: &str,
    start_time: SystemTime,
) -> Vec<SatellitePosition> {
    let mut positions = Vec::new();
    let mut failures = Vec::new();

    debug!(
        "Processing TLE for satellite {}: line1={}, line2={}",
        satellite_id, tle_line1, tle_line2
    );

    // Correct TLE line2 if necessary
    let eccentricity = parse_eccentricity(tle_line2);
    let inclination = parse_inclination(tle_line2);

    debug!(
        "Parsed TLE for satellite {}: eccentricity={}, inclination={}°",
        satellite_id, eccentricity, inclination
    );

    // Parse the TLE data
    let elements = match sgp4::Elements::from_tle(
        Some(satellite_id.to_owned()),
        tle_line1.as_bytes(),
        tle_line2.as_bytes(),
    ) {
        Ok(e) => e,
        Err(e) => {
            error!("Invalid TLE for satellite {}: {:?}", satellite_id, e);
            return positions;
        }
    };

    // Compute constants for propagation
    let constants = match Constants::from_elements(&elements) {
        Ok(c) => c,
        Err(e) => {
            error!("Failed to create constants for satellite {}: {:?}", satellite_id, e);
            return positions;
        }
    };

    let mut current_time = start_time;

    debug!(
        "Starting propagation for satellite {} from {:?}",
        satellite_id,
        system_time_to_iso8601(start_time)
    );

    while let Ok(duration_since_start) = current_time.duration_since(start_time) {
        let timestamp = duration_since_start.as_secs_f64();

        // Calculate minutes since start_time
        let minutes_since_start = MinutesSinceEpoch(timestamp / 60.0);

        match constants.propagate(minutes_since_start) {
            Ok(prediction) => {
                let (latitude, longitude, altitude) = calculate_lat_lon_alt(prediction.position);

                // Determine the duration and number of points based on altitude
                let (pass_duration_minutes, interval_seconds) = determine_sampling_params(altitude);

                positions.push(SatellitePosition {
                    id: satellite_id.to_string(),
                    timestamp: system_time_to_iso8601(current_time),
                    latitude,
                    longitude,
                    altitude,
                });

                current_time += Duration::from_secs(interval_seconds);

                // Stop propagating if the duration exceeds the pass duration
                if duration_since_start.as_secs() / 60 >= pass_duration_minutes {
                    break;
                }
            }
            Err(e) => {
                error!("Propagation error for satellite {}: {:?}", satellite_id, e);
                failures.push(current_time);
                break;
            }
        }
    }

    if failures.is_empty() {
        info!(
            "Successfully propagated satellite {}: {} positions generated.",
            satellite_id,
            positions.len()
        );
    } else {
        let total_failures = failures.len();
        error!(
            "Propagation failed for satellite {}: {} failures.",
            satellite_id, total_failures
        );
    }

    positions
}

/// Determine sampling parameters based on altitude.
fn determine_sampling_params(altitude_km: f64) -> (u64, u64) {
    if altitude_km < 2000.0 {
        // Low Earth Orbit (LEO)
        (1440, 60) // 1 day duration, 1-minute intervals
    } else if altitude_km < 35786.0 {
        // Medium Earth Orbit (MEO)
        (180, 216) // 3 hours duration, ~3.6-minute intervals
    } else {
        // Geostationary Earth Orbit (GEO)
        (1440, 8640) // 1 day duration, ~2.4-hour intervals
    }
}

/// Process a single TLE update, propagate positions, and store them in Redis.
pub async fn process_tle_update(mut con: redis::Connection, tle_update: TLEUpdate) {
    info!("Received TLE update for satellite {}", tle_update.id);

    let start_time = SystemTime::now();

    let positions = propagate_satellite_positions(
        &tle_update.id,
        &tle_update.line_1,
        &tle_update.line_2,
        start_time
    )
    .await;

    for position in &positions {
        // Log the position details
        debug!("Processing position for satellite {}: {:?}", tle_update.id, position);

        // Serialize the position into JSON
        let position_json = match serde_json::to_value(position) {
            Ok(json) => json,
            Err(e) => {
                error!("Failed to serialize position for satellite {}: {:?}, position: {:?}", tle_update.id, e, position);
                break; // Break on serialization error
            }
        };

        // Parse the timestamp as a score
        let score = match position.timestamp.parse::<DateTime<Utc>>() {
            Ok(dt) => dt.timestamp(), // Use seconds as an integer
            Err(e) => {
                error!(
                    "Invalid timestamp for satellite {}: {:?}, position: {:?}",
                    tle_update.id, e, position
                );
                break;
            }
        };

        // Attempt to store the position in Redis
        if let Err(e) = store_to_redis(
            &format!("satellite_positions:{}", tle_update.id),
            &position_json,
            score,
            &mut con,
        ) {
            error!(
                "Failed to store position in Redis for satellite {}: {:?}, position: {:?}",
                tle_update.id, e, position
            );
            break; // Break on Redis storage error
        }

    }

    let summary = json!({
        "event": "event_satellite_positions_updated",
        "satellite_id": tle_update.id,
        "start_time": positions.first().map(|p| &p.timestamp),
        "end_time": positions.last().map(|p| &p.timestamp),
        "positions_count": positions.len(),
    });

    match publish_to_redis("event_satellite_positions_updated", &summary, &mut con) {
        Ok(_) => {
            info!(
                "Successfully published update for satellite {} to Redis.",
                tle_update.id
            );
        }
        Err(e) => {
            error!("Failed to publish summary for satellite {}: {:?}", tle_update.id, e);
        }
    }
}

/// Subscribe to Redis channel for TLE updates and process them concurrently.
pub async fn subscribe_to_tle_updates(redis_client: Arc<RedisClient>, max_threads: usize) {
    let mut con = match redis_client.get_connection() {
        Ok(c) => c,
        Err(e) => {
            error!("Failed to connect to Redis: {:?}", e);
            return;
        }
    };

    let mut pubsub = con.as_pubsub();

    if let Err(e) = pubsub.subscribe("satellite_tle_updates") {
        error!("Failed to subscribe to channel: {:?}", e);
        return;
    }

    let (tx, mut rx) = mpsc::channel::<TLEUpdate>(max_threads);

    tokio::spawn(async move {
        while let Some(tle_update) = rx.recv().await {
            let redis_client = Arc::clone(&redis_client);
            tokio::spawn(async move {
                let mut con = match redis_client.get_connection() {
                    Ok(c) => c,
                    Err(e) => {
                        error!("Failed to connect to Redis: {:?}", e);
                        return;
                    }
                };

                process_tle_update(con, tle_update).await;
            });
        }
    });

    loop {
        if let Ok(message) = pubsub.get_message() {
            if let Ok(payload) = message.get_payload::<String>() {
                match serde_json::from_str::<TLEUpdate>(&payload) {
                    Ok(tle_update) => {
                        info!("Processing TLE update payload: {:?}", tle_update);
                        if tx.send(tle_update).await.is_err() {
                            error!("Failed to send TLE update to workers");
                        }
                    }
                    Err(e) => {
                        error!("Failed to deserialize TLE update: {:?}", e);
                    }
                }
            }
        }
    }
}
