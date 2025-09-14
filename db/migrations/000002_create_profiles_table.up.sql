-- public.profiles definition
-- Drop table
-- DROP TABLE public.profiles;
CREATE TABLE
    public.profiles (
        user_id int4 NULL,
        first_name varchar(100) NULL,
        last_name varchar(100) NULL,
        phone_number varchar(20) NULL,
        points int4 DEFAULT 0 NULL,
        image_path text NULL,
        CONSTRAINT profiles_user_id_key UNIQUE (user_id)
    );

-- public.profiles foreign keys
ALTER TABLE public.profiles ADD CONSTRAINT profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users (id) ON DELETE CASCADE;