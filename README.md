# Go vs Node Experiment Benchmark
Simple POST-based create-user endpoint inserting into Postgres.

## Running on AWS
### Requirements
- AWS account
- [AWS CLI](https://aws.amazon.com/cli/) (Installed and configured)
- [OpenTofu](https://opentofu.org/)
- Postgres RDS Instance
  - Create manually on AWS
  - Run the `tofu/create-db.sql` script on it
  - Create a .env file on `node-api` and `go-api` (you can copy the .env.example)
  - Populate the `POSTGRES_` variables

### How to run
```
cd tofu
tofu apply -auto-approve
```

This will create all the infrastructure required to spin up two Ubuntu servers (one to host the API, one to host the _gun_).
After created, it'll generate two files that can be used to connect via SSH to each server.

#### Running the API
SSH into the API by running 
```
./ssh_connect_api.sh
```

Then, inside the API server:

To start the Node API, run:
```
cd node-api
npm install
npm start
```

To start the Go API, run:
```
cd go-api
go build -o api
./api
```

Each API will output the **PID** (Process ID), and store it somewhere.

##### Monitoring the process
On a new terminal, reconnect to the API server.

Then, run:
```
./monitor_process.sh 2150 2000
```

If the PID is 2150 and you want to start with the 2,000 requests per second load

#### Running the GUN
Connect to the gun server by running:
```
./ssh_connect_gun.sh
```

Then, to start the test for 2,000 requests per second, run:
```
cd load-tester/vegeta
./metrics.sh 2000
```

This will run a stress test for 30s with 2,000 requests per second and it'll save some important metrics to a csv file.

## Running Locally
You can use docker-compose to start all the services.

First, start the postgres, node-api and go-api.

```
docker-compose up -d postgres node-api go-api
```

And you should see something like this
<img width="870" alt="image" src="https://github.com/ocodista/api-benchmark/assets/19851187/0aad0411-d171-415e-b2fd-c6c8cbad2222">

Then, after they're all up and running, start the **gun**:
```
docker-compose up gun
```

After the end of the test, you'll be able to see the core metrics on the terminal.
<img width="1595" alt="image" src="https://github.com/ocodista/api-benchmark/assets/19851187/50b146d4-201a-42fc-82f2-7167d1a3d82e">


# Remarks
- In Docker Compose, prefer `network_mode: host`, as the docker internal networking adds considerable bottleneck
- Adding Node.js Cluster can handle more requests/sec but also adds extra networking and coordination overhead
- Should you choose Go blindly? No, only if throughput is the only criterion, disregarding productivity, integration, team preference, toolset, etc.

