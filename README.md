# Traefik-U2F-SSO

```puml
@startuml
Browser --> Traefik : GET http://localhost:8080/whoami
Traefik -> TraefikU2F : GET http://authenticator:8080/verify
TraefikU2F -> TraefikU2F : check auth_session cookie: not logged in
Traefik <- TraefikU2F : 303 See Other /login (U2F_URL)
Browser <-- Traefik : 303 /login (U2F_URL)

Browser --> Traefik : GET /login (U2F_URL)
Browser --> Traefik : POST /webauthn/start?name=bob
Browser <-- Traefik: JSON WebAuthn meta data with challenge
Browser -> SecurityKey : Ask Authenticator with challenge to sign
Browser <- SecurityKey : Signed data
Browser --> Traefik : POST /webauthn/login/finish?name=bob with signed data
Browser <-- Traefik : JSON {"name": "bob"}

Browser -> Browser : reload
Browser --> Traefik : GET http://localhost:8080/whoami with auth_session cookie
Traefik -> TraefikU2F : GET http://authenticator:8080/verify with auth_session cookie
TraefikU2F -> TraefikU2F : check auth_session cookie: logged in
Traefik <- TraefikU2F : 200 OK

Traefik --> whoami : GET
Traefik <-- whoami : 200 Body
Browser <-- Traefik : 200 Body


@enduml
```

## Install

```
go get github.com/Tedyst/Traefik-U2F-SSO
```

## Running

```
go run .
```

This will start the server at port 8080, reachable at [localhost:8080](http://localhost:8080).

## License

MIT.
