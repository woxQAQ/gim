module github.com/woxQAQ/gim

go 1.21

replace (
	github.com/woxQAQ/gim/config => ./config
	github.com/woxQAQ/gim/docs => ./docs
	github.com/woxQAQ/gim/internal => ./internal
	github.com/woxQAQ/gim/internal/api => ./internal/api
	github.com/woxQAQ/gim/internal/api/friends => ./internal/api/friends
	github.com/woxQAQ/gim/internal/api/users => ./internal/api/users
	github.com/woxQAQ/gim/internal/db => ./internal/db
	github.com/woxQAQ/gim/internal/init => ./internal/init
	github.com/woxQAQ/gim/internal/middleware/jwt => ./internal/middleware/jwt
	github.com/woxQAQ/gim/internal/models => ./internal/models
	github.com/woxQAQ/gim/internal/server/gateway => ./internal/server/gateway
	github.com/woxQAQ/gim/internal/server/router => ./internal/server/router
	github.com/woxQAQ/gim/internal/server/transfer => ./internal/server/transfer
	github.com/woxQAQ/gim/pkg/server => ./pkg/server
	github.com/woxQAQ/gim/pkg/server/requests => ./pkg/server/requests
	github.com/woxQAQ/gim/pkg/util => ./pkg/util
)

require (
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.1.0
	github.com/panjf2000/gnet v1.6.7
	github.com/panjf2000/gnet/v2 v2.3.3
	github.com/swaggo/files v1.0.1
	github.com/swaggo/gin-swagger v1.6.0
	github.com/swaggo/swag v1.16.2
	github.com/vmihailenco/msgpack/v5 v5.4.1
	go.uber.org/zap v1.26.0
	gopkg.in/ini.v1 v1.67.0
	gorm.io/driver/mysql v1.5.2
	gorm.io/gorm v1.25.5
)

require (
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.29 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.22 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/bytedance/sonic v1.10.2 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/spec v0.20.9 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.16.0 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/panjf2000/ants/v2 v2.8.2 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.6.0 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.16.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
