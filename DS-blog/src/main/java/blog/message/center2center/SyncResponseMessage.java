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
    private List<EventRecord> logs;
    private TimeTable timeTable;

    /**
     * @param fromDataCenterName
     * @param toDataCenterName
     */
    public SyncResponseMessage(String fromDataCenterName, String toDataCenterName) {
        super(fromDataCenterName, toDataCenterName);
        // TODO Auto-generated constructor stub
    }

    public void setTimeTable(TimeTable timeTable){
        this.timeTable = timeTable;
    }

    public void setLog(List<EventRecord> logs){
        this.logs = logs;
    }

    public TimeTable getTimeTable(){
        return this.timeTable;
    }

    public List<EventRecord> getLog(){
        return this.logs;
    }
}
