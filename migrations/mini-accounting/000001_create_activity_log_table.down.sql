DO $$
BEGIN
   DROP INDEX IF EXISTS idx_trace_id;
   DROP TABLE IF EXISTS action_log;
END $$;
