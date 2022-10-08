# Go JWT Login Authentication Methods

1. Gin-Gonic
2. ServeMux

## Notes

- Find out how to redirect to "/" or "/login" if cookie is not found.
- In github.com/golang-jwt/jwt/v4 StandardClaims type is deprecated, you should replace StandardClaims with RegisteredClaims.