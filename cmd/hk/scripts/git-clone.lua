script.verb = "git:clone"
script.help = "clones a heroku app to your local machine at DIRECTORY (defaults to app name)"
script.author = "Xe"
script.version = "0.1"
script.usage = "[-app app-name] [DIRECTORY]"

local flag = require "flag"
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
  local directory = flags[1]

  if directory ~= nil then
    sh.git("clone", "https://git.heroku.com/" .. flags.app .. ".git", directory)
  else
    sh.git("clone", "https://git.heroku.com/" .. flags.app .. ".git")
  end
end
