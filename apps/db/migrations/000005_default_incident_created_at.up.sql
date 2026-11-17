-- The database records when each incident was created. The application does not send
-- created_at_utc, so the column takes this default on every insert.
-- Incidents created before this migration keep NULL: their real creation time is unknown.
ALTER TABLE incidents ALTER COLUMN created_at_utc SET DEFAULT now();
