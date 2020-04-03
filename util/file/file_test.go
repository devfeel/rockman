package file

import (
	"fmt"
	"path/filepath"
	"testing"
)

// 以下是功能测试

func Test_GetCurrentDirectory_1(t *testing.T) {
	thisDir := GetCurrentDirectory()
	t.Log(thisDir)
}

func Test_GetFileExt_1(t *testing.T) {
	fn := "/download/vagrant_1.9.2.dmg"
	fileExt := filepath.Ext(fn)
	if len(fileExt) < 1 {
		t.Error("fileExt null!")
	} else {
		t.Log(fileExt)
	}
}

func Test_GetFileExt_2(t *testing.T) {
	fn := "/download/vagrant_1.abc"
	fileExt := filepath.Ext(fn)
	if len(fileExt) < 1 {
		t.Error("fileExt null!")
	} else {
		t.Log(fileExt)
	}
}

func Test_Exist_1(t *testing.T) {
	fn := "testdata/file.test"
	//	fn := "/Users/kevin/Downloads/commdownload.dmg"
	isExist := Exists(fn)
	if isExist {
		t.Log(isExist)
	} else {
		t.Log("请修改测试代码中文件的路径！")
	}
}

func Test_ExistsInPath(t *testing.T) {
	fmt.Println(ExistsInPath("../../", "../../main.go"))
}

// 以下是性能测试
func BenchmarkExistsInPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExistsInPath("../../", "../../main.go")
	}
}
