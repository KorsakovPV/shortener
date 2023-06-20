alter table shortener.public.short_url DROP IF EXISTS created_by;

DROP INDEX IF EXISTS shortener.public.short_url_created_by_index;

