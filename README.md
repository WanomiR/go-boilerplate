# Go project

## Instructions
- To install tools, lint, format, build, and run the app use:
    ```bash
    task run
    # or just
    task
    ```
- To only build and run the app use:
    ```bash
    task docker:up
    ```
- To stop and remove docker containers:
    ```bash
    task docker:down
    ```
- To free up docker-used memory run:
    ```bash
    task docker:clean
    ```
- To see all tasks run:
    ```bash
    task --list
    ```
## Notes

- Service initialization and running logic is within `internal/app` package, `main` only creates and launches the app;
