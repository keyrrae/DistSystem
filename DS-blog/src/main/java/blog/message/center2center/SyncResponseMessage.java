package blog.message.center2center;

import blog.datacenter.TimeTable;
import blog.logs.EventRecord;
import blog.message.Message;

import java.sql.Time;
import java.util.List;

/**
 * @Project: ds-blog
 * @Title: SyncResponseMessage.java
 * @Package blog.message
 * @Description: TODO
 * @author YuesongWang
 * @date Apr 18, 2016 1:28:35 PM
 * @version V1.0
 */
public class SyncResponseMessage extends CrossDataCenterMessage {
    List<EventRecord> logs;
    TimeTable timeTable;

    /**
     * @param fromDataCenterName
     * @param toDataCenterName
     */
    public SyncResponseMessage(String fromDataCenterName, String toDataCenterName) {
        super(fromDataCenterName, toDataCenterName);
        // TODO Auto-generated constructor stub
    }

    public void setTimeTable(TimeTable timeTable){
        // TODO
    }

    public void setLog(List<EventRecord> logs){
        // TODO
    }

    public TimeTable getTimeTable(){
        //TODO
        return this.timeTable;
    }

    public List<EventRecord> getLog(){
        return this.logs;
    }
}
