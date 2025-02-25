-- Enable required extensions
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp') THEN 
        CREATE EXTENSION "uuid-ossp";
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'postgis') THEN 
        CREATE EXTENSION "postgis";
    END IF;
END $$;

-- Enable schema
CREATE SCHEMA IF NOT EXISTS config_schema;

-- Create Context table
CREATE TABLE config_schema.contexts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    name VARCHAR(255) UNIQUE NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    description VARCHAR(1024),
    max_satellite INT NOT NULL,
    max_tiles INT NOT NULL,
    activated_at TIMESTAMP NULL,
    desactivated_at TIMESTAMP NULL,
    trigger_generated_mapping_at TIMESTAMP NULL,
    trigger_imported_tle_at TIMESTAMP NULL,
    trigger_imported_satellite_at TIMESTAMP NULL,
    CONSTRAINT unique_tenant_context_name UNIQUE (tenant_id, name)
);

-- Create Satellite table
CREATE TABLE config_schema.satellites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    name VARCHAR(255) NOT NULL,
    space_id VARCHAR(255) UNIQUE NOT NULL,
    type VARCHAR(255),
    launch_date DATE,
    decay_date DATE,
    intl_designator VARCHAR(255),
    owner VARCHAR(255),
    object_type VARCHAR(255),
    period FLOAT,
    inclination FLOAT,
    apogee FLOAT,
    perigee FLOAT,
    rcs FLOAT,
    altitude FLOAT
);

-- Create TLE table
CREATE TABLE config_schema.tles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    space_id VARCHAR(255) NOT NULL,
    line1 VARCHAR(255) NOT NULL,
    line2 VARCHAR(255) NOT NULL,
    epoch TIMESTAMP NOT NULL
);

CREATE INDEX tle_space_id_idx ON config_schema.tles(space_id);

-- Create Tile table
CREATE TABLE config_schema.tiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    quadkey VARCHAR(256) UNIQUE NOT NULL,
    zoom_level INT NOT NULL,
    center_lat FLOAT NOT NULL,
    center_lon FLOAT NOT NULL,
    nb_faces INT NOT NULL,
    radius FLOAT NOT NULL,
    boundaries_json JSON,
    spatial_index GEOMETRY(POLYGON, 4326)
);

-- Create TileSatelliteMapping table
CREATE TABLE config_schema.tile_satellite_mappings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    space_id VARCHAR(255) NOT NULL,
    tile_id UUID NOT NULL,
    tle_id UUID NOT NULL,
    context_id UUID NOT NULL,
    intersection_latitude DOUBLE PRECISION NOT NULL,
    intersection_longitude DOUBLE PRECISION NOT NULL,
    intersected_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_tile FOREIGN KEY (tile_id) REFERENCES config_schema.tiles(id) ON DELETE CASCADE,
    CONSTRAINT fk_context FOREIGN KEY (context_id) REFERENCES config_schema.contexts(id) ON DELETE CASCADE
);

-- Create ContextSatellite table
CREATE TABLE config_schema.context_satellites (
    context_id UUID NOT NULL,
    satellite_id UUID NOT NULL,
    PRIMARY KEY (context_id, satellite_id),
    CONSTRAINT fk_context_satellite_context FOREIGN KEY (context_id) REFERENCES config_schema.contexts(id) ON DELETE CASCADE,
    CONSTRAINT fk_context_satellite_satellite FOREIGN KEY (satellite_id) REFERENCES config_schema.satellites(id) ON DELETE CASCADE
);

-- Create ContextTLE table
CREATE TABLE config_schema.context_tles (
    context_id UUID NOT NULL,
    tle_id UUID NOT NULL,
    PRIMARY KEY (context_id, tle_id),
    CONSTRAINT fk_context_tle_context FOREIGN KEY (context_id) REFERENCES config_schema.contexts(id) ON DELETE CASCADE,
    CONSTRAINT fk_context_tle_tle FOREIGN KEY (tle_id) REFERENCES config_schema.tles(id) ON DELETE CASCADE
);

-- Create ContextTile table
CREATE TABLE config_schema.context_tiles (
    context_id UUID NOT NULL,
    tile_id UUID NOT NULL,
    PRIMARY KEY (context_id, tile_id),
    CONSTRAINT fk_context_tile_context FOREIGN KEY (context_id) REFERENCES config_schema.contexts(id) ON DELETE CASCADE,
    CONSTRAINT fk_context_tile_tile FOREIGN KEY (tile_id) REFERENCES config_schema.tiles(id) ON DELETE CASCADE
);

-- Create Event Table
CREATE TABLE config_schema.events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_time_utc TIMESTAMP NOT NULL DEFAULT NOW(),
    event_uid VARCHAR(255) UNIQUE NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    payload JSON NOT NULL,
    comment TEXT,
    published_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_events_event_type ON config_schema.events(event_type);
CREATE INDEX idx_events_event_uid ON config_schema.events(event_uid);

-- Create EventHandler Table
CREATE TABLE config_schema.event_handlers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL,
    handler_name VARCHAR(255) NOT NULL,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP NULL,
    status VARCHAR(50) NOT NULL,
    error_message TEXT NULL,
    FOREIGN KEY (event_id) REFERENCES config_schema.events(id) ON DELETE CASCADE
);

CREATE INDEX idx_event_handlers_event_id ON config_schema.event_handlers(event_id);
CREATE INDEX idx_event_handlers_handler_name ON config_schema.event_handlers(handler_name);

-- Add indexes
CREATE INDEX idx_tile_satellite_mappings_space_id ON config_schema.tile_satellite_mappings(space_id);
CREATE INDEX idx_tile_satellite_mappings_tile_id ON config_schema.tile_satellite_mappings(tile_id);
CREATE INDEX idx_tile_satellite_mappings_tle_id ON config_schema.tile_satellite_mappings(tle_id);
CREATE INDEX idx_tile_satellite_mappings_context_id ON config_schema.tile_satellite_mappings(context_id);
