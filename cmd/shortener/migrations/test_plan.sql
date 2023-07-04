-- insert rows

do
$$
    begin
        for r in 1..100000
            loop
                --                 insert into schema_name.table_name(id) values(r);
                INSERT INTO short_url (id, original_url, created_by) VALUES (uuid_generate_v4(), uuid_generate_v4(), uuid_generate_v4());
            end loop;
    end;
$$;

DELETE FROM short_url;

-- index b-tree

create index short_url_pkey on public.short_url (id);

DROP INDEX IF EXISTS public.short_url_pkey;

-- index hash

create index short_url_pkey_hash on public.short_url using hash(id);

DROP INDEX IF EXISTS public.short_url_pkey_hash;

-- index b-tree

create index short_url_original_url_key on public.short_url (original_url);

DROP INDEX IF EXISTS public.short_url_original_url_key;

EXPLAIN ANALYZE SELECT * FROM short_url where id='1fe70fcd-c8e5-466c-a7f2-bd371755dab5';
EXPLAIN SELECT * FROM short_url where id='1fe70fcd-c8e5-466c-a7f2-bd371755dab5';

\di+

select pg_size_pretty(pg_total_relation_size('short_url'));

SELECT * FROM pg_stat_all_indexes where schemaname <> 'pg_toast' and  schemaname <> 'pg_catalog';

create index short_url_created_by_index_hash on public.short_url using hash (created_by);
create index short_url_created_by_index_hash on public.short_url (created_by);
DROP INDEX IF EXISTS public.short_url_created_by_index_hash;