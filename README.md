# Sum-Tel

A microservices-based system for parsing messages from Telegram channels and generating news summaries.

## Description

Sum-Tel is an early-stage project aimed at extracting messages from Telegram channels through parsing techniques and generating concise summaries of news content. The system is built using Go and follows a microservices architecture to ensure scalability and modularity.

## Features

- **Message Parsing**: Extract messages from specified Telegram channels.
- **News Summarization**: Generate AI-powered summaries of news articles.
- **Microservices Architecture**: Modular design with separate services for core functionality and parsing.
- **gRPC Communication**: Services communicate via gRPC for efficient inter-service calls.

## Architecture

The project consists of the following components:

- **Protos**: Protocol buffer definitions for gRPC services.
- **Services**:
  - `core`: Handles user management, channels, and subscriptions.
  - `parser`: Responsible for parsing Telegram messages and processing them.
- **Shared**: Common utilities and configurations.

## Getting Started

### Prerequisites

- Go 1.21
 or later
- PostgreSQL
- Docker (for containerized deployment)
- protoc (for generating gRPC code)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Negat1v9/sum-tel.git
   cd sum-tel
   ```

2. Generate gRPC code:
    ```bash
    make gen_proto
    ```

### Running

1. Start the services:
Use Docker Compose for full setup:
   ```bash
   docker-compose up --build -d
   ```

## License

This project is licensed under the MIT License. See the LICENSE file for details.
