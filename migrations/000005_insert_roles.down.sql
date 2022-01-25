ALTER TABLE users ALTER COLUMN role_id DROP NOT NULL;

UPDATE users SET role_id=NULL where role_id IN (select role_id from users limit 3);

DELETE FROM roles WHERE role_id IN (select role_id from users limit 3);
