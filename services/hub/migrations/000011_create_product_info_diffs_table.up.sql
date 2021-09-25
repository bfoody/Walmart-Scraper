CREATE TABLE IF NOT EXISTS product_info_diffs (
	id UUID PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL,
	product_id UUID NOT NULL,
	product_location_id UUID NOT NULL,
	old_timestamp TIMESTAMPTZ NOT NULL,
	new_timestamp TIMESTAMPTZ NOT NULL,
	old_product_info_id UUID NOT NULL,
	new_product_info_id UUID NOT NULL,
	old_in_stock_value BOOLEAN NOT NULL,
	new_in_stock_value BOOLEAN NOT NULL,
	CONSTRAINT fk_product
	FOREIGN KEY(product_id)
	REFERENCES products(id),
	CONSTRAINT fk_product_location
	FOREIGN KEY(product_location_id)
	REFERENCES product_locations(id),
	CONSTRAINT fk_old_product_info
	FOREIGN KEY(old_product_info_id)
	REFERENCES product_infos(id),
	CONSTRAINT fk_new_product_info
	FOREIGN KEY(new_product_info_id)
	REFERENCES product_infos(id)
);

CREATE INDEX index_product_info_diffs_by_product_location_id ON product_info_diffs USING btree(product_location_id);

CREATE INDEX index_product_info_diffs_sell_outs ON product_info_diffs USING btree(product_location_id) WHERE old_in_stock_value = TRUE AND new_in_stock_value = FALSE;

CREATE INDEX index_product_info_diffs_restocks ON product_info_diffs USING btree(product_location_id) WHERE old_in_stock_value = FALSE AND new_in_stock_value = TRUE;
