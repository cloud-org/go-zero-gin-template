<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [go-zero-gin-template](#go-zero-gin-template)
  - [test](#test)
  - [docker](#docker)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

### go-zero-gin-template

一个 go-zero gin web 模板

#### test

```
[GIN-debug] GET    /v1/hello                 --> go-zero-gin-template/api/internal/handler.(*Handler).ReflectHandler.func1 (4 handlers)
{"@timestamp":"2022-06-08T15:20:26.759+08:00","caller":"app/app.go:73","content":"server is start at: 0.0.0.0:9097","level":"info"}
```

```
$ curl http://127.0.0.1:9097/v1/hello\?name\=ashing                  
{"errMsg":"","errCode":0,"data":"hello, ashing"}% 
```

#### docker

```sh
make docker-build-api version=v0.0.1
```

```sh
make docker-run-api version=v0.0.1
```

