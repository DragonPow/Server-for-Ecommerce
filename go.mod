module Server-for-Ecommerce

go 1.18

replace google.golang.org/grpc => google.golang.org/grpc v1.36.0

replace github.com/spf13/viper => github.com/spf13/viper v1.7.1

replace gorm.io/driver/postgres => gorm.io/driver/postgres v1.3.10

require (
	github.com/allegro/bigcache/v3 v3.1.0
	github.com/envoyproxy/protoc-gen-validate v0.6.2
	github.com/go-logr/logr v1.2.2
	github.com/go-logr/zapr v1.2.3
	github.com/go-redis/redis/v8 v8.11.5
	github.com/gogo/protobuf v1.3.2
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/jmoiron/sqlx v1.3.1
	github.com/lib/pq v1.10.2
	github.com/segmentio/kafka-go v0.4.38
	github.com/spf13/viper v1.13.0
	github.com/tabbed/pqtype v0.1.1
	github.com/urfave/cli/v2 v2.20.3
	go.uber.org/zap v1.19.0
	golang.org/x/exp v0.0.0-20230113213754-f9f960f08ad4
	google.golang.org/genproto v0.0.0-20220822174746-9e6da59bd2fc
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.1 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
