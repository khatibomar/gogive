ALTER TABLE users ALTER COLUMN role_id DROP NOT NULL;

UPDATE users SET role_id=NULL where role_id IN (1,2,3);

DELETE FROM roles
	WHERE role_id=1;
DELETE FROM roles
	WHERE role_id=2;
DELETE FROM roles
	WHERE role_id=3;
