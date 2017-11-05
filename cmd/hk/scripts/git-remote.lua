script.verb = "git:remote"
script.help = "adds a git remote to an app repo"
script.author = "Xe"
script.version = "0.1"
script.usage = "[-app app-name]"

local flag = require "flag"
local question = require "question"
local sh = require "sh"

local fs = flag.new()

fs:string("app", "", "application name")

function run()
  if arg[1] == "-help" or arg[1] == "--help" then
    print(fs:usage())
    return
  end

  arg[0] = script.verb
  local flags = fs:parse(arg)
  if flags.app == "" then
    flags.app = question.ask "app name? "
  end

  sh.git("remote", "add", "heroku", "https://git.heroku.com/" .. flags.app .. ".git")
  print("heroku git remote added for app " .. flags.app)
end
