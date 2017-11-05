script.verb = "auth:logout"
script.help = "destroys credentials for heroku's API server."
script.author = "Xe"
script.version = "0.1"
script.usage = ""

local netrc = require "netrc"

function run()
  netrc.removeMachine("api.heroku.com")
  netrc.removeMachine("git.heroku.com")
  netrc.save()
end
