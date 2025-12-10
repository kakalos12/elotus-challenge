# elotus-challenge

## How to spin up the Docker Compose file

To get the application running using Docker Compose, follow these steps:

1.  Navigate to the `server` directory:
    ```bash
    cd server
    ```

2.  Build and start the services defined in `docker-compose.yml` in detached mode:
    ```bash
    docker compose up -d
    ```

    This will start two services:
    *   `db`: A MySQL 8.0 database, accessible on `localhost:3306`.
        *   **User**: `elotus_user`
        *   **Password**: `elotus_password`
        *   **Database Name**: `elotus`
        *   **Root Password**: `rootpassword`
    *   `server`: The Go application server, accessible on `localhost:8081`.

3.  To stop the services, run:
    ```bash
    docker-compose down
    ```

4.  To view the logs of all services:
    ```bash
    docker-compose logs -f
    ```
