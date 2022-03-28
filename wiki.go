package main

import (
	"log"
	"net/http"
	"os"
	"text/template"
)

type Page struct {
	Title string //标题
	Body  []byte //正文 类型定位为[]byte类型 而不是string，为了方便存储
}

//保存Page到一个文本文件
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600) //0600让当前用户拥有读写权限
}

//将保存的Page读取出来
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("viewHandler:", r.URL.Path)
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	t, err2 := template.ParseFiles("view.html")
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
	}
	err2 = t.Execute(w, p)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	t, err3 := template.ParseFiles("edit.html")
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}
	err3 = t.Execute(w, p)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Fprintf(w, "成功将 %s 保存!", title)
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}
