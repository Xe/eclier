script.verb = "apps:detect"
script.help = "prints the app name as detected by git commands"
script.author = "Xe"
script.version = "0.1"
script.usage = ""

local heroku = require "heroku"

function run(arg)
  print(heroku.detect_app())
end
