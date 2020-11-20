# hlc
Hybrid Logical Clock

Original Paper: [Logical Physical Clocksand Consistent Snapshots in Globally Distributed Databases](https://cse.buffalo.edu/tech-reports/2014-04.pdf)

Logical Clock is an important family of algorithms for ordering events into distributed system.

HLC presents enhanced Logic Clock's Lamport, which eliminates the gap between practical and theory of synchronised time.

# How use it

## Init HLC per Node

```go
lc := hlc.New(&NTPClock{})
```

## For every local event or sending, generate logical timestamp

```go
ts := lc.Now() // 1605806872706744321
```

**Notes**: *Timestamp compatible with NTP-format. First 48bit presents physical time and
other 16bit logical time (see chapter "6.2 Compact Timestamping using l and c" of paper).* 

## For every received event you should to update a local HLC instance

```go
receivedHlc := hlc.FromTimestamp(receivedTC)
_ = lc.Update(receivedHlc) // skip a new timestamp
```

# Useful links

1. [In-depth Analysis on HLC-based Distributed Transaction Processing](https://alibaba-cloud.medium.com/in-depth-analysis-on-hlc-based-distributed-transaction-processing-e75dad5f2af8)
2. [Hybrid Logical Clocks F# implementation](https://bartoszsypytkowski.com/hybrid-logical-clocks/)
