/**
 * @Author:huangshan
 * @Description:
 * @Date: 2022/11/3 10:08
 */
package snowflake

import (
	"fmt"
	"neuroxess-cloud/library/snowflake/workid"
	"testing"
)

func Test_turn(t *testing.T) {
	fmt.Println(turn(41, "毫秒"))
	fmt.Println(turn(41, "10毫秒"))
	fmt.Println(turn(41, "100毫秒"))
	fmt.Println(turn(41, "秒"))
	fmt.Println(turn(40, "毫秒"))
	fmt.Println(turn(40, "10毫秒"))
	fmt.Println(turn(40, "100毫秒"))
	fmt.Println(turn(40, "秒"))
	fmt.Println(turn(39, "毫秒"))
	fmt.Println(turn(39, "10毫秒"))
	fmt.Println(turn(39, "100毫秒"))
	fmt.Println(turn(39, "秒"))
	fmt.Println(turn(38, "毫秒"))
	fmt.Println(turn(38, "10毫秒"))
	fmt.Println(turn(38, "100毫秒"))
	fmt.Println(turn(38, "秒"))
	fmt.Println(turn(37, "毫秒"))
	fmt.Println(turn(37, "10毫秒"))
	fmt.Println(turn(37, "100毫秒"))
	fmt.Println(turn(37, "秒"))
	fmt.Println(turn(36, "毫秒"))
	fmt.Println(turn(36, "10毫秒"))
	fmt.Println(turn(36, "100毫秒"))
	fmt.Println(turn(36, "秒"))
	fmt.Println(turn(35, "毫秒"))
	fmt.Println(turn(35, "10毫秒"))
	fmt.Println(turn(35, "100毫秒"))
	fmt.Println(turn(35, "秒"))

}

func TestNextVal(t *testing.T) {
	if err := Init(&workid.Reg{}); err != nil {
		panic(err)
	}
	cnt := 10
	nums := make([]int64, cnt)
	for i := 0; i < cnt; i++ {
		nums[i] = NextVal()
	}
	ha := make(map[int64]int64)
	for i, v := range nums {
		if ha[v] > 0 {
			k, kk := i, ha[v]
			panic(fmt.Sprintf("%d %d", k, kk))
		} else {
			ha[v] = int64(i)
		}
	}
}

func TestInit(t *testing.T) {
	fmt.Println(Init(&workid.Random{}))
}

func TestGetWorkId(t *testing.T) {
	for {
		fmt.Println(GetWorkId(&workid.Reg{}))
	}
}
