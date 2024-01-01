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
./ssh_gun_server.sh
```

Then, to start the test for 2,000 requests per second, run:
```
cd load-tester/vegeta
./metrics.sh 2000
```

This will run a stress test for 30s with 2,000 requests per second and it'll save some important metrics to a csv file.


## Running Locally
To run the Go version:

```
cd go-api
docker compose up --build
```

To run the Node.js version:

```
cd node-api
docker compose up --build
```

To run the Vegeta stress test:

```
cd load-tester/vegeta
./start.sh 3000 # 3000 req/s, you can try with more
```

# Results

Running on a Ryzen 9 7950X3D with 96GB of DDR56000 and PCIe Gen 5 NVME.

Vegeta results for 3000 req/s for Go:

```
Starting Vegeta attack for 30s at 3000 requests per second...
Load test finished, generating reports...
Textual report generated: report_3000.txt


Requests      [total, rate, throughput]         89998, 2999.95, 2999.76
Duration      [total, attack, wait]             30.002s, 30s, 1.828ms
Latencies     [min, mean, 50, 90, 95, 99, max]  1.301ms, 2.348ms, 1.834ms, 2.909ms, 3ms, 3.823ms, 123.211ms
Bytes In      [total, mean]                     2429946, 27.00
Bytes Out     [total, mean]                     8279816, 92.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      201:89998
Error Set:
```

Vegeta results for 3000 req/s for Node.js (no cluster):

```
Starting Vegeta attack for 30s at 3000 requests per second...
Load test finished, generating reports...
Textual report generated: report_3000.txt


Requests      [total, rate, throughput]         89999, 3000.01, 2999.84
Duration      [total, attack, wait]             30.001s, 30s, 1.683ms
Latencies     [min, mean, 50, 90, 95, 99, max]  1.448ms, 4.458ms, 1.931ms, 3.097ms, 3.249ms, 97.854ms, 290.991ms
Bytes In      [total, mean]                     2339974, 26.00
Bytes Out     [total, mean]                     8279908, 92.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      201:89999
Error Set:
```

# Remarks

- In Docker Compose, prefer `network_mode: host`, as the docker internal networking adds considerable bottleneck
- Node.js will undoubtfully be "slower". Average, time per request will be slower than Go
- Adding Node.js Cluster can handle more requests/sec but also adds extra networking and coordination overhead
- Yes. Go will be faster
- Should you choose Go blindly? No, only if throughput is the only criteria, disregarding productivity, integration, team preference, toolset, etc.

