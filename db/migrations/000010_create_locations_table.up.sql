-- public.locations definition
-- Drop table
-- DROP TABLE public.locations;
CREATE TABLE
    public.locations (
        id serial4 NOT NULL,
        "name" varchar(100) NOT NULL,
        CONSTRAINT locations_pkey PRIMARY KEY (id)
    );