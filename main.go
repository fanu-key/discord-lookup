package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

const styles = `
    :root {
        --discord-blue: #5865F2;
        --discord-dark: #36393f;
        --discord-darker: #2f3136;
        --discord-light: #dcddde;
    }
    * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    }
    body {
        background-color: var(--discord-dark);
        color: var(--discord-light);
        line-height: 1.6;
    }
    .container {
        max-width: 800px;
        margin: 2rem auto;
        padding: 2rem;
        background-color: var(--discord-darker);
        border-radius: 8px;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }
    h1 {
        color: white;
        text-align: center;
        margin-bottom: 2rem;
        font-size: 2.5rem;
    }
    .form-group {
        margin-bottom: 1.5rem;
    }
    label {
        display: block;
        margin-bottom: 0.5rem;
        color: var(--discord-light);
        font-weight: 500;
    }
    input[type="text"] {
        width: 100%;
        padding: 0.8rem;
        border: none;
        border-radius: 4px;
        background-color: var(--discord-dark);
        color: white;
        font-size: 1rem;
        transition: all 0.3s ease;
    }
    input[type="text"]:focus {
        outline: none;
        box-shadow: 0 0 0 2px var(--discord-blue);
    }
    button {
        display: block;
        width: 100%;
        padding: 1rem;
        background-color: var(--discord-blue);
        color: white;
        border: none;
        border-radius: 4px;
        font-size: 1rem;
        font-weight: 600;
        cursor: pointer;
        transition: background-color 0.3s ease;
    }
    button:hover {
        background-color: #4752c4;
    }
    .info-table {
        width: 100%;
        border-collapse: collapse;
        margin-top: 1rem;
        background-color: var(--discord-dark);
        border-radius: 8px;
        overflow: hidden;
    }
    .info-table tr:hover {
        background-color: rgba(255, 255, 255, 0.1);
    }
    .info-table td {
        padding: 1rem;
        border-bottom: 1px solid var(--discord-darker);
    }
    .info-table td:first-child {
        font-weight: 600;
        width: 30%;
    }
    .back-button {
        display: inline-block;
        margin-top: 2rem;
        padding: 0.8rem 1.5rem;
        background-color: var(--discord-blue);
        color: white;
        text-decoration: none;
        border-radius: 4px;
        font-weight: 500;
        transition: background-color 0.3s ease;
    }
    .back-button:hover {
        background-color: #4752c4;
    }

	select {
		width: 100%;
		padding: 0.8rem;
		border: none;
		border-radius: 4px;
		background-color: var(--discord-dark);
		color: white;
		font-size: 1rem;
		transition: all 0.3s ease;
		appearance: none;
		background-image: url('data:image/svg+xml;charset=US-ASCII,%3Csvg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 4 5"%3E%3Cpath fill="white" d="M2 0L0 2h4zM2 5L0 3h4z"/%3E%3C/svg%3E');
		background-repeat: no-repeat;
		background-position: right 0.8rem center;
		background-size: 0.65em auto;
	}

	select:focus {
		outline: none;
		box-shadow: 0 0 0 2px var(--discord-blue);
	}

	option {
		background-color: var(--discord-dark);
		color: white;
	}
`

var flags = map[int]string{
	0:       "■ | Нет значков | 0",
	1:       "👨‍💻 | Сотрудник Discord | 1",
	2:       "🤝 | Партнер Discord | 2",
	4:       "🎉 | Участник HypeSquad Events | 4",
	8:       "🔍 | Охотник за багами уровня 1 | 8",
	64:      "🛡️ | Дом Отваги | 64",
	128:     "✨ | Дом Блеска | 128",
	256:     "⚖️ | Дом Баланса | 256",
	512:     "🛠️ | Ранний поддерживающий | 512",
	1024:    "📱 | Командный пользователь | 1024",
	16384:   "🔎 | Охотник за багами уровня 2 | 16384",
	65536:   "🤖 | Проверенный бот | 65536",
	131072:  "👨‍🎓 | Ранний разработчик проверенного бота | 131072",
	262144:  "🛡️ | Сертифицированный модератор Discord | 262144",
	4194304: "🧑‍💻 | Активный разработчик | 4194304",
}

var serverInfo map[string]interface{}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	}
	if err != nil {
		log.Printf("Ошибка при попытке открыть браузер: %v", err)
		log.Printf("Пожалуйста, зайдите на localhost:8080 в браузере!")
	}
}

func getInfo(ID string, token string, mode int) {
	var req *http.Request
	if mode == 1 {
		req, _ = http.NewRequest("GET", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/preview", ID), nil)
	}
	if mode == 2 {
		req, _ = http.NewRequest("GET", fmt.Sprintf("https://discord.com/api/v10/users/%s", ID), nil)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bot %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Ошибка при получении информации: %v", err)
		return
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&serverInfo); err != nil {
		log.Printf("Ошибка при декодировании информации: %v", err)
	}
}

func formatColor(color interface{}) string {
	if color == nil {
		return "N/A"
	}
	switch v := color.(type) {
	case float64:
		return fmt.Sprintf("#%06x", int(v))
	case string:
		return v
	default:
		return "N/A"
	}
}

func filterData(info map[string]interface{}) map[string]interface{} {
	filtered := make(map[string]interface{})

	for key, value := range info {
		if key == "stickers" || key == "emojis" || key == "features" ||
			key == "discovery_splash" || key == "splash" || key == "primary_guild" ||
			key == "clan" || key == "avatar_decoration_data" || key == "primary_guild" || key == "home_header" {
			continue
		}

		if (key == "icon" || key == "avatar") && value != nil {
			name := value.(string)
			nameURL := fmt.Sprintf("https://cdn.discordapp.com/%ss/%s/%s.webp?size=256", key, info["id"], name)
			filtered[key] = nameURL
		} else if key == "accent_color" || key == "banner_color" {
			filtered[key] = formatColor(value)
		} else {
			filtered[key] = value
		}
	}

	return filtered
}

func decodeFlags(flagValue int) []string {
	var decoded []string
	for bit, description := range flags {
		if flagValue&bit != 0 {
			decoded = append(decoded, description)
		}
	}
	if len(decoded) == 0 {
		decoded = append(decoded, flags[0])
	}
	return decoded
}

func generateFlagsHTML(flagValue int) string {
	decodedFlags := decodeFlags(flagValue)

	var sb strings.Builder
	sb.WriteString("<ul style='list-style: none; padding: 0;'>")
	for _, flag := range decodedFlags {
		parts := strings.Split(flag, "|")
		if len(parts) == 3 {
			sb.WriteString(fmt.Sprintf(`
				<li style="margin-bottom: 8px;">
					<span style="margin-right: 8px;">%s</span>
					<span style="font-weight: bold;">%s</span>
				</li>
			`, strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])))
		}
	}
	sb.WriteString("</ul>")

	return sb.String()
}

func getIndexTemplate() string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="ru">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>[TOPOL'] Lookup</title>
			<style>%s</style>
		</head>
		<body>
			<div class="container">
				<h1>[TOPOL'] Lookup</h1>
				<form method="POST" action="/getInfo">
					<div class="form-group">
						<label for="token">Токен бота</label>
						<input type="text" id="token" name="token" placeholder="Введите токен бота Discord" required>
					</div>
					<div class="form-group">
						<label for="lookupType">Выберите тип поиска</label>
						<select id="lookupType" name="lookupType" required>
							<option value="1">Поиск по серверу</option>
							<option value="2">Поиск по пользователю</option>
						</select>
					</div>
					<div class="form-group">
						<label for="ID">Введите ID</label>
						<input type="text" id="ID" name="ID" placeholder="Введите ID сервера или пользователя" required>
					</div>
					<button type="submit">Получить информацию</button>
				</form>
			</div>
		</body>
		</html>
	`, styles)
}

func generateHTMLPage(info map[string]interface{}) string {
	if info == nil {
		return "<h1>Ошибка: данные не найдены</h1>"
	}

	filteredInfo := filterData(info)

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="ru">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Discord Server Information</title>
			<style>%s</style>
		</head>
		<body>
			<div class="container">
				<h1>Полученная информация от Discord</h1>
				<table class="info-table">
	`, styles))

	for key, value := range filteredInfo {
		sb.WriteString("<tr><td>")
		sb.WriteString(key)
		sb.WriteString("</td><td>")

		if key == "icon" || key == "avatar" {
			sb.WriteString(fmt.Sprintf(`<img src="%s" alt="Icon" width="256" height="256">`, value))
		} else if key == "accent_color" || key == "banner_color" {
			color := value.(string)
			sb.WriteString(fmt.Sprintf(`
				<div style="display: flex; align-items: center;">
					<div style="width: 24px; height: 24px; background-color: %s; border: 1px solid #fff; margin-right: 8px;"></div>
					<span>%s</span>
				</div>
			`, color, color))
		} else if key == "flags" || key == "" {
			flagValue, ok := value.(float64)
			if ok {
				sb.WriteString(generateFlagsHTML(int(flagValue)))
			} else {
				sb.WriteString("N/A")
			}
		} else {
			sb.WriteString(fmt.Sprintf("%v", value))
		}
		sb.WriteString("</td></tr>")
	}

	sb.WriteString(`
				</table>
				<a href="/" class="back-button">← Назад</a>
			</div>
		</body>
		</html>
	`)

	return sb.String()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, getIndexTemplate())
}

func getInfoHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	lookupType := r.FormValue("lookupType")
	ID := r.FormValue("ID")

	var mode int
	if lookupType == "1" {
		mode = 1

	} else {
		mode = 2
	}

	getInfo(ID, token, mode)

	htmlContent := generateHTMLPage(serverInfo)
	fmt.Fprint(w, htmlContent)

	serverInfo = nil
}

func startWebServer() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/getInfo", getInfoHandler)

	go func() {
		log.Println("Запуск веб-сервера на порту 8080...")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}

func main() {

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}
	logger.SetLevel(logrus.InfoLevel)
	log.Println("\n\nTOPOL` LOOKUP v1.3\nСервер источник: https://discord.gg/UMbXmbEyF4\nПоявились вопросы: Discord - fanukey\n\nНекоторые предосторожности:\n\n1. Если вы боитесь что токен куда-либо попадёт:\n    Создайте нового бота Discord.\n    Запретите приглашать его кому-либо кроме владельца.\n    Не приглашайте его никуда.\n2. Помните, что в код не вшиты никакие запросы, кроме запросов к Discord API:\n    В случае если токен куда-либо попал, виноваты только вы.\n    К примеру так может случится, если вы используете незашифрованный Proxy-клиент,\n    от неизвестного провайдера.\n\n")
	startWebServer()

	openBrowser("http://localhost:8080")
	select {}
}
