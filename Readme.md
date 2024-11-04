## How to Run

### 1. Start the Database

- Open a terminal.
- Run `docker-compose up -d` to start the services in detached mode.

### 2. Start the Server

- Run `./bin/pbserver` to start the server on port 1234.

### 3. Start the Client UX

- Open another terminal.
- Run `./bin/pbclient http://localhost:1234`.

### 4. Available Commands

- Use the `help` command to see the available commands.

### 5. Run tests

- Run `go test -v -count=1 ./tests` to start tests.

---

Good luck!
