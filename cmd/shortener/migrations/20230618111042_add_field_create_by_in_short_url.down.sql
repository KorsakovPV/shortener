DROP INDEX IF EXISTS short_url_created_by_index;

alter table public.short_url DROP IF EXISTS created_by;