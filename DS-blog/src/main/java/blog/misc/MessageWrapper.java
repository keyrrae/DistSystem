package blog.misc;

import blog.message.Message;
//import messagefutures.message.*;

// Message is a wrapper class to abstract any message class.
public class MessageWrapper {
    private String serializedInnerMessage;
    private String messageClass;

    public MessageWrapper(String serializedMessage, Class className)
    {
        this.serializedInnerMessage = serializedMessage;
        this.messageClass = className.getName();
    }

    public MessageWrapper(String serializedMessage, String className)
    {
        this.serializedInnerMessage = serializedMessage;
        this.messageClass = className;
    }

    /**
     * 
     * Description: Get the content(in the form of supertype Message) of innerMessage. Cast it to ClientRequestMessage,
     * SyncRequestMessage etc. for further usage.
     * 
     * @return
     * @throws ClassNotFoundException
     *             Message
     */
    public Message getInnerMessage() throws ClassNotFoundException {
        Message innerMsg = Common
                .deserialize(this.serializedInnerMessage, Common.getClassFromString(this.messageClass));
        return innerMsg;
    }

    /**
     * 
     * Description: Get the Class(ClientRequest, SyncRequest, SyncResponse) of the innerMessage contained
     * 
     * @return
     * @throws ClassNotFoundException
     *             Class
     */
    public Class getInnerMessageClass() throws ClassNotFoundException {
        return Common.getClassFromString(this.messageClass);
    }
}
