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

    
    public Message getInnerMessage() throws ClassNotFoundException{
        Message innerMsg = Common.deserialize(this.serializedInnerMessage, Common.getClassFromString(this.messageClass));
        return innerMsg;
    }
}
