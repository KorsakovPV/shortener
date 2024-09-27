alter table public.short_url add IF NOT EXISTS created_by uuid;

create index short_url_created_by_index_hash on public.short_url using hash (created_by);
