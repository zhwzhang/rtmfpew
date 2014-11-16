## RTMFPew
##### RTMFP P2P transport implementation in Go
-----------------------------------------------
[![bountysource](https://www.bountysource.com/badge/team?team_id=47410)](https://www.bountysource.com/teams/rtmfpew)
[![gratipay](https://img.shields.io/gratipay/VoidNugget.svg)](https://gratipay.com/VoidNugget)

### Status

##### It's NOT working yet

#### What's TBD
 - [x] Data chunks processing & tests
 - [x] Session handling
 - [ ] Amf0 & Amf3 serialization with [Amfy](https://github.com/rtmfpew/amfy)
 - [ ] NetGroup & NetStream API
 - [ ] Data transmission tests
 - [ ] RFC7016 compliant tests
 - [ ] Echo testing with live flash client
 - [ ] Profiling & load-testing
 - [ ] Fuzzing

So in the end it should be very neat.

### Installation
Simple as
``` go get github.com/rtmfpew/rtmfpew ```

### Testing
```
./dependencies.sh
goconvey
```

Feel free to use [goop](https://github.com/nitrous-io/goop) as isolated environment provider.

### License
RTFMPew licensed under the Apache License 2.0
