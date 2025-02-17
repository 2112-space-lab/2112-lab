-- Insert test data into contexts
INSERT INTO config_schema.contexts (id, name, tenant_id, description, max_satellite, max_tiles)
VALUES 
    (uuid_generate_v4(), 'Context A', 'Tenant1', 'Test context A', 10, 50),
    (uuid_generate_v4(), 'Context B', 'Tenant2', 'Test context B', 20, 100),
    (uuid_generate_v4(), 'Context C', 'Tenant3', 'Test context C', 5, 25);

-- Insert test data into satellites
INSERT INTO config_schema.satellites (id, name, space_id, type, launch_date, owner, object_type, period, inclination, apogee, perigee, rcs, altitude)
VALUES 
    (uuid_generate_v4(), 'Satellite Alpha', '12345', 'Type1', '2022-01-01', 'NASA', 'LEO', 90.5, 45.6, 500, 400, 1.2, 450),
    (uuid_generate_v4(), 'Satellite Beta', '67890', 'Type2', '2023-02-01', 'ESA', 'GEO', 1436.0, 0.0, 35786, 35786, 2.5, 35786),
    (uuid_generate_v4(), 'Satellite Gamma', '54321', 'Type3', '2021-06-15', 'ISRO', 'MEO', 720.0, 20.1, 20000, 18000, 3.1, 19000);

-- Insert test data into TLE
INSERT INTO config_schema.tles (id, space_id, line1, line2, epoch)
VALUES 
    (uuid_generate_v4(), '12345', '1 25544U 98067A   23347.59899306  .00000123  00000-0  10270-4 0  9998', '2 25544  51.6446 200.5151 0001060  26.5673 333.5519 15.49397893260018', NOW()),
    (uuid_generate_v4(), '67890', '1 67890U 98067A   23347.59899306  .00000123  00000-0  10270-4 0  9998', '2 67890  51.6446 200.5151 0001060  26.5673 333.5519 15.49397893260018', NOW());

-- Insert test data into tiles
INSERT INTO config_schema.tiles (id, quadkey, zoom_level, center_lat, center_lon, nb_faces, radius, boundaries_json, spatial_index)
VALUES 
    (uuid_generate_v4(), '123123', 10, 37.7749, -122.4194, 4, 100.5, '{"boundary": "sample_data"}', ST_GeomFromText('POLYGON((-122.42 37.77, -122.41 37.78, -122.40 37.77, -122.42 37.77))', 4326)),
    (uuid_generate_v4(), '456456', 12, 48.8566, 2.3522, 6, 150.0, '{"boundary": "sample_data_2"}', ST_GeomFromText('POLYGON((2.35 48.85, 2.36 48.86, 2.37 48.85, 2.35 48.85))', 4326));

-- Insert test data into tile_satellite_mappings
INSERT INTO config_schema.tile_satellite_mappings (id, space_id, tile_id, tle_id, context_id, intersection_latitude, intersection_longitude, intersected_at)
VALUES 
    (uuid_generate_v4(), '12345', (SELECT id FROM config_schema.tiles LIMIT 1), (SELECT id FROM config_schema.tles WHERE space_id = '12345' LIMIT 1), (SELECT id FROM config_schema.contexts LIMIT 1), 37.7749, -122.4194, NOW()),
    (uuid_generate_v4(), '67890', (SELECT id FROM config_schema.tiles LIMIT 1 OFFSET 1), (SELECT id FROM config_schema.tles WHERE space_id = '67890' LIMIT 1), (SELECT id FROM config_schema.contexts LIMIT 1 OFFSET 1), 48.8566, 2.3522, NOW());

-- Insert test data into context_satellites (Many-to-Many relationship)
INSERT INTO config_schema.context_satellites (context_id, satellite_id)
VALUES 
    ((SELECT id FROM config_schema.contexts LIMIT 1), (SELECT id FROM config_schema.satellites LIMIT 1)),
    ((SELECT id FROM config_schema.contexts LIMIT 1 OFFSET 1), (SELECT id FROM config_schema.satellites LIMIT 1 OFFSET 1));

-- Insert test data into context_tles
INSERT INTO config_schema.context_tles (context_id, tle_id)
VALUES 
    ((SELECT id FROM config_schema.contexts LIMIT 1), (SELECT id FROM config_schema.tles LIMIT 1)),
    ((SELECT id FROM config_schema.contexts LIMIT 1 OFFSET 1), (SELECT id FROM config_schema.tles LIMIT 1 OFFSET 1));

-- Insert test data into context_tiles
INSERT INTO config_schema.context_tiles (context_id, tile_id)
VALUES 
    ((SELECT id FROM config_schema.contexts LIMIT 1), (SELECT id FROM config_schema.tiles LIMIT 1)),
    ((SELECT id FROM config_schema.contexts LIMIT 1 OFFSET 1), (SELECT id FROM config_schema.tiles LIMIT 1 OFFSET 1));