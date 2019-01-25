create or replace function fill_updated_at()
returns trigger as $$
begin
   new.updated_at = now();
   return new;
end;
$$ language 'plpgsql'
;
create sequence sq_entity_mock_id increment by 1 start with 1
;
create table entity_mock (
    id numeric(10) not null default nextval('sq_entity_mock_id'),
    name varchar(256) not null,
    age numeric(3) not null,
    data jsonb not null,
    deleted boolean default false,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    constraint pk_rarity primary key(id)
)
;
create trigger update_entity_mock_updated_at before update on entity_mock for each row execute procedure fill_updated_at()
;
create index ix_entity_mock_name on entity_mock (name asc)
;
