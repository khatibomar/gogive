CREATE INDEX IF NOT EXISTS items_name_idx ON items USING GIN (to_tsvector('simple', name));
