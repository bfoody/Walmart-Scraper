CREATE TABLE IF NOT EXISTS product_locations (
	id UUID PRIMARY KEY,
	location_id UUID NOT NULL,
	url TEXT,
	slug TEXT,
	category TEXT,
	CONSTRAINT fk_location
	FOREIGN KEY(location_id)
	REFERENCES locations(id)
);
