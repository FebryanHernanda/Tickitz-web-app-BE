-- public.casts definition
-- Drop table
-- DROP TABLE public.casts;
CREATE TABLE
    public.casts (
        id serial4 NOT NULL,
        "name" varchar(100) NOT NULL,
        CONSTRAINT casts_pkey PRIMARY KEY (id)
    );