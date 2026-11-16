-- An allocation links one resource to one incident. Resources from different agencies can be
-- allocated to the same incident. Rows are never deleted: releasing a resource sets
-- status = 'released', so the full history of who was sent where is kept.
CREATE TABLE allocations (
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id      uuid NOT NULL REFERENCES incidents (id),
    resource_id      uuid NOT NULL REFERENCES resources (id),

    status           text NOT NULL DEFAULT 'dispatched',
    eta_utc          timestamptz,          -- estimated time of arrival

    allocated_at_utc timestamptz NOT NULL DEFAULT now(),
    released_at_utc  timestamptz,

    CONSTRAINT allocations_status_valid
        CHECK (status IN ('dispatched', 'en_route', 'on_scene', 'released')),
    -- released_at_utc is set exactly when the allocation is released.
    CONSTRAINT allocations_released_pair
        CHECK ((status = 'released') = (released_at_utc IS NOT NULL))
);

-- A resource can only be on one incident at a time: at most one non-released allocation per
-- resource. This is enforced by the database, so two dispatchers allocating the same resource
-- at the same moment cannot both succeed.
CREATE UNIQUE INDEX allocations_one_active_per_resource
    ON allocations (resource_id)
    WHERE status <> 'released';

-- Speeds up "all allocations of incident X".
CREATE INDEX allocations_incident_id_idx ON allocations (incident_id);
