import java.util.ArrayList;
import java.util.List;

/**
 * Created by xuanwang on 4/12/16.
 */
class Log {
    private ArrayList<Event> eventList;

    Log(){
        eventList = new ArrayList<Event>();
    }

    Log(Log baseLog){
        eventList = new ArrayList<Event>();
        for(Event e: baseLog.getEvents()){
            eventList.add(e);
        }
    }

    public Event getEvent(){
        return null;
    }

    public ArrayList<Event> getEvents() { return eventList; }

    void printLog(){
        for(Event e: eventList){
            System.out.println(String.valueOf(e));
        }
    }
}
