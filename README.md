# Booking Service

A microservice for managing bookings, tickets, and events. This project is built with Golang, Gin, GORM, PostgreSQL, and Kafka.

---

## Documentation

### 1. High-Level Design
Refer to the [HighLevelDesignRFC.pdf](./HighLevelDesignRFC.pdf) for the architectural overview of the system.

### 2. Saga Pattern for Distributed Transactions
Refer to the [SagaPatternForDistributedTransactionRFC.pdf](./SagaPatternForDistributedTransactionRFC.pdf) for insights into how the Saga pattern is utilized to ensure consistency in distributed transactions across multiple services.

---

## How to Run the Project Locally

1. Clone the repository to your local machine:
    ```bash
    git clone <repository-url>
    cd booking-service
    ```

2. Start the project by running:
    ```bash
    docker-compose up
    ```
   This will spin up the entire stack, including:
   - **PostgreSQL**: For storing bookings, tickets, and events.
   - **Kafka**: For handling messaging between microservices.

3. Use the provided [Postman collection](./postman_collection.json) to interact with the APIs. Import the collection into Postman to test the endpoints.

---

## Key Features

### 1. Event Seeding
When the project is initialized for the first time, the database will be seeded with **Event 1**:
- `ID: 1`
- `Name: Event 1`
- `Date: 2024-11-16`
- `Location: Ho Chi Minh`
- `Capacity: 1000`

If the database is already initialized, the seed will not be applied again.

---

### 2. Mimicking Booking Service Receiving Kafka Messages
To simulate the booking service receiving a message to create tickets for an event, use the **Create Tickets for Event** API. This endpoint allows you to manually create tickets for a specific event ID, mimicking the functionality that would occur upon receiving a Kafka message.

---

## Areas for Improvement

The following enhancements are still pending:

1. **Finalize Kafka Message Consumers**:
   - Properly handle Kafka messages, especially for the payment domain.
   - Consume messages to ensure the system reacts to cross-domain events seamlessly.

2. **Fire Kafka Messages for Booking Status Changes**:
   - Emit Kafka events when bookings are created, confirmed, or canceled.
   - These messages should notify other domains (e.g., payments, notifications) to update their status.

3. **API and Validation Improvements**:
   - Enhance input validation for the APIs, especially for ticket and booking creation.
   - Add comprehensive error handling with detailed response codes and messages.

4. **Event Capacity Management**:
   - Introduce capacity checks for events to ensure tickets do not exceed event limits.
   - Add error responses if the capacity is exceeded when creating tickets.

5. **Testing Improvements**:
   - Expand test cases to include edge cases and negative scenarios.
   - Improve mocking for services like Kafka and TicketService to ensure proper isolation.

6. **Docker and Deployment Enhancements**:
   - Optimize the `docker-compose` setup to allow scaling services dynamically.
   - Use a proper Kafka setup for local development with retries and proper topic configurations.

7. **Documentation Updates**:
   - Add detailed documentation for all APIs, including request and response payloads.
   - Provide instructions for extending the service (e.g., adding new event domains).

---

## Contributing

We welcome contributions to improve this project. Please fork the repository and submit a pull request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

---

Happy coding! 🚀
