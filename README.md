# mgt-zookeeper

使用说明：
1.用户配置通过config.yml文件实现
2.配置源zk地址、目的zk地址
3.配置zk url黑白名单，白名单表示同步范围，不支持通配。黑名单表示禁止同步节点，支持通配

# 版本记录

##  2023-06-30-14-53

config.yml配置文件说明：
1.白名单列表中的节点，只同步白名单及其子节点。
2.白名单列表不能配置包含父子关系的节点，比如配了"/"表示同步所有，再配置其他子目录会多余
3.白名单配置格式：不支持以"/"结尾，不支持通配符，参考"/"、"/xxx"、"/xxx/yyy"
4.黑名单中的节点会被排除同步，相同节点黑名单优先级高于白名单
5.黑名单配置格式：不支持通配"*"或"/xxx*"，不支持以"/"结尾"/xxx/yyy/"，支持通配"/*"或"/xxx/*"，参考"/xxx"、"/xxx/yyy"、"/*"、"/xxx/*"
6.黑名单可以为空，白名单不可以为空