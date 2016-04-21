package blog.datacenter;

/**
 * @Project: ds-blog
 * @Title: Post.java
 * @Package blog.datacenter
 * @Description: TODO
 * @author YuesongWang
 * @date Apr 18, 2016 1:07:06 PM
 * @version V1.0
 */
public class Post {
    String content;

    public Post(String content) {
        super();
        this.content = content;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

}
