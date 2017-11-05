script.verb = "apps:create"
script.help = "create a new heroku application"
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

  local app = heroku.app_create(app)

  print("=== " .. tostring(app.name))
  print("Auto Cert Mgmt: " .. tostring(app.Acm))
  print("Git URL:        " .. tostring(app.GitURL))
  print("Owner:          " .. tostring(app.Owner.Email))
  print("Region:         " .. tostring(app.Region.Name))
  print("Repo Size:      " .. tostring(app.RepoSize))
  print("Slug Size:      " .. tostring(app.SlugSize))
  print("Stack:          " .. tostring(app.Stack.Name))
  print("Web URL:        " .. tostring(app.WebURL))
end
