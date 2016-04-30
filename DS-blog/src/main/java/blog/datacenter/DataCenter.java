package blog.datacenter;

import static java.lang.System.exit;

import java.io.IOException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.PriorityQueue;
import java.util.Scanner;
import java.util.concurrent.TimeoutException;

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

    TimeTable timeTable;

    // private PriorityQueue<Post> listOfPost;
    private PriorityQueue<EventRecord> logs;

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

        this.dataCenterName = dataCenterName;
        this.dataCenterNameToIndex = dataCenterNameToIndex;
        this.dataCenterIndex = dataCenterNameToIndex.get(dataCenterName);

        timeTable = new TimeTable(dataCenterNameToIndex, dataCenterName);
        logs = new PriorityQueue<EventRecord>();

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
     * 
     * Description: Used when decide what logs to send to other data center k.
     * 
     * @param table
     * @param e
     * @param kIndex
     * @return
     *         boolean
     */
    public boolean hasRec(EventRecord e, int kIndex) {
        int eventDatacenterIndex = dataCenterNameToIndex.get(e.getNodeName());
        return this.timeTable.ifDatacenterXKnowDatacenterYTillTimeT(kIndex, eventDatacenterIndex, e.getTimestamp());
    }

    /**
     * Description: TODO
     * 
     * @param message
     *            void
     */
    private void handleClientRequestPostMessage(ClientRequestPostMessage message) {
        timeTable.increaseLocalClock();
        logs.add(new EventRecord(timeTable.getLocalClock(), this.dataCenterName, message.getPost()));
    }

    /**
     * 
     * Description: TODO
     * 
     * @param message
     * @throws IOException
     *             void
     */
    private void handleClientRequestSyncMessage(ClientRequestSyncMessage message) throws IOException {
        SyncRequestMessage syncRequestMessage = new SyncRequestMessage(message.getFromDataCenterName(),
                message.getToDataCenterName());
        channel.basicPublish(
                Common.LOG_DIRECT_EXCHANGE_NAME,
                syncRequestMessage.getFromDataCenterName(),
                null,
                Common.serialize(
                        new MessageWrapper(Common.serialize(syncRequestMessage), syncRequestMessage.getClass()))
                        .getBytes());
        logger.info("Sending sync request to " + syncRequestMessage.getFromDataCenterName());
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
                this.logs);
        logger.info("Datacenter send lookup response");
        this.sendResponseToClient(responseMessage);
    }

    /**
     * 
     * Description: Datacenter handle other datacenter's sync request message. SyncRequestMessage's toDataCenterName is
     * the
     * data center that send SyncRequestMessage
     * 
     * @param message
     * @throws IOException
     *             void
     */
    private void handleDataCenterSyncRequestMessage(SyncRequestMessage message) throws IOException {

        String destDataCenterName = message.getToDataCenterName();
        int destDataCenterIndex = dataCenterNameToIndex.get(destDataCenterName);

        SyncResponseMessage responseMessage = new SyncResponseMessage(this.dataCenterName, destDataCenterName);

        responseMessage.setTimeTable(timeTable);

        // Filter what message to be sent out
        List<EventRecord> logsToBeSent = new ArrayList<EventRecord>();
        for (EventRecord e : logs) {
            if (!hasRec(e, destDataCenterIndex)) {
                logsToBeSent.add(e);
            }
        }
        responseMessage.setLogs(logsToBeSent);
        sendSyncResponseToDataCenter(responseMessage);
        logger.info("Sending sync response message to " + responseMessage.getToDataCenterName());
    }

    private void handleDataCenterResponseMessage(SyncResponseMessage message) {

        logger.info("Got message with " + message.getLogs().size() + " logs");
        // filter log to be added to own logs
        applyReceivedLogs(message.getLogs());
        // update time table
        updateTimeTable(dataCenterNameToIndex.get(message.getFromDataCenterName()), message.getTimeTable());

    }

    /**
     * Description: TODO
     * 
     * @param log
     *            void
     */
    private void applyReceivedLogs(List<EventRecord> log) {
        int count = 0;
        for (EventRecord l : log) {
            int logSourceDatacenterIndex = dataCenterNameToIndex.get(l.getNodeName());
            // Add logs that don't know
            if (timeTable.getTable()[dataCenterIndex][logSourceDatacenterIndex] < l.getTimestamp()) {
                logs.add(l);
                count++;
            }
        }
        logger.info(count + " logs are appended.");
    }

    /**
     * 
     * Description: Update local timeTable using the timeTabel received
     * 
     * @param recvTableDatacenterIndex
     * @param recv
     *            void
     */
    public void updateTimeTable(int recvTableDatacenterIndex, TimeTable recv) {
        // Should make local clock sync (one more than) received message source's clock?
        // this.increaseLocalClock();
        long[][] table = this.timeTable.getTable();
        table[dataCenterIndex][dataCenterIndex] = Math.max(
                recv.getTable()[recvTableDatacenterIndex][recvTableDatacenterIndex],
                table[dataCenterIndex][dataCenterIndex]) + 1;

        int n = table.length;
        int m = recv.size();

        if (n != m) {
            System.out.println("different time table dimensions");
            exit(1);
        }
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < n; j++) {
                table[i][j] = table[i][j] > recv.getTable()[i][j] ? table[i][j] : recv.getTable()[i][j];
            }
        }

        int recvIndex = recv.getDataCenterIndex();
        int thisIndex = dataCenterIndex;

        for (int j = 0; j < n; j++) {
            table[thisIndex][j] = table[thisIndex][j] > recv.getTable()[recvIndex][j] ?
                    table[thisIndex][j] : recv.getTable()[recvIndex][j];
        }
        String s = "";
        for (int i = 0; i < table.length; i++)
        {
            for (int j = 0; j < table[0].length; j++) {
                s += timeTable.getTable()[i][j];
            }
            s += "\n";
        }
        logger.info("Updated timetable: " + s);
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

    private void sendSyncResponseToDataCenter(SyncResponseMessage message) throws IOException {
        channel.basicPublish(Common.LOG_DIRECT_EXCHANGE_NAME, message.getToDataCenterName(), null,
                Common.serialize(new MessageWrapper(Common.serialize(message), message.getClass())).getBytes());
    }

    public static void main(String[] args) throws IOException, TimeoutException {
        Scanner scan = new Scanner(System.in);
        System.out.println("Please enter current datacenter name");
        String myName = scan.nextLine().trim();
        System.out.println("Current DataCenter Name:" + myName);

        System.out.println("Please enter the num of datacenters");
        int numOfDataCenters = Integer.parseInt(scan.nextLine().trim());

        HashMap<String, Integer> dataCenterNameIndexMap = new HashMap<String, Integer>();
        System.out
                .println("Please enter "
                        + numOfDataCenters
                        + " datacenter names(including yourself). ");
        System.out
                .println("Please make sure the names and their order is exactly the same at every datacenter console.");
        for (int i = 0; i < numOfDataCenters; i++) {
            System.out.print("Data Center " + i + ":");
            String dataCenterName = scan.nextLine().trim();
            dataCenterNameIndexMap.put(dataCenterName, i);
        }

        System.out.println("Please enter MQ address for communication(Or empty for default domain:rabbitmq)");
        String host = scan.nextLine().trim();
        if (host != null && host.length() != 0) {
            Common.MQ_HOST_NAME = host;
        }
        System.out.println("DataCenter " + myName + " started...");
        
        DataCenter dc = new DataCenter(myName, dataCenterNameIndexMap);
        new Thread(dc).start();
    }
}
