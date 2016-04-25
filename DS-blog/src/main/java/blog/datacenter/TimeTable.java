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

    public int size(){
        return table.length;
    }

    public void increaseLocalClock(){
        table[dataCenterIndex][dataCenterIndex]++;
    }

    public long getLocalClock(){
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

    public void upDateUponReceived(TimeTable recv){
        this.increaseLocalClock();
        int n = this.table.length;
        int m = recv.size();

        if(n != m){
            System.out.println("different time table dimensions");
            exit(1);
        }
        for(int i = 0; i < n; i++){
            for (int j = 0; j < n; j++){
                table[i][j] = table[i][j] > recv.getTable()[i][j] ? table[i][j] : recv.getTable()[i][j];
            }
        }

        int recvIndex = recv.getDataCenterIndex();
        int thisIndex = dataCenterIndex;

        for(int j = 0; j < n; j++){
            table[thisIndex][j] = table[thisIndex][j] > recv.getTable()[recvIndex][j] ?
                    table[thisIndex][j] : recv.getTable()[recvIndex][j];
        }
    }
}
