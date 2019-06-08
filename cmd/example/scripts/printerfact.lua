script.verb = "printerfact"
script.help = "displays useful and informative facts about printers"
script.author = "Xe"
script.version = "1.0"
script.usage = ""

local http = require "http"
local json = require "json"

function run(arg)
  local resp, err = http.get "https://xena.stdlib.com/printerfacts"

  if err ~= nil then
    error(err)
  end

  local obj = json.decode(resp.body)

  print(obj.facts[1])
end
