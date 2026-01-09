Welcome to the kafka cloud system

To run this system, feel free to simply activate

./dev_run.sh

if you wish to manually initialize the system, or only a select few services do the following:

docker compose up -d

this should activate the docker compose configuration that would activate:

- PostgreSQL for the User Verification Service

- Zookeeper, the Kafka Broker Orchestrator

- The Kafka Broker itself, on port 9092

To initialize a given service, simply enter its directory and activate it via the go cli, for example:

    cd consumer

    go mod download #only once

    go run .

Each service has it's own script file to test out client requests, for example:

    cd producer

    ./post_scripts.sh produce

To run the web client, both NodeJS and npm must be installed:

    curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
    sudo apt install nodejs

If you wish to use the client website, all the service must be active first, once they are all up, simply do :

    cd client
    npm install
    npm run dev

Or feel free to run the initialization script that activates the dockr composition, all services, and finally the client
