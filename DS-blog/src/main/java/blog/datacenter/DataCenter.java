package blog.datacenter;

import java.io.IOException;
import java.util.HashMap;
import java.util.List;
import java.util.PriorityQueue;
import java.util.concurrent.TimeoutException;

import org.apache.log4j.Logger;

import blog.logs.EventRecord;
import blog.message.center2center.SyncRequestMessage;
import blog.message.center2center.SyncResponseMessage;
import blog.message.client2center.CenterResponseLookUpMessage;
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
                        SyncRequestMessage message = (SyncRequestMessage) wrapper.getInnerMessage();
                        logger.info("SyncRequestMessage");
                    }
                    else if (classType.equals(SyncResponseMessage.class)) {
                        SyncResponseMessage message = (SyncResponseMessage) wrapper.getInnerMessage();
                        logger.info("SyncResponseMessage");
                    }
                    else if (classType.equals(ClientRequestLookUpMessage.class)) {
                        ClientRequestLookUpMessage message = (ClientRequestLookUpMessage) wrapper.getInnerMessage();
                        logger.info("ClientRequestLookUpMessage");
                        handleClientRequestLookUpMessage(message);

                    }
                    else if (classType.equals(ClientRequestPostMessage.class)) {
                        ClientRequestPostMessage message = (ClientRequestPostMessage) wrapper.getInnerMessage();
                        logger.info("ClientRequestPostMessage");
                    }
                    else if (classType.equals(ClientRequestSyncMessage.class)) {
                        ClientRequestSyncMessage message = (ClientRequestSyncMessage) wrapper.getInnerMessage();
                        logger.info("ClientRequestSyncMessage");
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
    private void handleClientRequestLookUpMessage(ClientRequestLookUpMessage message) {
        String clientName = message.getClientName();
        CenterResponseLookUpMessage responseMessage = new CenterResponseLookUpMessage(clientName, this.dataCenterName,
                this.listOfPost);

    }

    public void sendResponseToClient(Message m) {
        channel.basicPublish(Common.CLIENT_REQUEST_DIRECT_EXCHANGE_NAME, dataCenterName, null,
                Common.serialize(new MessageWrapper(Common.serialize(message), message.getClass())).getBytes());
    }

    public static void main(String[] args) throws IOException, TimeoutException {
        DataCenter dc = new DataCenter("dc1", new HashMap<String, Integer>());
        new Thread(dc).start();
    }

}
