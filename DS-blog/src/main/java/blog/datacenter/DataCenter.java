package blog.datacenter;

import java.io.IOException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.PriorityQueue;
import java.util.concurrent.TimeoutException;

import com.sun.tools.corba.se.idl.InterfaceGen;
import org.apache.log4j.BasicConfigurator;
import org.apache.log4j.Logger;

import blog.logs.EventRecord;
import blog.message.center2center.SyncRequestMessage;
import blog.message.center2center.SyncResponseMessage;
import blog.message.client2center.CenterResponseLookUpMessage;
import blog.message.client2center.ClientDataCenterMessage;
import blog.message.client2center.ClientRequestLookUpMessage;
import blog.message.client2center.ClientRequestPostMessage;
import blog.message.client2center.ClientRequestSyncMessage;
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

    private PriorityQueue<Post> listOfPost;
    private List<EventRecord> logs;
    private long timeStamp;

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

    public DataCenter(String dataCenterName, HashMap<String, Integer> dataCenterNameToIndex) throws IOException,
            TimeoutException {
        BasicConfigurator.configure();
        factory = new ConnectionFactory();
        factory.setHost(Common.MQ_HOST_NAME);
        connection = factory.newConnection();
        channel = connection.createChannel();
        timeStamp = 0;

        this.dataCenterName = dataCenterName;
        this.dataCenterNameToIndex = dataCenterNameToIndex;
        timeTable = new TimeTable(dataCenterNameToIndex, dataCenterName);
        listOfPost = new PriorityQueue<Post>();
        logs = new ArrayList<EventRecord>();

        this.logPropagationDirectQueueName = Common.getDatacenterLogPropagationDirectQueueName(dataCenterName);
        this.clientMessageDirectQueueName = Common.getClientMessageReceiverDirectQueueName(dataCenterName);
    }

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
                        SyncRequestMessage message = (SyncRequestMessage) wrapper.getInnerMessage();
                        logger.info("Received SyncRequestMessage");
                        handleDataCenterSyncRequestMessage(message);
                    }

                    else if (classType.equals(SyncResponseMessage.class)) {
                        SyncResponseMessage message = (SyncResponseMessage) wrapper.getInnerMessage();
                        logger.info("Received SyncResponseMessage");
                        handleDataCenterResponseMessage(message);
                    }
                    else if (classType.equals(ClientRequestLookUpMessage.class)) {
                        ClientRequestLookUpMessage message = (ClientRequestLookUpMessage) wrapper.getInnerMessage();
                        logger.info("ClientRequestLookUpMessage");
                        handleClientRequestLookUpMessage(message);

                    }
                    else if (classType.equals(ClientRequestPostMessage.class)) {
                        ClientRequestPostMessage message = (ClientRequestPostMessage) wrapper.getInnerMessage();
                        logger.info("ClientRequestPostMessage");
                        handleClientRequestPostMessage(message);
                    }
                    else if (classType.equals(ClientRequestSyncMessage.class)) {
                        ClientRequestSyncMessage message = (ClientRequestSyncMessage) wrapper.getInnerMessage();
                        logger.info("ClientRequestSyncMessage");
                        handleClientRequestSyncMessage(message);
                    }
                }
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    /**
     * Description: TODO
     * 
     * @param message
     *            void
     */
    private void handleClientRequestPostMessage(ClientRequestPostMessage message) {
        timeTable.increaseLocalClock();
        Post post = message.getPost();
        post.setTimeStamp(timeTable.getLocalClock());
        listOfPost.add(post);
    }

    /**
     * Description: TODO
     * 
     * @param message
     *            void
     * @throws IOException
     */
    private void handleClientRequestLookUpMessage(ClientRequestLookUpMessage message) throws IOException {
        String clientName = message.getClientName();
        CenterResponseLookUpMessage responseMessage = new CenterResponseLookUpMessage(clientName, this.dataCenterName,
                this.listOfPost);
        this.sendResponseToClient(responseMessage);
    }

     /**
     * 
     * Description: Send response message to client
     * 
     * @param message
     * @throws IOException
     *             void
     */
    public void sendResponseToClient(ClientDataCenterMessage message) throws IOException {
        channel.basicPublish(Common.CLIENT_REQUEST_DIRECT_EXCHANGE_NAME, message.getClientName(), null,
                Common.serialize(new MessageWrapper(Common.serialize(message), message.getClass())).getBytes());
    }

    private void handleDataCenterSyncRequestMessage(SyncRequestMessage message) throws IOException{

        String destDataCenterName = message.getToDataCenterName();
        SyncResponseMessage responseMessage = new SyncResponseMessage(this.dataCenterName, destDataCenterName);

        responseMessage.setTimeTable(timeTable);
        responseMessage.setLog(logs);
        sendSyncResponseToDataCenter(responseMessage);
        logger.info("Sending sync response message to " + responseMessage.getToDataCenterName());
    }

    private void handleDataCenterResponseMessage(SyncResponseMessage message){
        // TODO
        // update time table

        timeTable.upDateUponReceived(message.getTimeTable());
        // update log
    }

    private void sendSyncResponseToDataCenter(SyncResponseMessage message) throws IOException{
        channel.basicPublish(Common.LOG_DIRECT_EXCHANGE_NAME,  message.getToDataCenterName(), null,
                Common.serialize(new MessageWrapper(Common.serialize(message), message.getClass())).getBytes());
    }

    private void handleClientRequestSyncMessage(ClientRequestSyncMessage message) throws IOException{
        SyncRequestMessage syncRequestMessage = new SyncRequestMessage(message.getFromDataCenterName(),
                                                        message.getToDataCenterName());
        channel.basicPublish(Common.LOG_DIRECT_EXCHANGE_NAME,  syncRequestMessage.getFromDataCenterName(), null,
                Common.serialize(new MessageWrapper(Common.serialize(syncRequestMessage), syncRequestMessage.getClass())).getBytes());
        logger.info("Sending sync request to " + syncRequestMessage.getFromDataCenterName());
    }

    public static void main(String[] args) throws IOException, TimeoutException {

        HashMap<String, Integer> dataCenterNameIndexMap = new HashMap<String, Integer>();
        dataCenterNameIndexMap.put("dc1", 0);
        dataCenterNameIndexMap.put("dc2", 1);
        DataCenter dc1 = new DataCenter("dc1", new HashMap<String, Integer>());
        DataCenter dc2 = new DataCenter("dc2", new HashMap<String, Integer>());
        new Thread(dc1).start();
        new Thread(dc2).start();
    }

}
