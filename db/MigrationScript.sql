-- Role: cerdrifix
-- DROP ROLE IF EXISTS cerdrifix;

-- CREATE ROLE cerdrifix WITH
--   LOGIN
--   PASSWORD 'cerdrifix1234'
--   NOSUPERUSER
--   INHERIT
--   NOCREATEDB
--   NOCREATEROLE
--   NOREPLICATION;

-- Extension: "uuid-ossp"

DROP EXTENSION IF EXISTS "uuid-ossp";

CREATE EXTENSION "uuid-ossp"
    SCHEMA public
    VERSION "1.1";
	
  
-- Table: public.access_logs

DROP TABLE IF EXISTS public.access_logs;

CREATE TABLE public.access_logs
(
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    access_date timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    page text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT access_logs_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.access_logs
    OWNER to postgres;

GRANT INSERT, SELECT, UPDATE ON TABLE public.access_logs TO cerdrifix;

GRANT ALL ON TABLE public.access_logs TO postgres;
  
  
  
  -- PROCEDURE: public.sp_access_logs_insert(text)

DROP PROCEDURE IF EXISTS public.sp_access_logs_insert(text);

CREATE OR REPLACE PROCEDURE public.sp_access_logs_insert(
	page text)
LANGUAGE 'plpgsql'

AS $BODY$
BEGIN
    INSERT INTO access_logs ( page )
	VALUES ( page );
	
    COMMIT;
END;
$BODY$;

  