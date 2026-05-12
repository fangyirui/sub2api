-- Composite index to speed up GetAllGroupUsageSummary query
-- which aggregates actual_cost by group_id with a created_at filter for today's cost.
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_usage_logs_group_id_created_at ON usage_logs(group_id, created_at);
