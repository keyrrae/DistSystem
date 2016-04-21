package blog.client;

import java.io.IOException;
import java.util.concurrent.TimeoutException;

import org.apache.log4j.BasicConfigurator;
import org.apache.log4j.Logger;

import blog.datacenter.Post;
import blog.message.client2center.CenterResponseLookUpMessage;
import blog.message.client2center.ClientDataCenterMessage;
import blog.message.client2center.ClientRequestLookUpMessage;
import blog.message.client2center.ClientRequestPostMessage;
import blog.misc.Common;
import blog.misc.MessageWrapper;

import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;
import com.rabbitmq.client.QueueingConsumer;

/**
 * Created by xuanwang on 4/12/16.
 */
public class Client implements Runnable {
    private static Logger logger = Logger.getLogger(Client.class);
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

    public Client(String clientName) throws IOException, TimeoutException {
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
            System.exit(-1);
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
                    if (classType.equals(CenterResponseLookUpMessage.class)) {
                        CenterResponseLookUpMessage message = (CenterResponseLookUpMessage) wrapper.getInnerMessage();
                        logger.info("CenterResponseLookUpMessage");
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
        for (Post p : message.getListOfPost()) {
            System.out.println(p.getContent());
        }
    }

    public static void main(String[] args) throws IOException, TimeoutException, InterruptedException {
        Client c = new Client("client1");
        new Thread(c).start();
        logger.info("Send a post");
        c.sendMessageToDataCenter(new ClientRequestPostMessage(c.clientName, "dc1", new Post("FUCK")));
       
        Thread.sleep(3000);
        logger.info("Request look up");
        c.sendMessageToDataCenter(new ClientRequestLookUpMessage(c.clientName, "dc1"));

    }
}
