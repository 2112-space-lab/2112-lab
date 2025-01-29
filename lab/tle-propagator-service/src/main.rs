mod tle;
mod redis_utils;
mod propagation;

use std::sync::Arc;
use redis::Client as RedisClient;
use tokio::signal;
use tracing::{info, error};
use tracing_subscriber;

#[tokio::main]
async fn main() {
    // Initialize logging
    tracing_subscriber::fmt::init();
    info!("Starting tle-propagator-service...");

    // Load configuration
    let redis_host = std::env::var("REDIS_HOST").unwrap_or_else(|_| "redis-service".to_string());
    let redis_port = std::env::var("REDIS_PORT").unwrap_or_else(|_| "6379".to_string());
    let redis_url = format!("redis://{}:{}/", redis_host, redis_port);

    // Connect to Redis
    let redis_client = match RedisClient::open(redis_url.clone()) {
        Ok(client) => Arc::new(client),
        Err(e) => {
            error!("Failed to connect to Redis at {}: {:?}", redis_url, e);
            return;
        }
    };

    info!("Connected to Redis at {}", redis_url);

    let worker_count = 4; // Number of workers to process TLE updates

    // Run the subscription service and listen for termination signals
    tokio::select! {
        _ = propagation::subscribe_to_tle_updates(redis_client.clone(), worker_count) => {
            error!("Subscription service unexpectedly exited");
        }
        _ = signal::ctrl_c() => {
            info!("Termination signal received. Shutting down gracefully...");
        }
    }

    info!("tle-propagator-service exited.");
}
