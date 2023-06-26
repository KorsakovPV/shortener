DROP TABLE IF EXISTS public.short_url;

DROP INDEX IF EXISTS short_url_original_url_key;

DROP INDEX IF EXISTS short_url_pkey;

DROP EXTENSION IF EXISTS "uuid-ossp";