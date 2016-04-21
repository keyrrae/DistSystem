package blog.datacenter;

import java.util.Comparator;

/**
 * @Project: ds-blog
 * @Title: Post.java
 * @Package blog.datacenter
 * @Description: TODO
 * @author YuesongWang
 * @date Apr 18, 2016 1:07:06 PM
 * @version V1.0
 */
public class Post implements Comparable<Post>{
    private String content;
    private long timeStamp;

    public Post(String content) {
        super();
        this.content = content;
    }

    public int compareTo(Post o) {
        if(this.timeStamp > o.timeStamp){
            return 1;
        }
        else if(this.timeStamp < o.timeStamp){
            return -1;
        }
        return 0;
    }

    public String getContent() {
        return content;
    }

    public void setTimeStamp(long timeStamp){this.timeStamp = timeStamp;}

    public long getTimeStamp() {return this.timeStamp;}

    public void setContent(String content) {
        this.content = content;
    }

}
