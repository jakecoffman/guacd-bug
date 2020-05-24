# guacd-bug

## Error message:

```
19:55:16.113 [Thread-105] ERROR o.a.g.w.GuacamoleWebSocketTunnelEndpoint - Connection to guacd terminated abnormally: Element terminator of instruction was not ';' nor ','
```

Resulting in websocket disconnect.

## Steps to reproduce issue:

- Start guacd `docker run --name guacd -it --rm -p 4822:4822 guacamole/guacd`
- Start postgres `docker run --name some-postgres -e POSTGRES_PASSWORD=pw -e POSTGRES_USER=postgres -e POSTGRES_DB=guacamole_db -p 5432:5432 -d postgres`
- Generate init.db `docker run --rm guacamole/guacamole /opt/guacamole/bin/initdb.sh --postgres > initdb.sql`
- Run initdb.sql on guacamole_db database
- Start guacamole client `docker run --name my-guac --link guacd:guacd --link some-postgres:postgres -e POSTGRES_DATABASE=guacamole_db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=pw -p 8080:8080 guacamole/guacamole`
- Connect http://localhost:8080/guacamole/#/
- Login with guacadmin/guacadmin
- Configure an RDP session as usual, but be sure to set Maximum connections to 100 and Maximum per user to 100 too.
- Save and go Home
- Click on the connection to start the first connection
- In a new window/browser login to guacamole and go to settings
- Under active sessions you should see your other active session. Open developer tools and then click on the Connection Name link to connect to it.
- Copy the websocket connection from the network tab and paste it into the Go script (replace url).
- Also grab the Cookie header and paste it in the script (replace cookie). 
- Start the script `go run main.go`
- You should see "Diaing 0" etc for a while, sometimes it happens right away sometimes it takes a long time.
