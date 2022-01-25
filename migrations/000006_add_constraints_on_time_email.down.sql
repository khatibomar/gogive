ALTER TABLE users DROP CONSTRAINT unique_user_email;
ALTER TABLE users
ALTER COLUMN created_at TYPE DATE;

ALTER TABLE items 
ALTER COLUMN created_at TYPE DATE;

ALTER TABLE bans DROP CONSTRAINT unique_ban_email;

ALTER TABLE roles DROP CONSTRAINT unique_role;
ALTER TABLE permissions DROP CONSTRAINT unique_perm;
ALTER TABLE categories DROP CONSTRAINT unique_cat;
