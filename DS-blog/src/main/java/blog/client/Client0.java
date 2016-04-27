package blog.client;

import java.io.IOException;
import java.util.PriorityQueue;
import java.util.Scanner;
import java.util.concurrent.TimeoutException;

import blog.logs.EventRecord;
import blog.message.client2center.*;

import org.apache.log4j.BasicConfigurator;
import org.apache.log4j.Logger;

import blog.datacenter.Post;
import blog.misc.Common;
import blog.misc.MessageWrapper;

import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;
import com.rabbitmq.client.QueueingConsumer;

import static java.lang.System.exit;

/**
 * Created by xuanwang on 4/12/16.
 */
public class Client0 implements Runnable {
    private static Logger logger = Logger.getLogger(Client0.class);
    private ConnectionFactory factory;
    private Connection connection;
    private Channel channel;
    // UNIQUE routing key
    private String clientName;
    private String datacenterFeedbackMessageReceiverDirectQueueName;
    private QueueingConsumer consumer;
    private QueueingConsumer.Delivery delivery;

    public void bindToDatacenterMessageReceiverQueue() throws Exception {
        channel.exchangeDeclare(Common.CLIENT_REQUEST_DIRECT_EXCHANGE_NAME, "direct");
        channel.queueDeclare(this.datacenterFeedbackMessageReceiverDirectQueueName, false, false, false, null);
        // 使用datacenterFeedbackMessageReceiverDirectQueueName这个Queue绑定到Common.CLIENT_REQUEST_DIRECT_EXCHANGE_NAME这个exchange上，routing
        // key为clientName
        channel.queueBind(this.datacenterFeedbackMessageReceiverDirectQueueName,
                Common.CLIENT_REQUEST_DIRECT_EXCHANGE_NAME,
                this.clientName);
        if (consumer == null) {
            consumer = new QueueingConsumer(channel);
        }
        channel.basicConsume(this.datacenterFeedbackMessageReceiverDirectQueueName, true, consumer);
    }

    public Client0(String clientName) throws IOException, TimeoutException {
        super();
        BasicConfigurator.configure();
        this.clientName = clientName;
        factory = new ConnectionFactory();
        // NEED TO SETUP HOSTS FILE
        factory.setHost(Common.MQ_HOST_NAME);
        connection = factory.newConnection();
        channel = connection.createChannel();
        this.datacenterFeedbackMessageReceiverDirectQueueName = Common
                .getDatacenterFeedbackMessageReceiverDirectQueue(clientName);

    }

    public void sendMessageToDataCenter(ClientDataCenterMessage message) throws IOException {
        channel.basicPublish(Common.CLIENT_REQUEST_DIRECT_EXCHANGE_NAME, message.getDataCenterName(), null,
                Common.serialize(new MessageWrapper(Common.serialize(message), message.getClass())).getBytes());
    }

    public void postToDataCenter(String message) {

    }

    public void run() {
        MessageWrapper wrapper = null;
        QueueingConsumer.Delivery delivery;

        try {
            this.bindToDatacenterMessageReceiverQueue();
        } catch (Exception e) {
            logger.error("Client: " + this.clientName + " bind to client request exchange failed");
            e.printStackTrace();
            exit(-1);
        }

        logger.info("Client: " + this.clientName + " is Running");

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
                    System.out.println("Got a response");
                    if (classType.equals(CenterResponseLookUpMessage.class)) {
                        CenterResponseLookUpMessage message = (CenterResponseLookUpMessage) wrapper.getInnerMessage();
                        handleCenterResponseLookUpMessage(message);
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
    private void handleCenterResponseLookUpMessage(CenterResponseLookUpMessage message) {
        PriorityQueue<EventRecord> records = message.getListOfLogs();
        while (!records.isEmpty()) {
            println(records.poll().getContent());
        }
        printf("> ");
    }

    private static void println(String line) {
        System.out.println(line);
    }

    private static void printf(String line) {
        System.out.printf(line);
    }

    private static void printCommands() {
        println("===================================================");
        println("post(p) <message>");
        println("  - Post a message in DS-blog\n");

        println("lookup(l)");
        println("  - Display the posts in DS-blog in casual order\n");

        println("sync(s) <datacenter>");
        println("  - Synchronize with Datacenter");
        println("=================================================");
        printf("> ");
    }

    public static void main(String[] args) throws IOException, TimeoutException, InterruptedException {
        Client0 c = new Client0("client0");
        String dataCenterName = "dc0";
        new Thread(c).start();
        // logger.info("Send a post");
        // c.sendMessageToDataCenter(new ClientRequestPostMessage(c.clientName, "dc1", "FUCK"));
        //
        // Thread.sleep(3000);
        // logger.info("Request look up");
        // c.sendMessageToDataCenter(new ClientRequestLookUpMessage(c.clientName, "dc1"));

        printf("> ");
        while (true) {
            Scanner scan = new Scanner(System.in);

            String command = scan.nextLine();
            command = command.trim();
            String[] blogArgs = command.split("\\s+");

            blogArgs[0] = blogArgs[0].toLowerCase();

            if (blogArgs[0].equals("p") || blogArgs[0].equals("post")) {

                if (blogArgs.length == 1) {
                    println("Please enter your message");
                    continue;
                }
                else {
                    StringBuilder sb = new StringBuilder();
                    char[] commandChars = command.toCharArray();
                    int start;
                    if (blogArgs[0].equals("p")) {
                        start = 2;
                    }
                    else {
                        start = 4;
                    }
                    while (commandChars[start] == ' ' || commandChars[start] == '\t') {
                        start++;
                    }
                    for (int i = start; i < commandChars.length; i++) {

                        sb.append(commandChars[i]);
                    }
                    String message = sb.toString();
                    c.sendMessageToDataCenter(new ClientRequestPostMessage(c.clientName, dataCenterName, message));
                    printf("> ");
                }
            }

            else if (blogArgs[0].equals("l") || blogArgs[0].equals("lookup")) {
                c.sendMessageToDataCenter(new ClientRequestLookUpMessage(c.clientName, dataCenterName));
                printf("> ");
            }

            else if (blogArgs[0].equals("s") || blogArgs[0].equals("sync")) {
                if (blogArgs.length == 1) {
                    println("Please enter the hostname of the data center you want to sync with");
                    printf("> ");
                    continue;
                }
                else {
                    println("synchronizing with " + blogArgs[1]);
                    c.sendMessageToDataCenter(new ClientRequestSyncMessage(c.clientName, blogArgs[1], dataCenterName));
                    printf("> ");
                }
            }
            else if (blogArgs[0].equals("e") || blogArgs[0].equals("exit")) {
                println("exiting...");
                exit(0);
            }
            else {
                printCommands();
            }
        } // while
    }
}
