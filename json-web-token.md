# Introduction to JSON Web Tokens

[JSON Web Token](https://jwt.io/introduction)

## What is JSON Web Token?

- JSON Web Token (JWT) is an open standard (RFC 7519) that defines a compact and self-contained way for securely transmitting information between parties as a JSON object.
- This information can be verified and trusted because it is digitally signed.
- JWTs can be **signed using a secret**.

## Signed Tokens

- Signed Tokens can verify the *integrity* of the claims contained within it, while encrypted tokens hide those claims from other parties.
- When tokens are signed using public/private key pairs, the signature also certifies that only the party holding the private key is the one that signed it.

## When to use JSON Web Tokens?

---

- Authorization
    - Most common scenario for using JWT.
    - Once the user is logged in, each subsequent request will include the JWT, allowing the user to access routes, services and resources that are permitted with that token.
    - *Single Sign On* is a feature that widely uses JWT nowadays, because of its small overhead and its ability to be easily used across different domains.

---

- Information Exchange
    - JSON Web Tokens are a good way of securely transmitting information between parties.
    - Because JWTs can be signed - for example, using public/private key pairs - can be sure the senders are who they say they are.
    - Additionally, as the signature is calculated using the **header and the payload**, can also verify that the content hasn't been tampered with.

## What is the JSON Web Token Structure

- In its compact form, JSON Web Tokens consist of 3 parts separated by dots ( . ), which are: 
    - Header
    - Payload
    - Signature

### Header

---

- The header consists of 2 parts
    - The type of the token (which is JWT)
    - The signing algorithm being used such as HMAC, SHA256, or RSA.

---

- For Example, this JSON is **Base64Url** encoded to form the first part of the JWT.
```js
    {
        "alg"; "HS256",
        "typ": "JWT"
    }
```

### Payload 

---

- The second part of the token is the payload, which contains the **claims**.
- Claims are statements about an entity (typically, the user) and additional data.
- There are 3 types of claims
    - Registered Claims
        - A set of predefined claims which are not mandatory but recommended, to provide a set of useful, interoperable claims.
        - Some of them are: iss (issuer), exp (expiration time), sub (subject), aud (audience), and others.
    - Public Claims
        - Can be defined at will by those using JWTs.
        - To avoid collisions, they should be defined in the IANA JSON Web Token Registry or be defined as a URI that contains a collision resistant namespace.
    - Private Claims
        - Custom claims created to share information between parties that agree on using them.

---

- For exmaple, this payload is **Base64Url** encoded:
```js
{
  "sub": "1234567890",
  "name": "John Doe",
  "admin": true
}
```

### Signature

- To create the signature part, you have to take the encoded header, the encoded payload, a secret and the algorithm specified in the header and sign that.
- The signature is used to verify the message wasn't changed along the way and in the case of tokens signed with a private key, it can also very that the sender of the JWT is who it says it is.

---

- For example, if you want to use the HMAC SHA256 algorithm, the signature will be created in the following way:
```js
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  secret)
```

## Putting all together

- The output is 3 Base64-URL strings separated by dots that can be easily passed in HTML and HTTP environments, while being more compact when compared to XML-based standards such as SAML.

---

- The following shows a JWT that has the previous header and payload encoded, and it is signed with a secret.

```js
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.
SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

[JWT Encoded](https://jwt.io/#debugger-io)

## How do JSON Web Tokens work?

- In authentication, when the user successfully logs in using their credentials, a JSON Web Token will be returned. 
- Since tokens are credentials, great care must be taken to prevent security issues. 
- Should not keep tokens longer than required (set expiry)
- Should not store sensitive session data in browser storage due to lack of security.

---

- Whenever the user wants to access a protected route or resource, the user agent should sent the JWT, typically in the **Authorization** header using the **Bearer** schema.
- The content of the header should look like the following:
```
Authorization: Bearer <token>
```
- The server's protected routes will check for a valid JWT in the `Authorization` header, and if it's present, the user will be allowed to access protected resources.
- If the JWT contains the necessary data, the need to query the database for certain operations may be reduced, though this may not always be the case.

---

- If the token is sent in the `Authorization` header, Cross-Origin Resouce Sharing (CORS) won't be an issue as it doesn't use cookies.

