-- Resources are the units an agency can send to an incident, e.g. an ambulance or a fire truck.
--
-- status only says whether the agency has the resource in service. Whether it is busy on an
-- incident is NOT stored here: it is derived from having an active row in allocations, so the
-- two can never disagree.
CREATE TABLE resources (
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    canon_id       text UNIQUE,            -- canonical identifier shared with other systems
    agency_id      uuid NOT NULL REFERENCES agencies (id),

    name           text NOT NULL,          -- e.g. 'Ambulance 12'
    type_canon     text NOT NULL,          -- canonical type, e.g. 'ambulance'
    type_internal  text,                   -- type as named in the agency's own system
    status         text NOT NULL DEFAULT 'available',
    attributes     jsonb,                  -- extra details, e.g. capacity or equipment

    -- Home position (where the resource is stationed). Live positions are stored separately.
    home_latitude  double precision,
    home_longitude double precision,

    created_at_utc timestamptz NOT NULL DEFAULT now(),
    updated_at_utc timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT resources_status_valid CHECK (status IN ('available', 'out_of_service')),
    CONSTRAINT resources_home_pair CHECK ((home_latitude IS NULL) = (home_longitude IS NULL))
);

-- Speeds up "all resources of agency X".
CREATE INDEX resources_agency_id_idx ON resources (agency_id);
