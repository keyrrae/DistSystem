package blog.message.client2center;

import blog.message.Message;

/**
 * @Project: ds-blog
 * @Title: ClientDataCenterMessage.java
 * @Package blog.message
 * @Description: TODO
 * @author YuesongWang
 * @date Apr 20, 2016 7:30:59 PM
 * @version V1.0
 */
public class ClientDataCenterMessage extends Message {
    String clientName;
    String dataCenterName;

    
    public ClientDataCenterMessage(String clientName, String dataCenterName) {
        super();
        this.clientName = clientName;
        this.dataCenterName = dataCenterName;
    }

    public String getClientName() {
        return clientName;
    }

    public void setClientName(String clientName) {
        this.clientName = clientName;
    }

    public String getDataCenterName() {
        return dataCenterName;
    }

    public void setDataCenterName(String dataCenterName) {
        this.dataCenterName = dataCenterName;
    }

}
