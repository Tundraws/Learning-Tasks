CREATE TABLE IF NOT EXISTS page_visits (
    id SERIAL PRIMARY KEY,
    visit_count INTEGER NOT NULL DEFAULT 0
);

-- Инициализируем счетчик
INSERT INTO page_visits (visit_count) VALUES (0) ON CONFLICT DO NOTHING;