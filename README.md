# 🛠️ Koyjak Postgres Setup Guide

This guide walks you through setting up PostgreSQL for the **Koyjak** project on a Linux server.

---

## 📦 Step 1: Install PostgreSQL

```bash

sudo apt update

sudo apt install postgresql

## Enable and start the PostgreSQL service:

sudo systemctl enable postgresql

sudo systemctl start postgresql

```

## 🔐 Step 3: Login and Configure PostgreSQL Authentication
Login to postgresql and create db

```bash
sudo -u postgres psql
## if sudo -u postgres psql not working or ask for password and still wrong follow next step
```

Edit the pg_hba.conf file to allow password-based local access:

```bash
sudo nano /etc/postgresql/*/main/pg_hba.conf
```

Update or add this :
```bash
# Allow local connections with MD5 password authentication
local   all             all                                     md5
host    all             all             127.0.0.1/32            md5
host    all             all             ::1/128                 md5

# Replication settings (optional)
local   replication     all                                     peer
host    replication     all             127.0.0.1/32            scram-sha-256
host    replication     all             ::1/128                 scram-sha-256
```
restart postgresql
```bash
sudo systemctl restart postgresql
```

create database
```bash
## write 
CREATE DATABASE koyjak;
```

## 🔄 Step 4: Load Your Schema

Use pg_restore to import your SQL dump:
```bash
cd migrations

pg_restore --clean -U postgres -d koyjak ./db.sql

sudo -u postgres psql

\c koyjak
```

## 🧠 Step 5: Tune Linux Kernel for Max Connections

To increase system semaphore limits for PostgreSQL:
```bash

sudo nano /etc/sysctl.conf

and add this: 

kernel.sem = 250 32000 100 128

Then apply the changes:

sudo sysctl -p
```

## You’re now ready to run Koyjak with PostgreSQL
