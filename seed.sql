DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM events WHERE id = 1) THEN
        INSERT INTO events (id, name, date, location, capacity, created_at, updated_at)
        VALUES (1, 'Event 1', '2024-11-16 23:43:06.667', 'Ho Chi Minh', 1000, NOW(), NOW());
END IF;
END $$;
