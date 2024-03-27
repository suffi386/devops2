ALTER TABLE IF EXISTS projections.users11_notifications ADD COLUMN IF NOT EXISTS verified_email_lower TEXT GENERATED ALWAYS AS (lower(verified_email)) STORED;
CREATE INDEX IF NOT EXISTS users11_notifications_email_search ON projections.users11_notifications (instance_id, verified_email_lower);
