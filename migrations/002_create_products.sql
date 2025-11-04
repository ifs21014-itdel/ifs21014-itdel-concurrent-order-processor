CREATE TABLE public.products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL,
    price NUMERIC(12,2) NOT NULL CHECK (price >= 0),
    stock INTEGER NOT NULL CHECK (stock >= 0),
    weight NUMERIC(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT products_user_id_fkey 
        FOREIGN KEY (user_id) 
        REFERENCES public.users(id) 
        ON DELETE CASCADE
);
