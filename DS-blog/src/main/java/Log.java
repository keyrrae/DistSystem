import java.util.ArrayList;
import java.util.List;

import blog.logs.EventRecord;

/**
 * Created by xuanwang on 4/12/16.
 */
class Log {
    private ArrayList<EventRecord> eventList;

    Log(){
        eventList = new ArrayList<EventRecord>();
    }

    Log(Log baseLog){
        eventList = new ArrayList<EventRecord>();
        for(EventRecord e: baseLog.getEvents()){
            eventList.add(e);
        }
    }

    public EventRecord getEvent(){
        return null;
    }

    public ArrayList<EventRecord> getEvents() { return eventList; }

    void printLog(){
        for(EventRecord e: eventList){
            System.out.println(String.valueOf(e));
        }
    }
}
