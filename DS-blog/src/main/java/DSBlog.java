import java.util.Scanner;

import static java.lang.System.exit;


/**
 * Created by xuanwang on 4/12/16.
 */
public class DSBlog {


    private static void println(String line){
        System.out.println(line);
    }

    private static void printf(String line){
        System.out.printf(line);
    }

    private static void printUsage(){
        println("============================================");
        printf("Usage:");
        println("============================================");
    }

    private static void printCommands(){
        println("============================================");
        println("post(p) <message>");
        println("  - Post a message in DS-blog\n");

        println("lookup(l)");
        println("  - Display the posts in DS-blog in casual order\n");

        println("sync(s) <datacenter>");
        println("  - Synchronize with Datacenter");
        println( "============================================");
    }


    public static void main(String[] args){
        DataCenter dc = null;
        println(String.valueOf(args.length));

        if(args.length == 0){
            dc = new DataCenter("./config.cfg");
        }
        else if(args.length == 1){
            dc = new DataCenter(args[0]);
        }
        else{
            printUsage();
            exit(0);
        }

        while(true){
            Scanner scan = new Scanner(System.in);
            printf(">");
            String command = scan.nextLine();
            command = command.trim();
            String[] blogArgs = command.split("\\s+");

            blogArgs[0] = blogArgs[0].toLowerCase();

            if(blogArgs[0].equals("p") || blogArgs[0].equals("post")){
                println("post");
                if(blogArgs.length == 1){
                    println("Please enter your message");
                    continue;
                }
                else{
                    StringBuilder sb = new StringBuilder();
                    char[] commandChars = command.toCharArray();
                    for(int i = 4; i < commandChars.length; i++){
                        if(commandChars[i] == ' ' || commandChars[i] == '\t'){
                            continue;
                        }
                        else {
                            sb.append(commandChars[i]);
                        }
                    }
                    String message = sb.toString();
                    dc.post(message);
                }
            }
            else if(blogArgs[0].equals("l") || blogArgs[0].equals("lookup")){
                println("lookup");
                dc.lookup();
            }
            else if(blogArgs[0].equals("s") || blogArgs[0].equals("sync")){
                if(blogArgs.length == 1){
                    println("Please enter the hostname of the data center you want to sync with");
                    continue;
                }
                else {
                    println("synchronizing with " + blogArgs[1]);
                    dc.sync(blogArgs[1]);
                }
            }
            else if(blogArgs[0].equals("e") || blogArgs[0].equals("exit")){
                println("exiting...");
                exit(0);
            }
            else {
                printCommands();
            }
        } // while
    }
}
