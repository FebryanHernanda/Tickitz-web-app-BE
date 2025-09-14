-- public.movies definition
-- Drop table
-- DROP TABLE public.movies;
CREATE TABLE
    public.movies (
        id serial4 NOT NULL,
        title varchar(255) NOT NULL,
        poster_path text NULL,
        backdrop_path text NULL,
        synopsis text NULL,
        release_date date NULL,
        rating numeric(3, 1) NULL,
        age_rating varchar(50) NULL,
        duration int4 NULL,
        created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
        updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
        director_id int4 NULL,
        CONSTRAINT movies_pkey PRIMARY KEY (id)
    );

-- public.movies foreign keys
ALTER TABLE public.movies ADD CONSTRAINT movies_director_id_fkey FOREIGN KEY (director_id) REFERENCES public.directors (id) ON DELETE SET NULL;