# Go vs Node Experiment Benchmark

Simple POST based create-user endpoint inserting into Postgres.

Most of the bottleneck is in the network between the app and postgres handling the inserts.

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

