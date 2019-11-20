# go-structured-logging-demo

Sample app showing the power of structured logging.

1. Run application using docker compose.

   ```bash
   docker-compose up
   ```

   Alternatively run it with Seq <http://localhost:5341>.

   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.seq.yml up
   ```

   Alternatively run it with Greylog <http://localhost:9000>.

   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.greylog.yml up
   ```

1. Open <http://localhost:8080> in a browser several times.
