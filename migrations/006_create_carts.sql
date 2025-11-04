CREATE TABLE public.carts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    warehouse_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT carts_user_id_warehouse_id_unique UNIQUE(user_id, warehouse_id),
    CONSTRAINT carts_user_id_fkey
        FOREIGN KEY (user_id)
        REFERENCES public.users(id)
        ON DELETE CASCADE,
    CONSTRAINT carts_warehouse_id_fkey
        FOREIGN KEY (warehouse_id)
        REFERENCES public.warehouses(id)
        ON DELETE CASCADE
);