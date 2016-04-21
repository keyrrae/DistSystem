package blog.message.client2center;


/**
 * @Project: ds-blog
 * @Title: ClientRequestMessage.java
 * @Package blog.message
 * @Description: TODO
 * @author YuesongWang
 * @date Apr 18, 2016 1:28:44 PM
 * @version V1.0
 */
public class ClientRequestLookUpMessage extends ClientDataCenterMessage {

    /**
     * @param clientName
     * @param dataCenterName
     */
    public ClientRequestLookUpMessage(String clientName, String dataCenterName) {
        super(clientName, dataCenterName);
        // TODO Auto-generated constructor stub
    }
  
    
}
