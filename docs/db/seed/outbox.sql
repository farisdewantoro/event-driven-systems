DO $$ 
DECLARE 
    i INT;
    random_status TEXT;
    random_event TEXT;
    random_execute_at TIMESTAMP;
BEGIN
    FOR i IN 1..100000 LOOP
        -- Randomly select status
        SELECT CASE FLOOR(RANDOM() * 3)::INT
               WHEN 0 THEN 'PENDING'
               WHEN 1 THEN 'RETRYING'
               ELSE 'FAILED' 
               END INTO random_status;

        -- Random event_type from 1 to 100
        SELECT 'example_event_' || (FLOOR(RANDOM() * 100) + 1)::TEXT INTO random_event;

        -- Random execution timestamp between 2024-01-01 and 2026-12-31
        SELECT TIMESTAMP '2024-01-01 00:00:00' + (RANDOM() * INTERVAL '3 years') INTO random_execute_at;

        -- Insert into outbox table
        INSERT INTO outbox (event_type, status, attempt, destination_type, execute_at, payload, error_message)
        VALUES (random_event, random_status, 0, 'kafka', random_execute_at, '{}', NULL);
    END LOOP;
END $$;


-- CREATE PENDING MESSAGES
DO $$ 
DECLARE 
    i INT;
    random_status TEXT;
    random_event TEXT;
    random_execute_at TIMESTAMP;
BEGIN
    FOR i IN 1..3000 LOOP

        -- Random event_type from 1 to 100
        SELECT 'example_event_' || (FLOOR(RANDOM() * 100) + 1)::TEXT INTO random_event;

        -- Insert into outbox table
        INSERT INTO outbox (event_type, status, attempt, destination_type, execute_at, payload, error_message)
        VALUES (random_event, 'PENDING', 0, 'kafka', now(), '{}', NULL);
    END LOOP;
END $$;
