# Boardgame-SultansOfKaraya
Boardgame Sultans of Karaya

# 简介
本项目为业余项目，主要目的在于玩…
无明确任务，正在做游戏服务器

招募前端开发、客户端开发，有兴趣的同学可以加入~

```mermaid
graph TD
网页-->|httpproxy|webproxy
小程序-->|wsproxy|webproxy
App-->|tcp|proxy

webproxy-->logic
proxy-->logic

webproxy-->logic2(logic...)
proxy-->logic2(logic...)

logic-->svr
logic2-->svr

logic-->svr2(svr...)
logic2-->svr2(svr...)
```
