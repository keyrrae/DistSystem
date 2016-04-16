import org.json.JSONObject;

import java.util.ArrayList;

/**
 * Created by xuanwang on 4/14/16.
 */
public class Message {
    private ArrayList<ArrayList<Long>> timeTable;
    private Log log;

    Message(Log dcLog, ArrayList<ArrayList<Long>> dcTimeTable){
        log = new Log(dcLog);
        timeTable = new ArrayList<ArrayList<Long>>(dcTimeTable);

    }


    /*static<T> JSONObject toJSONObj(T objToConvert){
        ObjectMapper mapper = new ObjectMapper();

        JSONObject res = new JSONObject();
        return res;
    }
    */
}
