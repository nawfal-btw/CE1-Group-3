-- Incidents table for the MVP. Columns mirror incidentRow in
-- internal/repository/postgres/incidents.go. NULL means "not set" in the domain.
CREATE TABLE incidents (
    id                   uuid PRIMARY KEY,
    type                 smallint,
    status               smallint,
    priority             smallint,

    latitude             double precision,
    longitude            double precision,
    address              text,
    proximate_landmarks  jsonb,
    location_description text,

    description          text,

    created_at_utc       timestamptz,
    updated_at_utc       timestamptz,
    closed_at_utc        timestamptz,

    version              text NOT NULL DEFAULT '',

    -- Coordinates are either both present or both absent.
    CONSTRAINT incidents_gis_pair CHECK ((latitude IS NULL) = (longitude IS NULL))
);
