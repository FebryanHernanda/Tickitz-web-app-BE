-- public.schedules definition
-- Drop table
-- DROP TABLE public.schedules;
CREATE TYPE public.show_time AS ENUM ('10:00', '13:00', '16:00', '19:00');

CREATE TABLE
    public.schedules (
        id serial4 NOT NULL,
        "date" date NOT NULL,
        "time" public."show_time" NOT NULL,
        movie_id int4 NULL,
        CONSTRAINT schedules_pkey PRIMARY KEY (id),
        CONSTRAINT schedules_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies (id) ON DELETE CASCADE
    );