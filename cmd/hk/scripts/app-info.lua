script.verb = "app:info"
script.help = "looks up information on an application"
script.author = "Xe"
script.version = "0.1"
script.usage = "[-app app-name]"

local flag = require "flag"
local heroku = require "heroku"

local fs = flag.new()

fs:string("app", "", "application name")

script.usage = fs:usage()

function run(arg)
  if arg[1] == "-help" or arg[1] == "--help" then
    print(fs:usage())
    return
  end

  arg[0] = script.verb
  local flags = fs:parse(arg)

  if flags.app == "" then
    print("-app must be specified")
    return
  end

  local app = heroku.app_info(flags.app)

  print("app name: " .. tostring(app.name))
  print("app url: " .. tostring(app.weburl))
  print("app created at: " .. app.createdAt:string())
end
