ALTER TABLE users ALTER COLUMN role_id DROP NOT NULL;

DELETE FROM roles
	WHERE role_id=1;
DELETE FROM roles
	WHERE role_id=2;
DELETE FROM roles
	WHERE role_id=3;
