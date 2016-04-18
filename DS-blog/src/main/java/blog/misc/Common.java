package blog.misc;

import com.google.gson.Gson;

/**
 * @Project: helios
 * @Title: Common.java
 * @Package helios
 * @Description: TODO
 * @author YuesongWang
 * @date Feb 28, 2016 1:37:59 AM
 * @version V1.0
 */
public class Common {
    // For log propagation
    public final static String LOG_DIRECT_EXCHANGE_NAME = "logPropDirectExchange";
    // For client send request to data center, with routing key
    public final static String CLIENT_REQUEST_DIRECT_EXCHANGE_NAME = "clientRequestDirectExchange";
    // public final static String DC_RESPONSE_DIRECT_EXCHANGE_NAME = "dcResponseDirectExchange";
    public final static String MQ_HOST_NAME = "rabbithost";

    /**
     * Use json to serialize POJO message
     * Description: TODO
     * 
     * @param message
     * @return
     *         String
     */
    public static <T> String serialize(T message)
    {
        Gson gson = new Gson();
        return gson.toJson(message, message.getClass());
    }

    @SuppressWarnings({ "rawtypes", "unchecked" })
    public static <T> T deserialize(String json, Class className)
    {
        Gson gson = new Gson();
        return (T) gson.fromJson(json, className);
    }

    @SuppressWarnings("rawtypes")
    public static Class getClassFromString(String s) throws ClassNotFoundException
    {
        Class<?> cls = Class.forName(s);
        return cls;
    }

    // Datacenter as receiver
    public static String getClientMessageReceiverDirectQueueName(String dataCenterName) {
        return dataCenterName + ".clientMessageReceiver.direct.queue";
    }

    public static String getDatacenterLogPropagationDirectQueueName(String dataCenterName) {
        return dataCenterName + ".logPropagation.direct.queue";
    }

    // Client as a receiver
    public static String getDatacenterFeedbackMessageReceiverDirectQueue(String clientName) {
        return clientName + ".datacenterFeedbackMessageReceiver.direct.queue";
    }

}
