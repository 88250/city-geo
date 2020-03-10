## 简介

该仓库整理了中国城市县区的经纬度数据。

欢迎关注 B3log 开源社区微信公众号 `B3log开源`：

![image-d3c00d78](https://user-images.githubusercontent.com/873584/71566370-0d312c00-2af2-11ea-8ea1-0d45d6f0db20.png)

## 用法

[data.json](https://github.com/88250/city-geo/blob/master/data.json) 是已经整理好的数据，可直接使用。

代码中的 `baiduAK` 请勿在生产环境使用，可能会随时删除。

## 动机

整理这些数据的动机是满足[黑客派](https://hacpai.com)实现暗黑模式的需要，模式分为明亮、暗黑、随日出日落自动切换。

随日出日落自动切换特性需要知道日出日落时间，而不同地理位置的日出日落时间是不同的，但可以基于经纬度来进行计算，这就是制作该仓库的动机。

## 鸣谢

* [城市数据来源](https://github.com/modood/Administrative-divisions-of-China)
* [百度地图 API](http://lbsyun.baidu.com)
