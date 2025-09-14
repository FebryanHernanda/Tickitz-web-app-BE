-- public.seats definition
-- Drop table
-- DROP TABLE public.seats;
CREATE TABLE
    public.seats (
        id serial4 NOT NULL,
        seat_number varchar(5) NOT NULL,
        seat_type varchar(20) NULL,
        created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
        updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
        CONSTRAINT seats_pkey PRIMARY KEY (id)
    );