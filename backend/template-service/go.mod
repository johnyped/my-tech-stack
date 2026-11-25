// TODO: replace with the real module path.
// gitlab.awesome-poc-th.com/libero-engineering/<group>/<service-name>
// <group> = "core" for platform-wide services, "auxiliary-systems/<system>"
// for a system-specific one. See document/TECH_STACK.md "Module naming".
module gitlab.awesome-poc-th.com/libero-engineering/TODO-group/TODO-service-name

go 1.26.0

require (
	github.com/gin-gonic/gin v1.12.0
	github.com/go-playground/validator/v10 v10.30.1
	github.com/google/uuid v1.6.0
	github.com/spf13/viper v1.21.0
	github.com/stretchr/testify v1.12.1
	gitlab.awesome-poc-th.com/libero-engineering/core/common-lib.git/v2 v2.3.1
)

// Added by `script/new-service.sh` based on your infra checklist answers —
// see document/SCAFFOLD_CHECKLIST.md. Uncomment / `go get` the ones you need
// if you're hand-editing this file instead of using the generator.
//
// >>> CACHE
// require github.com/redis/go-redis/v9 v9.17.1 // --cache redis
// <<< CACHE
// >>> DB
// require github.com/jackc/pgx/v5 v5.11.0 // --db postgres
// <<< DB
// >>> KAFKA
// require github.com/segmentio/kafka-go v0.4.50 // --kafka yes
// <<< KAFKA
// >>> STORAGE
// require github.com/aws/aws-sdk-go-v2/service/s3 v1.96.0 // --storage s3
// <<< STORAGE
// >>> FIREBASE
// require firebase.google.com/go/v4 v4.23.0 // --firebase yes
// <<< FIREBASE
