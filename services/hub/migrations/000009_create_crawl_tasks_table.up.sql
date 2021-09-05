CREATE TABLE IF NOT EXISTS crawl_tasks (
	id UUID PRIMARY KEY,
	completed BOOLEAN,
	created_at TIMESTAMPTZ NOT NULL,
	origin_product_location_id UUID NOT NULL,
	CONSTRAINT fk_product_location
	FOREIGN KEY(origin_product_location_id)
	REFERENCES product_locations(id)
);

CREATE UNIQUE INDEX index_completed_crawl_tasks ON crawl_tasks USING btree(origin_product_location_id) WHERE completed = TRUE;
