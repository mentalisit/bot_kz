package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strconv"
	"ws/models"
)

func (s *Srv) getWsMatches(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization, content-type")
	limit := c.DefaultQuery("limit", "")
	page := c.DefaultQuery("page", "1")
	filter := c.DefaultQuery("filter", "")

	result := s.getMatchesAll(limit, page, filter)

	c.JSON(http.StatusOK, result)
}

func (s *Srv) getWsCorps(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization, content-type")
	limit := c.DefaultQuery("limit", "")
	page := c.DefaultQuery("page", "1")

	result := s.getCorps(limit, page)

	c.JSON(http.StatusOK, result)
}
func (s *Srv) docs(c *gin.Context) {
	htmlContent := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Docs</title>
    <style>
        html, body {
            height: 100%;
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            font-family: Arial, sans-serif;
            background-color: #f4f4f4; // Светлый фон всей страницы
        }
        .centered-content {
            width: 100%; // Ширина контента равна ширине страницы
            max-width: 600px; // Максимальная ширина контента
            text-align: center;
            box-shadow: 0 4px 8px rgba(0,0,0,0.1);
            background-color: white;
            padding: 20px;
            border-radius: 8px;
        }
        ul {
            list-style-type: none;
            padding: 0;
        }
        li {
            margin: 10px 0;
        }
        a {
            text-decoration: none;
            color: #3366cc;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="centered-content">
        <h1>Select endpoint:</h1>
        <ul>
            <li><a href="/corps">list corporations</a></li>
            <li><a href="/corps?limit=20">list corporations limit 20</a></li>
            <li><a href="/matches">list matches</a></li>
            <li><a href="/matches?limit=20">list matches limit 20</a></li>
        </ul>
    </div>
</body>
</html>
`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}

func (s *Srv) poll(c *gin.Context) {
	const pollTemplate = `
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>{{.Question}}</title>
		<style>
			html, body {
				height: 100%;
				margin: 0;
				padding: 0;
				display: flex;
				justify-content: center;
				align-items: center;
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
			}
			.centered-content {
				width: 100%;
				max-width: 600px;
				text-align: center;
				box-shadow: 0 4px 8px rgba(0,0,0,0.1);
				background-color: white;
				padding: 20px;
				border-radius: 8px;
			}
			.question {
				font-size: 24px;
				font-weight: bold;
				margin-bottom: 20px;
			}
			ul {
				list-style-type: none;
				padding: 0;
			}
			li {
				margin: 10px 0;
				font-size: 18px;
			}
			.votes {
				color: gray;
				font-size: 14px;
			}
			.user {
				margin-top: 20px;
				text-align: left;
				font-size: 16px;
			}
			h3 {
				font-size: 20px;
				margin-bottom: 10px;
			}
		</style>
	</head>
	<body>
		<div class="centered-content">
			<div class="question">{{.Question}}</div>

			<ul>
				{{range $index, $option := .Options}}
				<li>
					{{$option}} 
					<span class="votes">Голосов: {{GetVotesCount $.Votes $index}}</span>
				</li>
				{{end}}
			</ul>

			<div class="stats">
				<h3>Кто проголосовал:</h3>
				{{range .Votes}}
				<div class="user">
					{{.UserName}} - Вариант {{.Answer}} 
					({{if eq .Type "ds"}}Проголосовал в Discord{{else}}Проголосовал в Telegram{{end}})
				</div>
				{{end}}
			</div>
		</div>
	</body>
	</html>`

	// Указываем функцию для подсчета голосов
	tmpl := template.Must(template.New("poll").Funcs(template.FuncMap{
		"GetVotesCount": func(votes []models.Votes, optionIndex int) int {
			count := 0
			optionStr := strconv.Itoa(optionIndex + 1) // Преобразуем индекс в строку для сравнения с ответом
			for _, vote := range votes {
				if vote.Answer == optionStr {
					count++
				}
			}
			return count
		},
	}).Parse(pollTemplate))

	id := c.Param("id")
	file, err := os.ReadFile("docker/poll/" + id)
	if err != nil {
		c.JSON(http.StatusNotFound, "NotFound "+id)
		return
	}

	var p models.PollStruct
	err = json.Unmarshal(file, &p)
	if err != nil {
		c.JSON(http.StatusBadRequest, "BadRequest Unmarshal "+id)
		return
	}
	fmt.Printf("%s %s\n", p.Config.HostRelay, p.Question)

	sort.Slice(p.Votes, func(i, j int) bool {
		return p.Votes[i].Answer < p.Votes[j].Answer
	})

	// Рендеринг страницы с данными о голосовании
	c.Status(http.StatusOK)
	if err := tmpl.Execute(c.Writer, p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template: " + err.Error()})
		return
	}
}
