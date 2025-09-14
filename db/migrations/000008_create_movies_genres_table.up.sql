-- public.movies_genres definition
-- Drop table
-- DROP TABLE public.movies_genres;
CREATE TABLE
    public.movies_genres (
        movie_id int4 NOT NULL,
        genre_id int4 NOT NULL,
        CONSTRAINT movies_genres_pkey PRIMARY KEY (movie_id, genre_id),
        CONSTRAINT movies_genres_genre_id_fkey FOREIGN KEY (genre_id) REFERENCES public.genres (id) ON DELETE CASCADE,
        CONSTRAINT movies_genres_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies (id) ON DELETE CASCADE
    );