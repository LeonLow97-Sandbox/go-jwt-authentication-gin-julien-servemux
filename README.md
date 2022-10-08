# Go JWT Login Authentication Methods

1. Gin-Gonic
2. ServeMux

## Notes

- Find out how to redirect to "/" or "/login" if cookie is not found.
- In github.com/golang-jwt/jwt/v4 StandardClaims type is deprecated, you should replace StandardClaims with RegisteredClaims.

## Storing SubClaims

- Generating Claims (with subclaims) using `jwt.MapClaims`
```js
	tokenExpireTime := time.Now().Add(1 * time.Hour) // 1 hour expiry time from now
	claims := &jwt.MapClaims {
		"Issuer": dbUsername,
		"Expiry": tokenExpireTime,
		"data": map[string]string {
			"firstname": "Jie Wei",
			"lastname": "Low",
		},
	}
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
```

- Retrieving the Subclaims
```js
	token, _ := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	claims := token.Claims.(jwt.MapClaims)
	fmt.Println(`Retrieved Issuer using claims["Issuer"]:`, claims["Issuer"])
	c.Set("username", claims["Issuer"])

	data := claims["data"].(map[string]interface{})
	firstname := data["firstname"].(string)
	lastname := data["lastname"].(string)
	fmt.Println("Retrieved Subclaims firstname and lastname:", firstname, lastname)
```

## ServeMux

- `Content-Type: application/json` indicates that the request body format is JSON