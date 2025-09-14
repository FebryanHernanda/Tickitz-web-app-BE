-- public.cinemas definition
-- Drop table
-- DROP TABLE public.cinemas;
CREATE TABLE
    public.cinemas (
        id serial4 NOT NULL,
        "name" varchar(100) NOT NULL,
        prices numeric(10, 2) NULL,
        image_path text NULL,
        CONSTRAINT cinemas_pkey PRIMARY KEY (id)
    );