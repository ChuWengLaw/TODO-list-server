# TODO-list-server
A simple TODO-list server, where users can :-
1. Sign in using any one of: gmail, facebook or github login.
2. Add a TODO item.
3. Delete a TODO item.
4. List all TODO items.
5. Mark a TODO item as completed.

# Starting the service
## Navigate to the directory containing the docker-compose.yml file in your terminal and run the following command:-
**docker-compose up**
## Or to run the service in the background
**docker-compose up -d**

# Stoping the service
## In the same terminal
Ctrl+C to interrupt the service
## If it was running in daemon
**docker-compose down**

# Upon server running
## Invoke the following command to get an appropriate response object for each of the APIs (Add/Delete/List/Mark-complete)
**curl -H any_appropriate_auth_token http://host:port/api_path/params**
## Or simply insert http://host:port/ in browser url and follow the instructions there
