CREATE TABLE IF NOT EXISTS average_sellouts (
	id UUID PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
	product_id UUID NOT NULL UNIQUE,
	product_location_id UUID NOT NULL,
	average_availability_duration BIGINT NOT NULL,
	averaged_count BIGINT NOT NULL,
	CONSTRAINT fk_product
	FOREIGN KEY(product_id)
	REFERENCES products(id),
	CONSTRAINT fk_product_location
	FOREIGN KEY(product_location_id)
	REFERENCES product_locations(id)
);

CREATE INDEX index_average_sellouts_sorted ON average_sellouts(average_availability_duration ASC);
CREATE UNIQUE INDEX index_average_sellouts_on_product_id ON average_sellouts USING btree(product_id);
