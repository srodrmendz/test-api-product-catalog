# Authentication Microservice

The Authentication Microservice is a Go-based project aimed at providing authentication functionalities. It follows a modular structure and includes directories for repository, service, docs, utils, and a Makefile.

## Folder Structure

-   **repository**: Contains implementations of repositories for data storage and retrieval.
-   **service**: Holds business logic and services responsible for handling authentication operations.
-   **docs**: Swagger documentation.
-   **utils**: Utilities and helper functions shared across various parts of the project.
-   **Makefile**: Makefile for automating common development tasks, including testing and linting.

## Installation

To install dependencies execute `make deps`.

## Testing

To run tests, execute `make test` in the terminal. This command will execute all unit tests in the project.

## Linting

Run linting checks by executing `make lint` in the terminal. This command will check the codebase for style and formatting errors with **golangci-lint**.

## Create swagger docs

To create swagger docs run `make docs`.
