INSERT INTO categories(category_name)
	VALUES ('vehicules'),
	('mobile phones and accessories'),
	('electronics'),
	('fashion'),
	('pets'),
	('kids and babies'),
	('services'),
	('other')
	on conflict do nothing;
