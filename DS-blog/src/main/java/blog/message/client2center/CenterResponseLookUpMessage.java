package blog.message.client2center;

import java.util.PriorityQueue;

import blog.datacenter.Post;
import blog.logs.EventRecord;

/**
 * @Project: ds-blog
 * @Title: CenterResponseLookUpMessage.java
 * @Package blog.message
 * @Description: TODO
 * @author YuesongWang
 * @date Apr 18, 2016 1:28:44 PM
 * @version V1.0
 */
public class CenterResponseLookUpMessage extends ClientDataCenterMessage {
    PriorityQueue<EventRecord> listOfLogs;

    public CenterResponseLookUpMessage(String clientName, String dataCenterName, PriorityQueue<EventRecord> logs) {
        super(clientName, dataCenterName);
        this.listOfLogs = logs;
    }

    public PriorityQueue<EventRecord> getListOfLogs() {
        return listOfLogs;
    }

    public void setListOfLogs(PriorityQueue<EventRecord> listOfLogs) {
        this.listOfLogs = listOfLogs;
    }

}
