mail-provider
=============

把smtp封装为一个简单http接口，配置到sender中用来发送报警邮件

## 使用方法

```
curl http://$ip:4000/sender/mail -d "tos=a@a.com,b@b.com&subject=xx&content=yy"
```

zzlyzq修改为startls版本
=============
2017 06 16 修改明文smtp为startls，只是支持587端口的smtp发送。
zzlyzq修改为startls版本


wuxin4692的修改
=============
2018/10/23

1.增加html邮件模板
