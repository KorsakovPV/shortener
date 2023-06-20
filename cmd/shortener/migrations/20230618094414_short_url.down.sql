DROP EXTENSION IF EXISTS "uuid-ossp";

DROP TABLE IF EXISTS shortener.public.short_url;

DROP INDEX IF EXISTS shortener.public.short_url_original_url_key;

DROP INDEX IF EXISTS shortener.public.short_url_pkey;