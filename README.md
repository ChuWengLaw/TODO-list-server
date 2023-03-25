# TODO-list-server
A simple TODO-list server, where users can :-
1. Sign in using any one of: gmail, facebook or github login.
2. Add a TODO item.
3. Delete a TODO item.
4. List all TODO items.
5. Mark a TODO item as completed.

# Cloning the project
## Open git bash or any git gui extensions
**git clone https://github.com/ChuWengLaw/TODO-list-server.git**

# Starting the service
## Navigate to the directory containing the docker-compose.yml file from the cloned project in your terminal and run the following command:-
**docker-compose up**
## Or to run the service in the background
**docker-compose up -d**

# Stoping the service
## In the same terminal
**Ctrl+C to interrupt the service**
## If it was running in daemon
**docker-compose down**

# Upon server running
## Login to copy your access token
**curl -g -v http://localhost:8080/Login?method={1/2/3}**
## Invoke the following command to get an appropriate response object for each of the APIs depending on your system
**(powershell) curl -H @{'token' = 'your_token'} -v http://localhost:8080/(Add/Delete/List/Mark-complete)?params**
**(cmd) curl -H 'token: your_token' -g -v http://localhost:8080/(Add/Delete/List/Mark-complete)?params**