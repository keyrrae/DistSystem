package blog.message.client2center;

import blog.datacenter.Post;
import blog.message.Message;

/**
 * @Project: ds-blog
 * @Title: ClientRequestMessage.java
 * @Package blog.message
 * @Description: TODO
 * @author YuesongWang
 * @date Apr 18, 2016 1:28:44 PM
 * @version V1.0
 */
public class ClientRequestPostMessage extends ClientDataCenterMessage {
    String post;

    public ClientRequestPostMessage(String clientName, String dataCenterName, String post) {
        super(clientName, dataCenterName);
        this.post = post;
    }

    public String getPost() {
        return post;
    }

    public void setPost(String post) {
        this.post = post;
    }

}
