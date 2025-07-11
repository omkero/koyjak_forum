
-- koyjak db logs

alter table threads
add  column forum_title varchar(255) not null;


delete from threads;

select * from threads;

SELECT COUNT(*) AS total FROM forums WHERE forum_title = 'Frontend';

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

create table forums(
   forum_id serial primary key,
   forum_title varchar(255) not null unique,
   threads_count bigint default 0,
   posts_count bigint default 0,
   created_at timestamptz default now()
);

INSERT INTO forums (forum_title, threads_count) 
VALUES('Frontend',1) 
ON CONFLICT(forum_title) DO UPDATE SET threads_count = forums.threads_count + 1;

INSERT INTO forums (forum_title, posts_count) 
VALUES('Frontend',1) 
ON CONFLICT(forum_title) DO UPDATE SET posts_count = forums.posts_count + 1;

