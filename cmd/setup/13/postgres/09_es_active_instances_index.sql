CREATE INDEX IF NOT EXISTS es_active_instances ON eventstore.events (created_at DESC, instance_id);