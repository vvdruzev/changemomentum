package handlers

import (
	"github.com/heroku/changemomentum/schema"
	"time"
)

func validateOutParameters(participants []schema.Participant)  {
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
}