# library-management-app

This project aims to create a scalable backend system for a library management application using microservices architecture.

## Applications

The applications to be created include:
- Book Management
- Authors Management
- Book Categories Management
- Borrowing and Returning Books

## Requirements

- Golang
- PostgreSQL
- Redis

## Installation

1. **Clone Repository:**
2. **Setup Database:**
- Install PostgreSQL and Redis if not already installed.
- Configure PostgreSQL and Redis connection details in the application's configuration files.
3. **Install Dependencies:**


## Run

Use the Makefile provided to build and run the services. Ensure PostgreSQL and Redis are running before starting the services.

Example:
```bash
make run-user
