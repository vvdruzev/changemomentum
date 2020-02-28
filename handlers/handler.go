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
	participants,_ := db.List()

	c.HTML(200,"index.html", participants)
}


func (h *Handler) AddForm(c *gin.Context) {
	c.HTML(200,"create.html", nil)
}

func (h *Handler) Add(c *gin.Context) {
	// в целям упрощения примера пропущена валидация
	err := 	db.AddContact(c.PostForm("firstname"),	c.PostForm("lastname"))
	if err != nil {
		util.ResponseError(c.Writer,500,"Can't add contact")
	}
	c.Redirect(http.StatusFound,"/")
}

func (h *Handler) AddFormPhone(c *gin.Context) {
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
	//vars :=c.Request.URL.Query()
	id,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ResponseError(c.Writer,500,"Bad id")
	}

	// в целям упрощения примера пропущена валидация
	err = 	db.AddPhone(id,c.PostForm("PhoneNumber"))
	if err != nil {
		util.ResponseError(c.Writer,500,"Error add contact")
	}
	c.Redirect(http.StatusFound, "/")
}

func (h *Handler) Edit (c *gin.Context) {
	//vars := mux.Vars(r)
	id,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ResponseError(c.Writer,500,"Bad id")
	}
	contact, err := db.SelectItem(id)

	c.HTML(200,"edit.html", contact )
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Update(c *gin.Context) {
	fmt.Println(c.Param("phonenumber"))
	id,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ResponseError(c.Writer,500,"Bad id")
	}
	var contact schema.Contact
	contact.Id = id
	contact.FirstName = c.PostForm("firstname")
	contact.LastName  = c.PostForm("lastname")
	//var phonenumbers map[string][]string
	phonenumbers := c.PostFormArray("phonenumber")
	err = db.Update(contact,phonenumbers)
	if err != nil  {
		util.ResponseError(c.Writer,404,"Can't update contact")
	}

	c.Redirect( http.StatusFound ,"/" )

}

func (h *Handler) Delete(c *gin.Context) {
	//vars := c.PostFormMap("id")
	id, err := strconv.Atoi(c.Param("id"))
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

	//vars :=c.Request.URL.Query()
	a := (c.Query("search"))
	contacts, err := db.Search(a)

	c.HTML(200, "index.html", contacts)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

}