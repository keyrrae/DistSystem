package blog.datacenter;

import java.io.IOException;
import java.util.HashMap;
import java.util.List;
import java.util.concurrent.TimeoutException;

import org.apache.log4j.Logger;

import blog.logs.EventRecord;
import blog.message.ClientRequestMessage;
import blog.message.SyncRequestMessage;
import blog.message.SyncResponseMessage;
import blog.misc.Common;
import blog.misc.MessageWrapper;

import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;
import com.rabbitmq.client.QueueingConsumer;

/**
 * Created by xuanwang on 4/12/16.
 */
public class DataCenter extends Thread {

    private TimeTable timeTable;

    private List<Post> listOfPost;
    private List<EventRecord> logs;

    // DataCenter name for routing
    private String dataCenterName;
    // Index in timetable
    private int dataCenterIndex;

    // Mapping between dataCenterName to DataCenterIndex on timetable
    private HashMap<String, Integer> dataCenterNameToIndex;

    private String logPropagationDirectQueueName;
    private String clientMessageDirectQueueName;

    // Below for connecting MQ
    private ConnectionFactory factory;
    private Connection connection;
    private Channel channel;
    // Consumer for client message and log message
    private QueueingConsumer consumer;

    private static Logger logger = Logger.getLogger(DataCenter.class);

    /**
     * 
     * Description: Connect to log propagation exchange with routing key: this.dataCenterName
     * 
     * @throws Exception
     *             void
     */
    public void bindToLogExchange() throws Exception {
        channel.exchangeDeclare(Common.LOG_DIRECT_EXCHANGE_NAME, "direct");
        channel.queueDeclare(this.logPropagationDirectQueueName, false, false, false, null);
        // 使用logPropagationDirectQueueName这个Queue绑定到Common.CLIENT_REQUEST_DIRECT_EXCHANGE_NAME这个exchange上，routing
        // key为dataCenterName
        channel.queueBind(this.logPropagationDirectQueueName, Common.LOG_DIRECT_EXCHANGE_NAME,
                this.dataCenterName);
        if (consumer == null) {
            consumer = new QueueingConsumer(channel);
        }
        channel.basicConsume(this.logPropagationDirectQueueName, true, consumer);
    }

    /**
     * 
     * Description: Connect to client message exchange with routing key: this.dataCenterName
     * 
     * @throws Exception
     *             void
     */
    public void bindToClientExchange() throws Exception {
        channel.exchangeDeclare(Common.CLIENT_REQUEST_DIRECT_EXCHANGE_NAME, "direct");
        channel.queueDeclare(this.clientMessageDirectQueueName, false, false, false, null);

        channel.queueBind(this.clientMessageDirectQueueName, Common.CLIENT_REQUEST_DIRECT_EXCHANGE_NAME,
                this.dataCenterName);
        if (consumer == null) {
            consumer = new QueueingConsumer(channel);
        }
        channel.basicConsume(clientMessageDirectQueueName, true, consumer);

    }

    public DataCenter(String dataCenterName, HashMap<String, Integer> dataCenterNameToIndex) throws IOException,
            TimeoutException {
        factory = new ConnectionFactory();
        factory.setHost(Common.MQ_HOST_NAME);
        connection = factory.newConnection();
        channel = connection.createChannel();

        this.dataCenterName = dataCenterName;
        this.dataCenterNameToIndex = dataCenterNameToIndex;
        timeTable = new TimeTable(dataCenterNameToIndex);

        this.logPropagationDirectQueueName = Common.getDatacenterLogPropagationDirectQueueName(dataCenterName);
        this.clientMessageDirectQueueName = Common.getClientMessageReceiverDirectQueueName(dataCenterName);
    }

    /**
     * Datacenter run method, wait for incoming client request
     */
    public void run() {
        MessageWrapper wrapper = null;
        QueueingConsumer.Delivery delivery;
        logger.info("Data Center: " + this.dataCenterName + " is Running");
        try {
            this.bindToClientExchange();
        } catch (Exception e) {
            logger.error(this.dataCenterName + " binding to client exchange failed");
            e.printStackTrace();
            System.exit(-1);
        }

        try {
            this.bindToLogExchange();
        } catch (Exception e1) {
            logger.error(this.dataCenterName + " binding to log exchange failed");
            e1.printStackTrace();
            System.exit(-1);
        }

        // Receive Log
        try {
            while (true) {
                delivery = consumer.nextDelivery();
                if (delivery != null) {
                    String msg = new String(delivery.getBody());
                    wrapper = Common.deserialize(msg, MessageWrapper.class);
                }
                if (wrapper != null) {
                    Class classType = wrapper.getInnerMessageClass();
                    if (classType.equals(SyncRequestMessage.class)) {
                        SyncRequestMessage syncRequestMessage = (SyncRequestMessage) wrapper.getInnerMessage();
                        logger.info("SyncRequestMessage");
                    }
                    else if (classType.equals(SyncResponseMessage.class)) {
                        SyncResponseMessage syncResponseMessage = (SyncResponseMessage) wrapper.getInnerMessage();
                        logger.info("SyncResponseMessage");
                    }
                    else if (classType.equals(ClientRequestMessage.class)) {
                        ClientRequestMessage clientRequestMessage = (ClientRequestMessage) wrapper.getInnerMessage();
                        logger.info("ClientRequestMessage");
                    }
                }
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public static void main(String[] args) throws IOException, TimeoutException {
        DataCenter dc = new DataCenter("dc1", new HashMap<String, Integer>());
        new Thread(dc).start();
    }
    // public DataCenter(String configFile) {
    //
    // BufferedReader br;
    // try {
    // br = new BufferedReader(new FileReader(configFile));
    // StringBuilder sb = new StringBuilder();
    // String line = br.readLine();
    //
    // while (line != null) {
    // String[] config = line.trim().split("\\s+");
    //
    // line = br.readLine();
    // }
    //
    // } catch (FileNotFoundException e) {
    // System.out.println("cannot find file");
    // } catch (IOException e) {
    //
    // System.out.println("IO:cannot find file");
    //
    // }
    // }

    // void onReceive() {
    //
    // }
    //
    // void sync(String hostname) {
    // // check if hostname is available
    // if (!serverMap.containsKey(hostname)) {
    // System.out.println("hostname does not exist.");
    // return;
    // }
    //
    // // send sync request to hostname
    // // receive log items from DC
    // // update log and update timetable
    // onReceive();
    //
    // }
    //
    // public void lookup(int client) {
    // // Serialize log
    // // send to client
    //
    // }
    //
    // public void lookup() {
    // log.printLog();
    // }
    //
    // public void post(String msgStr) {
    // Message msg = new Message(log, timeTable);
    //
    // }

}
