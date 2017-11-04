script.verb = "auth:login"
script.help = "sets up initial credentials with heroku's API server."
script.author = "Xe"
script.version = "0.1"
script.usage = ""

local question = require "question"

function run()
  local email = question.ask "Heroku Email address: "
  local pass = question.secret "Heroku password (never stored): "

  set_heroku_userpass(email, pass)

  tkn, err = heroku_new_token()
  if err ~= nil then
    error(err:error())
  end

  netrc:machine("api.heroku.com"):set("login", email)
  netrc:machine("api.heroku.com"):set("password", tkn)
  netrc:machine("git.heroku.com"):set("login", email)
  netrc:machine("git.heroku.com"):set("password", tkn)

  netrc:save()

  print("Credentials saved for " .. email)
end