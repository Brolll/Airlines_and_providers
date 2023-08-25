package main

import (
	_ "database/sql/driver"
	"fmt"
	_ "github.com/golang-migrate/migrate"
	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"html/template"
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

func showAirlines(w http.ResponseWriter, db *sqlx.DB) {
	airlines := []Airline{}
	err := db.Select(&airlines, `SELECT * FROM "Airline"`)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	for _, airline := range airlines {
		fmt.Fprintf(w, "<p>id: %s, name: %s\n<p>", airline.Id, airline.Name)
	}
}

func showAirlinesAndProviders(w http.ResponseWriter, db *sqlx.DB) {
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
}

func showSchemas(w http.ResponseWriter, db *sqlx.DB) {
	schemas := []SSchema{}
	err := db.Select(&schemas, `SELECT * FROM "SSchema"`)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	for _, schema := range schemas {
		fmt.Fprintf(w, "<p>id: %d, name: %s\n<p>", schema.Id, schema.Name)
	}
}

func showAccounts(w http.ResponseWriter, db *sqlx.DB) {
	accounts := []Account{}
	err := db.Select(&accounts, `SELECT * FROM "Account"`)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	for _, account := range accounts {
		fmt.Fprintf(w, "<p>id: %d Schemaid: %d name: %s<p>", account.Id, account.SchemaId, account.Name)
	}
}

func showProviders(w http.ResponseWriter, db *sqlx.DB) {
	providers := []Provider{}
	err := db.Select(&providers, `SELECT * from "Provider"`)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	for _, provider := range providers {
		fmt.Fprintf(w, "<p>id: %s, name: %s, code: %v<p>", provider.Id, provider.Name, provider.Code)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	tp, err := template.ParseFiles("./ui/html/mainPage.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "mainPage.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
}

func handlerAddAirline(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	tp, err := template.ParseFiles("./ui/html/addAirline.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "addAirline.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		showAirlines(w, db)
		fmt.Fprintln(w, "<h1>Список поставщиков и компаний которыми они владеют:</h1>")
		showAirlinesAndProviders(w, db)
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
		return
	}
}

func handlerDeleteAirline(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	tp, err := template.ParseFiles("./ui/html/deleteAirline.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "deleteAirline.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		showAirlines(w, db)
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
		return
	}
}

func handlerAddProvider(w http.ResponseWriter, r *http.Request, db *sqlx.DB) { //Возможно надо доработать
	tp, err := template.ParseFiles("./ui/html/addProvider.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "addProvider.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	if r.Method == "POST" {
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
		return
	}
}

func handlerDeleteProvider(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	tp, err := template.ParseFiles("./ui/html/deleteProvider.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "deleteProvider.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		showAirlinesAndProviders(w, db)
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
		return
	}
}

func handlerModifyProvider(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	tp, err := template.ParseFiles("./ui/html/modifyProvider.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "modifyProvider.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		showAirlinesAndProviders(w, db)
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
		_, err = db.Exec(`DELETE FROM "ProviderAirline" WHERE provider_code = $1 AND airline_id = $2`, provider_code, id)
		if err != nil {
			fmt.Fprintln(w, http.StatusBadRequest)
			return
		}
		return
	}
}

func handlerAddSchema(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	tp, err := template.ParseFiles("./ui/html/addSchema.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "addSchema.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		showSchemas(w, db)
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
		return
	}
}

func handlerSchemaSearch(w http.ResponseWriter, r *http.Request, db *sqlx.DB) { //Вывод схемы не массивом, а строкой
	tp, err := template.ParseFiles("./ui/html/schemaSearch.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "schemaSearch.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	if r.Method == "POST" {
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
		fmt.Fprintf(w, `<p>%v</p>`, sch)
		return
	}
}

func handlerModifySchema(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	tp, err := template.ParseFiles("./ui/html/modifySchema.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "modifySchema.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		showSchemas(w, db)
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
		return
	}
}

func handlerDeleteSchema(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	tp, err := template.ParseFiles("./ui/html/deleteSchema.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "deleteSchema.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		showSchemas(w, db)
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
		return
	}
}

func handlerAccountAdd(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	tp, err := template.ParseFiles("./ui/html/accountAdd.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "accountAdd.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		showSchemas(w, db)
		fmt.Fprintln(w, "<h1>Аккаунты:</h1>")
		showAccounts(w, db)
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
		return
	}
}

func handlerModifyAccsSchema(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	tp, err := template.ParseFiles("./ui/html/modifyAccsSchema.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "modifyAccsSchema.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		showSchemas(w, db)
		fmt.Fprintln(w, "<h1>Аккаунты:</h1>")
		showAccounts(w, db)
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
		return
	}
}

func handlerDeleteAccount(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	tp, err := template.ParseFiles("./ui/html/deleteAccount.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "deleteAccount.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		showAccounts(w, db)
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
		return
	}
}

func handlerAviaListAcc(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	tp, err := template.ParseFiles("./ui/html/aviaListAcc.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "aviaListAcc.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		showAccounts(w, db)
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
		for _, acc := range accounts {
			fmt.Fprintf(w, "<p><p>name: %s", acc.Name)
		}
		return
	}
}

func handlerAviaListProvider(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	tp, err := template.ParseFiles("./ui/html/aviaListProvider.html")
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	err = tp.ExecuteTemplate(w, "aviaListProvider.html", r)
	if err != nil {
		fmt.Fprintln(w, http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		showProviders(w, db)
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
		fmt.Fprintf(w, `<h2>%s</h2>`, airlines)
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
