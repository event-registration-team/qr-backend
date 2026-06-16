DROP TRIGGER IF EXISTS set_updated_at_admins ON admins;
DROP TABLE IF EXISTS admins CASCADE;
DROP FUNCTION IF EXISTS trigger_set_updated_at;