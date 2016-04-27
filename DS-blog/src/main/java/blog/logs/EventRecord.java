package blog.logs;

/**
 * Created by xuanwang on 4/12/16.
 */
public class EventRecord implements Comparable<EventRecord> {
    private OperationType op;
    private long timestamp;

    // Equals to NodeId on paper
    private String nodeName;
    // POST's content
    private String content;

    public EventRecord(long timestamp, String nodeName, String content) {
        this.op = OperationType.POST;
        this.timestamp = timestamp;
        this.nodeName = nodeName;
        this.content = content;
    }

    public OperationType getOp() {
        return op;
    }

    public void setOp(OperationType op) {
        this.op = op;
    }

    public long getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(long timestamp) {
        this.timestamp = timestamp;
    }

    public String getNodeName() {
        return nodeName;
    }

    public void setNodeName(String nodeName) {
        this.nodeName = nodeName;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    /*
     * (non-Javadoc)
     * 
     * @see java.lang.Comparable#compareTo(java.lang.Object)
     */
    public int compareTo(EventRecord o) {
        if (timestamp > o.timestamp)
            return 1;
        else if (timestamp == o.timestamp)
            return 0;
        else
            return -1;
    }

}
