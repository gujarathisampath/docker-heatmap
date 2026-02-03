-- Safe cleanup script
-- This will not crash if tables are missing

DO $$ 
BEGIN
    -- Disable triggers
    SET session_replication_role = 'replica';

    -- Clean Activity Events
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'activity_events') THEN
        TRUNCATE TABLE activity_events CASCADE;
        IF EXISTS (SELECT FROM pg_class WHERE relname = 'activity_events_id_seq') THEN
            ALTER SEQUENCE activity_events_id_seq RESTART WITH 1;
        END IF;
    END IF;

    -- Clean Docker Accounts
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'docker_accounts') THEN
        TRUNCATE TABLE docker_accounts CASCADE;
        IF EXISTS (SELECT FROM pg_class WHERE relname = 'docker_accounts_id_seq') THEN
            ALTER SEQUENCE docker_accounts_id_seq RESTART WITH 1;
        END IF;
    END IF;

    -- Clean Users
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users') THEN
        TRUNCATE TABLE users CASCADE;
        IF EXISTS (SELECT FROM pg_class WHERE relname = 'users_id_seq') THEN
            ALTER SEQUENCE users_id_seq RESTART WITH 1;
        END IF;
    END IF;

    -- Re-enable triggers
    SET session_replication_role = 'origin';
END $$;

SELECT 'Reset complete. Missing tables were skipped.' as status;
