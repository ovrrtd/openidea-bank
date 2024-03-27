# to run server
.PHONY: runServerLinux
runServerLinux:
	cd cmd && GOOS=linux GOARCH=amd64 go build -o main && ./main

.PHONY: runServerMac
runServerMac:
	cd cmd && go build -o main && ./main

.PHONY: migrateUp
migrateUp:
	migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)"  -path db/migrations up

# to run rollback migration
.PHONY: migrateDown
migrateDown:
	migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)"  -path db/migrations down
	

.PHONY: buildProd
buildProd:
	cd cmd && GOOS=linux GOARCH=amd64 go build -o main_nu

.PHONY: scpProd
scpProd:
	 scp -i ./Project-Sprint-Key.pem ./cmd/main ubuntu@54.255.241.225:~/

.PHONY: migrateUpProd
migrateUpProd:
	migrate -path ./db/migrations -database "postgres://postgres:quiuxi9paeGh2EiS2Kiesh9euh2Ing4je@projectsprint-db.cavsdeuj9ixh.ap-southeast-1.rds.amazonaws.com:5432/postgres?sslmode=verify-full&sslrootcert=ap-southeast-1-bundle.pem" up

.PHONY: migrateDownProd
migrateDownProd:
	migrate -path ./db/migrations -database "postgres://postgres:quiuxi9paeGh2EiS2Kiesh9euh2Ing4je@projectsprint-db.cavsdeuj9ixh.ap-southeast-1.rds.amazonaws.com:5432/postgres?sslmode=verify-full&sslrootcert=ap-southeast-1-bundle.pem" drop
