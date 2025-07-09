# Create

## production

### admin

```sql
create user datapipe_production_admin with password '';
create database datapipe_production;
grant all privileges on database datapipe_production to datapipe_production_admin;
alter database datapipe_production owner to datapipe_production_admin;

-- connect to target db.
\c datapipe_production

create extension postgis;
create extension postgis_topology;

-- Create database objects for permissions handling during migrations.

create table zz_table_permissions (
    schema_name text,
    table_name  text,
    processed   boolean default false,
    primary key (schema_name, table_name)
);

grant all on zz_table_permissions to datapipe_production_admin;

create or replace function record_new_table()
returns event_trigger as $$
declare
    obj record;
begin
    for obj in
        select *
        from pg_event_trigger_ddl_commands()
        where command_tag = 'CREATE TABLE'
    loop
        insert into zz_table_permissions (schema_name, table_name)
        values (
            (
                select pg_namespace.nspname
                from pg_class
                join pg_namespace
                    on pg_namespace.oid = pg_class.relnamespace
                where pg_class.oid = obj.objid
            ),
            (
                select relname
                from pg_class
                where oid = obj.objid
            )
        )
        on conflict do nothing;
    end loop;
end;
$$ language plpgsql;

create event trigger record_new_table_trigger on ddl_command_end
when tag in ('CREATE TABLE')
execute procedure record_new_table();

create or replace function grant_select_on_tables(target_role text)
returns void as $$
declare
    tbl record;
begin
    for tbl in
        select schema_name, table_name
        from zz_table_permissions
        where not processed
    loop
        begin
            execute format(
                'grant select on table %s.%s to %s',
                tbl.schema_name, tbl.table_name, target_role
            );

            update zz_table_permissions
            set processed = true
            where
                schema_name = tbl.schema_name
                and table_name = tbl.table_name
            ;
        end;
    end loop;
end;
$$ language plpgsql;
```

### reader

```sql
create user datapipe_production_reader with password 'HzPLVDwt2ftMZymq9kbUJ498';
alter user datapipe_production_reader set default_transaction_read_only = on;

-- connect to target db.
\c datapipe_production

grant connect on database datapipe_production to datapipe_production_reader;
grant usage on schema public to datapipe_production_reader;
grant select on all tables in schema public to datapipe_production_reader;
```

## testing

### admin

```sql
create user datapipe_testing_admin with password '';
create database datapipe_testing;
grant all privileges on database datapipe_testing to datapipe_testing_admin;
alter database datapipe_testing owner to datapipe_testing_admin;

-- connect to target db.
\c datapipe_testing

create extension postgis;
create extension postgis_topology;

-- Create database objects for permissions handling during migrations.

create table zz_table_permissions (
    schema_name text,
    table_name  text,
    processed   boolean default false,
    primary key (schema_name, table_name)
);

grant select on zz_table_permissions to datapipe_testing_admin;

create or replace function record_new_table()
returns event_trigger as $$
declare
    obj record;
begin
    for obj in
        select *
        from pg_event_trigger_ddl_commands()
        where command_tag = 'CREATE TABLE'
    loop
        insert into zz_table_permissions (schema_name, table_name)
        values (
            (
                select pg_namespace.nspname
                from pg_class
                join pg_namespace
                    on pg_namespace.oid = pg_class.relnamespace
                where pg_class.oid = obj.objid
            ),
            (
                select relname
                from pg_class
                where oid = obj.objid
            )
        )
        on conflict do nothing;
    end loop;
end;
$$ language plpgsql;

create event trigger record_new_table_trigger on ddl_command_end
when tag in ('CREATE TABLE')
execute procedure record_new_table();

create or replace function grant_select_on_tables(target_role text)
returns void as $$
declare
    tbl record;
begin
    for tbl in
        select schema_name, table_name
        from zz_table_permissions
        where not processed
    loop
        begin
            execute format(
                'grant select on table %s.%s to %s',
                tbl.schema_name, tbl.table_name, target_role
            );

            update zz_table_permissions
            set processed = true
            where
                schema_name = tbl.schema_name
                and table_name = tbl.table_name
            ;
        end;
    end loop;
end;
$$ language plpgsql;
```

### reader

```sql
create user datapipe_testing_reader with password '';
alter user datapipe_testing_reader set default_transaction_read_only = on;

-- connect to target db.
\c datapipe_testing

grant connect on database datapipe_testing to datapipe_testing_reader;
grant usage on schema public to datapipe_testing_reader;
grant select on all tables in schema public to datapipe_testing_reader;
```
