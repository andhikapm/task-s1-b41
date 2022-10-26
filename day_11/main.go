package main

import (
	"context"
	"day_11/connection"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type MetaData struct {
	Title     string
	StatLogin bool
	UserName  string
	FlashData string
}

var DataAcc = MetaData{
	Title: "Personal Web",
}

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
	Duration     string
	StatLogin    bool
	NodeJS       bool
	NextJS       bool
	ReactJS      bool
	TypeScript   bool
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

var dataBlogs = []Blog{}

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")

	route.HandleFunc("/blog-detail/{id}", blogDetail).Methods("GET")
	route.HandleFunc("/form-blog", formAddBlog).Methods("GET")
	route.HandleFunc("/add-blog", addBlog).Methods("POST")
	route.HandleFunc("/delete-blog/{id}", deleteBlog).Methods("GET")

	route.HandleFunc("/form-edit-blog/{id}", formEditBlog).Methods("GET")
	route.HandleFunc("/edit-blog/{id}", editBlog).Methods("POST")

	route.HandleFunc("/form-register", formRegister).Methods("GET")
	route.HandleFunc("/register", register).Methods("POST")

	route.HandleFunc("/form-login", formLogin).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")

	route.HandleFunc("/logout", logout).Methods("GET")

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

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["StatLogin"] != true {
		DataAcc.StatLogin = false
	} else {
		DataAcc.StatLogin = session.Values["StatLogin"].(bool)
		DataAcc.UserName = session.Values["Name"].(string)
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

		each.NodeJS = false
		each.NextJS = false
		each.ReactJS = false
		each.TypeScript = false

		for _, tList := range each.Technologies {
			if tList == "Node JS" {
				each.NodeJS = true
			}
			if tList == "Next JS" {
				each.NextJS = true
			}
			if tList == "React JS" {
				each.ReactJS = true
			}
			if tList == "TypeScript" {
				each.TypeScript = true
			}
		}

		if session.Values["StatLogin"] != true {
			each.StatLogin = false
		} else {
			each.StatLogin = session.Values["StatLogin"].(bool)
		}

		dataBlogs = append(dataBlogs, each)
	}

	respData := map[string]interface{}{
		"Blogs":  dataBlogs,
		"Status": DataAcc,
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

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["StatLogin"] != true {
		DataAcc.StatLogin = false
	} else {
		DataAcc.StatLogin = session.Values["StatLogin"].(bool)
		DataAcc.UserName = session.Values["Name"].(string)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, DataAcc)
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

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	var BlogDetail = Blog{}

	if session.Values["StatLogin"] != true {
		DataAcc.StatLogin = false
		BlogDetail.StatLogin = false
	} else {
		DataAcc.StatLogin = session.Values["StatLogin"].(bool)
		DataAcc.UserName = session.Values["Name"].(string)
		BlogDetail.StatLogin = session.Values["StatLogin"].(bool)
	}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id=$1", id).Scan(
		&BlogDetail.Id, &BlogDetail.Name, &BlogDetail.StartDate, &BlogDetail.EndDate, &BlogDetail.Description, &BlogDetail.Technologies, &BlogDetail.Image,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	data := map[string]interface{}{
		"Blog":   BlogDetail,
		"Status": DataAcc,
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

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["StatLogin"] != true {
		DataAcc.StatLogin = false
	} else {
		DataAcc.StatLogin = session.Values["StatLogin"].(bool)
		DataAcc.UserName = session.Values["Name"].(string)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, DataAcc)
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

	var EditProject = Blog{}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["StatLogin"] != true {
		DataAcc.StatLogin = false
	} else {
		DataAcc.StatLogin = session.Values["StatLogin"].(bool)
		DataAcc.UserName = session.Values["Name"].(string)
	}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id=$1", id).Scan(
		&EditProject.Id, &EditProject.Name, &EditProject.StartDate, &EditProject.EndDate, &EditProject.Description, &EditProject.Technologies, &EditProject.Image,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	EditProject.NodeJS = false
	EditProject.NextJS = false
	EditProject.ReactJS = false
	EditProject.TypeScript = false

	for _, tList := range EditProject.Technologies {
		if tList == "Node JS" {
			EditProject.NodeJS = true
		}
		if tList == "Next JS" {
			EditProject.NextJS = true
		}
		if tList == "React JS" {
			EditProject.ReactJS = true
		}
		if tList == "TypeScript" {
			EditProject.TypeScript = true
		}

	}

	t := time.Now()

	currTime := t.Format("2006-01-02")
	formStart := EditProject.StartDate.Format("2006-01-02")
	formEnd := EditProject.EndDate.Format("2006-01-02")

	if currTime >= formStart {
		currTime = formStart
	}

	data := map[string]interface{}{
		"dataEdit":  EditProject,
		"Status":    DataAcc,
		"timeNow":   currTime,
		"timeStart": formStart,
		"timeEnd":   formEnd,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
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

	fmt.Println(id)
	fmt.Println("Name : " + name)
	fmt.Println("Start Date : " + start_date)
	fmt.Println("End Date : " + end_date)
	fmt.Println("Description : " + description)
	fmt.Println("Image : " + image)
	fmt.Println(technologies)

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name=$2, start_date=$3, end_date=$4, description=$5, technologies=$6, image='../public/img/indprof.jpg' WHERE id=$1", id, name, start_date, end_date, description, technologies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func formRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var nameAcc = r.PostForm.Get("dataNameAcc")
	var email = r.PostForm.Get("dataEmail")
	var pass = r.PostForm.Get("dataPass")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(pass), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES($1, $2, $3)", nameAcc, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
}

func formLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, DataAcc)
}

func login(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("dataEmail")
	pass := r.PostForm.Get("dataPass")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	session.Values["StatLogin"] = true
	session.Values["Name"] = user.Name
	session.Options.MaxAge = 7200

	session.AddFlash("Successfully Login!", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout!")
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
