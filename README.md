# temporal-dispatcher-demo

## 定义
- master workflow, 主workflow，用于启动child workflow。其中包括producer和consumer2个主workflow
- child workflow, 这里指代负责收发消息的workflow

## 关键思路
1. 服务启动的时候，同时启动master和worker，master负责启动scheduler,worker负责执行workflow，其中master启动加zk全局锁，只有一个master
可以加锁成功。其他节点阻塞等待，一旦master退出，其他节点重新竞争锁。
    > 这里有个问题，启用child workflow的主workflow并没有限制在master节点上，与题目要求不符。
    >
    > 可能的解决方案有：
    > 1. 主workflow只在 master节点上注册，限制调度。但是个人感觉没什么意义。
    > 2. 主workflow不使用temporal的schedule模式，自己启用cornjob来实现。
2. worker注册用于收发消息的child workflow,以及用于调度的 master workflow
3. schedule 触发的时候，利用temporal的overlap特性关闭掉，之前没有运行完的workflow
4. master如果退出，不会关闭schedule。而是交由新的master去关闭，并启用新的schedule。(这里delete再删除，主要是语义收敛在master负责schedule的生命周期)
5. 每个节点只有一个consumer和producer。producer可以并发访问
6. consumer节点不易过多，以防超过kafka partition数量，导致部分节点接受不到消息。所以这里设置的consumer节点为一个。
7. consumer以group形式消费，因为单节点只有一个client，但是可能有多个workflow来访问，所以将callback消费修改为ping-pong模式。使用channel控制并发。
8. consumer生命周期和服务节点一致，不跟随workflow。
9. consumer消费到的消息使用了temporal的sideeffect模式，防止workflow重试的时候，重复消费消息。（暂没有验证）

## TODO
- [ ] 验证sideeffect模式的正确性
- [ ] 优化日志输出，当前日志是直接打印到终端的，不优雅
- [ ] 丰富日志监控，当前各个节点状态有点像黑盒子 (需要再熟悉下temporal，看下是否有更多可以展示的内容)
- [ ] dockerfile 配置文件改为外部传入
- [ ] 优化helm chart，支持添加依赖