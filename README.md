# Event-driven

Project Setup Guide

This guide explains how to set up and run the system step by step. The system uses Kafka and multiple microservices: order, inventory, and payment.

1. Run Kafka from the root directory

In the project root directory, there is a docker-compose.yml file for Kafka.
Start Kafka with the following command:

docker compose up -d

Verify that Kafka containers are running:

docker ps

2. Create Kafka Topics via Kafka UI

After Kafka is running, open Kafka UI in your browser.

Go to the Kafka UI address (e.g., http://localhost:8080
).

Navigate to the Topics section.

Create the following topics manually:

order_topic

inventory_topic

3. Prepare .env files for each service

Each service requires a .env file. If the file does not exist, create it in the service directory with the following content:

Inventory Service (inventory/.env)

MONGO_CONN_URL="mongodb://admin:secret123@mongo-inventory:27017/mydb?authSource=admin"
KAFKA_BROKER_ADDRESS="kafka:9092"
KAFKA_ORDER_TOPIC="order_topic"
KAFKA_GROUP="my_group"
KAFKA_INVENTORY_TOPIC="inventory_topic"

Order Service (order/.env)

COUCHBASE_URL="couchbase://couchbase-order"
COUCHBASE_USERNAME="admin"
COUCHBASE_PASSWORD="123456"
KAFKA_BROKER_ADDRESS="kafka:9092"
KAFKA_TOPIC="order_topic"

Payment Service (payment/.env)

MONGO_CONN_URL="mongodb://admin:secret123@mongo-payment:27018/mydb?authSource=admin"
KAFKA_BROKER_ADDRESS="kafka:9092"
KAFKA_INVENTORY_TOPIC="inventory_topic"
KAFKA_GROUP="my_group"

Note: Make sure the .env files are in the root of their respective service directories and are not ignored by .dockerignore.
All services are connected via a Docker network, so they can reach each other using service names defined in Docker Compose.

4. Run the Services

Each microservice has its own docker-compose.yml. Start them one by one in the following order:

Start Order Service
cd order
docker compose up -d

Start Inventory Service
cd ../inventory
docker compose up -d

Start Payment Service
cd ../payment
docker compose up -d

5. Verify Services

Check all running containers:

docker ps

You should see:

Kafka broker(s)

Order service

Inventory service

Payment service

MongoDB / Couchbase (if included)