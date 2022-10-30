package main

import (
	"context"
	"day_12/connection"
	"day_12/middleware"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
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
	UserId       int
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

var tempImg string

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")

	route.HandleFunc("/blog-detail/{id}", blogDetail).Methods("GET")
	route.HandleFunc("/form-blog", formAddBlog).Methods("GET")
	route.HandleFunc("/add-blog", middleware.UploadFile(addBlog)).Methods("POST")
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

	var rows pgx.Rows

	if session.Values["StatLogin"] != true {

		DataAcc.StatLogin = false
		rows, _ = connection.Conn.Query(context.Background(), "SELECT tb_projects.*, tb_user.name FROM tb_projects LEFT JOIN tb_user ON tb_projects.user_id = tb_user.id ORDER BY id ASC")

	} else {

		DataAcc.StatLogin = session.Values["StatLogin"].(bool)
		DataAcc.UserName = session.Values["Name"].(string)
		rows, _ = connection.Conn.Query(context.Background(), "SELECT tb_projects.*, tb_user.name FROM tb_projects LEFT JOIN tb_user ON tb_projects.user_id = tb_user.id WHERE tb_user.name=$1 ORDER BY id ASC", DataAcc.UserName)

	}

	//rows, _ = connection.Conn.Query(context.Background(), "SELECT tb_projects.*, tb_user.name FROM tb_projects LEFT JOIN tb_user ON tb_projects.user_id = tb_user.id ORDER BY id ASC")

	dataBlogs = []Blog{}

	mon30 := []string{"4", "6", "9", "11"}

	for rows.Next() {
		var each = Blog{}
		err := rows.Scan(&each.Id, &each.Name, &each.StartDate, &each.EndDate, &each.Description, &each.Technologies, &each.Image, &each.UserId, &each.Author)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Post_date = each.StartDate.Format("2006-01-02")

		splStart := strings.Split(each.StartDate.Format("2006-01-02"), "-")

		feb29, _ := strconv.Atoi(splStart[0])
		feb29 = feb29 % 4

		resDur := each.EndDate.Sub(each.StartDate).Hours()
		resDur = resDur / 24
		//fmt.Println(resDur)
		each.Duration = strconv.Itoa(int(resDur)) + " Days"

		if resDur >= 365 {
			if feb29 == 0 {
				resDur = resDur / 366
			} else {
				resDur = resDur / 365
			}
			each.Duration = strconv.Itoa(int(resDur)) + " Years"

		} else if (resDur >= 28) && (splStart[1] == "02") {
			if feb29 == 0 {
				resDur = resDur / 29
			} else {
				resDur = resDur / 28
			}

			each.Duration = strconv.Itoa(int(resDur)) + " Months"

		} else if resDur >= 30 {
			tanda31 := true

			for _, strMon := range mon30 {
				if strMon == splStart[1] {
					tanda31 = false
					resDur = resDur / 30
					break
				}
			}

			if tanda31 == true {
				resDur = resDur / 31
			}

			each.Duration = strconv.Itoa(int(resDur)) + " Months"
		}

		//fmt.Println(each.Duration)

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

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	DataAcc.FlashData = strings.Join(flashes, "")

	//fmt.Println(dataBlogs)
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
	} else {
		DataAcc.StatLogin = session.Values["StatLogin"].(bool)
		DataAcc.UserName = session.Values["Name"].(string)
	}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id=$1", id).Scan(
		&BlogDetail.Id, &BlogDetail.Name, &BlogDetail.StartDate, &BlogDetail.EndDate, &BlogDetail.Description, &BlogDetail.Technologies, &BlogDetail.Image,
	)

	mon30 := []string{"4", "6", "9", "11"}

	splStart := strings.Split(BlogDetail.StartDate.Format("2006-01-02"), "-")

	feb29, _ := strconv.Atoi(splStart[0])
	feb29 = feb29 % 4

	resDur := BlogDetail.EndDate.Sub(BlogDetail.StartDate).Hours()
	resDur = resDur / 24
	//fmt.Println(resDur)
	BlogDetail.Duration = strconv.Itoa(int(resDur)) + " Days"

	if resDur >= 365 {
		if feb29 == 0 {
			resDur = resDur / 366
		} else {
			resDur = resDur / 365
		}
		BlogDetail.Duration = strconv.Itoa(int(resDur)) + " Years"

	} else if (resDur >= 28) && (splStart[1] == "02") {
		if feb29 == 0 {
			resDur = resDur / 29
		} else {
			resDur = resDur / 28
		}

		BlogDetail.Duration = strconv.Itoa(int(resDur)) + " Months"

	} else if resDur >= 30 {
		tanda31 := true

		for _, strMon := range mon30 {
			if strMon == splStart[1] {
				tanda31 = false
				resDur = resDur / 30
				break
			}
		}

		if tanda31 == true {
			resDur = resDur / 31
		}

		BlogDetail.Duration = strconv.Itoa(int(resDur)) + " Months"
	}

	BlogDetail.NodeJS = false
	BlogDetail.NextJS = false
	BlogDetail.ReactJS = false
	BlogDetail.TypeScript = false

	for _, tList := range BlogDetail.Technologies {
		if tList == "Node JS" {
			BlogDetail.NodeJS = true
		}
		if tList == "Next JS" {
			BlogDetail.NextJS = true
		}
		if tList == "React JS" {
			BlogDetail.ReactJS = true
		}
		if tList == "TypeScript" {
			BlogDetail.TypeScript = true
		}
	}

	formStart := BlogDetail.StartDate.Format("2 January 2006")
	formEnd := BlogDetail.EndDate.Format("2 January 2006")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	data := map[string]interface{}{
		"Blog":      BlogDetail,
		"Status":    DataAcc,
		"startDate": formStart,
		"startEnd":  formEnd,
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
		DataAcc.FlashData = "Need Login!!"
	} else {
		DataAcc.StatLogin = session.Values["StatLogin"].(bool)
		DataAcc.UserName = session.Values["Name"].(string)
	}

	t := time.Now()

	currTime := t.Format("2006-01-02")

	data := map[string]interface{}{
		"Status":  DataAcc,
		"timeNow": currTime,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func addBlog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	dataContex := r.Context().Value("dataFile")

	name := r.PostForm.Get("dataName")
	start_date := r.PostForm.Get("dataStartDate")
	end_date := r.PostForm.Get("dataEndDate")
	description := r.PostForm.Get("dataDescription")
	technologies := r.Form["dataTechnologies"]
	image := dataContex.(string)

	userId := session.Values["Id"].(int)

	fmt.Println("Name : " + name)
	fmt.Println("Start Date : " + start_date)
	fmt.Println("End Date : " + end_date)
	fmt.Println("Description : " + description)
	fmt.Println("Image : " + image)
	fmt.Println(technologies)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name, start_date, end_date, description, technologies, image, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", name, start_date, end_date, description, technologies, image, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func deleteBlog(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var deletoImg = Blog{}
	err := connection.Conn.QueryRow(context.Background(), "SELECT image FROM tb_projects WHERE id=$1", id).Scan(
		&deletoImg.Image,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	e := os.Remove("uploads/" + deletoImg.Image)
	if e != nil {
		log.Fatal(e)

	}

	_, err = connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)

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

	tempImg = EditProject.Image

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

	file, handler, err := r.FormFile("dataImage")
	name := r.PostForm.Get("dataName")
	start_date := r.PostForm.Get("dataStartDate")
	end_date := r.PostForm.Get("dataEndDate")
	description := r.PostForm.Get("dataDescription")
	technologies := r.Form["dataTechnologies"]

	if file != nil {

		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode("Error Retrieving the File")
			return
		}
		defer file.Close()

		tempFile, err := ioutil.TempFile("uploads", "*"+handler.Filename)
		if err != nil {
			fmt.Println(err)
			fmt.Println("path upload error")
			json.NewEncoder(w).Encode(err)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		tempFile.Write(fileBytes)

		dataE := tempFile.Name()
		filename := dataE[8:]

		ctx := context.WithValue(r.Context(), "dataFile", filename)

		dataContex := ctx.Value("dataFile")
		image := dataContex.(string)

		e := os.Remove("uploads/" + tempImg)
		if e != nil {
			log.Fatal(e)

		}

		fmt.Println("Image : " + image)

		_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name=$2, start_date=$3, end_date=$4, description=$5, technologies=$6, image=$7 WHERE id=$1", id, name, start_date, end_date, description, technologies, image)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}

	} else {

		_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name=$2, start_date=$3, end_date=$4, description=$5, technologies=$6 WHERE id=$1", id, name, start_date, end_date, description, technologies)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}

	}

	fmt.Println(id)
	fmt.Println("Name : " + name)
	fmt.Println("Start Date : " + start_date)
	fmt.Println("End Date : " + end_date)
	fmt.Println("Description : " + description)
	fmt.Println(technologies)

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
	DataAcc.FlashData = strings.Join(flashes, "")
	session.Options.MaxAge = -1
	session.Save(r, w)
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
		session.AddFlash("Wrong Email!!", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
		/*
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("message : " + err.Error()))
			return*/
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	if err != nil {
		session.AddFlash("Wrong Password!!", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
		/*
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("message : " + err.Error()))
			return*/
	}

	session.Values["StatLogin"] = true
	session.Values["Name"] = user.Name
	session.Values["Id"] = user.Id
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
