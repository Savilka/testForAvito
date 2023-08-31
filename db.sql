create table "user"
(
    id serial not null
        constraint user_pk
            primary key
);

create table segment
(
    id   serial not null
        constraint segment_pk
            primary key,
    name varchar(200) unique
);

create table user_segment
(
    user_id      serial       not null
        constraint user_segment_user_id_fk
            references "user"
            on delete cascade,
    segment_name varchar(200) not null
        constraint user_segment_segment_name_fk
            references segment (name)
            on delete cascade,
    ttl          bigint,
    date_insert  bigint,
    constraint user_segment_pk
        primary key (user_id, segment_name)
);

