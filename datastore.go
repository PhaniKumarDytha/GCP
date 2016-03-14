package datastore

import (
	"fmt"
	"html/template" //We can use html/template to keep the HTML in a separate file, allowing us to change the layout of our edit page without modifying the underlying Go code
	"io/ioutil"     //used to read from and write to files
	"net/http"
	"time"

	"database/sql" //These packages help connect with the CloudSQL

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/ziutek/mymysql/godrv"

	"appengine"
	"appengine/datastore"
)

type Stocks struct {
	ScripID   string
	Timestamp time.Time
	Price     float32 //Make sure the field names start with caps. Otherwise they don't get exposed
	High52    float32
	Low52     float32
	Ucircuit  bool
	Lcircuit  bool
}

func init() {
	http.HandleFunc("/", write_form)
	http.HandleFunc("/submit", save_data)
	http.HandleFunc("/connectsql", connect_sql)
}

//Example with Execute method
//Another simple example is here - http://blog.joshsoftware.com/2014/03/14/learn-to-build-and-deploy-simple-go-web-apps-part-three/
func get_data(w http.ResponseWriter, r *http.Request) {

	p, err := ioutil.ReadFile("form.html")
	if err != nil {
		fmt.Fprintf(w, err.Error()) //the Error method helps display the error as a string
	}
	t, _ := template.ParseFiles("form.html")
	t.Execute(w, p)

}

//Example with ExecuteTemplate
//FInally found the solution here http://stackoverflow.com/questions/29605632/template-execute-invalid-memory-address-or-nil-pointer-dereference
func write_form(w http.ResponseWriter, r *http.Request) {
	s := Stocks{
		ScripID:   "New",
		Timestamp: time.Now(),
		Price:     240,
		High52:    300,
		Low52:     120,
		Ucircuit:  false,
		Lcircuit:  false, //It is important to put this comma
	}

	var mytemplate = template.Must(template.New("form.html").ParseFiles("form.html")) //important : in both cases use the file name "form.html" - it doesn't work otherwise. You'll get an error saying it is an incomplete template
	if err := mytemplate.ExecuteTemplate(w, "form.html", s); err != nil {
		fmt.Fprintf(w, err.Error())

	}
}

//This function connects to a regular CloudSQL instance
// Package cloudsql exposes access to Google Cloud SQL databases.
// This package is intended for MySQL drivers to make App Engine-specific connections.
// Applications should use this package through database/sql:
// Select a pure Go MySQL driver that supports this package, and use sql.Open with protocol "cloudsql" and an address of the Cloud SQL instance.
// A Go MySQL driver that has been tested to work well with Cloud SQL is the go-sql-driver:
func connect_sql(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Attempting to connect to Google Cloud SQL...")
	//db, err := sql.Open("mysql", "user@cloudsql(project-id:instance-name)/dbname")
	//_, err := sql.Open("mysql", "root:admin@cloudsql(project-id:myappmulti:us-central1:myappmulti)/demo")
	db, err := sql.Open("mymysql", "cloudsql:myappmulti*demo/root/adin")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	//Use regular queries to get results and output them
	query_results, qry_err := db.Query("select scripid from stocks")
	for query_results.Next() {
		var scripid string
		qry_err = query_results.Scan(&scripid)
		if qry_err != nil {
			fmt.Fprintf(w, qry_err.Error())
		}
		fmt.Fprintf(w, scripid)
	}
}

func save_data(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	s := Stocks{
		ScripID:   "New",
		Timestamp: time.Now(),
		Price:     240,
		High52:    300,
		Low52:     120,
		Ucircuit:  false,
		Lcircuit:  false, //It is important to put this comma
	}

	key, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Stocks", nil), &s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Stock added to datastore. Key was %q Last updated:%q ", key, s.Timestamp)
}
