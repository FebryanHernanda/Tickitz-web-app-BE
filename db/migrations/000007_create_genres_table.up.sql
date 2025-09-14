-- public.genres definition
-- Drop table
-- DROP TABLE public.genres;
CREATE TABLE
    public.genres (
        id serial4 NOT NULL,
        "name" varchar(50) NOT NULL,
        CONSTRAINT genres_pkey PRIMARY KEY (id)
    );