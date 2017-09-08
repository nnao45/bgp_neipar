# bgp_neipar
This script is Parsing Router'OS "show bgp ipv4 unicast summary" Results.  
And sort, AS number, LastUP/Down, Pfx/Stat, and You Setting(peer/peer.go)discription.  

## Dependencies
shell -> colordiff  
go    -> github.com/codeskyblue/go-sh  
         github.com/golang/glog  
         github.com/google/goexpect  
         github.com/google/goterm/term  
         github.com/ziutek/telnet  
         golang.org/x/crypto/ssh/terminal  

## Demo(sort by Pfx or Status)
```bash
[nao4arale@bgp ~]$ bgp_neipar -p
Getting show bgp ipv4 unicast summary
Login Router Address: gw-router-tokyo-1  
Login Username: araisan
Login Password: 

################## Sort by Pfx or Status ##################

Peer             AS      LastUP      Pfx/Stat  Description   
---------------  ------  ----------  ------    ------------  
10.0.1.1　　　　　　　  65200   240h0m0s    645352    IX-EX       
10.0.2.1　　　　　　　  65100   6960h0m0s   638534    IXA-user          
10.0.3.1　　　　　　　  65032   6120h0m0s   642       IXB-user          
10.0.4.1　　　　　　　  65001   240h0m0s    33        IXB-user       
10.0.5.1　　　　　　　  65024   240h0m0s    30        IXB-user       
10.0.6.1　　　　　　　  65102   3192h0m0s   4         IXA-user          
10.0.7.1　　　　　　　  65003   0s          Active    IXB-user      
10.0.8.1　　　　　　　  65307   0s          Idle      Route-server        
10.0.9.1　　　　　　　  65007   5064h0m0s   Idle      IXA-user      
10.1.1.1　　　　　　　  65009   6960h0m0s   Idle      IXA-user       
10.1.2.1　　　　　　　  65010   6960h0m0s   Idle      IXA-user        

################ diff Now and Last show cmd ###############

--- /usr/local/bgp_neipar/diff/lastdiff.txt     2017-09-08 21:47:53.272941246 +0900
+++ /usr/local/bgp_neipar/diff/diff.txt 2017-09-08 21:50:21.029792034 +0900
@@ -1,13 +1,13 @@
-10.0.1.1 65200 645352 IX-EX
+10.0.1.1 65200 645352 IX-EX
 10.0.5.1 65024 30 IXB-user
-10.0.8.1 65307 660839 Route-server 
+10.0.8.1 65307 Idle Route-server 
 10.1.2.1 65010 Idle IXA-user
 10.1.1.1 65009 Idle IXA-user
 ```
 
