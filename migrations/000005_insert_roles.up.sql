INSERT INTO public.roles(
	role_name, role_description)
	VALUES ('user', 'for normal user'),
	('admin','for admins'),
	('analytic','for analytics');
commit;

UPDATE users SET role_id=1 WHERE role_id IS NULL;
alter table users alter column role_id set not null;
