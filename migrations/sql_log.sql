
-- koyjak db logs



---------------------------------

create table Threads (
  thread_id SERIAL primary key,
  user_id bigint not null,
  thread_title varchar(255) not null,
  thread_content text not null,
  created_at timestamptz default now()
);

insert into Threads (user_id, thread_title, thread_content)
values (500, 'Hello koyjacks', 'this is the first koyjak');

select * from Threads;


----------------------------------

create table Users (
  user_id SERIAL primary key,
  username varchar(20) unique not null,
  email_address varchar(255) unique not null,
  pwd text not null,
  created_at timestamptz default now()
);

insert into Users (username, email_address, pwd)
values ('koyjak', 'koyjak@soyo.org', 'soyjackos');

select * from Users;

------------------------------------------------------------

CREATE table forums (
  forum_id SERIAL primary key,
  forum_category varchar(255) not null,
  forum_title varchar(255) not null,
  forum_description text not null,
  threads_count bigint not null,
  posts_count bigint not null,
  created_at timestamptz default now()
);

select * from forums;