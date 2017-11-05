script.verb = "apps:destroy"
script.help = "permanently destroy a heroku application"
script.author = "Xe"
script.version = "0.1"
script.usage = "[app-name]"

local heroku = require "heroku"
local question = require "question"

function run(arg)
  local app = ""
  if #arg == 0 then
    app = question.ask "app name? (leave blank for auto-generated name) "
  else
    app = arg[1]
  end

  while true do
    local confirm = question.ask "Please confirm the application name: "
    if confirm == app then
      break
    end
  end

  heroku.app_destroy(app)

  print("application " .. app .. " destroyed")
end
