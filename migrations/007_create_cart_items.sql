CREATE TABLE public.cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    sub_total NUMERIC(12,2) NOT NULL CHECK (sub_total >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT cart_items_cart_id_fkey
        FOREIGN KEY (cart_id)
        REFERENCES public.carts(id)
        ON DELETE CASCADE,
    CONSTRAINT cart_items_product_id_fkey
        FOREIGN KEY (product_id)
        REFERENCES public.products(id)
        ON DELETE CASCADE
);