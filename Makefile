createmigration:
	./migrate create -ext=sql -dir=migrations/init -seq init

migrateup:
	./migrate -path=migrations/init -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up

migratedown:
	./migrate -path=migrations/init -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose down

# .PHONY: migrate migratedown createmigration