# Docker Installation

Details about how to use Janus in Docker can be found on the Quay.io repo hosting the image: [janus](https://quay.io/repository/hellofresh/janus). We also have a some cool examples with [Docker Compose template](https://github.com/hellofresh/janus/blob/master/examples) with built-in orchestration and scalability.

Here is a quick example showing how to link a Janus container to a Cassandra or PostgreSQL container:

1. **Start your database:**

    If you wish to use a database instead of a file system based configuration just start the mongodb container:

    ```bash
    $ docker run -d --name janus-database \
                  -p 27017:27017 \
                  mongo:3.0
    ```

2. **Start your key value storage:**

    You can choose to use redis as your key value store instead of the in memory one. If you choose to do so you just start the container:
You will also need to set the storage dsn using the `STORAGE_DNS` to something like `memory://localhost` or if you use redis `redis://janus-storage:6379`

    ```bash
    $ docker run -d --name janus-storage \
                  -p 6379:6379 \
                  redis:3.0
    ```


3. **Start Janus:**

    Start a Janus container and link it to your database container (if you are using it), configuring the `DATABASE_DSN` environment variable with the connection string like `mongodb://janus-database:27017/janus`:
    

    ```bash
    $ docker run -d --name janus \
                  --link janus-database:janus-database \
                  --link janus-storage:janus-storage \
                  -e "DATABASE_DSN=mongodb://janus-database:27017/janus" \
                  -e "STORAGE_DNS=redis://janus-storage:6379" \
                  -p 8080:8080 \
                  -p 8443:8443 \
                  -p 8081:8081 \
                  -p 8444:8444 \
                  quay.io/hellofresh/janus
    ```

3. **Janus is running:**

    ```bash
    $ curl http://127.0.0.1:8081/
    ```
