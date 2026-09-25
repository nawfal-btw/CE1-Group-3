-- name: GetAllIncidents :many
SELECT
    *
FROM
    incidents;

-- name: GetIncidentByID :one
SELECT
    *
FROM
    incidents
WHERE
    id = $1
LIMIT 1;

-- name: AddIncident :one
INSERT INTO incidents (latitude, longitude)
    VALUES ($1, $2)
RETURNING
    *;

-- name: DeleteIncidentByID :one
DELETE FROM incidents
WHERE id = $1
RETURNING
    *;

