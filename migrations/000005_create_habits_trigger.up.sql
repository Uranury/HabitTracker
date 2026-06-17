CREATE TRIGGER IF NOT EXISTS update_habits_updated_at
AFTER UPDATE ON habits
FOR EACH ROW
BEGIN
UPDATE habits SET updated_at = strftime('%Y-%m-%dT%H:%M:%SZ', 'now') WHERE id = OLD.id;
END;