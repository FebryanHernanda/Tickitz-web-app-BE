-- public.orders definition
-- Drop table
-- DROP TABLE public.orders;
CREATE TABLE
    public.orders (
        id serial4 NOT NULL,
        qr_code varchar(255) NULL,
        ispaid bool DEFAULT false NULL,
        isactive bool DEFAULT true NULL,
        total_prices numeric(12, 2) NULL,
        created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
        updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
        user_id int4 NULL,
        cinemas_schedule_id int4 NULL,
        payment_method_id int4 NULL,
        CONSTRAINT orders_pkey PRIMARY KEY (id),
        CONSTRAINT orders_cinemas_schedule_id_fkey FOREIGN KEY (cinemas_schedule_id) REFERENCES public.cinemas_schedules (id) ON DELETE CASCADE,
        CONSTRAINT orders_payment_method_id_fkey FOREIGN KEY (payment_method_id) REFERENCES public.payment_methods (id) ON DELETE SET NULL,
        CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users (id) ON DELETE CASCADE
    );