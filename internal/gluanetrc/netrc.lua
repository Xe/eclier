--- Module netrc offers a simple interface to a user's netrc file in their home directory.
-- @module netrc

local netrc = {}

--- add_machine adds a machine to the netrc manifest with a username and password.
-- @param name string the domain name of the machine
-- @param login string the user name to log in as
-- @param password string the password or similar secret for the machine
-- @return Machine
function netrc.add_machine(name, login, password)
end

--- machine loads netrc data for a given machine by domain name.
-- Any changes made with the `set` method of a machine will be saved to the disk
-- when the module's `save` function is called. If the given machine does not
-- exist in the netrc file, this function will return nil.
-- @param name string
-- @return Machine
-- @usage local creds = netrc.machine("api.foobar.com")
-- @usage print(creds:get("username"), creds:get("password"))
function netrc.machine(name)
  return nil
end

--- remove_machine removes a single machine from the netrc manifest by name.
-- @param name string the name of the machine to remove from the netrc manifest
-- @usage netrc.remove_machine("api.digg.com")
function netrc.remove_machine(name)
end

--- save writes all changes made in machine `set` methods to the disk at $HOME/.netrc.
-- This function will raise a lua error if the save fails. This function should
-- not fail in the course of normal operation.
-- @usage netrc.save()
function netrc.save()
end

--- Machine is a userdata wrapper around the go netrc.Machine type.
-- https://godoc.org/github.com/dickeyxxx/netrc#Machine
-- @type Machine

local Machine = {}

--- get gets a Machine value by key.
-- @param key the netrc key to get
-- @return string the value from the netrc
-- @usage local cli = api.new(m:get("login"), m:get("password"))
function Machine:get(key)
end

--- set updates information in this Machine by a key, value pair.
-- @param key the netrc key to set
-- @param value the value to set the above key to
-- @usage m:set("password", "hunter2")
function Machine:set(key, value)
end

return netrc
