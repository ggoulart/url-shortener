CREATE TABLE urls
(
    encoded_key VARCHAR(255) PRIMARY KEY,
    long_url    TEXT UNIQUE NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Index for faster search by long_url
CREATE INDEX idx_urls_long_url ON urls (long_url);
