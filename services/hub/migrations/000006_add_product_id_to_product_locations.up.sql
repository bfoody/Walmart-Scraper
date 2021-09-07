ALTER TABLE product_locations ADD product_id UUID;
ALTER TABLE product_locations ADD CONSTRAINT fk_product
	FOREIGN KEY(product_id)
	REFERENCES products(id);
