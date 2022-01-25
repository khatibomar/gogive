ALTER TABLE users ADD CONSTRAINT unique_user_email UNIQUE (email);
ALTER TABLE users
ALTER COLUMN created_at TYPE TIMESTAMP(0) with time zone,
ALTER COLUMN created_at SET DEFAULT NOW();

ALTER TABLE items 
ALTER COLUMN created_at TYPE TIMESTAMP(0) with time zone,
ALTER COLUMN created_at SET DEFAULT NOW();

ALTER TABLE bans ADD CONSTRAINT unique_ban_email UNIQUE (email);

ALTER TABLE roles ADD CONSTRAINT unique_role UNIQUE (role_name);
ALTER TABLE permissions ADD CONSTRAINT unique_perm UNIQUE (code);
ALTER TABLE categories ADD CONSTRAINT unique_cat UNIQUE (category_name);
