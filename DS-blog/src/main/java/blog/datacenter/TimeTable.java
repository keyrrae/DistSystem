package blog.datacenter;

import java.util.HashMap;

import static java.lang.System.exit;

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
    private long[][] table;
    private int dataCenterIndex;

    public TimeTable(HashMap<String, Integer> dataCenterNameToIndex, String dataCenterName) {
        table = new long[dataCenterNameToIndex.size()][dataCenterNameToIndex.size()];
        dataCenterIndex = dataCenterNameToIndex.get(dataCenterName);

    }

    public long getLocalClock() {
        return table[dataCenterIndex][dataCenterIndex];
    }

    public long[][] getTable() {
        return table;
    }

    public void setTable(long[][] table) {
        this.table = table;
    }

    public int getDataCenterIndex() {
        return dataCenterIndex;
    }

    public void setDataCenterIndex(int dataCenterIndex) {
        this.dataCenterIndex = dataCenterIndex;
    }

    public int size() {
        return table.length;
    }

    public void increaseLocalClock() {
        table[dataCenterIndex][dataCenterIndex]++;
    }

    /**
     * 
     * Description: Check if datacenter X(index) know datacenter Y(index) till timestamp T.
     * 
     * @param X
     * @param Y
     * @param T
     * @return
     *         boolean
     */
    public boolean ifDatacenterXKnowDatacenterYTillTimeT(int X, int Y, long T) {
        return this.table[X][Y] >= T;
    }

  
}
