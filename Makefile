include app.env

ifeq (${NODE_ENV}, DEV)
	DB_NAME=${DEV_POSTGRES_DB}
	DB_USER=${DEV_POSTGRES_USER}
	DB_PASSWORD=${DEV_POSTGRES_PASSWORD}
	DB_HOST=${DEV_POSTGRES_HOST}
	DB_PORT=${DEV_POSTGRES_PORT}
	SSL_MODE=${DEV_SSL_MODE}
endif

ifeq (${NODE_ENV}, STAGE)
	DB_NAME=${STAGE_POSTGRES_DB}
	DB_USER=${STAGE_POSTGRES_USER}
	DB_PASSWORD=${STAGE_POSTGRES_PASSWORD}
	DB_HOST=${STAGE_POSTGRES_HOST}
	DB_PORT=${STAGE_POSTGRES_PORT}
	SSL_MODE=${STAGE_SSL_MODE}
endif

ifeq (${NODE_ENV}, PROD)
	DB_NAME=${PROD_POSTGRES_DB}
	DB_USER=${PROD_POSTGRES_USER}
	DB_PASSWORD=${PROD_POSTGRES_PASSWORD}
	DB_HOST=${PROD_POSTGRES_HOST}
	DB_PORT=${PROD_POSTGRES_PORT}
	SSL_MODE=${PROD_SSL_MODE}
endif


migrate_init:
	migrate create -ext sql -dir ./db/schema -seq init_schema

migrate_create:
	migrate create -ext sql -dir ./db/schema -seq ${MIGRATE_NAME}_schema

migration_up:
	migrate -path ./db/schema -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${SSL_MODE}" -verbose up

migration_down_all:
	migrate -path ./db/schema -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${SSL_MODE}" -verbose down

migration_down_by_id:
	migrate -path ./db/schema -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${SSL_MODE}" -verbose down ${VERSION}


migration_fix:
	migrate -path ./db/schema -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${SSL_MODE}" force ${VERSION}

mockgen:
	mockgen -source=D:\development\backend\callboard_microservices\users_mrc\internal\usecases\usecases_interfaces.go -destination=D:\development\backend\callboard_microservices\users_mrc\internal\db\mocks\mock_usecases.go -package=mocks

sqlc:
	sqlc compile
	sqlc generate