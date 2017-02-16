import std.stdio;
import server.server;

void main()
{
    auto server = new Server();
    server.run();
	writeln("Edit source/app.d to start your project.");
}
