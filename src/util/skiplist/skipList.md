
 - 版本1对应的skipList结构
![skipList](https://static.studygolang.com/190805/198f69e95d7af643ca8b8fe893d2e0e8.png)

 - 版本1对应的插入描述
  ![insert](https://static.studygolang.com/190805/90f667174c19fe92933482e14a29dee9.png)
  
  
 - 版本2对应的skipList结构
 ![](https://upload-images.jianshu.io/upload_images/10335749-e7a205a8c9d9604a.png)
 
  - 版本2数据插入
  ![](https://upload-images.jianshu.io/upload_images/10335749-f9357f6f81d44cd5.png)
  ```text
   我们需要插入一个3，并且调用randomLevel得到的层数为3，那么插入3需要如下几部：
   
       1、沿着LEVEL3的链查找第一个比3大的节点或者TAIL节点，记录下该节点的前一个节点——HEAD和层数——3。
   
       2、沿着LEVEL2的链查找第一个比3大的节点或者TAIL节点，记录下该节点的前一个节点——&1和层数——2。
   
       3、沿着LEVEL1的链查找第一个比3大的节点或者TAIL节点，记录下该节点的前一个节点——&2和层数——1。
   
       4、生成一个新的节点newNode，key赋值为3，将newNode插入HEAD、&1、&2之后，即HEAD.next[3]=&3，&1.next[2]=&3，&2.next[1]=&3。
   
       5、给newNode的next赋值，即步骤4中HEAD.next[3]、&1.next[2]、2.next[1]原本的值。
   
       注意：为了易于理解，上述步骤中所有索引均从1开始，而代码中则从0开始，所以代码中均有索引=层数-1的关系。
   
```