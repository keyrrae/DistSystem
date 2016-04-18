package blog.datacenter;

import java.io.IOException;
import java.util.List;
import java.util.concurrent.TimeoutException;

import blog.logs.EventRecord;
import blog.misc.Common;

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

    private String dataCenterName;
    private String logPropagationDirectQueueName;
    private String clientMessageDirectQueueName;

    // Below for connecting MQ
    private ConnectionFactory factory;
    private Connection connection;
    private Channel channel;
    // Consumer for client message and log message
    private QueueingConsumer consumer;

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

    public DataCenter(String dataCenterName) throws IOException, TimeoutException {
        factory = new ConnectionFactory();
        factory.setHost(Common.MQ_HOST_NAME);
        connection = factory.newConnection();
        channel = connection.createChannel();
        
        this.dataCenterName = dataCenterName;
        this.logPropagationDirectQueueName = Common.getDatacenterLogPropagationDirectQueueName(dataCenterName);
        this.clientMessageDirectQueueName = Common.getClientMessageReceiverDirectQueueName(dataCenterName);
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

    public void run() {

    }

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
