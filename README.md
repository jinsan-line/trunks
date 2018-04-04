# Trunks [![Build Status](https://travis-ci.org/straightdave/trunks.svg?branch=master)](https://travis-ci.org/straightdave/trunks)

Trunks, like every son, is derived from the father Vegeta with some enhanced skills. Same as Vegeta, Trunks is yet another HTTP load testing tool built out of a need to drill
HTTP services with a constant request rate.
It can be used both as a command line utility and a library.

![Trunks](http://images2.wikia.nocookie.net/__cb20100725123520/dragonballfanon/images/5/52/Future_Trunks_SSJ2.jpg)

## Usage manual

for original usage of Vegeta, please refer to [vegeta' readme](https://github.com/tsenart/vegeta/blob/master/README.md)

## More functionalities

### dump http response to file
```console
(add one more option '-respf')
-respf string
      Dump responses to file
```

### gRPC perf test (as a lib)
Sample file:
```go
package main

import (
  "fmt"
  "log"
  "time"

  trunks "github.com/straightdave/trunks/lib"
)

func main() {
  log.Println("hello")
  tgt := &trunks.GTargeter{
    Target:     ":8087",
    IsEtcd:     false,
    MethodName: "/myapppb.MyApp/Hello",
    Request:    &HelloRequest{Name: "dave"},
    Response:   &HelloResponse{},
  }

  b, err := tgt.GenBurner()
  if err != nil {
    fmt.Println(err)
    return
  }
  defer b.Conn.Close()

  var metrics trunks.Metrics

  startT := time.Now()
  for res := range b.Burn(tgt, uint64(5), 10*time.Second) {
    metrics.Add(res)
  }
  dur := time.Since(startT)

  metrics.Close()

  fmt.Printf("dur: %v\n", dur.Seconds())
  fmt.Printf("earliest: %v\n", metrics.Earliest.Sub(startT).Nanoseconds())
  fmt.Printf("latest: %v\n", metrics.Latest.Sub(startT).Nanoseconds())
  fmt.Printf("end: %v\n", metrics.End.Sub(startT).Nanoseconds())
  fmt.Printf("reqs: %d\n", metrics.Requests)
  fmt.Printf("50th: %s\n", metrics.Latencies.P50)
  fmt.Printf("95th: %s\n", metrics.Latencies.P95)
  fmt.Printf("99th: %s\n", metrics.Latencies.P99)
  fmt.Printf("mean: %s\n", metrics.Latencies.Mean)
  fmt.Printf("max: %s\n", metrics.Latencies.Max)
}
```

For details please read the code (currently smell but promised to be better)
