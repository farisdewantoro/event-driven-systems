
-- CREATE PENDING MESSAGES
DO $$ 
DECLARE 
    i INT;
    random_status TEXT;
    random_event TEXT;
    random_execute_at TIMESTAMP;
BEGIN
    FOR i IN 1..3000 LOOP


        -- Insert into outbox table
        INSERT INTO outbox (event_type, status, attempt, destination_type, execute_at, payload, error_message)
        VALUES ('email:send_notification', 'PENDING', 0, 'ASYNQ', now(), '{
  "user_id": "c0eee119-3f01-4ac9-a2a7-a1576f6185a5",
  "notification_id": "1f99669d-85fc-4a94-be90-46c59a026477",
  "notification_type": "USER_REGISTRATION"
}', NULL);
    END LOOP;
END $$;
