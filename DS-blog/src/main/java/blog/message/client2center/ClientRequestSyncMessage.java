
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
public class ClientRequestSyncMessage extends ClientDataCenterMessage{
    private String fromDataCenterName;


    public ClientRequestSyncMessage(String clientName, String fromDataCenterName, String toDataCenterName){
        super(clientName, toDataCenterName);
        this.fromDataCenterName = fromDataCenterName;
    }

    public String getClientName(){
        return this.clientName;
    }

    public String getFromDataCenterName(){
        return this.fromDataCenterName;
    }

    public String getToDataCenterName(){
        return this.dataCenterName;
    }
}
