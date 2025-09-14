-- public.payment_methods definition
-- Drop table
-- DROP TABLE public.payment_methods;
CREATE TABLE
    public.payment_methods (
        id serial4 NOT NULL,
        "name" varchar(50) NOT NULL,
        provider varchar(50) NOT NULL,
        CONSTRAINT payment_methods_pkey PRIMARY KEY (id)
    );