package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"ws/models"

	"github.com/gin-gonic/gin"
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
	// Определяем префикс из текущего пути
	prefix := ""
	if strings.HasPrefix(c.Request.URL.Path, "/ws") {
		prefix = "/ws"
	}

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Docs</title>
    <style>
        html, body {
            height: 100%%;
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
        }
        .centered-content {
            width: 100%%;
            max-width: 600px;
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
            <li><a href="%s/corps">list corporations</a></li>
            <li><a href="%s/corps?limit=20">list corporations limit 20</a></li>
            <li><a href="%s/matches">list matches</a></li>
            <li><a href="%s/matches?limit=20">list matches limit 20</a></li>
        </ul>
    </div>
</body>
</html>
`, prefix, prefix, prefix, prefix)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}

// Структура для передачи данных в шаблон
type PollTemplateData struct {
	models.PollStruct
	TotalVotes int
	// Сгруппированные голоса: ключ - индекс варианта (например, "1", "2"), значение - список проголосовавших
	GroupedVotes map[string][]models.Votes
}

func (s *Srv) poll(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "not found id")
		return
	}
	targetURL := fmt.Sprintf("https://mentalisit.myds.me/web/poll.html?id=%s", id)

	// Выполняем переадресацию (302 Found)
	c.Redirect(http.StatusFound, targetURL)
}

func (s *Srv) pollAPI(c *gin.Context) {
	id := c.Param("id")

	var p models.PollStruct

	file, err := os.ReadFile("docker/poll/" + id)
	if err == nil {
		if err = json.Unmarshal(file, &p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "BadRequest Unmarshal " + id})
			return
		}
	} else {
		poll2Struct, votes, err := s.Db.GetPollById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "NotFound " + id})
			return
		}
		p = poll2Struct
		p.Votes = votes
	}

	fmt.Printf("%s %s\n", p.Config.HostRelay, p.Question)

	// Group votes by answer
	groupedVotes := make(map[string][]models.Votes)
	for _, vote := range p.Votes {
		groupedVotes[vote.Answer] = append(groupedVotes[vote.Answer], vote)
	}

	// Prepare API response
	response := gin.H{
		"question":      p.Question,
		"options":       p.Options,
		"total_votes":   len(p.Votes),
		"grouped_votes": groupedVotes,
	}

	c.JSON(http.StatusOK, response)
}
