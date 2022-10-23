package main

import (
	"context"
	"day_9/connection"
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
	route.HandleFunc("/blog-detail/{index}", blogDetail).Methods("GET")
	route.HandleFunc("/form-blog", formAddBlog).Methods("GET")
	route.HandleFunc("/add-blog", addBlog).Methods("POST")
	route.HandleFunc("/delete-blog/{index}", deleteBlog).Methods("GET")

	var giga = 100
	fmt.Println(giga)
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

	rows, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects")

	for rows.Next() {
		var each = Blog{} // manggil struct

		err := rows.Scan(&each.Id, &each.Name, &each.StartDate, &each.EndDate, &each.Description, &each.Technologies, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Author = "rudi"
		each.Post_date = each.StartDate.Format("2006-01-02")

		dataBlogs = append(dataBlogs, each)
	}

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

	// var query = "SELECT id, title, content FROM tb_blog"

	rows, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image"+" FROM tb_projects")

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
		dataBlogs = append(dataBlogs, each)
	}

	//dataBlogs[0].Technologies = []string{"fa-brands fa-google", "fa-brands fa-github", "fa-brands fa-windows", "fa-brands fa-android"}

	//fmt.Println(dataBlogs[0].Technologies[0])

	respData := map[string]interface{}{
		"Blogs": dataBlogs,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func blogDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/blog-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var BlogDetail = Blog{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range dataBlogs {
		if index == i {
			BlogDetail = Blog{
				Name:        data.Name,
				Description: data.Description,
				Post_date:   data.Post_date,
				Author:      data.Author,
			}
		}
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

	var title = r.PostForm.Get("dataTitle")
	var description = r.PostForm.Get("dataDescription")
	var image = r.PostForm.Get("dataImage")
	var postDate = strconv.Itoa(time.Now().Day()) + "-" + strconv.Itoa(int(time.Now().Month())) + "-" + strconv.Itoa(time.Now().Year())

	fmt.Println("Title : " + title)
	fmt.Println("Content : " + description)
	fmt.Println("Image : " + image) //Image : 101166621_p0.jpg
	fmt.Println(postDate)

	var newBlog = Blog{
		Name:        title,
		Description: description,
		Author:      "rudi",
		Post_date:   postDate,
	}

	dataBlogs = append(dataBlogs, newBlog)

	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

func deleteBlog(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	fmt.Println(index)

	dataBlogs = append(dataBlogs[:index], dataBlogs[index+1:]...)
	fmt.Println(dataBlogs)

	http.Redirect(w, r, "/blog", http.StatusFound)
}
