# Go Web Template

This is base template for doing spikes for smaller applications. It uses the following:

- Golang Templates for UI
- sqlite3 for datastore (can use Postgres)
- Generic oauth2 login (tested with Google)

## Quick Start

Copy the `.env.template` file and change any values you like. You can see the possible values by looking at `internals/models/config.go`.

Once the settings are in place you can do:

```
make start
```

and then browse to `http://localhost:3000`.

```
.
├── cmd
│   └── server               <= Main entry point
├── docker-compose.yaml
├── Dockerfile
├── internals                <= Your App code
│   ├── handlers             <= HTTP handlers
│   ├── models               <= App structs
│   └── repository           <= Database queries
├── migrations               <= SQL migrations
│   └── 000000-init.sql
├── static                   <= Images, css, etc
│   └── plant-research.png
├── templates                <= HTML Pages
│   └── home.html
└── upload-temp              <= For uploaded files

```

## OAuth2 Setup

### Google

- [Setup on Google Console](https://console.cloud.google.com/apis/dashboard)
- [How to](https://support.google.com/cloud/answer/6158849?hl=en)

```bash
WB_AUTH_REDIRECT_URL=http://localhost:3000/callback
WB_AUTH_CLIENT_ID=xxxxxxxxxxxxxxxxxxxxxxx.apps.googleusercontent.com
WB_AUTH_CLIENT_SECRET=xxxxx-xx-xxxxxxxxxx
WB_AUTH_SCOPES=email openid https://www.googleapis.com/auth/userinfo.email
WB_AUTH_AUTH_URL=https://accounts.google.com/o/oauth2/auth
WB_AUTH_TOKEN_URL=https://oauth2.googleapis.com/token
WB_AUTH_AUTH_STYLE=1
WB_AUTH_ACCESS_TOKEN_URL=https://www.googleapis.com/oauth2/v2/userinfo?access_token=
```

### Auth0

```bash
WB_AUTH_REDIRECT_URL=http://localhost:3000/callback
WB_AUTH_CLIENT_ID=xxxxxxxxxxxxxxxxxxxxxxxx
WB_AUTH_CLIENT_SECRET=xxxxxxxxxxxxxxxxxxxxx
WB_AUTH_SCOPES=email openid
WB_AUTH_AUTH_URL=https://{auth0_domain}/authorize
WB_AUTH_TOKEN_URL=https://{auth0_domain}/oauth/token
WB_AUTH_AUTH_STYLE=1
WB_AUTH_ACCESS_TOKEN_URL=https://{auth0_domain}/userinfo?access_token=
```

## sqlite3 vs Postgres

The code can support both sqlite3 or postgres. By default it uses sqlite3, but if you look at `start_db` in the Makefile and the example values in `.env.template` you can see how to get Postgres working.
