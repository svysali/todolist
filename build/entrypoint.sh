#!/bin/sh

# Wait for database to be ready
until nc -w 1 -z $DB_HOST $DB_PORT; do
    echo "Database is unavailable -- sleeping..."
    sleep 2
done

/app/todolist
