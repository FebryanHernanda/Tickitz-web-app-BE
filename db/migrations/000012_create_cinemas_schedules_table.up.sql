-- public.cinemas_schedules definition
-- Drop table
-- DROP TABLE public.cinemas_schedules;
CREATE TABLE
    public.cinemas_schedules (
        id serial4 NOT NULL,
        cinemas_id int4 NULL,
        schedules_id int4 NULL,
        locations_id int4 NULL,
        CONSTRAINT cinemas_schedules_pkey PRIMARY KEY (id),
        CONSTRAINT cinemas_schedules_cinema_id_fkey FOREIGN KEY (cinemas_id) REFERENCES public.cinemas (id) ON DELETE CASCADE,
        CONSTRAINT cinemas_schedules_location_id_fkey FOREIGN KEY (locations_id) REFERENCES public.locations (id) ON DELETE CASCADE
    );