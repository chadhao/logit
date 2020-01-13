# logit

## Troubleshoots

Error: cannot unmarshal DNS messsage
Reason: dns nameserver can not resolve the request
Solution: change default dns server from local to public such as 8.8.8.8 
```
panic: error parsing uri: lookup logit-dev-cluster-h9szy.gcp.mongodb.net on 127.0.0.53:53: cannot unmarshal DNS message
        panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x2d8 pc=0xbf2fca]

goroutine 1 [running]:
go.mongodb.org/mongo-driver/mongo.(*Client).endSessions(0x0, 0x10e52a0, 0xc0003441c0)
        /home/shien/go/pkg/mod/go.mongodb.org/mongo-driver@v1.2.0/mongo/client.go:286 +0x3a
go.mongodb.org/mongo-driver/mongo.(*Client).Disconnect(0x0, 0x10e52a0, 0xc0003441c0, 0xc0003441c0, 0xc0003b8150)
        /home/shien/go/pkg/mod/go.mongodb.org/mongo-driver@v1.2.0/mongo/client.go:190 +0x56
github.com/chadhao/logit/modules/user/model.Close()
        /home/shien/Development/logit/modules/user/model/client.go:47 +0x67
github.com/chadhao/logit/modules/user.ShutdownModule()
        /home/shien/Development/logit/modules/user/main.go:20 +0x20
main.shutdownModules()
        /home/shien/Development/logit/module.go:44 +0x49
panic(0xd911e0, 0xc0003b8140)
        /usr/local/go/src/runtime/panic.go:679 +0x1b2
main.main()
        /home/shien/Development/logit/main.go:27 +0x230
```