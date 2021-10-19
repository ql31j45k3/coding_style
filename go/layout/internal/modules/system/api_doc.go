package system

import (
	"net/http"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func RegisterRouterAPIDoc(condAPI APIDocCond) {
	condAPI.R.GET("/api-doc", getAPIHtml)
}

func getAPIHtml(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("/srv/layout/api/api.html"))
	if err := tmpl.Execute(c.Writer, nil); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("getAPIHtml - tmpl.Execute")

		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
}
