-- Agencies are the organisations that own resources, e.g. a fire department,
-- an ambulance service or the police.
CREATE TABLE agencies (
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    canon_id       text UNIQUE,            -- canonical identifier shared with other systems
    name           varchar(100) NOT NULL,
    agency_type    text,                   -- e.g. 'fire', 'ems', 'police'
    contact_info   jsonb,

    created_at_utc timestamptz NOT NULL DEFAULT now()
);
