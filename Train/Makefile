postgres: #Запускает контейнер
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15

createdb: #Создаёт базу данных
	docker exec -it postgres15 createdb --username=root --owner=root airline

dropdb: #Удаляет базу данных
	docker exec -it postgres15 dropdb airline

migrateup: #Загружает базу данных из файла
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/airline?sslmode=disable" -verbose up

test:
	docker exec -it postgres15 psql -U root -d airline
migratedown: #Удаляет базу данных из файла
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/airline?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup test migratedown
