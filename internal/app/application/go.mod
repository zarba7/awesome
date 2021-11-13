module src

go 1.16

require (
	dubbo.apache.org/dubbo-go/v3 v3.0.0-rc2
	github.com/apache/rocketmq-client-go/v2 v2.1.0-rc3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/gorm v1.9.16
	github.com/json-iterator/go v1.1.10
	github.com/onsi/ginkgo v1.15.2 // indirect
	github.com/onsi/gomega v1.11.0 // indirect
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.38.0
	pbRole v0.0.0
)

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4
	github.com/golang/protobuf => github.com/golang/protobuf v1.3.2
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	model => ./../../../model
	pbRole => ./../../../proto/pbRole
)
