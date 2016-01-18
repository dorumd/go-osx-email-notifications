# OS X Email Notifications v0.1.0 built using GoLang.

# Running
- Install GoLang (with v.1.4 tested)
- Install required packaged by running:
`go get golang.org/x/oauth2`
`go get golang.org/x/net/context`
`go get google.golang.org/api/gmail/v1`
`go get github.com/deckarep/gosx-notifier`
- Create `client_secret.json` file inside config directory
- And execute `go build main.go && ./main`
- At the first run you'll have to open the generated link and then paste your token to allow using your gmail account

# TODO
- Fix email url
- Run it on startup
- Use GoRoutines for processing at the same time multiple emails
