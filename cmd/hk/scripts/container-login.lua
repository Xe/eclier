script.verb = "container:login"
script.help = "logs the currently logged in heroku user to the heroku docker registry"
script.author = "Xe"
script.version = "0.1"
script.usage = ""

local netrc = require "netrc"
local sh = require "sh"

function run()
  local pass = netrc.machine("api.heroku.com"):get("password")

  sh.docker("login", "--username=_", "--password", pass, "registry.heroku.com")

  print("logged into registry.heroku.com sucessfully")
end
