alter table public.short_url DROP IF EXISTS created_by;

DROP INDEX IF EXISTS public.short_url_created_by_index;

