package blog.datacenter;

import static org.junit.Assert.*;

import java.io.IOException;
import java.util.HashMap;
import java.util.concurrent.TimeoutException;

import org.junit.Before;
import org.junit.Test;

import blog.logs.EventRecord;

/**
 * @Project: ds-blog
 * @Title: DataCenterUnitTest.java
 * @Package blog.datacenter
 * @Description: TODO
 * @author YuesongWang
 * @date Apr 26, 2016 1:10:36 AM
 * @version V1.0
 */
public class DataCenterUnitTest {

    DataCenter dc0;
    DataCenter dc1;
    DataCenter dc2;

    @Before
    public void initalized() throws IOException, TimeoutException {
        HashMap<String, Integer> dataCenterNameToIndex = new HashMap<String, Integer>();
        dataCenterNameToIndex.put("dc0", 0);
        dataCenterNameToIndex.put("dc1", 1);
        dataCenterNameToIndex.put("dc2", 2);
        dc0 = new DataCenter("dc0", dataCenterNameToIndex);
        dc1 = new DataCenter("dc1", dataCenterNameToIndex);
        dc2 = new DataCenter("dc2", dataCenterNameToIndex);
    }

    @Test
    public void testHasRec() {
        // Set dc0 know dc1 know dc2 till 3
        dc0.timeTable.getTable()[1][2] = 3;
        EventRecord e1 = new EventRecord(2, "dc2", "DC2 time 2 event");
        EventRecord e2 = new EventRecord(4, "dc2", "DC2 time 4 event");
        // Expect dc0 knows dc1 know dc2's time 2 event
        assertTrue(dc0.hasRec(e1, 1));
        // Expect dc0 not knows dc1 know dc2's time 4 event
        assertFalse(dc0.hasRec(e2, 1));
    }
}
