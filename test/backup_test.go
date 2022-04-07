package test

import (
	"fmt"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
)

func TestRel(t *testing.T) {
	_, caches, err := nginx.Load("./config_test/nginx.conf")
	if err != nil {
		t.Log(err)
	}

	//fileList, err := conf.List()
	//if err != nil {
	//	t.Log(err)
	//}

	//for _, s := range fileList {
	for s := range caches {
		t.Log(filepath.Rel("F:\\GO_Project\\src\\bifrost\\test", s))
	}
}

func checkBackups(name, dir string, saveTime, cycle int, now time.Time) (bool, error) {
	needBackup := true
	saveDate := now.Add(-24 * time.Hour * time.Duration(saveTime))
	cycleDate := now.Add(-24 * time.Hour * time.Duration(cycle))
	bakFilePattern := fmt.Sprintf(`^%s\.(\d{8})\.tgz$`, name)
	bakFileReg := regexp.MustCompile(bakFilePattern)

	baks, gErr := filepath.Glob(filepath.Join(dir, fmt.Sprintf("%s.*.tgz", name)))
	if gErr != nil {
		return false, gErr
	}

	for i := 0; i < len(baks) && needBackup; i++ {
		bakName := filepath.Base(baks[i])
		if isBak := bakFileReg.MatchString(bakName); isBak {
			bakDate, tpErr := time.Parse("20060102", bakFileReg.FindStringSubmatch(bakName)[1])
			if tpErr != nil {
				return false, tpErr
			}

			if bakDate.Unix() < saveDate.Unix() {
				//rmErr := os.Remove(baks[i])
				//if rmErr != nil {
				//	return false, rmErr
				//}
				fmt.Println("remove:", baks[i])
			}

			if bakDate.Unix() > cycleDate.Unix() || bakDate.Format("20060102") == now.Format("20060102") {
				fmt.Println("no backup required, cased by", baks[i], "is exist.")
				needBackup = false
			}

		}
	}

	return needBackup, nil
}

func TestCheckBackups(t *testing.T) {
	name := "nginx.conf"
	dir := `F:/Code_Path/src/bifrost/test/tgz_test`
	saveTime := 7
	cycle := 2
	now := time.Now()
	needBak, checkErr := checkBackups(name, dir, saveTime, cycle, now)
	if checkErr != nil {
		t.Log(checkErr)
	}
	fmt.Println(needBak)
}
