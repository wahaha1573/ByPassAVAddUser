# 绕过杀软添加用户 By Golang

ByPassAVAddUser.exe会默认添加一个用户（testnb$/123@pass）

目前在火绒和360上进行测试，可以绕过

ByPassAVAddUserCli.exe支持自定义用户名和密码

用法：

ByPassAVAddUserCli.exe -u admin -p password

```
Usage of ByPassAVAddUserCli.exe:
  -p string
        密码
  -u string
        用户名
```

![](PictureFile\360.jpg)

![](PictureFile\huorong.jpg)

![](PictureFile\ByPass360.jpg)

![](\PictureFile\ByPassHuoRong.jpg)

## 参考：

https://docs.microsoft.com/zh-cn/windows/win32/api/lmaccess/nf-lmaccess-netuseradd/
https://github.com/CodyGuo/xcgui/blob/master/doc/cToGo.go
https://github.com/iamacarpet/go-win64api/blob/master/users.go
https://git.itch.ovh/itchio/ox/-/blob/ec75be15423d72ab9691a3318ecce3feee67b19b/syscallex/netapi32_windows.go
https://pkg.go.dev/github.com/itchio/ox/syscallex#NetUserAdd