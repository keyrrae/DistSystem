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
class DataCenter extends Thread{
    private HashMap<Long, Long> timeVector;
    private ArrayList<ArrayList<Long>> timeTable;
    private Log log;
    private String localIP;
    private HashMap<String, InetSocketAddress> serverMap;
    private JSONsender jsonSender;
    private JSONreceiver jsoReceiver;

    DataCenter(String configFile){

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

    public void run(){

    }

    void onReceive(){

    }

    void sync(String hostname){
        // check if hostname is available
        if(!serverMap.containsKey(hostname)){
            System.out.println("hostname does not exist.");
            return;
        }





        // send sync request to hostname
        // receive log items from DC
        // update log and update timetable
        onReceive();

    }

    public void lookup(int client){
        // Serialize log
        // send to client

    }

    public void lookup(){
        log.printLog();
    }


    public void post(String msgStr){
        Message msg = new Message(log, timeTable);


    }



}
