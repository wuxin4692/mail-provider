mail-provider
=============

把smtp封装为一个简单http接口，配置到sender中用来发送报警邮件

## 使用方法

```
curl http://$ip:4000/sender/mail -d "tos=a@a.com,b@b.com&subject=xx&content=yy"
```

=============
2017 06 16 修改明文smtp为startls，只是支持587端口的smtp发送。
