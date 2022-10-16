package controllers

import (
	"fmt"
	"goblog/admin/helpers"
	"goblog/admin/models"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gosimple/slug"
	"github.com/julienschmidt/httprouter"
)

type Dashboard struct{}

func (dashboard Dashboard) Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	view, err := template.ParseFiles(helpers.Include("dashboard/list")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Posts"] = models.Post{}.GetAll()
	data["Alert"] = helpers.GetAlert(w, r)
	view.ExecuteTemplate(w, "index", data)
}
func (dashboard Dashboard) NewItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	view, err := template.ParseFiles(helpers.Include("dashboard/add")...)
	if err != nil {
		return
	}
	view.ExecuteTemplate(w, "index", nil)
}
func (dashboard Dashboard) Add(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	title := r.FormValue("blog-title")
	slug := slug.Make(title)
	description := r.FormValue("blog-desc")
	categoryId, _ := strconv.Atoi(r.FormValue("blog-category"))
	content := r.FormValue("blog-content")

	//upload
	r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("blog-pic")
	if err != nil {
		fmt.Println(err)
		return
	}
	f, err := os.OpenFile("uploads/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = io.Copy(f, file)
	//upload end
	if err != nil {
		fmt.Println(err)
		return
	}
	models.Post{
		Title:       title,
		Slug:        slug,
		Description: description,
		Content:     content,
		CategoryID:  categoryId,
		Picture_url: "uploads/" + header.Filename,
	}.Add()
	helpers.SetAlert(w, r, "Kayıt İşlemi Başarıyla Gerçekleşti")
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
func (dashboard Dashboard) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	post := models.Post{}.Get(p.ByName("id"))
	post.Delete()
	helpers.SetAlert(w, r, "Kayıt Silme İşlemi Başarıyla Gerçekleşti!")
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
func (dashboard Dashboard) Edit(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	view, err := template.ParseFiles(helpers.Include("dashboard/edit")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Post"] = models.Post{}.Get(p.ByName("id"))
	view.ExecuteTemplate(w, "index", data)
}
func (dashboard Dashboard) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	post := models.Post{}.Get(p.ByName("id"))
	title := r.FormValue("blog-title")
	slug := slug.Make(title)
	description := r.FormValue("blog-desc")
	categoryId, _ := strconv.Atoi(r.FormValue("blog-category"))
	content := r.FormValue("blog-content")
	is_selected := r.FormValue("is_selected")
	var picture_url string
	if is_selected == "1" {
		r.ParseMultipartForm(10 << 20)
		file, header, err := r.FormFile("blog-pic")
		if err != nil {
			fmt.Println(err)
			return
		}
		f, err := os.OpenFile("uploads/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, _ = io.Copy(f, file)
		picture_url = "uploads/" + header.Filename
		os.Remove(post.Picture_url)
	} else {
		picture_url = post.Picture_url
	}
	post.Updates(models.Post{
		CategoryID:  categoryId,
		Title:       title,
		Slug:        slug,
		Description: description,
		Content:     content,
		Picture_url: picture_url,
	})
	helpers.SetAlert(w, r, "Kayıt Düzenleme Başarıyla Gerçekleşti")
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
