alter table users alter column role_id set not null;

INSERT INTO public.roles(
	role_name, role_description)
	VALUES ('user', 'for normal user'),
	('admin','for admins'),
	('analytic','for analytics');
