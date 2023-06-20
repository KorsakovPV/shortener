alter table public.short_url add IF NOT EXISTS created_by uuid not null;

create index short_url_created_by_index on public.short_url (created_by);