script.verb = "auth:token"
script.help = "prints your current authentication token."
script.author = "Xe"
script.version = "0.1"
script.usage = ""

function run()
  print(netrc:machine("api.heroku.com"):get("password"))
end
