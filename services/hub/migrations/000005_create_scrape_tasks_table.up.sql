CREATE TABLE IF NOT EXISTS scrape_tasks (
	id UUID PRIMARY KEY,
	completed BOOLEAN,
	created_at TIMESTAMP NOT NULL,
	scheduled_for TIMESTAMP NOT NULL,
	product_location_id UUID NOT NULL,
	repeat BOOLEAN,
	interval INTERVAL,
	CONSTRAINT fk_product_location
	FOREIGN KEY(product_location_id)
	REFERENCES product_locations(id)
);

CREATE INDEX index_upcoming_tasks ON scrape_tasks USING btree(scheduled_for) WHERE completed = FALSE;
