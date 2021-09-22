# Nats-resender

[![CI](https://github.com/a1fred/nats-resender/actions/workflows/ci.yml/badge.svg)](https://github.com/a1fred/nats-resender/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/a1fred/nats-resender/badge.svg)](https://coveralls.io/github/a1fred/nats-resender)

Resend messages from one nats to another

# Usage
```sh
$ ./build/nats-resender --help
nats-resender version 9410787-master-20210922-00:29:55
Usage of ./build/nats-resender:
  -debug
        Debug

  -from-url string
        nats source url (default "nats://127.0.0.1:4222")
  -from-subj string
        nats source subject (default "*")

  -to-url string
        nats destination url (default "nats://127.0.0.1:4222")
  -to-subj string
        nats destination subject (default "nats-resender")

  -queue string
        queue, disabled by default
  -pendingMsgLimit int
        pending subscription message limit (default 524288)
  -pendingByteLimit int
        pending subscription byte limit (default 67108864)

  -print-period uint
        Print period seconds (default 10)
```

# Example
```shell
$ ./build/nats-resender
nats-resender version 9410787-master-20210921-18:26:04
2021/09/21 18:26:16 0 messages processed, elapsed 10.00s, 0.00 msg/sec
2021/09/21 18:26:26 114043 messages processed, elapsed 10.00s, 11404.28 msg/sec
2021/09/21 18:26:36 228342 messages processed, elapsed 10.00s, 22834.15 msg/sec
2021/09/21 18:26:46 228038 messages processed, elapsed 10.00s, 22803.75 msg/sec
2021/09/21 18:26:56 204887 messages processed, elapsed 10.00s, 20488.60 msg/sec
```


# TODO
 * fix goveralls
 * better test coverage
