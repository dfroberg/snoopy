# Snoopy
Snoopy subscribes to events on the Ethereum network you specify and spits out stats and information about blocks it has gathered since it started.

*NOTES:*
* There is no persistance layer yet.
* Snoopy uses infura API endpoints to access the network.

|Endpoint|Port|Function|Method|Auth Reqd.|
|---|---|---|---|---|
|/health|9080|Cluster liveness probe endpoint|GET|No|
|/ping|9080|Pongs if alive|GET|No|
|/|9080|returns stats on snooped blocks|GET|Token|
|/blockid|9080|Return dump of block with internal id|POST|Token|
|/blockhash|9080|Return dump of block with hash|POST|Token|
|/blocknumber|9080|Return dump of block with number|POST|Token|
|/metrics|2112|Prometheus metrics endpoint|GET|No|

# Some ideas:
* To make this little widget useful, add filters for To/From so you can capture events related to something that matters to YOU, i.e. your wallet address or some payment destination.
## Set it up
Register for an [INFURA Project ID](https://infura.io/register)
You will have to use this key for subsequent requests to INFURA endpoints,
as briefly shown in the [Choose a network](https://infura.io/docs/gettingStarted/chooseaNetwork) section of the site

Clone this repo and proceed to the steps below;

# Tests
~~~
export SNOOPY_PROJECT_ID=<INFURA PROJECT ID>
export SNOOPY_NETWORK_NAME=ropsten
export SNOOPY_API_TOKEN=TestToken
~~~
An very minimalistic test;
~~~
$ go test -cover
~~~
~~~
2022/05/04 16:05:02 Success! you connected to the ropsten Network
2022/05/04 16:05:02 Success! you connected to the ropsten Network
2022/05/04 16:05:10 Got: {"Id":1,"BlockHash":"0x90d337977aa098f7f69b19fe29e09464486d725f58fa84b1ccdcb04246d74ada","BlockNumber":"14711389","BlockTime":1651673102,"BlockNonce":7604351258204595666,"BlockNumTransactions":134}
PASS
coverage: 33.1% of statements
ok      snoopy/v2       11.100s
~~~
Quickly build and run the executable;
~~~
$ cd src/
$ go build . && ./snoopy
~~~
Starts the snooping on whatever net you've specified above in the ENV's
Use the curl commands below in Endpoints to test.

# Deploy
~~~
export SNOOPY_PROJECT_ID=<INFURA PRODUCTION PROJECT ID>
export SNOOPY_NETWORK_NAME=mainnet
export SNOOPY_API_TOKEN=TestToken
export SNOOPY_NAMESPACE=snoopy
export SNOOPY_INGRESS_CLASS=traefik
export SNOOPY_INGRESS_HOST=snoopy.local
export SNOOPY_VERSION=v0.6.4
~~~
Create the snoopy namespace;
~~~
envsubst < <(cat manifests/snoopy-namespace.yaml) | kubectl apply -f -
~~~
~~~
namespace/snoopy created
~~~
Since secrets should be actually be applied via Vault / SOPS / Age etc,
this is a workaround;
~~~
envsubst < <(cat manifests/snoopy-*.yaml) | kubectl apply -f -
~~~
~~~
secret/common-snoopy-secret created
configmap/snoopy-config-map created
deployment.apps/snoopy created
ingress.networking.k8s.io/snoopy-ingress created
namespace/snoopy unchanged
service/snoopy-service created
~~~
*NOTE:*
Please adjust the ingress!
If it's a local domain adjust your /etc/hosts to point to your ingress controller IP or annotations to match your environment otherwise.

# Load Test
This app manages a measly 34 rps (mean) on the /blocks and  389.98 rps (mean) on the /health endpoint in kuberenetes, this is more than likely du to my local infrastructure rather than anything else, as it gets 1912.50 rps(mean) on localhost on the same endpoint and 18313.34 rps (mean) on / (stats) and 23187.33 rps (mean) on /health;

## Localhost
~~~
ab -n2000 -c100 -H "X-Token: TestToken" http://localhost:9080/blocks
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 200 requests
Completed 400 requests
Completed 600 requests
Completed 800 requests
Completed 1000 requests
Completed 1200 requests
Completed 1400 requests
Completed 1600 requests
Completed 1800 requests
Completed 2000 requests
Finished 2000 requests


Server Software:
Server Hostname:        localhost
Server Port:            9080

Document Path:          /blocks
Document Length:        2225 bytes

Concurrency Level:      100
Time taken for tests:   1.046 seconds
Complete requests:      2000
Failed requests:        0
Total transferred:      4626000 bytes
HTML transferred:       4450000 bytes
Requests per second:    1912.50 [#/sec] (mean)
Time per request:       52.288 [ms] (mean)
Time per request:       0.523 [ms] (mean, across all concurrent requests)
Transfer rate:          4319.93 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.3      0       3
Processing:     0   51  20.9     47     140
Waiting:        0   51  20.8     47     140
Total:          0   51  20.8     47     140

Percentage of the requests served within a certain time (ms)
  50%     47
  66%     55
  75%     59
  80%     63
  90%     71
  95%     94
  98%    117
  99%    123
 100%    140 (longest request)
~~~
## K8s:
~~~
ab -n2000 -c100 -H "X-Token: TestToken" https://snoopy.local/blocks
~~~
The test below is with ~500 blocks deployed in a k8s cluster.
Shows quite a bit of improvement needed, but it keeps working along while under load.
~~~
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking snoopy.local (be patient)
Completed 200 requests
Completed 400 requests
Completed 600 requests
Completed 800 requests
Completed 1000 requests
Completed 1200 requests
Completed 1400 requests
Completed 1600 requests
Completed 1800 requests
Completed 2000 requests
Finished 2000 requests


Server Software:
Server Hostname:        snoopy.local
Server Port:            443
SSL/TLS Protocol:       TLSv1.2,ECDHE-RSA-AES128-GCM-SHA256,2048,128
Server Temp Key:        X25519 253 bits
TLS Server Name:        snoopy.local

Document Path:          /blocks
Document Length:        85066 bytes

Concurrency Level:      100
Time taken for tests:   57.467 seconds
Complete requests:      2000
Failed requests:        1938
   (Connect: 0, Receive: 0, Length: 1938, Exceptions: 0)
Non-2xx responses:      730
Total transferred:      108921400 bytes
HTML transferred:       108717660 bytes
Requests per second:    34.80 [#/sec] (mean)
Time per request:       2873.367 [ms] (mean)
Time per request:       28.734 [ms] (mean, across all concurrent requests)
Transfer rate:          1850.94 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        8   43  74.7     15     621
Processing:     2 2826 2412.0   3385   19972
Waiting:        2 2801 2396.0   3364   19873
Total:         18 2869 2397.2   3400   19989

Percentage of the requests served within a certain time (ms)
  50%   3400
  66%   3814
  75%   4193
  80%   4402
  90%   5212
  95%   6698
  98%   8297
  99%   9500
 100%  19989 (longest request)
 ~~~
# Endpoints
## Getting cumulative stats
~~~
curl -s -X GET -H "X-Token: TestToken" http://localhost:9080/ | jq
~~~
~~~
{
  "NumBlocks": 5,
  "NumTx": 95,
  "NumAuthRequests": 6,
  "NumUnAuthRequests": 3,
  "NumSystemRequests": 2,
  "NumApiConns": 2
}
~~~
## Dumping blocks in memory:
~~~
curl -s -X GET -H "X-Token: TestToken" http://localhost:9080/blocks | jq
~~~
~~~
{
    "1": {
        "Id": 1,
        "BlockHash": "0xfdca36c087e54707159880b4f9d3f5d8e507f1bf79f7ab4f4204299c123b2ea8",
        "BlockNumber": 12234036,
        "BlockTime": 1651564772,
        "BlockNonce": 14773452359893719969,
        "BlockNumTransactions": 32
    },
    "2": {
        "Id": 2,
        "BlockHash": "0x966ad79c014cd222980c5d825f404451cc7fe71e58421cf885b4f9d1164b29de",
        "BlockNumber": 12234038,
        "BlockTime": 1651564871,
        "BlockNonce": 10722525007728445607,
        "BlockNumTransactions": 18
    }
}
~~~
## Get Block Data by Hash
~~~
curl -s -H "X-Token: TestToken" -d '{"Id": 1}' http://localhost:9080/blockid | jq
~~~
~~~
{
  "Id": 1,
  "BlockHash": "0xfdca36c087e54707159880b4f9d3f5d8e507f1bf79f7ab4f4204299c123b2ea8",
  "BlockNumber": 12234036,
  "BlockTime": 1651564772,
  "BlockNonce": 14773452359893720000,
  "BlockNumTransactions": 32
}
~~~
## Get Block Data by Hash
~~~
curl -s -H "X-Token: TestToken" -d '{"Hash": "0xd92a881dd3e68c25fc78e9495f200e4e2bb26370bc5ec9c02c809700fa2354c8"}' http://localhost:9080/blockhash | jq
~~~
~~~
[
  {
    "Id": 1,
    "BlockHash": "0x1f19810bc87163449aed4ebc32b2d542d285c39e300954b4b51af2c4eee6f7f4",
    "BlockNumber": "12234088",
    "BlockTime": 1651566840,
    "BlockNonce": 7343538744684689000,
    "BlockNumTransactions": 28
  }
]
~~~
## Get Block Data by Number
~~~
curl -s -H "X-Token: TestToken" -d '{"Number": "14711278"}' http://localhost:9080/blocknumber | jq
~~~
~~~
[
  {
    "Id": 4,
    "BlockHash": "0xf2075ccaa5e1f77b063fa58e06479678681ed60409356a65a28e8a21b20e2b67",
    "BlockNumber": "14711278",
    "BlockTime": 1651671625,
    "BlockNonce": 16791431146648115000,
    "BlockNumTransactions": 32
  }
]
~~~
## Healtcheck
~~~
curl -s -X GET -H "X-Token: TestToken" http://localhost:9080/health | jq
~~~
~~~
{
  "alive": "true"
}
~~~
## Ping
~~~
curl -s -X GET -H "X-Token: TestToken" http://localhost:9080/ping | jq
~~~
~~~
{
  "ping": "pong"
}
~~~
## Unknown url
~~~
curl -s -X GET -H "X-Token: TestToken" http://localhost:9080/weirdurl
~~~
~~~
404 page not found
~~~

## Wrong or no token
~~~
curl -s -X GET -H "X-Token: WrongToken" http://localhost:9080/
~~~
~~~
Forbidden
~~~
# Prometheus Instrumentation
~~~
curl -s -X GET http://localhost:2112/metrics
~~~
Prometheus endpoint metrics
~~~
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 4.79e-05
go_gc_duration_seconds{quantile="0.25"} 4.79e-05
go_gc_duration_seconds{quantile="0.5"} 8.02e-05
go_gc_duration_seconds{quantile="0.75"} 8.02e-05
go_gc_duration_seconds{quantile="1"} 8.02e-05
go_gc_duration_seconds_sum 0.0001281
go_gc_duration_seconds_count 2
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 17
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.18.1"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 1.722832e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 6.946744e+06
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.448599e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 41541
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 2.0503668434863447e-06
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 5.092192e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 1.722832e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 4.66944e+06
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 3.260416e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 5831
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 2.4576e+06
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 7.929856e+06
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.6515687268106744e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 47372
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 4800
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15600
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 68408
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 81600
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.194304e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 784993
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 458752
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 458752
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 1.5811592e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 8
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 0.06
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 4096
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 12
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 2.3465984e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.65156860536e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.188114432e+09
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes -1
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 2
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# HELP snoopy_processed_apicalls_total The total number of processed api calls
# TYPE snoopy_processed_apicalls_total counter
snoopy_processed_apicalls_total 2
# HELP snoopy_processed_blocks_total The total number of processed blocks
# TYPE snoopy_processed_blocks_total counter
snoopy_processed_blocks_total 8
# HELP snoopy_processed_events_total The total number of processed events
# TYPE snoopy_processed_events_total counter
snoopy_processed_events_total 8
# HELP snoopy_processed_requests_total The total number of processed requests
# TYPE snoopy_processed_requests_total counter
snoopy_processed_requests_total 3
# HELP snoopy_processed_transactions_total The total number of processed transactions
# TYPE snoopy_processed_transactions_total counter
snoopy_processed_transactions_total 1543
~~~
