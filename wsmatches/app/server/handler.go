package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
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

//func (s *Srv) pollOLD(c *gin.Context) {
//	const pollTemplate = `
//	<!DOCTYPE html>
//	<html lang="ru">
//	<head>
//		<meta charset="UTF-8">
//		<meta name="viewport" content="width=device-width, initial-scale=1.0">
//		<title>{{.Question}}</title>
//		<style>
//			html, body {
//				height: 100%;
//				margin: 0;
//				padding: 0;
//				display: flex;
//				justify-content: center;
//				align-items: center;
//				font-family: Arial, sans-serif;
//				background-color: #f4f4f4;
//			}
//			.centered-content {
//				width: 100%;
//				max-width: 600px;
//				text-align: center;
//				box-shadow: 0 4px 8px rgba(0,0,0,0.1);
//				background-color: white;
//				padding: 20px;
//				border-radius: 8px;
//			}
//			.question {
//				font-size: 24px;
//				font-weight: bold;
//				margin-bottom: 20px;
//			}
//			ul {
//				list-style-type: none;
//				padding: 0;
//			}
//			li {
//				margin: 10px 0;
//				font-size: 18px;
//			}
//			.votes {
//				color: gray;
//				font-size: 14px;
//			}
//			.user {
//				margin-top: 20px;
//				text-align: left;
//				font-size: 16px;
//			}
//			h3 {
//				font-size: 20px;
//				margin-bottom: 10px;
//			}
//		</style>
//	</head>
//	<body>
//		<div class="centered-content">
//			<div class="question">{{.Question}}</div>
//
//			<ul>
//				{{range $index, $option := .Options}}
//				<li>
//					{{$option}}
//					<span class="votes">Голосов: {{GetVotesCount $.Votes $index}}</span>
//				</li>
//				{{end}}
//			</ul>
//
//			<div class="stats">
//				<h3>Кто проголосовал:</h3>
//				{{range .Votes}}
//				<div class="user">
//					{{.UserName}} - Вариант {{.Answer}}
//					({{if eq .Type "ds"}}Проголосовал в Discord{{else}}Проголосовал в Telegram{{end}})
//				</div>
//				{{end}}
//			</div>
//		</div>
//	</body>
//	</html>`
//
//	// Указываем функцию для подсчета голосов
//	tmpl := template.Must(template.New("poll").Funcs(template.FuncMap{
//		"GetVotesCount": func(votes []models.Votes, optionIndex int) int {
//			count := 0
//			optionStr := strconv.Itoa(optionIndex + 1) // Преобразуем индекс в строку для сравнения с ответом
//			for _, vote := range votes {
//				if vote.Answer == optionStr {
//					count++
//				}
//			}
//			return count
//		},
//	}).Parse(pollTemplate))
//
//	id := c.Param("id")
//	file, err := os.ReadFile("docker/poll/" + id)
//	if err != nil {
//		c.JSON(http.StatusNotFound, "NotFound "+id)
//		return
//	}
//
//	var p models.PollStruct
//	err = json.Unmarshal(file, &p)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, "BadRequest Unmarshal "+id)
//		return
//	}
//	fmt.Printf("%s %s\n", p.Config.HostRelay, p.Question)
//
//	sort.Slice(p.Votes, func(i, j int) bool {
//		return p.Votes[i].Answer < p.Votes[j].Answer
//	})
//
//	// Рендеринг страницы с данными о голосовании
//	c.Status(http.StatusOK)
//	if err := tmpl.Execute(c.Writer, p); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template: " + err.Error()})
//		return
//	}
//}

// Структура для передачи данных в шаблон
type PollTemplateData struct {
	models.PollStruct
	TotalVotes int
	// Сгруппированные голоса: ключ - индекс варианта (например, "1", "2"), значение - список проголосовавших
	GroupedVotes map[string][]models.Votes
}

func (s *Srv) poll(c *gin.Context) {
	const pollTemplate = `
    <!DOCTYPE html>
    <html lang="ru">
    <head>
       <meta charset="UTF-8">
       <meta name="viewport" content="width=device-width, initial-scale=1.0">
       <title>Опрос: {{.Question}}</title>
       <style>
          /* Общие стили (без изменений) */
          html, body { min-height: 100%; margin: 0; padding: 0; display: flex; justify-content: center; align-items: flex-start; padding-top: 40px; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f0f2f5; }
          .centered-content { width: 100%; max-width: 650px; background-color: white; padding: 30px; border-radius: 12px; box-shadow: 0 6px 15px rgba(0,0,0,0.1); }
          .question { font-size: 28px; font-weight: 600; color: #333; margin-bottom: 30px; text-align: center; border-bottom: 2px solid #007bff; padding-bottom: 15px; }
          
          /* Секция Вариантов Ответа (Результаты) */
          h3 { font-size: 22px; color: #007bff; margin-top: 25px; margin-bottom: 15px; padding-bottom: 5px; border-bottom: 1px solid #e0e0e0; }
          ul { list-style-type: none; padding: 0; }
          
          /* Стили li */
          li { background-color: #f9f9f9; border: 1px solid #eee; border-radius: 6px; margin-bottom: 10px; font-size: 16px; color: #2c3e50; position: relative; overflow: hidden; cursor: pointer; padding: 0; }
          li:hover { background-color: #f0f0f0; }

          .progress-container { position: relative; width: 100%; min-height: 40px; overflow: hidden; border-radius: 6px 6px 0 0; }
          .progress-bar { position: absolute; top: 0; left: 0; height: 100%; background-color: #d1e7ff; transition: width 0.5s ease-out; z-index: 1; }
          
          /* Контент (текст и счетчик голосов) */
          .option-content { position: relative; z-index: 2; display: flex; justify-content: space-between; align-items: center; padding: 10px 15px; }

          /* Счетчики Голосов */
          .votes { font-weight: bold; color: #0056b3; white-space: nowrap; }
          
          /* Секция Проголосовавших */
          .user-list { display: none; padding: 10px 15px 5px 15px; margin: 0 -1px -1px -1px; border-top: 1px solid #ddd; background-color: #fff; border-radius: 0 0 6px 6px; }
          
          /* Стиль заголовка списка и кнопки */
          .voters-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 5px; }
          
          /* Кнопка копирования */
          .copy-button { background-color: #007bff; color: white; border: none; border-radius: 4px; padding: 5px 10px; font-size: 13px; cursor: pointer; transition: background-color 0.2s; }
          .copy-button:hover { background-color: #0056b3; }
          .copy-button:active { background-color: #003d80; }

          .user { text-align: left; font-size: 15px; color: #555; padding: 5px 0; border-bottom: 1px dotted #ccc; }
          .user:last-child { border-bottom: none; }
          .user-source { font-style: italic; color: #888; margin-left: 10px; }

          /* Элемент для копирования (скрыт) */
          .copy-target { position: absolute; left: -9999px; opacity: 0; }
       </style>
    </head>
    <body>
       <div class="centered-content">
          <div class="question">{{.Question}}</div>

          <h3>Результаты голосования ({{.TotalVotes}} {{if eq .TotalVotes 1}}голос{{else}}голосов{{end}})</h3>
          
          {{$totalVotes := .TotalVotes}}

          <ul>
             {{range $index, $option := .Options}}
             {{$count := GetVotesCount $.GroupedVotes (Itoa (Add $index 1))}}
             {{$percentage := 0.0}}
             {{if gt $totalVotes 0}}
                {{$percentage = CalculatePercentage $count $totalVotes}}
             {{end}}
             
             <li onclick="toggleVoters('voters-{{$index}}')">
                <div class="progress-container">
                    <div class="progress-bar" style="width: {{$percentage}}%;"></div>
                    
                    <div class="option-content">
                        <div>{{$option}}</div> 
                        <span class="votes">{{$count}} голосов ({{$percentage}}%)</span>
                    </div>
                </div>
                
                <textarea class="copy-target" id="full-copy-target-{{$index}}" readonly>
*** Результаты опроса: {{$.Question}} ***
Вариант: {{$option}} | Голосов: {{$count}} ({{$percentage}}%)
--- Проголосовали ({{$count}}): ---
{{range GetGroupedVotes $.GroupedVotes (Itoa (Add $index 1))}}{{.UserName}} ({{if eq .Type "ds"}}Discord{{else}}Telegram{{end}})
{{end}}
</textarea>

                <div class="user-list" id="voters-{{$index}}">
                    <div class="voters-header">
                       <strong>Проголосовали ({{$count}}):</strong>
                       
                       <button class="copy-button" 
                               onclick="event.stopPropagation(); copyToClipboard('full-copy-target-{{$index}}', this);">
                               Копировать результат
                       </button>
                    </div>
                    <hr style="border: none; border-top: 1px dashed #eee; margin: 5px 0;">
                    
                    {{range GetGroupedVotes $.GroupedVotes (Itoa (Add $index 1))}}
                       <div class="user">
                          <strong>{{.UserName}}</strong>
                          <span class="user-source">
                            ({{if eq .Type "ds"}}Discord{{else}}Telegram{{end}})
                          </span>
                       </div>
                    {{end}}
                </div>
             </li>
             {{end}}
          </ul>

          <script>
            // --- ФУНКЦИЯ JAVASCRIPT ДЛЯ ПЕРЕКЛЮЧЕНИЯ СПИСКОВ ---
            let activeId = null;

            function toggleVoters(targetId) {
                const targetElement = document.getElementById(targetId);
                
                if (activeId === targetId) {
                    targetElement.style.display = 'none';
                    activeId = null;
                    return;
                }

                if (activeId) {
                    const activeElement = document.getElementById(activeId);
                    if (activeElement) {
                        activeElement.style.display = 'none';
                    }
                }

                targetElement.style.display = 'block';
                activeId = targetId;
            }
            
            // --- УНИВЕРСАЛЬНАЯ ФУНКЦИЯ JAVASCRIPT ДЛЯ КОПИРОВАНИЯ ---
            function copyToClipboard(targetId, buttonElement) {
                const target = document.getElementById(targetId);
                
                if (navigator.clipboard) {
                    navigator.clipboard.writeText(target.value).then(() => {
                        const originalText = buttonElement.textContent;
                        buttonElement.textContent = 'Скопировано!';
                        setTimeout(() => {
                            buttonElement.textContent = originalText;
                        }, 1500);
                    }).catch(err => {
                        console.error('Не удалось скопировать текст: ', err);
                        alert('Не удалось скопировать. Попробуйте вручную.');
                    });
                } else {
                    target.select();
                    try {
                        document.execCommand('copy');
                        const originalText = buttonElement.textContent;
                        buttonElement.textContent = 'Скопировано!';
                        setTimeout(() => {
                            buttonElement.textContent = originalText;
                        }, 1500);
                    } catch (err) {
                        console.error('Не удалось скопировать текст: ', err);
                        alert('Не удалось скопировать. Попробуйте вручную.');
                    }
                }
            }
          </script>
       </div>
    </body>
    </html>`

	// --- Структуры и логика Go (без изменений) ---

	type PollTemplateData struct {
		models.PollStruct
		TotalVotes   int
		GroupedVotes map[string][]models.Votes
	}

	groupVotes := func(votes []models.Votes) map[string][]models.Votes {
		grouped := make(map[string][]models.Votes)
		for _, vote := range votes {
			grouped[vote.Answer] = append(grouped[vote.Answer], vote)
		}
		return grouped
	}

	tmpl := template.Must(template.New("poll").Funcs(template.FuncMap{
		"GetVotesCount": func(grouped map[string][]models.Votes, optionKey string) int {
			return len(grouped[optionKey])
		},
		"GetGroupedVotes": func(grouped map[string][]models.Votes, optionKey string) []models.Votes {
			return grouped[optionKey]
		},
		"CalculatePercentage": func(count int, total int) string {
			if total == 0 {
				return "0"
			}
			percent := (float64(count) / float64(total)) * 100
			return fmt.Sprintf("%.1f", percent)
		},
		"Itoa": strconv.Itoa,
		"Add":  func(a, b int) int { return a + b },
	}).Parse(pollTemplate))

	id := c.Param("id")
	file, err := os.ReadFile("docker/poll/" + id)
	if err != nil {
		c.JSON(http.StatusNotFound, "NotFound "+id)
		return
	}

	var p models.PollStruct
	if err = json.Unmarshal(file, &p); err != nil {
		c.JSON(http.StatusBadRequest, "BadRequest Unmarshal "+id)
		return
	}
	fmt.Printf("%s %s\n", p.Config.HostRelay, p.Question)

	groupedVotes := groupVotes(p.Votes)
	totalVotes := len(p.Votes)

	data := PollTemplateData{
		PollStruct:   p,
		TotalVotes:   totalVotes,
		GroupedVotes: groupedVotes,
	}

	c.Status(http.StatusOK)
	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template: " + err.Error()})
		return
	}
}
