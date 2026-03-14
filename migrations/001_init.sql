CREATE TABLE IF NOT EXISTS links (
    id           SERIAL PRIMARY KEY,
    short_code   VARCHAR(20) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at   TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_links_short_code ON links(short_code);

CREATE TABLE IF NOT EXISTS clicks (
    id         SERIAL PRIMARY KEY,
    link_id    INT REFERENCES links(id) ON DELETE CASCADE,
    user_agent TEXT,
    ip_address VARCHAR(45),
    clicked_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_clicks_link_id ON clicks(link_id);
CREATE INDEX IF NOT EXISTS idx_clicks_clicked_at ON clicks(clicked_at);