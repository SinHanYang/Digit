CREATE OR REPLACE FUNCTION hook_update() RETURNS trigger AS $hook_update$
DECLARE
    x smallint;
BEGIN
    SELECT d INTO x FROM flag;
    IF TG_OP = 'DELETE' THEN
        RAISE EXCEPTION triggered_action_exception;
    ELSIF x=1 THEN
        RETURN NEW;
    ELSE
        NEW.digitsnid:= nextval(pg_get_serial_sequence(TG_TABLE_NAME,'digitsnid'));
        EXECUTE format(
            'INSERT INTO %I select $1.*'
            ,TG_TABLE_NAME) USING NEW;
        RAISE NOTICE 'here';
        OLD.digitstatus=0;
        RETURN OLD;
    END IF;
EXCEPTION
    WHEN triggered_action_exception THEN
        EXECUTE format (
            'UPDATE flag SET d = 1;'
        );
        EXECUTE format (
            'UPDATE %I SET digitstatus = 0 
            WHERE digitsnid=%L;'
        ,TG_TABLE_NAME,OLD.digitsnid); 
        EXECUTE format(
            'UPDATE flag SET d = 0;'
        );
        RETURN NULL;
END;
$hook_update$ LANGUAGE plpgsql;

CREATE TABLE flag (d int);
--CREATE TRIGGER hook_update BEFORE UPDATE OR DELETE on t__backend FOR EACH ROW EXECUTE FUNCTION hook_update();

