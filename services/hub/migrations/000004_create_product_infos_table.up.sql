CREATE TABLE IF NOT EXISTS product_infos (
	id UUID PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL,
	product_id UUID NOT NULL,
	product_location_id UUID NOT NULL,
	price DECIMAL,
	availability_status TEXT,
	in_stock BOOLEAN,
	CONSTRAINT fk_product
	FOREIGN KEY(product_id)
	REFERENCES products(id),
	CONSTRAINT fk_product_location
	FOREIGN KEY(product_location_id)
	REFERENCES product_locations(id)
);

CREATE INDEX index_product_infos_by_product_id ON product_infos USING btree(product_id);

CREATE INDEX index_product_infos_by_product_location_id ON product_infos USING btree(product_location_id);
