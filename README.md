# Simple Web Service
This project is a service designed to display order data through a simple interface. 

It integrates with PostgreSQL for data storage, NATS Streaming for real-time data updates, and provides a basic web interface for data visualization.

## Features
PostgreSQL Integration: Utilizes PostgreSQL for persistent storage of order data.

NATS Streaming Subscription: Subscribes to a NATS Streaming channel to receive real-time updates on order data.

In-Memory Caching: Implements in-memory caching to speed up data retrieval and ensure service resilience.

Data Persistence: Ensures data integrity and service continuity by restoring cache from the database in case of service failure.

Web Interface: Provides a basic web interface for displaying order data by ID.

## Getting Started
## Prerequisites
Go 

PostgreSQL 13 or later

NATS Streaming Server

## Installation
Clone the repository:

git clone https://github.com/KarinaRdk/Web-server.git
cd DemonstrationService
Install dependencies:

go mod download

Set up your PostgreSQL database and NATS Streaming Server.

Configure your application by editing the config/local.yaml file. Ensure the database connection string and NATS Streaming Server details are correctly set.

Build and run the application

## Usage
HTTP Endpoints:
/get_order: Retrieves an order by its ID.
/: Serves the HTML page for displaying order data.

## Subscriber:
Listens for messages on the NATS Streaming Server and updates the database and cache with new orders.

## To publish a message
go mod download github.com/nats-io/nats.go
go mod download github.com/nats-io/stan.go
go run internal/publisher/main.go

if you wish to edit the message you can do it by editing internal/publisher/model.json