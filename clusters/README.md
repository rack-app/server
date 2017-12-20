# Target Goal
```
$ wrk -t12 -c400 -d30s http://localhost:9292
Running 30s test @ http://localhost:9292
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.33ms    6.65ms 180.74ms   93.02%
    Req/Sec    16.15k    11.46k   58.06k    53.57%
  2871948 requests in 30.07s, 109.56MB read
Requests/sec:  95517.29
Transfer/sec:      3.64MB
```