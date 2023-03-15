CREATE TABLE IF NOT EXISTS public.cryptocurrency (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS public.request (
    user_id BIGINT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    currency_id UUID NOT NULL,
    constraint currency_id FOREIGN KEY (currency_id) REFERENCES public.cryptocurrency (id)
);
