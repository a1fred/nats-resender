# Nats-resender

[![CI](https://github.com/a1fred/nats-resender/actions/workflows/ci.yml/badge.svg)](https://github.com/a1fred/nats-resender/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/a1fred/nats-resender/badge.svg)](https://coveralls.io/github/a1fred/nats-resender)

Resend messages from one nats to another

# Usage
```sh
$ ./nats-resender --help
nats-resender version v1.0.29-3-g47b65aa-tls-20220111-16:31:16
Usage of ./nats-resender:
  -debug
        Debug
  -from-creds string
        nats source user credentials file
  -from-nkey string
        nats source nkey seed file
  -from-subj string
        nats source url (default "*")
  -from-tlscacrt string
        nats source tls ca cert file
  -from-tlscrt string
        nats source tls cert file
  -from-tlsinsecure
        nats source disable tls cert verification
  -from-tlskey string
        nats source tls key file
  -from-url string
        nats source url (default "nats://127.0.0.1:4222")
  -pendingByteLimit int
        pending subscription byte limit (default 67108864)
  -pendingMsgLimit int
        pending subscription message limit (default 524288)
  -print-period uint
        Print period seconds (default 10)
  -queue string
        queue, disabled by default
  -to-creds string
        nats source user credentials file
  -to-nkey string
        nats source nkey seed file
  -to-subj string
        nats source url (default "*")
  -to-tlscacrt string
        nats source tls ca cert file
  -to-tlscrt string
        nats source tls cert file
  -to-tlsinsecure
        nats source disable tls cert verification
  -to-tlskey string
        nats source tls key file
  -to-url string
        nats source url (default "nats://127.0.0.1:4222")
```

# Example
```shell
$ ./nats-resender --to-url=nats://127.0.0.1:4223
nats-resender version v1.0.29-3-g47b65aa-tls-20220111-16:35:20
Resending nats://127.0.0.1:4222/* -> nats://127.0.0.1:4223/*
2022/01/11 16:35:35 973972 processed, 10.00s elapsed, 97396.64 msg/sec, buffer:85539msgs/10948992bytes
2022/01/11 16:35:45 1870045 processed, 10.00s elapsed, 187003.72 msg/sec, buffer:209215msgs/26779520bytes
```
