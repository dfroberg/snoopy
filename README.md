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
|/blocks|9080|Return dump of block|GET|Token|
|/blockid|9080|Return dump of block with internal id|POST|Token|
|/blockhash|9080|Return dump of block with hash|POST|Token|
|/blocknumber|9080|Return dump of block with number|POST|Token|
|/txs|9080|Return dump of transactions|GET|Token|
|/txid|9080|Return dump of transaction with internal id|POST|Token|
|/txnumber|9080|Return dump of transaction in blocknumber number|POST|Token|
|/filters|9080|Return dump of filters|GET|Token|
|/filteradd|9080|Add a TxTo address filter|POST|Token|
|/filterdelete|9080|Remove a TxTo address filter|POST|Token|
|/filterid|9080|Return filter matching filter id|POST|Token|
|/filterto|9080|Return filter matching TxTo|POST|Token|
|/metrics|2112|Prometheus metrics endpoint|GET|No|

# Some ideas:
* WIP: To make this little widget useful, add filters for To/From so you can capture events related to something that matters to YOU, i.e. your wallet address or some payment destination.
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
curl -s -H "X-Token: TestToken" -d '{"Number": 14711278}' http://localhost:9080/blocknumber | jq
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
## Get Transactions
~~~
curl -s -H "X-Token: TestToken" http://localhost:9080/txs | jq
~~~
~~~
{
  "1": {
    "Id": 1,
    "TxBlockId": 22,
    "TxBlockNumber": 14717116,
    "TxHash": "0xda6a3cce91baf1a4d8648fb659c6630078fae6a72ea12caff6aff975d92617c7",
    "TxValue": 1.5e+19,
    "TxGas": 393786,
    "TxGasPrice": 35792179968,
    "TxCost": 15014094459380880000,
    "TxNonce": 1018,
    "TxTo": "0x75A6787C7EE60424358B449B539A8b774c9B4862",
    "TxReceiptStatus": 1
  },
  ...
  "99": {
    "Id": 99,
    "TxBlockId": 22,
    "TxBlockNumber": 14717116,
    "TxHash": "0x5f1406e75002398534d874d0392ad52a14cb983aab213856b91f2ed757a9fa7c",
    "TxValue": 15375963677623800,
    "TxGas": 21000,
    "TxGasPrice": 56000000000,
    "TxCost": 16551963677623800,
    "TxNonce": 1,
    "TxTo": "0xA090e606E30bD747d4E6245a1517EbE430F0057e",
    "TxReceiptStatus": 1
  }
}
~~~
## Get Transaction Data by Id
~~~
curl -s -H "X-Token: TestToken" -d '{"Id": 1}' http://localhost:9080/txid | jq
~~~
~~~
{
  "Id": 1,
  "TxBlockId": 20,
  "TxBlockNumber": 14717114,
  "TxHash": "0x55bcc6f4fdf880ff81612da123c1795e798cc82cb559bdd8a70a347100820533",
  "TxGas": 50000,
  "TxGasPrice": 50000000000,
  "TxCost": 2500000000000000,
  "TxNonce": 4,
  "TxTo": "0xdAC17F958D2ee523a2206206994597C13D831ec7",
  "TxReceiptStatus": 1
}
~~~
## Get Transaction Data by Block Number
~~~
curl -s -H "X-Token: TestToken" -d '{"Number": 14717097}' http://localhost:9080/txnumber | jq
~~~
~~~
[
  {
    "Id": 1,
    "TxBlockId": 3,
    "TxBlockNumber": 14717097,
    "TxHash": "0x9c628cb1623b22d4886786f61f23dedd950d25a3897f5852b61cf863b64a401e",
    "TxGas": 700000,
    "TxGasPrice": 51200000000,
    "TxCost": 35840000000000000,
    "TxNonce": 147495,
    "TxTo": "0x0000006daea1723962647b7e189d311d757Fb793",
    "TxReceiptStatus": 1
  },
  ...
  {
    "Id": 211,
    "TxBlockId": 3,
    "TxBlockNumber": 14717097,
    "TxHash": "0x43608849e0ffff749426b70a68700377308647df69bfdb8a915786ff92c94d0e",
    "TxValue": 33389713214078600,
    "TxGas": 80000,
    "TxGasPrice": 58988117698,
    "TxCost": 38108762629918600,
    "TxNonce": 1,
    "TxTo": "0x2a67035357C3045438F3A92E46870a9E48e5AAB7",
    "TxReceiptStatus": 1
  }
]
~~~
## Get Filters
~~~
curl -s -H "X-Token: TestToken" http://localhost:9080/filters | jq
~~~
~~~
{
  "0": {
    "TxTo": "0xA090e606E30bD747d4E6245a1517EbE430F0057e"
  }
}
~~~
## Add Filter
~~~
curl -s -H "X-Token: TestToken" -d '{"To": "0xA090e606E30bD747d4E6245a1517EbE430F0057e"}' http://localhost:9080/filteradd | jq
~~~
~~~
[
  {
    "TxTo": "0xA090e606E30bD747d4E6245a1517EbE430F0057e"
  }
]
~~~
## Get Filter by To
return null on not found
~~~
curl -s -H "X-Token: TestToken" -d '{"To": "0xA090e606E30bD747d4E6245a1517EbE430F0057e"}' http://localhost:9080/filterto | jq
~~~
~~~
[
  {
    "TxTo": "0xA090e606E30bD747d4E6245a1517EbE430F0057e"
  }
]
~~~
## Get Filter by ID
return null on not found
~~~
curl -s -H "X-Token: TestToken" -d '{"Id": 1}' http://localhost:9080/filterid | jq
~~~
~~~
[
  {
    "TxTo": "0xA090e606E30bD747d4E6245a1517EbE430F0057e"
  }
]
~~~
## Delete Filter
~~~
curl -s -H "X-Token: TestToken" -d '{"Id": 1}' http://localhost:9080/filterdelete | jq
~~~
~~~
{
  "result": "true"
}
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
# HELP go_gc_cycles_automatic_gc_cycles_total Count of completed GC cycles generated by the Go runtime.
# TYPE go_gc_cycles_automatic_gc_cycles_total counter
go_gc_cycles_automatic_gc_cycles_total 91
# HELP go_gc_cycles_forced_gc_cycles_total Count of completed GC cycles forced by the application.
# TYPE go_gc_cycles_forced_gc_cycles_total counter
go_gc_cycles_forced_gc_cycles_total 0
# HELP go_gc_cycles_total_gc_cycles_total Count of all completed GC cycles.
# TYPE go_gc_cycles_total_gc_cycles_total counter
go_gc_cycles_total_gc_cycles_total 91
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 3.2e-05
go_gc_duration_seconds{quantile="0.25"} 6.02e-05
go_gc_duration_seconds{quantile="0.5"} 7.64e-05
go_gc_duration_seconds{quantile="0.75"} 9.57e-05
go_gc_duration_seconds{quantile="1"} 0.0005015
go_gc_duration_seconds_sum 0.007921
go_gc_duration_seconds_count 91
# HELP go_gc_heap_allocs_by_size_bytes_total Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_allocs_by_size_bytes_total histogram
go_gc_heap_allocs_by_size_bytes_total_bucket{le="8.999999999999998"} 3833
go_gc_heap_allocs_by_size_bytes_total_bucket{le="24.999999999999996"} 396617
go_gc_heap_allocs_by_size_bytes_total_bucket{le="64.99999999999999"} 756520
go_gc_heap_allocs_by_size_bytes_total_bucket{le="144.99999999999997"} 908358
go_gc_heap_allocs_by_size_bytes_total_bucket{le="320.99999999999994"} 915154
go_gc_heap_allocs_by_size_bytes_total_bucket{le="704.9999999999999"} 919215
go_gc_heap_allocs_by_size_bytes_total_bucket{le="1536.9999999999998"} 921959
go_gc_heap_allocs_by_size_bytes_total_bucket{le="3200.9999999999995"} 923398
go_gc_heap_allocs_by_size_bytes_total_bucket{le="6528.999999999999"} 923797
go_gc_heap_allocs_by_size_bytes_total_bucket{le="13568.999999999998"} 924101
go_gc_heap_allocs_by_size_bytes_total_bucket{le="27264.999999999996"} 924274
go_gc_heap_allocs_by_size_bytes_total_bucket{le="+Inf"} 925131
go_gc_heap_allocs_by_size_bytes_total_sum 2.28943952e+08
go_gc_heap_allocs_by_size_bytes_total_count 925131
# HELP go_gc_heap_allocs_bytes_total Cumulative sum of memory allocated to the heap by the application.
# TYPE go_gc_heap_allocs_bytes_total counter
go_gc_heap_allocs_bytes_total 2.28943952e+08
# HELP go_gc_heap_allocs_objects_total Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_allocs_objects_total counter
go_gc_heap_allocs_objects_total 925131
# HELP go_gc_heap_frees_by_size_bytes_total Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_frees_by_size_bytes_total histogram
go_gc_heap_frees_by_size_bytes_total_bucket{le="8.999999999999998"} 1519
go_gc_heap_frees_by_size_bytes_total_bucket{le="24.999999999999996"} 385937
go_gc_heap_frees_by_size_bytes_total_bucket{le="64.99999999999999"} 739035
go_gc_heap_frees_by_size_bytes_total_bucket{le="144.99999999999997"} 887315
go_gc_heap_frees_by_size_bytes_total_bucket{le="320.99999999999994"} 893147
go_gc_heap_frees_by_size_bytes_total_bucket{le="704.9999999999999"} 896751
go_gc_heap_frees_by_size_bytes_total_bucket{le="1536.9999999999998"} 899294
go_gc_heap_frees_by_size_bytes_total_bucket{le="3200.9999999999995"} 900620
go_gc_heap_frees_by_size_bytes_total_bucket{le="6528.999999999999"} 900980
go_gc_heap_frees_by_size_bytes_total_bucket{le="13568.999999999998"} 901260
go_gc_heap_frees_by_size_bytes_total_bucket{le="27264.999999999996"} 901430
go_gc_heap_frees_by_size_bytes_total_bucket{le="+Inf"} 902276
go_gc_heap_frees_by_size_bytes_total_sum 2.2493568e+08
go_gc_heap_frees_by_size_bytes_total_count 902276
# HELP go_gc_heap_frees_bytes_total Cumulative sum of heap memory freed by the garbage collector.
# TYPE go_gc_heap_frees_bytes_total counter
go_gc_heap_frees_bytes_total 2.2493568e+08
# HELP go_gc_heap_frees_objects_total Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_frees_objects_total counter
go_gc_heap_frees_objects_total 902276
# HELP go_gc_heap_goal_bytes Heap size target for the end of the GC cycle.
# TYPE go_gc_heap_goal_bytes gauge
go_gc_heap_goal_bytes 4.194304e+06
# HELP go_gc_heap_objects_objects Number of objects, live or unswept, occupying heap memory.
# TYPE go_gc_heap_objects_objects gauge
go_gc_heap_objects_objects 22855
# HELP go_gc_heap_tiny_allocs_objects_total Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size.
# TYPE go_gc_heap_tiny_allocs_objects_total counter
go_gc_heap_tiny_allocs_objects_total 134842
# HELP go_gc_pauses_seconds_total Distribution individual GC-related stop-the-world pause latencies.
# TYPE go_gc_pauses_seconds_total histogram
go_gc_pauses_seconds_total_bucket{le="-5e-324"} 0
go_gc_pauses_seconds_total_bucket{le="9.999999999999999e-10"} 0
go_gc_pauses_seconds_total_bucket{le="9.999999999999999e-09"} 0
go_gc_pauses_seconds_total_bucket{le="9.999999999999998e-08"} 0
go_gc_pauses_seconds_total_bucket{le="1.0239999999999999e-06"} 0
go_gc_pauses_seconds_total_bucket{le="1.0239999999999999e-05"} 18
go_gc_pauses_seconds_total_bucket{le="0.00010239999999999998"} 176
go_gc_pauses_seconds_total_bucket{le="0.0010485759999999998"} 182
go_gc_pauses_seconds_total_bucket{le="0.010485759999999998"} 182
go_gc_pauses_seconds_total_bucket{le="0.10485759999999998"} 182
go_gc_pauses_seconds_total_bucket{le="+Inf"} 182
go_gc_pauses_seconds_total_sum NaN
go_gc_pauses_seconds_total_count 182
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 15
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.18.1"} 1
# HELP go_memory_classes_heap_free_bytes Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory.
# TYPE go_memory_classes_heap_free_bytes gauge
go_memory_classes_heap_free_bytes 606208
# HELP go_memory_classes_heap_objects_bytes Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector.
# TYPE go_memory_classes_heap_objects_bytes gauge
go_memory_classes_heap_objects_bytes 4.008272e+06
# HELP go_memory_classes_heap_released_bytes Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory.
# TYPE go_memory_classes_heap_released_bytes gauge
go_memory_classes_heap_released_bytes 6.643712e+06
# HELP go_memory_classes_heap_stacks_bytes Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use.
# TYPE go_memory_classes_heap_stacks_bytes gauge
go_memory_classes_heap_stacks_bytes 458752
# HELP go_memory_classes_heap_unused_bytes Memory that is reserved for heap objects but is not currently used to hold heap objects.
# TYPE go_memory_classes_heap_unused_bytes gauge
go_memory_classes_heap_unused_bytes 865968
# HELP go_memory_classes_metadata_mcache_free_bytes Memory that is reserved for runtime mcache structures, but not in-use.
# TYPE go_memory_classes_metadata_mcache_free_bytes gauge
go_memory_classes_metadata_mcache_free_bytes 10800
# HELP go_memory_classes_metadata_mcache_inuse_bytes Memory that is occupied by runtime mcache structures that are currently being used.
# TYPE go_memory_classes_metadata_mcache_inuse_bytes gauge
go_memory_classes_metadata_mcache_inuse_bytes 4800
# HELP go_memory_classes_metadata_mspan_free_bytes Memory that is reserved for runtime mspan structures, but not in-use.
# TYPE go_memory_classes_metadata_mspan_free_bytes gauge
go_memory_classes_metadata_mspan_free_bytes 6528
# HELP go_memory_classes_metadata_mspan_inuse_bytes Memory that is occupied by runtime mspan structures that are currently being used.
# TYPE go_memory_classes_metadata_mspan_inuse_bytes gauge
go_memory_classes_metadata_mspan_inuse_bytes 91392
# HELP go_memory_classes_metadata_other_bytes Memory that is reserved for or used to hold runtime metadata.
# TYPE go_memory_classes_metadata_other_bytes gauge
go_memory_classes_metadata_other_bytes 5.15144e+06
# HELP go_memory_classes_os_stacks_bytes Stack memory allocated by the underlying operating system.
# TYPE go_memory_classes_os_stacks_bytes gauge
go_memory_classes_os_stacks_bytes 0
# HELP go_memory_classes_other_bytes Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more.
# TYPE go_memory_classes_other_bytes gauge
go_memory_classes_other_bytes 934823
# HELP go_memory_classes_profiling_buckets_bytes Memory that is used by the stack trace hash map used for profiling.
# TYPE go_memory_classes_profiling_buckets_bytes gauge
go_memory_classes_profiling_buckets_bytes 1.485345e+06
# HELP go_memory_classes_total_bytes All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes.
# TYPE go_memory_classes_total_bytes gauge
go_memory_classes_total_bytes 2.026804e+07
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 4.008272e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 2.28943952e+08
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.485345e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 1.037118e+06
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 0
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 5.15144e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 4.008272e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 7.24992e+06
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 4.87424e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 22855
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 6.643712e+06
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 1.212416e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.6517366291967328e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 1.059973e+06
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 4800
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15600
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 91392
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 97920
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.194304e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 934823
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 458752
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 458752
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 2.026804e+07
# HELP go_sched_goroutines_goroutines Count of live goroutines.
# TYPE go_sched_goroutines_goroutines gauge
go_sched_goroutines_goroutines 15
# HELP go_sched_latencies_seconds Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running.
# TYPE go_sched_latencies_seconds histogram
go_sched_latencies_seconds_bucket{le="-5e-324"} 0
go_sched_latencies_seconds_bucket{le="9.999999999999999e-10"} 1117
go_sched_latencies_seconds_bucket{le="9.999999999999999e-09"} 1117
go_sched_latencies_seconds_bucket{le="9.999999999999998e-08"} 1117
go_sched_latencies_seconds_bucket{le="1.0239999999999999e-06"} 3522
go_sched_latencies_seconds_bucket{le="1.0239999999999999e-05"} 3586
go_sched_latencies_seconds_bucket{le="0.00010239999999999998"} 3788
go_sched_latencies_seconds_bucket{le="0.0010485759999999998"} 3811
go_sched_latencies_seconds_bucket{le="0.010485759999999998"} 3811
go_sched_latencies_seconds_bucket{le="0.10485759999999998"} 3811
go_sched_latencies_seconds_bucket{le="+Inf"} 3811
go_sched_latencies_seconds_sum NaN
go_sched_latencies_seconds_count 3811
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 8
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 2.4
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 4096
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 12
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 2.6042368e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.6517349603e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.188737024e+09
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 3
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# HELP snoopy_processed_apicalls_total The total number of processed api calls
# TYPE snoopy_processed_apicalls_total counter
snoopy_processed_apicalls_total 2
# HELP snoopy_processed_blocks_total The total number of processed blocks
# TYPE snoopy_processed_blocks_total counter
snoopy_processed_blocks_total 123
# HELP snoopy_processed_events_total The total number of processed events
# TYPE snoopy_processed_events_total counter
snoopy_processed_events_total 124
# HELP snoopy_processed_requests_total The total number of processed requests
# TYPE snoopy_processed_requests_total counter
snoopy_processed_requests_total 5
# HELP snoopy_processed_transactions_total The total number of processed transactions
# TYPE snoopy_processed_transactions_total counter
snoopy_processed_transactions_total 22452
~~~
