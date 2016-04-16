import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;
import org.json.JSONObject;

import java.io.IOException;

/**
 * Created by xuanwang on 4/14/16.
 */
public class JSONsender extends Thread {
    private final static String QUEUE_NAME = "jsonMsgQueue";
    private ConnectionFactory factory = null;

    public JSONsender() {
        // TODO Auto-generated constructor stub
    }

    @Override
    public void run() {

        factory = new ConnectionFactory();
        factory.setHost("localhost");
        factory.setPort(5676);

        System.out.println("connected to rabbitMQ on localhost ...");
        Connection connection = factory.newConnection();
        Channel channel = null;
        try{
            channel = connection.createChannel();
        }
        catch(IOException e){
            e.printStackTrace();
        }

        channel.queueDeclare(QUEUE_NAME, false, false, false, null);
        for (int i = 1; i <= 10; i++) {
            JSONObject obj = new JSONObject();

            obj.put("name", String.format("Person%s", i));
            obj.put("age", new Integer(37));

            channel.basicPublish("", QUEUE_NAME, null, obj.toString().getBytes());
            System.out.println(" [x] Sent '" + obj.toString() + "'");
        }

        channel.close();
        connection.close();
    }
}
