#!/bin/bash
echo "Creating tables..."
psql -U jetspotter -h localhost -p 5432 <<EOF
CREATE TABLE IF NOT EXISTS aircraft (
    aircraft_id SERIAL PRIMARY KEY,
    callsign VARCHAR(50),
    type VARCHAR(100),
    description VARCHAR(255),
    tail_number VARCHAR(50),
    icao VARCHAR(50),
    image_url VARCHAR(255),
    military BOOLEAN
);

CREATE TABLE IF NOT EXISTS spot_configurations (
    spot_configuration_id SERIAL PRIMARY KEY,
    lattitude NUMERIC,
    longitude NUMERIC,
    max_range_kilometers INTEGER
);

CREATE TABLE IF NOT EXISTS spots (
    spot_id SERIAL PRIMARY KEY,
    aircraft_id int REFERENCES aircraft(aircraft_id),
    spot_configuration_id int REFERENCES spot_configurations(spot_configuration_id),
    spot_date DATE
);

CREATE TABLE IF NOT EXISTS notification_types (
    notification_type_id SERIAL PRIMARY KEY,
    name VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS notifications (
    notification_id SERIAL PRIMARY KEY,
    spot_id int REFERENCES spots(spot_id),
    notification_type_id int REFERENCES notification_types(notification_type_id)
);

CREATE TABLE IF NOT EXISTS notification_configurations (
    notification_configuration_id SERIAL PRIMARY KEY,
    notification_type_id int REFERENCES notification_types(notification_type_id),
    webhook_url VARCHAR(255),
    gotify_token VARCHAR(255),
    ntfy_topic VARCHAR(255),
    ntfy_server VARCHAR(255),
    max_range_kilometers INTEGER,
    max_altitude_feet INTEGER
);

CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email VARCHAR(100),
    password VARCHAR(255),
    notification_configuration_id int REFERENCES notification_configurations(notification_configuration_id)
);

CREATE EXTENSION IF NOT EXISTS pgcrypto;
EOF

echo "Inserting aircraft test data..."
psql -U jetspotter -h localhost -p 5432 <<EOF
INSERT INTO aircraft (
    callsign,
    type,
    description,
    tail_number,
    icao,
    image_url,
    military
)
VALUES
    (
        'APEX-11',
        'F16',
        'Multirole fighter jet',
        'J-146',
        'ABCDEF1234',
        'https://www.f-16.net/g3/f-16-photos/album37/album13/baw',
        TRUE
    ),
    (
        'GRIZZLY-21',
        'A400',
        'Military transport aircraft',
        'G-273',
        'GHIJKL5678',
        'https://www.airhistory.net/photo/713163/CT-03',
        TRUE
    ),
    (
        'CES001',
        'CESSNA',
        'Light aircraft',
        'N12345',
        'MNOPQR9012',
        'https://www.airhistory.net/photo/713163/CT-03',
        FALSE
    );

EOF
