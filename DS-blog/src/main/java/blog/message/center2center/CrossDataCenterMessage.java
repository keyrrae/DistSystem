package blog.message.center2center;

import blog.message.Message;

/**
 * @Project: ds-blog
 * @Title: CrossDataCenterMessage.java
 * @Package blog.message
 * @Description: TODO
 * @author YuesongWang
 * @date Apr 20, 2016 7:31:21 PM
 * @version V1.0
 */
public class CrossDataCenterMessage extends Message {
    String fromDataCenterName;
    String toDataCenterName;

    
    public CrossDataCenterMessage(String fromDataCenterName, String toDataCenterName) {
        super();
        this.fromDataCenterName = fromDataCenterName;
        this.toDataCenterName = toDataCenterName;
    }

    public String getFromDataCenterName() {
        return fromDataCenterName;
    }

    public void setFromDataCenterName(String fromDataCenterName) {
        this.fromDataCenterName = fromDataCenterName;
    }

    public String getToDataCenterName() {
        return toDataCenterName;
    }

    public void setToDataCenterName(String toDataCenterName) {
        this.toDataCenterName = toDataCenterName;
    }

}
