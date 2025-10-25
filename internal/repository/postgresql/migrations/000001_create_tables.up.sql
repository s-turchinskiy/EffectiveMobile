create table if not exists effectivemobile.users (
    id SERIAL primary key,
    uuid UUID not null unique
);

create table if not exists effectivemobile.services (
    id SERIAL primary key,
    name TEXT not null unique
);

create table if not exists effectivemobile.subscriptions (
    id SERIAL primary key,
    user_id INTEGER not null references effectivemobile.users (id),
    service_id INTEGER not null references effectivemobile.services (id),
    begin_date timestamp not null,
    end_date timestamp,
    sum numeric(15, 0) not null,
    UNIQUE (user_id, service_id, begin_date)
);