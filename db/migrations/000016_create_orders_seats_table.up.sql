-- public.orders_seats definition

-- Drop table

-- DROP TABLE public.orders_seats;

CREATE TABLE public.orders_seats (
	id serial4 NOT NULL,
	status varchar(20) NULL,
	order_id int4 NULL,
	seat_id int4 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT orders_seats_pkey PRIMARY KEY (id),
	CONSTRAINT orders_seats_status_check CHECK (((status)::text = ANY ((ARRAY['booked'::character varying, 'available'::character varying])::text[])))
);


-- public.orders_seats foreign keys

ALTER TABLE public.orders_seats ADD CONSTRAINT orders_seats_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;
ALTER TABLE public.orders_seats ADD CONSTRAINT orders_seats_seat_id_fkey FOREIGN KEY (seat_id) REFERENCES public.seats(id) ON DELETE CASCADE;