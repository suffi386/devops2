CREATE INDEX CONCURRENTLY IF NOT EXISTS user_sessions_by_user ON auth.user_sessions (instance_id, user_id);