import java.io.BufferedReader;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.net.InetSocketAddress;
import java.util.ArrayList;
import java.util.HashMap;

import com.rabbitmq.client.ConnectionFactory;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.Channel;


/**
 * Created by xuanwang on 4/12/16.
 */
public class DataCenter extends Thread{
    private HashMap<Long, Long> timeVector;
    private ArrayList<ArrayList<Long>> timeTable;
    private Log log;
    private String localIP;
    private HashMap<String, InetSocketAddress> serverMap;

    public DataCenter(String configFile){

        serverMap = new HashMap<String, InetSocketAddress>();

        BufferedReader br;
        try {
            br = new BufferedReader(new FileReader(configFile));
            StringBuilder sb = new StringBuilder();
            String line = br.readLine();

            while (line != null) {
                String[] config = line.trim().split("\\s+");
                final InetSocketAddress server = new InetSocketAddress(config[1], Integer.parseInt(config[2]));
                serverMap.put(config[0], server);

                line = br.readLine();
            }

        } catch(FileNotFoundException e){
            System.out.println("cannot find file");
        } catch(IOException e){

            System.out.println("IO:cannot find file");

        }

        timeVector = new HashMap<Long, Long>();
        timeTable = new ArrayList<ArrayList<Long>>();
        log = new Log();
    }

    public void sync(String hostname){
        //send sync request to hostname
        //receive log items from DC
        //update log
        //update timetable

    }

    public void lookup(int client){
        // Serialize log
        // send to client
    }

    public void lookup(){

    }

    public void run(){

    }

    public void post(String message){

    }

    private static void processMessage(){

    }

}
