[![baby-gopher](http://www.babygopher.org/images/babygopher-badge.png)](http://www.babygopher.org)

RTMFPew
========

RTMFP transport implementation in Go.

### Status
What's TBD
 - [x] Data chunks processing & tests
 - [ ] Session handling
 - [ ] NetGroup & NetStream API
 - [ ] Data transmission & tests
 - [ ] RFC7016 compilant tests
 - [ ] Tesing in live flash client
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
