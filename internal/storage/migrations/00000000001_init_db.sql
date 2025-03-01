create table public.user
(
    id         uuid        not null constraint user_pk primary key,
    login      varchar     not null unique,
    password   varchar     not null,
    created_at timestamptz not null
);

create type public.secret_kind as enum('credentials', 'note', 'blob', 'bank_card');

create table public.secret
(
    id           uuid         not null primary key,
    user_id      uuid         not null references public.user (id) on delete cascade,
    name         varchar      not null,
    kind         secret_kind  not null,
    is_encrypted bool         not null,
    unique (user_id, name)
);

create table public.secret_credentials
(
    id       uuid    not null primary key references secret (id) on delete cascade,
    login    varchar not null,
    password varchar not null
);

create table public.secret_note
(
    id   uuid    not null primary key references secret (id) on delete cascade,
    body varchar not null
);

create table public.secret_blob
(
    id   uuid    not null primary key references secret (id) on delete cascade,
    body varchar not null
);

create table public.secret_bank_card
(
    id     uuid    not null primary key references secret (id) on delete cascade,
    name   varchar not null,
    number varchar not null,
    date   varchar not null,
    cvv    varchar not null
);

create table public.tag
(
    secret_id uuid    not null references secret (id) on delete cascade,
    text      varchar not null,
    primary key (secret_id, text)
);

---- create above / drop below ----

drop table public.user;
drop type  public.secret_kind;
drop table public.secret;
drop table public.secret_credentials;
drop table public.secret_note;
drop table public.secret_blob;
drop table public.secret_bank_card;
drop table public.tag;
