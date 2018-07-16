# TeaWeb - 可视化智能Web服务
TeaWeb集静态资源、缓存、代理、统计、监控于一体的可视化智能WebServer。

# 架构
~~~
            |---------|       |---------------------------| 
Client  ->  | TeaWeb  |  <->  | Nginx, Apache, Tomcat ... |
            |---------|       |---------------------------|
               Web
               Proxy
               Statistics
               Monitor
               Security
               Log
~~~
