create table statuses(
  id serial primary key,
  name text not null,
  priority integer not null,
  owner_id varchar(255) not null
);

create table tasks(
  id serial primary key,
  content text not null,
  priority integer not null,
  owner_id varchar(255) not null,
  created_at date not null,
  status_id integer references statuses(id) on delete cascade
);
