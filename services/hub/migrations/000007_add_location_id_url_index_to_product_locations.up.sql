CREATE INDEX index_product_locations_location_id_url ON product_locations USING btree(location_id, url);
