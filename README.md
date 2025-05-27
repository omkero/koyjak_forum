## use this design 
https://www.webdesignerforum.co.uk/


## setup database 
sudo apt install postgresql

sudo systemctl start postgresql
sudo systemctl enable postgresql

-- now you have to login postgres cli and create database

sudo psql -U postgres

sql: CREATE DATABASE koyjak;

-- exit and load the schema
pg_restore --clean -U postgres -d koyjak ./db.sql

sudo psql -U postgres

\c koyjak


## and now you configure postgresql to accept connections use like this


sudo nano /etc/postgresql/*/main/pg_hba.conf

local   all             postgres                                md5

# TYPE  DATABASE        USER            ADDRESS                 METHOD

# "local" is for Unix domain socket connections only
local   all             all                                     md5
# IPv4 local connections:
host    all             all             127.0.0.1/32            md5
# IPv6 local connections:
host    all             all             ::1/128                 md5
# Allow replication connections from localhost, by a user with the
# replication privilege.
local   replication     all                                     peer
host    replication     all             127.0.0.1/32            scram-sha-256
host    replication     all             ::1/128                 scram-sha-256

