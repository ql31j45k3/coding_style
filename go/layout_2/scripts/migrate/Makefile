MYSQL_URL="mysql://root:123456@tcp(127.0.0.1:3306)/demo_2?charset=utf8&parseTime=True&loc=Local"
MYSQL_SQL_PATH="migrations"

db-init:
	migrate create -ext sql -dir ${MYSQL_SQL_PATH} -seq init_schema

db-up:
	migrate --path ${MYSQL_SQL_PATH} -database ${MYSQL_URL} up

db-down:
	migrate --path ${MYSQL_SQL_PATH} -database ${MYSQL_URL} down
