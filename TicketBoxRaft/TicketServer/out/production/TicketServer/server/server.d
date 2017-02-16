module server.server;

import std.stdio;
import server.config;

class Server {
    Config config;

    this() {
       this.config = new Config();
    }
    void run(){
        writeln("run from Server");
    }

}