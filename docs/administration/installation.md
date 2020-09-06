# Installation

At the end of this section, Cell should be installed. It will still require configuration.

## Installing Snap

[Snap](https://snapcraft.io) is a great package manager for Linux; it'll make this process a whole lot easier. You might already have it, but you're able to use `apt`:

```console
$ sudo apt update
$ sudo apt install snapd
```

## Installing Go

[Go](https://golang.org) is the programming language used by Cell. It's available as a Snap:

```console
$ sudo snap install go --classic
```

## Cloning Cell

Since there aren't any pre-built executables and migrations are stored in the repo, you need to clone the source. This can be done with [git](https://git-scm.com):

```console
$ git clone https://github.com/open-slicer/cell
$ cd cell
```

Don't worry about installing git; you likely already have it.

## Setting up PostgreSQL

Cell relies on [PostgreSQL](https://www.postgresql.org) to keep state for all userland structures.

### Installation

The Snap for PostgreSQL is outdated.

```console
$ sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
$ sudo apt update
$ curl -L https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
$ sudo apt install postgresql
```

### Creating a database

First, log in as `postgres` and open a PostgreSQL prompt:

```console
$ sudo su postgres
$ psql
```

Now this is open, you can create a database called `cell`:

```sql
CREATE DATABASE cell;
```

To exit, type `exit`.

### Migrating the database

After creating a DB, you must migrate it. This can be done with [`golang-migrate`](https://github.com/golang-migrate/migrate).

#### Installing `golang-migrate`

This, again, isn't available through Snap:

```console
$ echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
$ sudo apt update
$ curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | sudo apt-key add -
$ sudo apt install migrate
```

#### Running migrations

The `up` command can be used to migrate the DB:

```console
$ migrate -database 'postgres://username:password@localhost/cell' -path migrations up
```

Make sure to replace the placeholders here!

## Installing Redis

[Redis](https://redis.io) is used to keep track of Lockets - WebSocket nodes - and publish messages to clients. It is available as a Snap:

```console
$ sudo snap install redis --classic
```
