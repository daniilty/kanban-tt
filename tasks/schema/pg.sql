create table if not exists statuses(
  id serial primary key,
  name text not null,
  owner_id varchar(255) not null,
  parent_id integer not null,
  child_id integer not null
);

create table if not exists tasks(
  id serial primary key,
  content text not null,
  priority integer not null,
  owner_id varchar(255) not null,
  created_at date not null,
  status_id integer references statuses(id) on delete cascade
);

create index statuses_owner_idx on statuses(owner_id);
