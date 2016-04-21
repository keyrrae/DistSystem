package blog.message.client2center;

import java.util.PriorityQueue;

import blog.datacenter.Post;

/**
 * @Project: ds-blog
 * @Title: CenterResponseLookUpMessage.java
 * @Package blog.message
 * @Description: TODO
 * @author YuesongWang
 * @date Apr 18, 2016 1:28:44 PM
 * @version V1.0
 */
public class CenterResponseLookUpMessage extends ClientDataCenterMessage {

    public CenterResponseLookUpMessage(String clientName, String dataCenterName, PriorityQueue<Post> listOfPost) {
        super(clientName, dataCenterName);
        this.listOfPost = listOfPost;
    }

    PriorityQueue<Post> listOfPost;

    public PriorityQueue<Post> getListOfPost() {
        return listOfPost;
    }

    public void setListOfPost(PriorityQueue<Post> listOfPost) {
        this.listOfPost = listOfPost;
    }

}
