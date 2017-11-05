script.verb = "container:push"
script.help = "pushes a prevously tagged docker image to a heroku app"
script.author = "Xe"
script.version = "0.1"
script.usage = "[-app app-name]"

local flag = require "flag"
local question = require "question"
local sh = require "sh"

local fs = flag.new()
fs:string("app", "", "heroku application name")

function run(arg)
  if arg[1] == "-help" or arg[1] == "--help" then
    print(fs:usage())
    return
  end

  arg[0] = script.verb
  local flags = fs:parse(arg)

  if flags.app == "" then
    flags.app = question.ask "Heroku app name? "
  end

  sh.docker("push", "registry.heroku.com/" .. flags.app .. "/web")
end
