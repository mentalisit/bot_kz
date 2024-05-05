package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
func (s *Srv) getWsCorpsCount(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Authorization, content-type")
	limit := c.DefaultQuery("limit", "")
	page := c.DefaultQuery("page", "1")

	result := s.getCorpsCount(limit, page)

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
            <li><a href="/corps2">list corporations v2</a></li>
            <li><a href="/corps2?limit=20">list corporations v2 limit 20</a></li>
            <li><a href="/matches">list matches</a></li>
            <li><a href="/matches?limit=20">list matches limit 20</a></li>
        </ul>
    </div>
</body>
</html>
`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}
