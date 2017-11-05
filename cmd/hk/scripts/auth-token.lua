script.verb = "auth:token"
script.help = "prints your current authentication token."
script.author = "Xe"
script.version = "0.1"
script.usage = ""

local heroku = require "heroku"

function run()
  local user, pass = heroku.get_userpass()
  print(pass)
end
