# event-driven

Project Setup Guide

This guide explains how to set up and run the system step by step. The system uses Kafka and multiple Go-based microservices: order, inventory, and payment.

1. Run Kafka from the root directory

In the project root directory, there is a docker-compose.yml file for Kafka.
Start Kafka with the following command:

docker compose up -d

Verify that Kafka containers are running with:

docker ps

2. Create Kafka Topics via Kafka UI

After Kafka is running, open Kafka UI in your browser.

Go to the configured Kafka UI address (for example: http://localhost:8080
).

Navigate to the Topics section.

Create the following topics manually:

order_topic

inventory_topic

3. Check .env files

Each service (order, inventory, payment) requires a .env file for configuration.

If the .env file does not exist inside a service directory, create one.

At minimum, it should contain the Kafka broker connection details and any database URIs.

Example values:
KAFKA_BROKER=localhost:9092
MONGO_URI=mongodb://admin:secret123@mongo-service:27017/mydb

4. Run the Services

Each microservice (order, inventory, payment) has its own docker-compose.yml.
Start them one by one in the following order:

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

You should see Kafka broker(s), order service, inventory service, payment service, and MongoDB (if included).