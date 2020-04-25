CREATE DATABASE grpc_2;

-- Extension: "uuid-ossp"
CREATE EXTENSION "uuid-ossp";

-- FUNCTION: public.trigger_update_timestamp()
CREATE FUNCTION public.trigger_update_timestamp()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $BODY$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;

$BODY$;

ALTER FUNCTION public.trigger_update_timestamp()
    OWNER TO postgres;

-- Table: public."user"
CREATE TABLE public."user"
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    username character varying COLLATE pg_catalog."default" NOT NULL,
    state smallint NOT NULL DEFAULT 1,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    CONSTRAINT user_pkey PRIMARY KEY (id)
);

-- Index: idx_user_username
CREATE UNIQUE INDEX idx_user_username
    ON public."user" USING btree
    (username COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

-- Trigger: update_user_updated_at
CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE 
    ON public."user"
    FOR EACH ROW
    EXECUTE PROCEDURE public.trigger_update_timestamp();