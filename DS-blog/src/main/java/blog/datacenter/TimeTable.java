package blog.datacenter;

import java.util.HashMap;

/**
 * @Project: ds-blog
 * @Title: TimeTable.java
 * @Package blog.datacenter
 * @Description: TODO
 * @author YuesongWang
 * @date Apr 18, 2016 12:58:33 PM
 * @version V1.0
 */
public class TimeTable {
    long[][] table;

    public TimeTable(HashMap<String, Integer> dataCenterNameToIndex) {
        table = new long[dataCenterNameToIndex.size()][dataCenterNameToIndex.size()];
        
    }
}
