/**
 * Created by xuanwang on 4/12/16.
 */
class Event {
    private long timestamp;
    private String message;

    public Event(long time, String msg){
        timestamp = time;
        message = msg;
    }

    public long getTime(){
        return timestamp;
    }

    public String getMessage(){
        return message;
    }



}
