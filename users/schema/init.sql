create table if not exists users(
  id serial primary key,
  name varchar not null,
  email varchar not null,
  password_hash varchar not null,
  email_confirmed boolean not null,
  created_at date not null
);

alter table users add column task_ttl integer;

update users set task_ttl=365 where task_ttl is null;
