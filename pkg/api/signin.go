package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

/*
Почему-то с jwt.RegisteredClaims возвращался пустой СheckSum уже при получении. Не стал убирать. Закоментировал.
type Claims struct {
	СheckSum             [32]byte // чексум пароля
	jwt.RegisteredClaims          // базовый тип
}*/

type Pwd struct {
	Pwd string `json:"password"`
}

const MySecret = "my_secret_key"

func taskSignInHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var pwd Pwd
	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	// десериализуем JSON
	if err = json.Unmarshal(buf.Bytes(), &pwd); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	pass := os.Getenv("TODO_PASSWORD")
	if pass != pwd.Pwd {
		writeJson(w, map[string]string{"error": "Неверный пароль"})
		return
	}
	//claims := Claims{}

	checksum := sha256.Sum256([]byte(pwd.Pwd))
	// Claims хранит преобразованный в строку hash пароля в токене
	claims := jwt.MapClaims{
		"checksum": hex.EncodeToString(checksum[:]), // преобразуем в hex-строку
	}

	//claims.СheckSum = sha256.Sum256([]byte(pwd.Pwd))
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// получаем подписанный токен
	signedToken, err := jwtToken.SignedString([]byte(MySecret))
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]string{"token": signedToken})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var tokenString string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenString = cookie.Value
			} else {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
			// здесь код для валидации и проверки JWT-токена
			/*			var claims Claims
						//claims := Claims{}
						token, err := jwt.ParseWithClaims(tokenString, &claims,
							func(token *jwt.Token) (interface{}, error) {
								// желательно проверять используемый метод
								if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
									return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
								}
								return []byte(MySecret), nil
							})
						if err != nil {
							writeJson(w, map[string]string{"error": err.Error()})
						}
						if claims.СheckSum == sha256.Sum256([]byte(pass)) {
							http.Error(w, "Authentification required", http.StatusUnauthorized)
							return
						}*/
			jwtToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return []byte(MySecret), nil
			})
			if err != nil {
				writeJson(w, map[string]string{"error": err.Error()})
				return
			}
			if !jwtToken.Valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
			// приводим поле Claims к типу jwt.MapClaims
			res, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				writeJson(w, map[string]string{"error": "failed to type assertion of jwt.MapCalims"})
				return
			}
			checksumStr, ok := res["checksum"].(string)
			if !ok {
				writeJson(w, map[string]string{"error": "Invalid checksum format"})
				return
			}
			// Преобразуем hex-строку обратно в байтовый массив
			checksumFromToken, err := hex.DecodeString(checksumStr)
			if err != nil {
				writeJson(w, map[string]string{"error": "Failed to decode checksum"})
				return
			}
			// Получаем ожидаемый checksum
			expectedChecksum := sha256.Sum256([]byte(pass))
			// Сравниваем байты по одному
			if !bytes.Equal(checksumFromToken, expectedChecksum[:]) {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
