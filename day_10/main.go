package main

import (
	"context"
	"day_10/connection"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Blog struct {
	Id           int
	Name         string
	StartDate    time.Time
	EndDate      time.Time
	Description  string
	Technologies []string
	Image        string
	Post_date    string
	Author       string
}

var dataBlogs = []Blog{
	/*{
		Name:        "Alpha",
		Post_date:   "20-06-1999",
		Author:      "Beta",
		Description: "Test",
	},
	{
		Name:        "Gamma",
		Post_date:   "17-09-2000",
		Author:      "beta",
		Description: "Hellol",
	},*/
}

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/blog", blog).Methods("GET")
	route.HandleFunc("/blog-detail/{id}", blogDetail).Methods("GET")
	route.HandleFunc("/form-blog", formAddBlog).Methods("GET")
	route.HandleFunc("/add-blog", addBlog).Methods("POST")
	route.HandleFunc("/delete-blog/{id}", deleteBlog).Methods("GET")
	route.HandleFunc("/form-edit-blog/{id}", formEditBlog).Methods("GET")
	route.HandleFunc("/edit-blog/{id}", editBlog).Methods("POST")

	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", route)

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	rows, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects ORDER BY id ASC")

	dataBlogs = []Blog{}

	for rows.Next() {
		var each = Blog{}
		err := rows.Scan(&each.Id, &each.Name, &each.StartDate, &each.EndDate, &each.Description, &each.Technologies, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Author = "rudi"
		each.Post_date = each.StartDate.Format("2006-01-02")

		dataBlogs = append(dataBlogs, each)
	}

	//fmt.Println(dataBlogs)
	respData := map[string]interface{}{
		"Blogs": dataBlogs,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func blog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/datablog.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	rows, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects ORDER BY id ASC")

	var result []Blog

	for rows.Next() {
		var each = Blog{} // manggil struct

		err := rows.Scan(&each.Id, &each.Name, &each.StartDate, &each.EndDate, &each.Description, &each.Technologies, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Author = "rudi"
		each.Post_date = each.StartDate.Format("2006-01-02")

		//fmt.Println(each.Technologies)
		result = append(result, each)
	}

	//dataBlogs[0].Technologies = []string{"fa-brands fa-google", "fa-brands fa-github", "fa-brands fa-windows", "fa-brands fa-android"}

	//fmt.Println(dataBlogs[0].Technologies[0])
	fmt.Println(result)

	respData := map[string]interface{}{
		"Blogs": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func blogDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/blog-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var BlogDetail = Blog{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id=$1", id).Scan(
		&BlogDetail.Id, &BlogDetail.Name, &BlogDetail.StartDate, &BlogDetail.EndDate, &BlogDetail.Description, &BlogDetail.Technologies, &BlogDetail.Image,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	data := map[string]interface{}{
		"Blog": BlogDetail,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func formAddBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/addproject.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func addBlog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("dataName")
	start_date := r.PostForm.Get("dataStartDate")
	end_date := r.PostForm.Get("dataEndDate")
	description := r.PostForm.Get("dataDescription")
	image := r.PostForm.Get("dataImage")
	technologies := r.Form["dataTechnologies"]

	fmt.Println("Name : " + name)
	fmt.Println("Start Date : " + start_date)
	fmt.Println("End Date : " + end_date)
	fmt.Println("Description : " + description)
	fmt.Println("Image : " + image) //Image : 101166621_p0.jpg
	fmt.Println(technologies)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name, start_date, end_date, description, technologies, image) VALUES ($1, $2, $3, $4, $5, '../public/img/indprof.jpg')", name, start_date, end_date, description, technologies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func deleteBlog(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func formEditBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/editproject.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func editBlog(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("dataName")
	start_date := r.PostForm.Get("dataStartDate")
	end_date := r.PostForm.Get("dataEndDate")
	description := r.PostForm.Get("dataDescription")
	image := r.PostForm.Get("dataImage")
	technologies := r.Form["dataTechnologies"]

	fmt.Println("Name : " + name)
	fmt.Println("Start Date : " + start_date)
	fmt.Println("End Date : " + end_date)
	fmt.Println("Description : " + description)
	fmt.Println("Image : " + image)
	fmt.Println(technologies)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name, start_date, end_date, description, technologies, image) VALUES ($1, $2, $3, $4, $5, '../public/img/indprof.jpg')", name, start_date, end_date, description, technologies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
