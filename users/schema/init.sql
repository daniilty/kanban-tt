create table users(
  id serial primary key,
  name varchar not null,
  email varchar not null,
  password_hash varchar not null,
  email_confirmed boolean not null,
  created_at date not null
)
