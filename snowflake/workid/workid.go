/**
 * @Author:huangshan
 * @Description:
 * @Date: 2022/11/3 16:21
 */
package workid

type WorkId interface {
	GetWorkId() (int64, error)
}
