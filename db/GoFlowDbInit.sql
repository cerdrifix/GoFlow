-- Table: public.maps

DROP TABLE IF EXISTS public.instances;
DROP TABLE IF EXISTS public.states;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.maps;


CREATE TABLE IF NOT EXISTS public.maps
(
	id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
	name varchar(255) COLLATE pg_catalog."default" NOT NULL,
	version integer NOT NULL,
	creation_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
	data json NOT NULL,
	CONSTRAINT maps_pkey PRIMARY KEY (id),
	CONSTRAINT maps_unique_name_version UNIQUE (name, version)

)

TABLESPACE pg_default;

ALTER TABLE public.maps
	OWNER to cerdrifix;

-- Table: users
CREATE TABLE public.users
(
    username character varying(64) NOT NULL,
    name character varying(255) NOT NULL,
    surname character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    PRIMARY KEY (username)
);

ALTER TABLE public.users
    OWNER to postgres;

-- Table: public.states

CREATE TABLE public.states
(
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    map_id uuid NOT NULL,
    node_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    creator_id character varying(64) COLLATE pg_catalog."default" NOT NULL,
    owner_id character varying(64) COLLATE pg_catalog."default",
    enter_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    leave_date timestamp without time zone,
    CONSTRAINT states_pk PRIMARY KEY (id),
    CONSTRAINT states_fk_map_id FOREIGN KEY (map_id)
        REFERENCES public.maps (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT states_fk_users_creator FOREIGN KEY (creator_id)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT states_fk_users_owner FOREIGN KEY (owner_id)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE public.states
    OWNER to postgres;

-- Table: instances

CREATE TABLE public.instances
(
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    map_id uuid NOT NULL,
    start_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_date timestamp without time zone,
    current_node uuid NOT NULL,
    CONSTRAINT instances_pk PRIMARY KEY (id),
    CONSTRAINT instances_fk_states FOREIGN KEY (current_node)
        REFERENCES public.states (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT instances_fk_maps FOREIGN KEY (map_id)
        REFERENCES public.maps (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

ALTER TABLE public.instances
    OWNER to postgres;

------ Stored Procedures

-- Procedure sp_map_insert
CREATE OR REPLACE PROCEDURE sp_map_insert (
	_data json
)
AS
$$
DECLARE 
	_name		VARCHAR(255);
	_version 	INTEGER;
BEGIN

	_name := _data->'name';
	_name := trim(both '"' from _name);
	
	if _name is null then
		raise exception 'Errore! JSON non contenente il nome del workflow';
	end if;
	
	SELECT 	COUNT(name) + 1 INTO _version
	FROM	public.maps
	WHERE	name = _name;

	INSERT INTO public.maps ( name, version, data)
	VALUES ( _name, _version, _data );
	
	raise notice 'Inserita mappa % - versione %', _name, _version;

END 
$$ LANGUAGE plpgsql;

-- Procedure sp_maps_getlatestbyname
CREATE OR REPLACE FUNCTION fn_maps_getlatestbyname (
	_name varchar(255)
)
RETURNS TABLE (
	name 	varchar(255),
	version	int,
	data	json
)
AS
$$
DECLARE 
BEGIN
	
	RETURN QUERY
	SELECT 		M.id, M.name, M.version, M.data
	FROM		public.maps M
	WHERE		M.name = _name
	ORDER BY	version desc
	LIMIT 1;

END
$$ LANGUAGE plpgsql;

-- Procedure sp_maps_getbynameandversion
CREATE OR REPLACE FUNCTION fn_maps_getbynameandversion (
	_name 		varchar(255),
	_version 	int
)
RETURNS TABLE (
	name 	varchar(255),
	version	int,
	data	json
)
AS
$$
DECLARE
BEGIN

	RETURN QUERY
	SELECT 		M.id, M.name, M.version, M.data
	FROM		public.maps M
	WHERE		name = _name
	and			version = _version;

END 
$$ LANGUAGE plpgsql;


DO $$
DECLARE
	data json := '{"name":"richiesta_con_approvazione","description":"Richiesta con approvazione","nodes":[{"name":"start_1","type":"start","events":{"pre":[{"type":"validator","name":"CheckInputVariable","parameters":[{"name":"inputVariableName","type":"variable","value":"nome"}]}],"post":[{"type":"function","name":"CopyVariable","parameters":[{"name":"srcVariable","type":"variable","value":"nome"},{"name":"dstVariable","type":"variable","value":"NOMINATIVO"}]}]},"triggers":[{"name":"auto","after":{"unit":"seconds","value":0},"transaction":"start_to_task_approvativo"}],"transactions":[{"name":"start_to_task_approvativo","description":"Eseguito","to":"task_approvativo","events":{"pre":[],"post":[]}}]},{"name":"task_approvativo","type":"task","events":{"pre":[],"post":[]},"triggers":[{"name":"auto_approve","after":{"unit":"days","value":10},"transaction":"task_approvativo_cancel"}],"transactions":[{"name":"task_approvativo_ok","description":"Approva","visible":true,"to":"end_ok","events":{"pre":[],"post":[]}},{"name":"task_approvativo_ko","description":"Rifiuta","visible":true,"to":"end_ko","events":{"pre":[],"post":[]}},{"name":"task_approvativo_cancel","description":"Annulla","visible":false,"to":"end_canceled","events":{"pre":[],"post":[]}}]},{"name":"end_ok","type":"end","description":"Richiesta terminata con successo"},{"name":"end_ko","type":"end","description":"Richiesta rifiutata"},{"name":"end_canceled","type":"end","description":"Richiesta annullata da sistema (tempo massimo raggiunto)"}]}';
BEGIN

	call sp_map_insert(data);

END $$;
