CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS shortener.public.short_url (id uuid PRIMARY KEY, original_url TEXT NOT NULL UNIQUE);