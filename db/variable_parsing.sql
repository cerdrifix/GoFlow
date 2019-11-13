DO $$
DECLARE
	map_id uuid := '231330bc-0668-11ea-9071-ffe17b0f393e';
	start_node varchar(32) := 'start_1';
	creator_id varchar(64) := 'cerdrifix';
	variables json := '{"NOMINATIVO":"Davide","cognome":"Ceretto","dataCreazione":"20191027T22:02:54.254","nome":"Davide","testoRichiesta":"Necessario nuovo PC"}';
   _instance_id uuid;
   _key   text;
   _value text;
BEGIN

	select public.fn_instance_new(map_id,start_node,creator_id,variables) INTO _instance_id;
	
	RAISE NOTICE 'Instance id: %', _instance_id;

    FOR _key, _value IN
       SELECT * FROM json_each_text(variables)
    LOOP
       RAISE NOTICE '%: %', _key, _value;
	   
-- 	   INSERT INTO public.variables ( state_id, )
    END LOOP;

END $$;

-- SELECT public.fn_instance_new('231330bc-0668-11ea-9071-ffe17b0f393e', 'start_1', 'cerdrifix', '{"cognome":"Ceretto","dataCreazione":"20191027T22:02:54.254","nome":"Davide","testoRichiesta":"Necessario nuovo PC"}')


-- select * from public.maps;