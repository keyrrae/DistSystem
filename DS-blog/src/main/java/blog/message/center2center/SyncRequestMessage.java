
package blog.message.center2center;

import blog.message.Message;

/**  
 * @Project: ds-blog
 * @Title: SyncRequestMessage.java
 * @Package blog.message
 * @Description: TODO
 * @author YuesongWang
 * @date Apr 18, 2016 1:28:26 PM
 * @version V1.0  
 */
public class SyncRequestMessage extends CrossDataCenterMessage{

    /**
     * @param fromDataCenterName
     * @param toDataCenterName
     */
    public SyncRequestMessage(String fromDataCenterName, String toDataCenterName) {
        super(fromDataCenterName, toDataCenterName);
        // TODO Auto-generated constructor stub
    }

}
