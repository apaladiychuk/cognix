How to start with Socket messagnger:

SERVER
1. Copy "Server" folder to a server
2. Edit Server.exe.config file to set up correct ip and port:
  <appSettings>
    <add key="ip" value="127.0.0.1" />
    <add key="port" value="8086" />
  </appSettings>
3. Open port on a server for incoming TCP connections
4. Run Server.exe

Clients
1. Copy "Client" folder to a client environment
2. Edit Client.exe.config file to set up correct ip and port - it should machh to the server configuration:
  <appSettings>
    <add key="ip" value="127.0.0.1" />
    <add key="port" value="8086" />
  </appSettings>
3. Run Client.exe

