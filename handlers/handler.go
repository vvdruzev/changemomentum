package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/heroku/changemomentum/db"
	"github.com/heroku/changemomentum/schema"
	"github.com/heroku/changemomentum/util"
	"html/template"
	"net/http"
	"strconv"
)

type Handler struct {
	Tmpl *template.Template
}

func NewHandler() *Handler  {
	return &Handler{
	}
}

func (h *Handler) List(c *gin.Context) {
	contacts,_ := db.List()

	//err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
	//	Contacts map[int]schema.Contact
	//	Phones map[int][]string
	//	schema.Search
	//}{
	//	contacts,
	//	phones,
	//	schema.Search{"",},
	//})

	c.HTML(200,"index.html", contacts)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

}

func (h *Handler) AddForm(c *gin.Context) {
	//err := h.Tmpl.ExecuteTemplate(c.Writer, "create.html", nil)
	c.HTML(200,"create.html", nil)
	//if err != nil {
	//	http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	//	return
	//}
}

func (h *Handler) Add(c *gin.Context) {
	// в целям упрощения примера пропущена валидация
	//err := 	db.AddContact(r.FormValue("firstname"),	r.FormValue("lastname"))
	err := 	db.AddContact(c.Param("firstname"),	c.Param("lastname"))
	if err != nil {
		util.ResponseError(c.Writer,500,"Can't add contact")
	}
	c.Redirect(http.StatusFound,"/")
}

func (h *Handler) AddFormPhone(c *gin.Context) {
	//vars := mux.Vars(r)
	item  := &schema.Phone{}

	id,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ResponseError(c.Writer,500,"Bad id")
	}
	item.ContactId = id
	//err = h.Tmpl.ExecuteTemplate(w, "addphone.html", item)
	c.HTML(200,"addphone.html", item )
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AddPhone(c *gin.Context) {
	//vars := mux.Vars(r)
	id,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ResponseError(c.Writer,500,"Bad id")
	}

	// в целям упрощения примера пропущена валидация
	err = 	db.AddPhone(id,c.Param("PhoneNumber"))
	if err != nil {
		util.ResponseError(c.Writer,500,"Error add contact")
	}
	c.Redirect(http.StatusFound, "/")
	//http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Edit (c *gin.Context) {
	//vars := mux.Vars(r)
	id,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ResponseError(c.Writer,500,"Bad id")
	}
	contact, err := db.SelectItem(id)

	//err = h.Tmpl.ExecuteTemplate(c.Writer, "edit.html", contact)
	c.HTML(200,"edit.html", contact )
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Update(c *gin.Context) {
	//vars := mux.Vars(r)
	//r.ParseForm()
	fmt.Println(c.Param("phonenumber"))
	id,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ResponseError(c.Writer,500,"Bad id")
	}
	var contact schema.Contact
	contact.Id = id
	contact.FirstName = c.Param("firstname")
	contact.LastName  = c.Param("lastname")
	//var phonenumbers map[string][]string
	phonenumbers := c.GetStringMapStringSlice("phonenumber")
	err = db.Update(contact,phonenumbers["phonenumber"])
	if err != nil  {
		util.ResponseError(c.Writer,404,"Can't update contact")
	}

	c.Redirect( http.StatusFound ,"/" )

}

func (h *Handler) Delete(c *gin.Context) {
	vars := c.PostFormMap("id")
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.ResponseError(c.Writer,500,"Bad id")
	}

	err = db.Delete(id)
	if err != nil {
		util.ResponseError(c.Writer,500,"Error delete contact")
	}

	c.Writer.Header().Set("Content-type", "application/json")
	resp :=[]byte(`{"affected": ` + strconv.Itoa(int(id)) + `}`)
	c.Writer.Write(resp)

}


func (h *Handler) Search(c *gin.Context) {
	field :=c.Param("field")

	contacts, err := db.Search(field)

	c.HTML(200, "index.html", contacts)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

}