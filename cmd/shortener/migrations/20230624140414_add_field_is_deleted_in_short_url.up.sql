alter table public.short_url add if not exists is_deleted bool default false not null;
