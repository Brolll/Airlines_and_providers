package main

import (
	_ "database/sql/driver"
	"fmt"
	_ "github.com/golang-migrate/migrate"
	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

type Provider struct {
	Id   string `db:"id"`
	Name string `db:"name"`
	Code int    `db:"code"`
}

type SSchema struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

type Airline struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}

type ProviderAirline struct {
	Provider_code int    `db:"provider_code"`
	Airline_id    string `db:"airline_id"`
}

type ProviderSchema struct {
	Schema_id   string `db:"schema_id"`
	Provider_id int    `db:"provider_id"`
}

type Account struct {
	Id       int    `db:"id"`
	SchemaId int    `db:"SchemaId"`
	Name     string `db:"name"`
}

type Airpv struct {
	Code     int
	Id       string
	Name     string
	Airlines string
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<a href="http://localhost:8080/1">1.   Добавить авиакомпанию</a><p><p>
				<a href="http://localhost:8080/2">2.   Удалить авиакомпанию по id</a><p><p>
				<a href="http://localhost:8080/3">3.   Добавить поставщика</a><p><p>
				<a href="http://localhost:8080/4">4.   Удалить поставщика по Id</a><p><p>
				<a href="http://localhost:8080/5">5.   Изменить список поставщиков</a><p><p>
				<a href="http://localhost:8080/6">6.   Добавить схему</a><p><p>
				<a href="http://localhost:8080/7">7.   Искать схему по названию</a><p><p>
				<a href="http://localhost:8080/8">8.   Изменить схему</a><p><p>
				<a href="http://localhost:8080/9">9.   Удалить схему</a><p><p>
				<a href="http://localhost:8080/10">10. Добавить аккаунт</a><p><p>
				<a href="http://localhost:8080/11">11. Изменение схемы аккаунта</a><p><p>
				<a href="http://localhost:8080/12">12. Удалить аккаунт</a><p><p>
				<a href="http://localhost:8080/13">13. Получить список авиакомпаний по Id аккаунта</a><p><p>
				<a href="http://localhost:8080/14">14. Получить список авиакомпаний по Id поставщика</a><p><p>
			</body>
			</html>
		`)
}

func handlerAddAirline(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>1 - Добавление авиакомпании<h1>
				<form action="/1" method="post">
					<label for="idInput">Введите Id:</label>
					<input type="text" id="idInput" name="id">
					<label for="nameInput">Введите Name:</label>
					<input type="text" id="nameInput" name="name">
					<label for="nameInput">Введите Provider_code:</label>
					<input type="text" id="provider_codeInput" name="provider_code">
					<button type="submit">Отправить</button>
				</form>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
		`)
		airlines := []Airline{}
		err := db.Select(&airlines, `SELECT * FROM "Airline"`)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, "<p>Авиакомпании:<p>")
		for _, airline := range airlines {
			fmt.Fprintf(w, "<p>id: %s, name: %s\n<p>", airline.Id, airline.Name)
		}
		fmt.Fprintln(w, "Список поставщиков и компаний которыми они владеют:")
		airpvs := []Airpv{}
		err = db.Select(&airpvs, `
		SELECT
		     "Provider"."code",
			 "Provider"."id",
			 "Provider"."name",
			 ARRAY_AGG("ProviderAirline"."airline_id") AS "airlines"
		FROM
			 "Provider"
			 left JOIN "ProviderAirline" ON "Provider"."code" = "ProviderAirline"."provider_code"
		GROUP BY
			 "Provider"."id",
			 "Provider"."name"
			 `)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		for _, airpv := range airpvs {
			fmt.Fprintf(w, "<p>code: %d id: %s, name: %s, airlines: %v\n<p>", airpv.Code, airpv.Id, airpv.Name, airpv.Airlines)
		}
	} else if r.Method == "POST" {
		id, name := r.FormValue("id"), r.FormValue("name")
		provider_code, err := strconv.Atoi(r.FormValue("provider_code"))
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		if id == "" || name == "" || provider_code == 0 {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`INSERT INTO "Airline" ("id", "name") VALUES ($1,$2)`, id, name)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`INSERT INTO "ProviderAirline" (provider_code, airline_id) VALUES ($1,$2)`, provider_code, id)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
			<body>
				<h1>Успешно</h1>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
			`)
		return
	}
}

func handlerDeleteAirline(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>2 - Удаление авиакомпании<h1>
				<form action="/2" method="post">
					<label for="idInput">Введите Id:</label>
					<input type="text" id="idInput" name="id">
					<button type="submit">Отправить</button>
				</form>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
		`)
		airlines := []Airline{}
		err := db.Select(&airlines, `SELECT * FROM "Airline"`)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, "<p>Авиакомпании:<p>")
		for _, airline := range airlines {
			fmt.Fprintf(w, "<p>id: %s, name: %s\n<p>", airline.Id, airline.Name)
		}
	} else if r.Method == "POST" {
		id := r.FormValue("id")
		if id == "" {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		_, err := db.Exec(`DELETE FROM "ProviderAirline" WHERE airline_id = $1`, id)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`DELETE FROM "Airline" WHERE id = $1`, id)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
			<body>
				<h1>Успешно</h1>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
			`)
		return
	}
}

func handlerAddProvider(w http.ResponseWriter, r *http.Request, db *sqlx.DB) { //Возможно надо доработать
	if r.Method == "GET" {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>1 - Добавление поставщика<h1>
				<form action="/3" method="post">
					<label for="idInput">Введите Id:</label>
					<input type="text" id="idInput" name="id">
					<label for="nameInput">Введите Name:</label>
					<input type="text" id="nameInput" name="name">
					<button type="submit">Отправить</button>
				</form>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
		`)
	} else if r.Method == "POST" {
		id, name := r.FormValue("id"), r.FormValue("name")
		if id == "" || name == "" {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		_, err := db.Exec(`INSERT INTO "Provider" ("id", "name") VALUES ($1,$2)`, id, name)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
			<body>
				<h1>Успешно</h1>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
			`)
		return
	}
}

func handlerDeleteProvider(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>4 - Удаление поставщика по Id<h1>
				<form action="/4" method="post">
					<label for="idInput">Введите Id:</label>
					<input type="text" id="idInput" name="id">
					<button type="submit">Отправить</button>
				</form>
				<a href="http://localhost:8080">Вернутся на главную</a>
				<h2><p>Список поставщиков и компаний которыми они владеют:<p><h2>
			</body>
			</html>
		`)
		airpvs := []Airpv{}
		err := db.Select(&airpvs, `
		SELECT
		     "Provider"."code",
			 "Provider"."id",
			 "Provider"."name",
			 ARRAY_AGG("ProviderAirline"."airline_id") AS "airlines"
		FROM
			 "Provider"
			 left JOIN "ProviderAirline" ON "Provider"."code" = "ProviderAirline"."provider_code"
		GROUP BY
			 "Provider"."id",
			 "Provider"."name"
			 `)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		for _, airpv := range airpvs {
			fmt.Fprintf(w, "<p>code: %d id: %s, name: %s, airlines: %v\n<p>", airpv.Code, airpv.Id, airpv.Name, airpv.Airlines)
		}
	} else if r.Method == "POST" {
		id := r.FormValue("id")
		if id == "" {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		var prv_code int
		var prv_codeStr string
		err := db.Get(&prv_codeStr, `SELECT "Provider"."code" FROM "Provider" WHERE "id" = $1`, id)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		prv_code, err = strconv.Atoi(prv_codeStr)
		_, err = db.Exec(`DELETE FROM "ProviderAirline" WHERE "provider_code" = $1`, prv_code)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`DELETE FROM "ProviderSchema" WHERE "provider_id" = $1`, id)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`DELETE FROM "Provider" WHERE "id" = ($1)`, id)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
			<body>
				<h1>Успешно</h1>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
			`)
		return
	}
}

func handlerModifyProvider(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>5 - Изменения списка поставщиков <p>Введите код поставщика и Id авикомпании которую хотите убрать<p><h1>
				<form action="/5" method="post">
					<label for="idInput">Введите Id:</label>
					<input type="text" id="idInput" name="id">
					<label for="nameInput">Введите Provider_code:</label>
					<input type="text" id="provider_codeInput" name="provider_code">
					<button type="submit">Отправить</button>
				</form>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
		`)
		fmt.Fprintln(w, "<p>Список поставщиков и компаний которыми они владеют:<p>")
		airpvs := []Airpv{}
		err := db.Select(&airpvs, `
		SELECT
		     "Provider"."code",
			 "Provider"."id",
			 "Provider"."name",
			 ARRAY_AGG("ProviderAirline"."airline_id") AS "airlines"
		FROM
			 "Provider"
			 left JOIN "ProviderAirline" ON "Provider"."code" = "ProviderAirline"."provider_code"
		GROUP BY
			 "Provider"."id",
			 "Provider"."name"
			 `)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		for _, airpv := range airpvs {
			fmt.Fprintf(w, "<p>code: %d id: %s, name: %s, airlines: %v\n<p>", airpv.Code, airpv.Id, airpv.Name, airpv.Airlines)
		}
	} else if r.Method == "POST" {
		id := r.FormValue("id")
		provider_code, err := strconv.Atoi(r.FormValue("provider_code"))
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		if id == "" || provider_code == 0 {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`DELETE FROM "ProviderAirline" WHERE provider_code = $1 AND airline_id = $2`, provider_code, id)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
			<body>
				<h1>Успешно</h1>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
			`)
		return
	}
}

func handlerAddSchema(w http.ResponseWriter, r *http.Request, db *sqlx.DB) { //Возможно надо доработать
	if r.Method == "GET" {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>6 - Добавление схемы<h1>
				<form action="/6" method="post">
					<label for="nameInput">Введите Name:</label>
					<input type="text" id="nameInput" name="name">
					<button type="submit">Отправить</button>
				</form>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
		`)
		schemas := []SSchema{}
		err := db.Select(&schemas, `SELECT * FROM "SSchema"`)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, "<p>Схемы:<p>")
		for _, schema := range schemas {
			fmt.Fprintf(w, "<p>id: %d, name: %s\n<p>", schema.Id, schema.Name)
		}
	} else if r.Method == "POST" {
		name := r.FormValue("name")
		if name == "" {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		_, err := db.Exec(`INSERT INTO "SSchema" ("name") VALUES ($1)`, name)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
			<body>
				<h1>Успешно</h1>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
			`)
		return
	}
}

func handlerSchemaSearch(w http.ResponseWriter, r *http.Request, db *sqlx.DB) { //Вывод схемы не массивом, а строкой
	if r.Method == "GET" {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>7 - Поиск схемы по названию?<h1>
				<form action="/7" method="post">
					<label for="nameInput">Введите Name:</label>
					<input type="text" id="nameInput" name="name">
					<button type="submit">Отправить</button>
				</form>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
		`)
	} else if r.Method == "POST" {
		name := r.FormValue("name")
		if name == "" {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		var sch SSchema
		err := db.Get(&sch, `SELECT * FROM "SSchema" WHERE name = $1`, name)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
			<body>
				<h1>Вот ваша схема %v</h1>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
			`, sch)
		return
	}
}

func handlerModifySchema(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>8 - Изменение схемы<h1>
				<form action="/8" method="post">
					<label for="nameInput">Введите Id:</label>
					<input type="text" id="nameInput" name="id">
					<label for="nameInput">Введите новое имя для схемы:</label>
					<input type="text" id="nameInput" name="name">
					<button type="submit">Отправить</button>
				</form>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
		`)
		var sch []SSchema
		err := db.Select(&sch, `SELECT * FROM "SSchema"`)
		for _, schema := range sch {
			fmt.Fprintf(w, "<p>id: %d name: %s <p>", schema.Id, schema.Name)
		}
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
	} else if r.Method == "POST" {
		name := r.FormValue("name")
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`UPDATE "SSchema" SET "name" = $1 WHERE id = $2`, name, id)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
			<body>
				<h1>Успешно</h1>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
			`)
		return
	}
}

func handlerDeleteSchema(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>9 - Удаление схемы<h1>
				<form action="/9" method="post">
					<label for="nameInput">Введите Id:</label>
					<input type="text" id="idInput" name="id">
					<button type="submit">Отправить</button>
				</form>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
		`)
		schemas := []SSchema{}
		err := db.Select(&schemas, `SELECT * FROM "SSchema"`)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, "<p>Схемы:<p>")
		for _, schema := range schemas {
			fmt.Fprintf(w, "<p>id: %d, name: %s\n<p>", schema.Id, schema.Name)
		}
	} else if r.Method == "POST" {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		if id == 0 {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`DELETE FROM "SSchema"
								WHERE id = $1
  								AND NOT EXISTS (
    							SELECT 1 FROM "Account" WHERE "SchemaId" = $1)`, id)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
			<body>
				<h1>Успешно</h1>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
			`)
		return
	}
}

func handlerAccountAdd(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>10 - Добавление аккаунта<h1>
				<form action="/10" method="post">
					<label for="nameInput">Введите SchemaId:</label>
					<input type="text" id="idInput" name="schemaId">
					<label for="nameInput">Введите Name:</label>
					<input type="text" id="idInput" name="name">
					<button type="submit">Отправить</button>
				</form>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
		`)
		schemas := []SSchema{}
		err := db.Select(&schemas, "SELECT * FROM \"SSchema\"")
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, "<p>Схемы:<p>")
		for _, schema := range schemas {
			fmt.Fprintf(w, "<p>id: %d, name: %s\n<p>", schema.Id, schema.Name)
		}
		fmt.Fprintln(w, "<p>Аккаунты:<p>")
		accounts := []Account{}
		err = db.Select(&accounts, `SELECT * FROM "Account"`)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		for _, account := range accounts {
			fmt.Fprintf(w, "<p>id: %d Schemaid: %d name: %s<p>", account.Id, account.SchemaId, account.Name)
		}
	} else if r.Method == "POST" {
		name := r.FormValue("name")
		id, err := strconv.Atoi(r.FormValue("schemaId"))
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		if id == 0 {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`INSERT INTO "Account"("SchemaId", "name")  values ($1,$2)`, id, name)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
			<body>
				<h1>Успешно</h1>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
			`)
		return
	}
}

func handlerModifyAccsSchema(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>11 - Изменение схемы у аккаунта<h1>
				<h2><p>Введите Id аккаунта у которого хотите поменять схему и Id схемы на которую поменять<p><h2>
				<form action="/11" method="post">
					<label for="nameInput">Введите Id аккаунта:</label>
					<input type="text" id="nameInput" name="idacc">
					<label for="nameInput">Введите Id схемы:</label>
					<input type="text" id="nameInput" name="idsch">
					<button type="submit">Отправить</button>
				</form>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
		`)
		var sch []SSchema
		err := db.Select(&sch, `SELECT * FROM "SSchema"`)
		fmt.Fprintln(w, "<p>Схемы:<p>")
		for _, schema := range sch {
			fmt.Fprintf(w, "<p>id: %d name: %s <p>", schema.Id, schema.Name)
		}
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, "<p>Аккаунты:<p>")
		accounts := []Account{}
		err = db.Select(&accounts, `SELECT * FROM "Account"`)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		for _, account := range accounts {
			fmt.Fprintf(w, "<p>id: %d Schemaid: %d name: %s<p>", account.Id, account.SchemaId, account.Name)
		}
	} else if r.Method == "POST" {
		idacc, err := strconv.Atoi(r.FormValue("idacc"))
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		idsch, err := strconv.Atoi(r.FormValue("idsch"))
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`UPDATE "Account" SET "SchemaId" = $1 WHERE id = $2`, idsch, idacc)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
			<body>
				<h1>Успешно</h1>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
			`)
		return
	}
}

func handlerDeleteAccount(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>12 - Удаление аккаунта<h1>
				<form action="/12" method="post">
					<label for="nameInput">Введите Id аккаунта:</label>
					<input type="text" id="nameInput" name="id">
					<button type="submit">Отправить</button>
				</form>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
		`)
		fmt.Fprintln(w, "<p>Аккаунты:<p>")
		accounts := []Account{}
		err := db.Select(&accounts, `SELECT * FROM "Account"`)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		for _, account := range accounts {
			fmt.Fprintf(w, "<p>id: %d Schemaid: %d name: %s<p>", account.Id, account.SchemaId, account.Name)
		}
	} else if r.Method == "POST" {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`DELETE from "Account" where id = $1`, id)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
			<body>
				<h1>Успешно</h1>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
			`)
		return
	}
}

func handlerAviaListAcc(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>13 - Поиск авиакомпаний по аккаунту<h1>
				<form action="/13" method="post">
					<label for="nameInput">Введите Id аккаунта:</label>
					<input type="text" id="nameInput" name="id">
					<button type="submit">Отправить</button>
				</form>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
		`)
		fmt.Fprintln(w, "<p>Аккаунты:<p>")
		accounts := []Account{}
		err := db.Select(&accounts, `SELECT * FROM "Account"`)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		for _, account := range accounts {
			fmt.Fprintf(w, "<p>id: %d Schemaid: %d name: %s<p>", account.Id, account.SchemaId, account.Name)
		}
	} else if r.Method == "POST" {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		type Accoun struct {
			Name string
		}
		accounts := []Accoun{}
		err = db.Select(&accounts, `SELECT 
								"Airline"."name"
								FROM "Account"
								LEFT JOIN "ProviderSchema"
								ON "Account"."SchemaId" = "ProviderSchema"."schema_id"
								LEFT JOIN "Provider"
								ON "ProviderSchema"."provider_id" = "Provider"."id"
								LEFT JOIN "ProviderAirline"
								ON "Provider"."code" = "ProviderAirline"."provider_code"
								LEFT JOIN "Airline"
								ON "Airline"."id" = "ProviderAirline"."airline_id"
								WHERE 
								"Account"."id" = $1
								GROUP BY
								"Airline"."name"`, id)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
			<body>
				<h1>Список авиакомпаний у аккаунта:</h1>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
			`)
		for _, acc := range accounts {
			fmt.Fprintf(w, "<p><p>name: %s", acc.Name)
		}
		return
	}
}

func handlerAviaListProvider(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>14 - Просмотр авиакомпаний которыми владеет поставщик<h1>
				<form action="/14" method="post">
					<label for="nameInput">Введите Id поставщика:</label>
					<input type="text" id="nameInput" name="id">
					<button type="submit">Отправить</button>
				</form>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
		`)
		fmt.Fprintln(w, "<p>Поставщики:<p>")
		providers := []Provider{}
		err := db.Select(&providers, `SELECT * from "Provider"`)
		if err != nil {
			panic(err)
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		for _, provider := range providers {
			fmt.Fprintf(w, "<p>id: %s, name: %s, code: %v<p>", provider.Id, provider.Name, provider.Code)
		}
	} else if r.Method == "POST" {
		id := r.FormValue("id")
		var airlines string
		err := db.Get(&airlines, `SELECT
			 ARRAY_AGG("ProviderAirline"."airline_id") AS "airlines"
		FROM
			 "Provider"
			 left JOIN "ProviderAirline" ON "Provider"."code" = "ProviderAirline"."provider_code"
			 where "Provider"."id"  = $1
		`, id)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
			<html>
			<body>
				<h1>Вот список компаний которыми владеет поставщик с Id: %s</h1>
				<h2>%s<h2>
				<a href="http://localhost:8080">Вернутся на главную</a>
			</body>
			</html>
			`, id, airlines)
		return
	}
}

func main() {
	connStr := "host=localhost port=5432 user=root dbname=airline password=secret sslmode=disable"
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		fmt.Println(http.StatusBadRequest)
		return
	}
	defer db.Close()
	http.HandleFunc("/", handler)
	http.HandleFunc("/1", func(w http.ResponseWriter, r *http.Request) { handlerAddAirline(w, r, db) })
	http.HandleFunc("/2", func(w http.ResponseWriter, r *http.Request) { handlerDeleteAirline(w, r, db) })
	http.HandleFunc("/3", func(w http.ResponseWriter, r *http.Request) { handlerAddProvider(w, r, db) })
	http.HandleFunc("/4", func(w http.ResponseWriter, r *http.Request) { handlerDeleteProvider(w, r, db) })
	http.HandleFunc("/5", func(w http.ResponseWriter, r *http.Request) { handlerModifyProvider(w, r, db) })
	http.HandleFunc("/6", func(w http.ResponseWriter, r *http.Request) { handlerAddSchema(w, r, db) })
	http.HandleFunc("/7", func(w http.ResponseWriter, r *http.Request) { handlerSchemaSearch(w, r, db) })
	http.HandleFunc("/8", func(w http.ResponseWriter, r *http.Request) { handlerModifySchema(w, r, db) })
	http.HandleFunc("/9", func(w http.ResponseWriter, r *http.Request) { handlerDeleteSchema(w, r, db) })
	http.HandleFunc("/10", func(w http.ResponseWriter, r *http.Request) { handlerAccountAdd(w, r, db) })
	http.HandleFunc("/11", func(w http.ResponseWriter, r *http.Request) { handlerModifyAccsSchema(w, r, db) })
	http.HandleFunc("/12", func(w http.ResponseWriter, r *http.Request) { handlerDeleteAccount(w, r, db) })
	http.HandleFunc("/13", func(w http.ResponseWriter, r *http.Request) { handlerAviaListAcc(w, r, db) })
	http.HandleFunc("/14", func(w http.ResponseWriter, r *http.Request) { handlerAviaListProvider(w, r, db) })
	http.ListenAndServe(":8080", nil)
}
