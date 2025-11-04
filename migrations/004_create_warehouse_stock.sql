CREATE TABLE public.warehouse_stock (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    CONSTRAINT fk_product 
        FOREIGN KEY (product_id) 
        REFERENCES public.products(id),
    CONSTRAINT fk_warehouse 
        FOREIGN KEY (warehouse_id) 
        REFERENCES public.warehouses(id)
);
