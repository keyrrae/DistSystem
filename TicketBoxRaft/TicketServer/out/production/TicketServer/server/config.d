module server.config;

import std.stdio;
import std.json;

class Config {
    string address;

    this(string filename){

        File file = File(filename, "r");
        string lines;

        while (!file.eof()) {
            lines = lines ~ (file.readln());
        }
        file.close();

        JSONValue json = parseJSON(lines);
        this.address = json["address"].str;
        writeln(this.address);
    }
}
