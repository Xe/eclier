script.verb = "auth:logout"
script.help = "destroys credentials for heroku's API server."
script.author = "Xe"
script.version = "0.1"
script.usage = ""

function run()
  netrc:removeMachine("api.heroku.com")
  netrc:removeMachine("git.heroku.com")
end
