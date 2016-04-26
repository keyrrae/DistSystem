package blog.datacenter;

import static org.junit.Assert.*;

import java.util.HashMap;

import org.junit.Before;
import org.junit.Test;

/**
 * @Project: ds-blog
 * @Title: TimeTableUnitTest.java
 * @Package blog.datacenter
 * @Description: TODO
 * @author YuesongWang
 * @date Apr 26, 2016 1:54:22 AM
 * @version V1.0
 */

public class TimeTableUnitTest {
    TimeTable tb0;
    TimeTable tb1;
    TimeTable tb2;

    @Before
    public void init() {
        HashMap<String, Integer> dataCenterNameToIndex = new HashMap<String, Integer>();
        dataCenterNameToIndex.put("dc0", 0);
        dataCenterNameToIndex.put("dc1", 1);
        dataCenterNameToIndex.put("dc2", 2);

        tb0 = new TimeTable(dataCenterNameToIndex, "dc0");
        tb1 = new TimeTable(dataCenterNameToIndex, "dc1");
        tb2 = new TimeTable(dataCenterNameToIndex, "dc2");

    }

    @Test
    public void testUpdateUponReceived() {
        tb1.updateTimeTable(0, tb0);

        long[][] mergedTb0 = new long[][] {
                { 3, 0, 0 }, { 3, 4, 0 }, { 0, 0, 0 }
        };
        for (int i = 0; i < mergedTb0.length; i++)
            for (int j = 0; j < mergedTb0[0].length; j++) {
                assertTrue(mergedTb0[i][j] == tb1.getTable()[i][j]);
            }
    }
}
