package common

func Ternary(cond bool, res1, res2 any) any {
	if cond {
		return res1
	} else {
		return res2
	}
}
