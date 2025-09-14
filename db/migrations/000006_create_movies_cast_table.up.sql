-- public.movies_cast definition
-- Drop table
-- DROP TABLE public.movies_cast;
CREATE TABLE
    public.movies_cast (
        movie_id int4 NOT NULL,
        cast_id int4 NOT NULL,
        CONSTRAINT movies_cast_pkey PRIMARY KEY (movie_id, cast_id)
    );

-- public.movies_cast foreign keys
ALTER TABLE public.movies_cast ADD CONSTRAINT movies_cast_cast_id_fkey FOREIGN KEY (cast_id) REFERENCES public.casts (id) ON DELETE CASCADE;

ALTER TABLE public.movies_cast ADD CONSTRAINT movies_cast_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies (id) ON DELETE CASCADE;