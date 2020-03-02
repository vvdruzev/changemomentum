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
	"time"
)

type Handler struct {
	Tmpl *template.Template
	hashToken map[string]int
}

func NewHandler() *Handler  {
	return &Handler{
	}
}

func (h *Handler) List(c *gin.Context) {
	participants,_ := db.List()
	var data struct{ PartNotRegistred []schema.Participant
		PartRegistred []schema.Participant
	}
	for _, val := range participants {
		date,_ := time.Parse("2006-01-02 15:04:05",val.Date)
		if  date.Unix()< 1262304000 {   //01 Jan 2010 00:00:00
			val.Date = ""
			data.PartNotRegistred = append(data.PartNotRegistred, val)
		}else {
			data.PartRegistred = append(data.PartRegistred, val)
		}
	}
	c.HTML(200,"index.html", data)
}


func (h *Handler) AddForm(c *gin.Context) {
	c.HTML(200,"create.html", nil)
}

func (h *Handler) Add(c *gin.Context) {
	// в целям упрощения примера пропущена валидация
	err := 	db.AddContact(c.PostForm("firstname"),	c.PostForm("lastname"), c.PostForm("command"), 1)
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
	participant, err := db.SelectItem(id)

	c.HTML(200,"edit.html", participant )
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
	var participant schema.Participant
	participant.Id = id
	participant.FirstName = c.PostForm("firstname")
	participant.LastName  = c.PostForm("lastname")
	participant.Command = c.PostForm("command")
	//var phonenumbers map[string][]string
	//phonenumbers := c.PostFormArray("phonenumber")
	err = db.Update(participant)
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

func (h *Handler) Registration(c *gin.Context) {
	//vars := c.PostFormMap("id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ResponseError(c.Writer,500,"Bad id")
	}

	err = db.Registration(id)
	if err != nil {
		util.ResponseError(c.Writer,500,"Error delete contact")
	}

	c.Writer.Header().Set("Content-type", "application/json")
	resp :=[]byte(`{"affected": ` + strconv.Itoa(int(id)) + `}`)
	c.Writer.Write(resp)

}

func (h *Handler) UnRegistration(c *gin.Context) {
	//vars := c.PostFormMap("id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ResponseError(c.Writer,500,"Bad id")
	}

	err = db.UnRegistration(id)
	if err != nil {
		util.ResponseError(c.Writer,500,"Error delete contact")
	}

	participants,_ := db.List()
	var data struct{ PartNotRegistred []schema.Participant
		PartRegistred []schema.Participant
	}
	for _, val := range participants {
		date,_ := time.Parse("2006-01-02 15:04:05",val.Date)
		if  date.Unix()< 1262304000 {   //01 Jan 2010 00:00:00
			val.Date = ""
			data.PartNotRegistred = append(data.PartNotRegistred, val)
		}else {
			data.PartRegistred = append(data.PartRegistred, val)
		}
	}

	c.Writer.Header().Set("Content-type", "application/json")
	resp :=[]byte(`{"affected": ` + strconv.Itoa(int(id)) + `}`)
	c.Writer.Write(resp)

}


func (h *Handler) Search(c *gin.Context) {

	//vars :=c.Request.URL.Query()
	a := (c.Query("search"))
	participants, err := db.Search(a)
	var data struct{ PartNotRegistred []schema.Participant
		PartRegistred []schema.Participant
	}
	for _, val := range participants {
		date,_ := time.Parse("2006-01-02 15:04:05",val.Date)
		if  date.Unix()< 1262304000 {   //01 Jan 2010 00:00:00
			val.Date = ""
			data.PartNotRegistred = append(data.PartNotRegistred, val)
		}else {
			data.PartRegistred = append(data.PartRegistred, val)
		}
	}

	c.HTML(200, "index.html", data)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

}