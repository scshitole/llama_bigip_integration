# llama_bigip_integration

![image](https://github.com/user-attachments/assets/4b61edc9-5de7-45dc-9716-050603b289af)
![image](https://github.com/user-attachments/assets/ee38afc7-702f-4081-af1e-9c3a309bbc71)
## Some example which were tried
xuser@ubuntu:~/llama_bigip_integration$ go run main.go "How many virtual servers"
Based on the provided query and virtual server information from F5 BIG-IP, here are the key details:

* Total Virtual Servers: 4
* Names:
 + NEwVSServer
 + VS_app1
 + VS_app2
 + VS_app3
* Destinations:
 + /Common/200.1.1.1:80 (NEwVSServer)
 + /Common/192.168.1.101:80 (VS_app1)
   + /Common/192.168.1.102:80 (VS_app2)
    + /Common/192.168.1.103:80 (VS_app3)
* Associated Pools:
   + NEwVSServer - No pool associated
    + VS_app1 - No pool associated
     + VS_app2 - No pool associated
      + VS_app3 - No pool associated
* Enabled Status:
  All virtual servers are enabled.
xuser@ubuntu:~/llama_bigip_integration$ go run main.go "How many pools"
The provided query and pool information does not directly ask for the count of pools. However, based on the provided data, there is only 1 pool available:

- Name: ext_gw_pool
- Description: This pool is defined as round-robin with "pass-through" tos settings.

Therefore, there is only one pool specified in the provided information.
xuser@ubuntu:~/llama_bigip_integration$ go run main.go "How many pool members"
There is 1 pool member in the specified pool 'ext_gw_pool'. The pool member's details are:

- Name: 10.1.20.1:0
- Address: 10.1.20.1
- State: up
xuser@ubuntu:~/llama_bigip_integration$ go run main.go "How many pool members and what is the status of the pool members is it UP or Down"
Based on the provided query and pool member information from F5 BIG-IP, here are the key points:

* Number of pool members: 1
* Status of the pool member: UP (specifically, the member "10.1.20.1:0" is up)

There is only one pool member in the specified pool named "ext_gw_pool". This member is currently operating correctly and is marked as "up", indicating that it is available for use.
xuser@ubuntu:~/llama_bigip_integration$ go run main.go "How many virtual servers are there on this BIG-IP, can you provide the IP addresses of these Virtual servers"
Based on the provided virtual server information, here are the key details:

Number of Virtual Servers: 4

IP Addresses of Virtual Servers:
- 200.1.1.1
- 192.168.1.101
- 192.168.1.102
- 192.168.1.103

Associated Pools:
- New VSServer: no pool specified (empty)
- VS_app1: /Common/192.168.1.101:80
- VS_app2: /Common/192.168.1.102:80
- VS_app3: /Common/192.168.1.103:80

Enabled Status: All virtual servers are enabled (true).
xuser@ubuntu:~/llama_bigip_integration$ 

