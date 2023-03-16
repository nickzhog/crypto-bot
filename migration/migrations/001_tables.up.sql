CREATE TABLE IF NOT EXISTS public.cryptocurrency (
    currency_name TEXT PRIMARY KEY,
    last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    price TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS public.request (
    user_id BIGINT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    currency_name TEXT NOT NULL,
    price TEXT NOT NULL,
    constraint currency_name FOREIGN KEY (currency_name) REFERENCES public.cryptocurrency (currency_name)
);
