ALTER TABLE items ADD CONSTRAINT categories_length_check CHECK (array_length(categories, 1) BETWEEN 1 AND 5);
