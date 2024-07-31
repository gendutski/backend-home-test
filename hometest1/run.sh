#!/bin/bash

# http port
read -p "Enter http port to use: " HTTP_PORT
if [ "$HTTP_PORT" == "" ]; then
    echo "Error: Invalid input"
    exit 
fi

# get mysql credential
read -p "Enter your MySQL host (don't use localhost): " MYSQL_HOST
read -p "Enter your MySQL username: " MYSQL_USERNAME
read -sp "Enter the MySQL password: " MYSQL_PASSWORD
echo
# Trying to connect to MySQL
mysql -h"$MYSQL_HOST" -u"$MYSQL_USERNAME" -p"$MYSQL_PASSWORD" -e "quit" 2>/dev/null
if [ $? -eq 0 ]; then
    echo "Connection to MySQL successful."
else
    echo "Error: Connection to MySQL failed. Check your username and password."
    exit
fi

# Validate database name
read -p "Enter the name of the database used: " MYSQL_DB_NAME

# Tried to checking for database existence
DB_EXISTS=$(mysql -u"$MYSQL_USERNAME" -p"$MYSQL_PASSWORD" -e "SHOW DATABASES LIKE '$MYSQL_DB_NAME';" 2>/dev/null | grep "^$MYSQL_DB_NAME$")

if [ "$DB_EXISTS" == "" ]; then
    # if exists ok
    echo "Error: Database '$MYSQL_DB_NAME' doesn't exists."
    exit
fi

# create docker image
read -p "Enter your docker image name: " DOCKER_NAME
if [ "$DOCKER_NAME" == "" ]; then 
    echo "Error: Invalid input"
    exit
fi
docker build -t "$DOCKER_NAME" .

# get container name
read -p "Enter your docker container name: " CONTAINER_NAME
if [ "$CONTAINER_NAME" == "" ]; then 
    echo "Error: Invalid input"
    exit
fi
# delete existing container 
echo "Delete existing container"
docker container rm "$CONTAINER_NAME"
# create container
echo "Creating container"
docker container create --name "$CONTAINER_NAME" -e HTTP_PORT=$HTTP_PORT -e MYSQL_HOST="$MYSQL_HOST" -e MYSQL_USERNAME="$MYSQL_USERNAME" -e MYSQL_DB_NAME="$MYSQL_DB_NAME" -e MYSQL_PASSWORD="$MYSQL_PASSWORD" -p $HTTP_PORT:$HTTP_PORT $DOCKER_NAME

# start container
echo "Start container"
docker container start "$CONTAINER_NAME"

echo "Finish!"


