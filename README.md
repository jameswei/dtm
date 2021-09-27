![license](https://img.shields.io/github/license/yedf/dtm)
[![Build Status](https://travis-ci.com/yedf/dtm.svg?branch=main)](https://travis-ci.com/yedf/dtm)
[![Coverage Status](https://coveralls.io/repos/github/yedf/dtm/badge.svg?branch=main)](https://coveralls.io/github/yedf/dtm?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/yedf/dtm)](https://goreportcard.com/report/github.com/yedf/dtm)
[![Go Reference](https://pkg.go.dev/badge/github.com/yedf/dtm.svg)](https://pkg.go.dev/github.com/yedf/dtm)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/avelino/awesome-go#database)

# [中文文档](http://dtm.pub)

# A Cross Language Transaction Manager

## Who's using DTM

[Ivydad](https://ivydad.com)

[Eglass](https://epeijing.cn)

[Jiou](http://jiou.me)

[GoldenData]()

## What is DTM

DTM is the first distributed transaction management framework in Golang. Unlike other frameworks, DTM provides extremely easy access interfaces of HTTP and gRPC, supports multiple language bindings, and handles tricky problems of unordered sub-transactions at the framework level.

## Features

* Extremely easy to adapt
  - Support HTTP and gRPC, provide easy-to-use programming interfaces, lower substantially the barrier of getting started with distributed transactions. Newcomers can adapt quickly.

* Easy to use
  - Relieving developers from worrying about suspension, null compensation, idempotent transaction, and other tricky problems, the framework layer handles them all.

* Language-agnostic
  - Suit for companies with multiple-language stacks.
    Easy to write bindings for Go, Python, PHP, Node.js, Ruby, and other languages.

* Easy to deploy, easy to extend
  - DTM depends only on MySQL, easy to deploy, cluster, and scale horizontally.

* Support for multiple distributed transaction protocol
  - TCC, SAGA, XA, Transactional messages.

## DTM vs. others

There is no mature open-source distributed transaction framework for non-Java languages.
Mature open-source distributed transaction frameworks for Java language include Ali's Seata, Huawei's ServiceComb-Pack, Jingdong's shardingsphere, himly, tcc-transaction, ByteTCC, and so on, of which Seata is most widely used.

The following is a comparison of the main features of dtm and Seata.


| Features                | DTM                                                                                           | Seata                                                                                            | Remarks                                                             |
| :-----:                 | :----:                                                                                        | :----:                                                                                           | :----:                                                              |
| Supported languages     | <span style="color:green">Golang, Python, PHP,  and others</span>                               | <span style="color:orange">Java</span>                                                           | dtm allows easy access from a new language                            |
| Exception handling      | [Sub-transaction barrier](https://zhuanlan.zhihu.com/p/388444465)                             | <span style="color:orange">manual</span>                                                         | dtm solves idempotent transaction, hanging, null compensation                   |
| TCC                     | <span style="color:green">✓</span>                                                            | <span style="color:green">✓</span>                                                               |                                                                     |
| XA                      | <span style="color:green">✓</span>                                                            | <span style="color:green">✓</span>                                                               |                                                                     |
| AT                      | <span style="color:red">✗</span>                                                              | <span style="color:green">✓</span>                                                               | AT is similar to XA with better performance but with dirty rollback |
| SAGA                    | <span style="color:orange">Simple mode</span>                                                 | <span style="color:green">complicated state-machine mode</span>                                   | dtm's state-machine mode is being planned                         |
| Transactional Messaging | <span style="color:green">✓</span>                                                            | <span style="color:red">✗</span>                                                                 | dtm provides Transactional Messaging similar to RocketMQ               |
| Communication protocols | <span style="color:green">HTTP, gRPC</span>                                                   | <span style="color:green">Dubbo, no HTTP</span>                                             |                                                                     |
| Star count              | <img src="https://img.shields.io/github/stars/yedf/dtm.svg?style=social" alt="github stars"/> | <img src="https://img.shields.io/github/stars/seata/seata.svg?style=social" alt="github stars"/> | dtm 0.1 is released from 20210604 and under fast development                    |

From the features' comparison above, if your language stack includes languages other than Java, then dtm is the one for you.
If your language stack is Java, you can also choose to access dtm and use sub-transaction barrier technology to simplify your business development.

## [Cook Book](https://en.dtm.pub)

# Quick start

### Install

`git clone https://github.com/yedf/dtm`

### Configure Mysql

`cp conf.sample.yml conf.yml # Modify conf.yml`

### Start the example
`go run app/main.go`

# Code

### Use
``` go
  // business micro-service address
  const qsBusi = "http://localhost:8081/api/busi_saga"
  // The address where DtmServer serves DTM, which is a url
  DtmServer := "http://localhost:8080/api/dtmsvr"
  req := &gin.H{"amount": 30} // micro-service payload
	// DtmServer is the address of DTM micro-service
	saga := dtmcli.NewSaga(DtmServer, dtmcli.MustGenGid(DtmServer)).
		// add a TransOut subtraction，forward operation with url: qsBusi+"/TransOut", reverse compensation operation with url: qsBusi+"/TransOutCompensate"
		Add(qsBusi+"/TransOut", qsBusi+"/TransOutCompensate", req).
		// add a TransIn subtraction, forward operation with url: qsBusi+"/TransIn", reverse compensation operation with url: qsBusi+"/TransInCompensate"
		Add(qsBusi+"/TransIn", qsBusi+"/TransInCompensate", req)
	// submit the created saga transaction，dtm ensures all subtractions either complete or get revoked
	err := saga.Submit()
```
### Complete example

Refer to [examples/quick_start.go](./examples/quick_start.go).

### Slack

You can join the [DTM slack channel here](https://join.slack.com/t/dtm-w6k9662/shared_invite/zt-vkrph4k1-eFqEFnMkbmlXqfUo5GWHWw).

### Wechat

Add wechat friend with id yedf2008, or scan the OR code. Fill in dtm as verification.

![yedf2008](http://service.ivydad.com/cover/dubbingb6b5e2c0-2d2a-cd59-f7c5-c6b90aceb6f1.jpeg)

### Give a star! ⭐

If you think this project is good, or helpful to you, please give a star!

### Who is using
<div style='vertical-align: middle'>
    <img alt='Ivydad' height='40'  src='https://www.ivydad.com/_nuxt/img/header-logo.5b3eb96.png'>
    <img alt='Eglass' height='40'  src='https://img.epeijing.cn/official-website/assets/logo.png'>
    <img alt='Jiou' height='40'  src='http://www.siqitech.com.cn/img/logo.3f6c2914.png'>
    <img alt='GoldenData' height='40'  src='https://pic1.zhimg.com/80/v2-dc1d0cef5f7b72be345fc34d768e69e3_1440w.png'>
</div>
