create table if not exists subscription
(
    id           serial primary key,
    service_name varchar not null,
    price        int     not null,
    user_id      varchar not null,
    start_date   date    not null,
    end_date     date
);

create index if not exists idx_user_interval_price on subscription(user_id, start_date, end_date);
