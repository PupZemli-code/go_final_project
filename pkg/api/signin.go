package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Секретный ключ для подписи
var secretKey = "qwerty"
var SignedToken string

// todo
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	// Принемает пароль из переменной окружения "TODO_PASSWORD"
	passTODO := os.Getenv("TODO_PASSWORD")
	if len(passTODO) > 0 {

		var data map[string]interface{}

		// Декодируем тело запроса
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			Logger.Printf("ошибка декодирования JSON: %v", err)
			sendError(w, http.StatusInternalServerError, fmt.Errorf("ошибка декодирования JSON: %v", err))
			return
		}

		// Получаем пароль
		password, ok := data["password"].(string)
		if !ok {
			Logger.Printf("неверный формат пароля: %v", err)
			sendError(w, http.StatusInternalServerError, fmt.Errorf("неверный формат пароля: %v", err))
			return
		}

		// Проверяет пароль
		if passTODO == password {

			// Получает hash из пароля, переводит в
			// HEX (от англ. hexadecimal — «шестнадцатеричный») представление
			hashPassTODO := sha256.Sum256([]byte(passTODO))
			hashString := hex.EncodeToString(hashPassTODO[:])

			// Создаем полезную нагрузку
			claims := jwt.MapClaims{
				"hash": hashString,
			}

			// Создаем токен
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			// Подписываем токен хешом из пароля и секретного ключа
			SignedToken, err = token.SignedString([]byte(secretKey))
			if err != nil {
				Logger.Printf("ошибка подписи jwt: %v", err)
				sendError(w, http.StatusInternalServerError, fmt.Errorf("ошибка подписи jwt: %v", err))
				return
			}
			writeJson(w, map[string]any{"token": SignedToken})
			return
		}
		Logger.Printf("неверный пароль")
		sendError(w, http.StatusUnauthorized, fmt.Errorf("неверный пароль"))
	}
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		passTODO := os.Getenv("TODO_PASSWORD")
		if len(passTODO) > 0 {
			var token string // JWT-токен из куки

			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				Logger.Println("cookie.Value идентифицирован")
				token = cookie.Value
			}

			// парсит токен
			jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})
			if err != nil {
				fmt.Printf("ошибка парсинга токена: %s\n", err)
				return
			}

			// проверка валидности
			if !jwtToken.Valid {
				Logger.Println("ошибка авторизации, токен не валиден")
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			// Извлекает Claims из токена
			res, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				Logger.Println("ошибка чтения jwt.MapCalims")
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			// Извлекает hash из Claims
			hashRaw := res["hash"]
			var hash string
			hash, ok = hashRaw.(string)
			if !ok {
				Logger.Println("ошибка приведения hashRaw к типу string jwt.MapCalims")
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			// Получает hash из пароля, переводит в
			// HEX (от англ. hexadecimal — «шестнадцатеричный») представление
			bytes := sha256.Sum256([]byte(passTODO))
			hashString := hex.EncodeToString(bytes[:])

			// Сверяет хиши, необходимо при изменении пароля,
			// чтобы токен утратил валидность
			if hashString != hash {
				Logger.Println("ошибка авторизации, hash не совпадает")
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
