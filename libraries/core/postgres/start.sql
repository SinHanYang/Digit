CREATE OR REPLACE FUNCTION hook_create_table() RETURNS event_trigger AS $$
DECLARE
    r record;
BEGIN
    FOR r IN SELECT * FROM pg_event_trigger_ddl_commands() 
    LOOP
        IF r.command_tag LIKE 'CREATE TABLE' THEN
            RAISE NOTICE 'caught % %', r.command_tag, r.object_identity;
            EXECUTE 'ALTER TABLE '
                    || substring(r.object_identity from 8)
                    || ' RENAME TO '
                    || substring(r.object_identity from 8)||'__backend';
            EXECUTE 'ALTER TABLE '
                    || substring(r.object_identity from 8)||'__backend'
                    || ' ADD COLUMN digitSNID SERIAL ';
            EXECUTE 'ALTER TABLE '
                    || substring(r.object_identity from 8)||'__backend'
                    || ' ADD COLUMN digitStatus INTEGER NOT NULL DEFAULT 1';
            EXECUTE 'CREATE OR REPLACE VIEW '
                    || substring(r.object_identity from 8)
                    || ' AS SELECT '
                    || ARRAY_TO_STRING(ARRAY(
                        SELECT column_name::text FROM information_schema.columns 
                        WHERE table_schema = 'public' 
                        AND table_name = substring(r.object_identity from 8)||'__backend' 
                        AND column_name::text NOT LIKE 'digit%'
                        ),',')
                    || ' FROM '
                    || substring(r.object_identity from 8)||'__backend'
                    || ' WHERE digitstatus=1';
            EXECUTE 'CREATE TRIGGER '
                    || 'hook_update_'||substring(r.object_identity from 8)
                    || ' BEFORE UPDATE OR DELETE ON '
                    || substring(r.object_identity from 8)||'__backend'
                    || ' FOR EACH ROW EXECUTE FUNCTION hook_update()';
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

create event trigger hook_create on ddl_command_end 
when tag in ('create table') execute function hook_create_table();