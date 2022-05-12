package api

import (
	"encoding/json"
	"fmt"
	"github.com/Issei0804-ie/who-is-in-a-lab/domain"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"os"
)

func InitAPI(members *[]domain.Member) {
	r := gin.Default()
	h := `{{ define "index.tmpl" }}
	<html lang="jp">
	<head>
	<meta charset="UTF-8">
	<title>kono-lab</title>
	</head>
	<body>
	<h2> 在学者 </h2>
	{{ range .}}
	{{ if .IsLab}}
	<li>{{ .Name }} </li>
	{{ end }}
	{{ end }}

	<h2>帰宅者</h2>

	{{ range .}}
	{{ if eq .IsLab  false}}
	<li>{{ .Name }} </li>
	{{ end }}
	{{ end }}

	</body>
	</html>
	{{end}}`

	r.GET("/", func(c *gin.Context) {
		limit := 30
		for i := 0; i < len(*members); i++ {
			(*members)[i].SetIsLab(limit)
		}
		t := template.Must(template.New("index.tmpl").Parse(h))
		err := t.Execute(c.Writer, *members)
		if err != nil {
			log.Fatal(err.Error())
			c.JSON(http.StatusInternalServerError, map[string]string{"message": "server error"})
			return
		}
	})

	r.POST("/register", func(c *gin.Context) {

		newMember := domain.Member{}
		err := c.BindJSON(&newMember)
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Printf("%v, \n", newMember)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"message": "invalid json format."})
			return
		}

		if newMember.Name == "" || newMember.Addresses == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"message": "invalid body. check parameter."})
			return
		}

		// 送られてきた mac address が既に登録されていないか確認
		for _, member := range *members {
			for _, newMemberAddress := range newMember.Addresses {
				for _, address := range member.Addresses {
					if address == newMemberAddress {
						c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"message": "this mac address is already used. if you want to remove stored mac address, you need to ask issei."})
						return
					}
				}
			}
		}

		didAdd := false
		// 同じ名前なら mac address を追加
		for i, member := range *members {
			if newMember.Name == member.Name {
				(*members)[i].Addresses = append(member.Addresses, newMember.Addresses...)
				didAdd = true
			}
		}

		if !didAdd {
			*members = append(*members, newMember)
		}
		jsonMembers, err := json.Marshal(members)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadGateway, map[string]string{"message": "server error"})
			return
		}
		file, err := os.Create("./address.json")
		defer file.Close()
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadGateway, map[string]string{"message": "server error"})
			return
		}

		file.Write(jsonMembers)
		file.Sync()
		c.JSON(http.StatusOK, map[string]string{"message": "ok"})
		return
	})
	r.Run(":80")
}
