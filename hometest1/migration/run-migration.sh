#!/bin/bash

# get mysql credential
read -p "Enter your MySQL host (eg: localhost): " MYSQL_HOST
read -p "Enter your MySQL username: " MYSQL_USERNAME
read -sp "Enter the MySQL password: " MYSQL_PASSWORD
echo

# Trying to connect to MySQL
mysql -h"$MYSQL_HOST" -u"$MYSQL_USERNAME" -p"$MYSQL_PASSWORD" -e "quit" 2>/dev/null
if [ $? -eq 0 ]; then
    echo "Connection to MySQL successful."
else
    echo "Connection to MySQL failed. Check your username and password."
    exit
fi

# Validate database name
read -p "Enter the name of the database used: " MYSQL_DB_NAME

# Tried to checking for database existence
DB_EXISTS=$(mysql -u"$MYSQL_USERNAME" -p"$MYSQL_PASSWORD" -e "SHOW DATABASES LIKE '$MYSQL_DB_NAME';" 2>/dev/null | grep "^$MYSQL_DB_NAME$")

if [ "$DB_EXISTS" == "$MYSQL_DB_NAME" ]; then
    # if exists ok
    echo "Database '$MYSQL_DB_NAME' exists."
else
    # if not exists, attempt to crete one
    echo "Database '$MYSQL_DB_NAME' not exists."
    read -p "Create database '$MYSQL_DB_NAME'? (yes/no): " CREATE_DB
    if [ "$CREATE_DB" != "yes" ]; then
        echo "Database must exists"
        exit
    else
        mysql -u"$MYSQL_USERNAME" -p"$MYSQL_PASSWORD" -e "CREATE DATABASE $MYSQL_DB_NAME;" 2>/dev/null
        if [ $? -eq 0 ]; then
            echo "Database created."
        else
            echo "Failed to create database."
            exit
        fi
    fi
fi

# create table if not exists
TABLES=("product" "product_quantity" "promotion")
COUNT=1

for TABLE_NAME in "${TABLES[@]}"; do
    # check table if exists
    TABLE_EXISTS=$(mysql -u"$MYSQL_USERNAME" -p"$MYSQL_PASSWORD" -D "$MYSQL_DB_NAME" -e "SHOW TABLES LIKE '$TABLE_NAME';" 2>/dev/null | grep "^$TABLE_NAME$")
    if [ "$TABLE_EXISTS" == "$TABLE_NAME" ]; then
        echo "Tabel '$TABLE_NAME' is exists in '$MYSQL_DB_NAME'."
    else
        echo "Tabel '$TABLE_NAME' is not exists in database '$MYSQL_DB_NAME'."
        echo "Creating..."

        mysql -u"$MYSQL_USERNAME" -p"$MYSQL_PASSWORD" $MYSQL_DB_NAME <./0$COUNT-$TABLE_NAME.sql
    fi

    COUNT=$((COUNT + 1))
done

# run seed data
echo
echo "Do you want to fill example data?"
echo "WARNING!!!"
echo "If you choose yes, it will truncate all tables"
read -p "(yes/no): " SEED_DB
if [ "$SEED_DB" == "yes" ]; then
    echo "Seeding tables"
    mysql -u"$MYSQL_USERNAME" -p"$MYSQL_PASSWORD" $MYSQL_DB_NAME <./04-seed-data.sql
fi

echo "Finish"