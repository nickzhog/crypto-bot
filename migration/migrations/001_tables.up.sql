CREATE TABLE IF NOT EXISTS public.cryptocurrency (
    name TEXT PRIMARY KEY,
    last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    price TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS public.request (
    user_id BIGINT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    currency_name UUID NOT NULL,
    constraint currency_name FOREIGN KEY (currency_name) REFERENCES public.cryptocurrency (name)
);
