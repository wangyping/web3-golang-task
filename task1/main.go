package main

import (
	"fmt"
	"strconv"
)

func twoSum(nums []int, target int) []int {
	var result []int = make([]int, 2)
	length := len(nums)
	for i := 0; i < length-1; i++ {
		for j := i + 1; j < length; j++ {
			if nums[i]+nums[j] == target {
				result[0] = i
				result[1] = j
				return result
			}
		}
	}
	return result

	//使用哈希表，可以将寻找 target - x 的时间复杂度降低到从 O(N) 降低到 O(1)。
	//这样我们创建一个哈希表，对于每一个 x，我们首先查询哈希表中是否存在 target - x，然后将 x 插入到哈希表中，即可保证不会让 x 和自己匹配。

	// hashTable := map[int]int{}
	// for i, x := range nums {
	//     if p, ok := hashTable[target-x]; ok {
	//         return []int{p, i}
	//     }
	//     hashTable[x] = i
	// }
	// return nil

}

func singleNumber(nums []int) int {
	m := make(map[int]int)
	for _, v := range nums {
		m[v]++
	}
	for k, v := range m {
		if v == 1 {
			return k
		}
	}
	return -1

	//最优解
	//异或运算有以下三个性质。

	//任何数和 0 做异或运算，结果仍然是原来的数，即 a⊕0=a。
	//任何数和其自身做异或运算，结果是 0，即 a⊕a=0。
	//异或运算满足交换律和结合律，即 a⊕b⊕a=b⊕a⊕a=b⊕(a⊕a)=b⊕0=b。

	// single := 0
	// for _, num := range nums {
	//     single ^= num
	// }
	// return single

}

func isPalindrome(x int) bool {
	str := strconv.Itoa(x)
	strLength := len(str)
	for i := 0; i < strLength/2; i++ {
		if str[i] != str[strLength-1-i] {
			return false
		}
	}
	return true

	// 最优解
	// // 特殊情况：
	// // 如上所述，当 x < 0 时，x 不是回文数。
	// // 同样地，如果数字的最后一位是 0，为了使该数字为回文，
	// // 则其第一位数字也应该是 0
	// // 只有 0 满足这一属性
	// if x < 0 || (x % 10 == 0 && x != 0) {
	//     return false
	// }

	// revertedNumber := 0
	// for x > revertedNumber {
	//     revertedNumber = revertedNumber * 10 + x % 10
	//     x /= 10
	// }

	// // 当数字长度为奇数时，我们可以通过 revertedNumber/10 去除处于中位的数字。
	// // 例如，当输入为 12321 时，在 while 循环的末尾我们可以得到 x = 12，revertedNumber = 123，
	// // 由于处于中位的数字不影响回文（它总是与自己相等），所以我们可以简单地将其去除。
	// return x == revertedNumber || x == revertedNumber / 10

}

func isValid(s string) bool {
	sLength := len(s)
	if sLength%2 != 0 {
		return false
	}
	stack := make([]string, 0)
	var m map[string]string = map[string]string{"(": ")", "[": "]", "{": "}"}
	for i := 0; i < sLength; i++ {
		//入栈
		if m[string(s[i])] != "" {
			//fmt.Println("入栈: ", string(s[i]))
			stack = append(stack, string(s[i]))
		} else {
			//出栈
			if len(stack) == 0 {
				return false
			}
			//fmt.Println("出栈: ", stack[len(stack)-1])
			if m[stack[len(stack)-1]] != string(s[i]) {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

func longestCommonPrefix(strs []string) string {
	fmt.Println("longestCommonPrefix strs: ", strs)
	if len(strs) == 0 || len(strs[0]) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	str := ""
	for i := 1; i <= len(strs[0]); i++ {
		str = strs[0][:i]
		for j := 0; j < len(strs); j++ {
			if len(strs[j]) < len(str) || strs[j][:i] != str {
				return strs[0][:i-1]
			}
		}
	}
	return str
}

func plusOne(digits []int) []int {
	fmt.Println("plusOne digits: ", digits)
	digitsLen := len(digits)
	if digitsLen == 0 {
		return []int{1}
	}
	for i := digitsLen - 1; i >= 0; i-- {
		if digits[i]+1 > 9 {
			digits[i] = 0
		} else {
			digits[i] = digits[i] + 1
			return digits
		}
		if i == 0 {
			digits = append([]int{1}, digits...)
		}
	}
	return digits

}

func removeDuplicates(nums []int) int {
	fmt.Println("before removeDuplicates nums: ", nums)
	numsLen := len(nums)
	if numsLen == 0 {
		return 0
	}
	newNums := nums[:1]
	for i := 1; i < numsLen; i++ {
		if nums[i] == newNums[len(newNums)-1] {
			continue
		}
		newNums = append(newNums, nums[i])
	}
	nums = newNums
	fmt.Println("after removeDuplicates nums: ", nums)
	return len(nums)
}

func merge(intervals [][]int) [][]int {
	fmt.Println("before merge intervals: ", intervals)

	result := [][]int{}
	if len(intervals) == 0 {
		return result
	}
	if len(intervals) == 1 {
		return intervals
	}

	result = intervals[0:1]
	for i := 1; i < len(intervals); i++ {
		result = mergeOther(result, intervals[i])
		fmt.Println(result)
	}
	return result
}

func mergeOther1(intervals [][]int, other []int) [][]int {
	result := [][]int{}
	for i := 0; i < len(intervals); i++ {
		if intervals[i][0] > other[1] {
			result = append(result, intervals[0:i]...)
			result = append(result, [][]int{other}...)
			result = append(result, intervals[i:]...)
		} else if intervals[i][1] < other[0] {
			if i < len(intervals)-1 {
				continue
			}
			result = append(result, intervals[0:i+1]...)
			result = append(result, [][]int{other}...)
		} else {
			intervals[i][0] = min(intervals[i][0], other[0])
			intervals[i][1] = max(intervals[i][1], other[1])
			result = append(result, intervals[0:i+1]...)
			nextIndex := i + 1
			for j := i + 1; j < len(intervals); j++ {
				nextIndex = j
				if intervals[j][0] > result[len(result)-1][1] {
					break
				}
				result = mergeOther(result, intervals[j])
				if j == len(intervals)-1 {
					return result
				}
			}
			if nextIndex > len(intervals)-1 {
				return result
			}
			result = append(result, intervals[nextIndex:]...)
			return result
		}
		return result
	}
	return result
}
func mergeOther(intervals [][]int, other []int) [][]int {
	result := [][]int{}
	startIndex := -1
	step := 0
	for i := 0; i < len(intervals); i++ {
		if startIndex == -1 {
			if intervals[i][0] > other[1] {
				result = append(result, intervals[0:i]...)
				result = append(result, [][]int{other}...)
				result = append(result, intervals[i:]...)
				return result
			} else if intervals[i][1] < other[0] {
				if i < len(intervals)-1 {
					continue
				}
				result = append(result, intervals[0:i+1]...)
				result = append(result, [][]int{other}...)
				return result
			} else {
				if step == 0 {
					startIndex = i
				}
				intervals[i][0] = min(intervals[i][0], other[0])
				intervals[i][1] = max(intervals[i][1], other[1])
				other = intervals[i]
				step++
				continue
			}
		} else {
			if intervals[i][0] > other[1] {
				if step > 1 {
					result = append(result, intervals[0:startIndex]...)
					result = append(result, intervals[i-1:]...)
				} else {
					result = append(result, intervals[0:startIndex+1]...)
					result = append(result, intervals[i:]...)
				}
				return result
			} else {
				if step == 0 {
					startIndex = i
				}
				intervals[i][0] = min(intervals[i][0], other[0])
				intervals[i][1] = max(intervals[i][1], other[1])
				other = intervals[i]
				step++
				// fmt.Println("intervals", intervals)
				// fmt.Println("other", other)
				// fmt.Println("startIndex", startIndex)
				// fmt.Println("step", step)
				continue
			}
		}
	}
	if startIndex >= 0 {
		if step > 1 {
			result = append(result, intervals[0:startIndex]...)
			result = append(result, intervals[startIndex+step-1:]...)
		} else {
			result = append(result, intervals[0:startIndex+1]...)
			result = append(result, intervals[startIndex+step:]...)
		}
	}
	return result
}
func main() {
	//控制流程
	fmt.Println(singleNumber([]int{1, 2, 3, 4, 5, 6, 6, 5, 4, 2, 1}))

	//回文
	fmt.Printf("1.回文：%d %t\n", 121, isPalindrome(121))
	fmt.Printf("2.回文：%d %t\n", -121, isPalindrome(-121))

	//有效的括号
	fmt.Printf("1. 有效括号%s： %t\n", "()", isValid("()"))
	fmt.Printf("2. 有效括号%s： %t\n", "()[]{}", isValid("()[]{}"))
	fmt.Printf("3. 有效括号%s： %t\n", "(]", isValid("(]"))
	fmt.Printf("4. 有效括号%s： %t\n", "([])", isValid("([])"))
	fmt.Printf("5. 有效括号%s： %t\n", "([{)]})", isValid("([{)]})"))
	fmt.Printf("6. 有效括号%s： %t\n", "((", isValid("(("))
	fmt.Printf("7. 有效括号%s： %t\n", "(([]){})", isValid("(([]){})"))
	fmt.Printf("8. 有效括号%s： %t\n", "[({(())}[()])]", isValid("[({(())}[()])]"))

	//最长公共前缀
	fmt.Printf("1. 最长公共前缀：%s\n", longestCommonPrefix([]string{"flower", "flow", "flight"}))
	fmt.Printf("2. 最长公共前缀：%s\n", longestCommonPrefix([]string{"dog", "racecar", "car"}))
	fmt.Printf("3. 最长公共前缀：%s\n", longestCommonPrefix([]string{"flower", "flower", "flower", "flower"}))
	fmt.Printf("4. 最长公共前缀：%s\n", longestCommonPrefix([]string{"", "b"}))

	//加1
	fmt.Printf("1. 加一：%v\n", plusOne([]int{1, 2, 3}))
	fmt.Printf("2. 加一：%v\n", plusOne([]int{4, 3, 2, 9}))
	fmt.Printf("3. 加一：%v\n", plusOne([]int{0}))
	fmt.Printf("4. 加一：%v\n", plusOne([]int{9, 9, 9, 9}))

	//删除有序数组中的重复项
	fmt.Printf("1. 删除有序数组中的重复项：%v\n", removeDuplicates([]int{1, 1, 2}))
	fmt.Printf("2. 删除有序数组中的重复项：%v\n", removeDuplicates([]int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}))
	fmt.Printf("3. 删除有序数组中的重复项：%v\n", removeDuplicates([]int{1, 1, 1, 2, 2, 3}))

	//合并区间
	fmt.Printf("1. 合并区间：%v\n", merge([][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}))
	fmt.Printf("2. 合并区间：%v\n", merge([][]int{{1, 4}, {4, 5}}))
	fmt.Printf("3. 合并区间：%v\n", merge([][]int{{1, 4}, {4, 8}, {5, 10}}))
	fmt.Printf("4. 合并区间：%v\n", merge([][]int{{1, 4}, {5, 6}}))
	fmt.Printf("5. 合并区间：%v\n", merge([][]int{{1, 4}, {0, 6}}))
	fmt.Printf("6. 合并区间：%v\n", merge([][]int{{1, 4}, {0, 1}}))
	fmt.Printf("7. 合并区间：%v\n", merge([][]int{{1, 4}, {0, 0}}))
	fmt.Printf("8. 合并区间：%v\n", merge([][]int{{1, 4}, {0, 2}, {3, 5}}))
	fmt.Printf("9. 合并区间：%v\n", merge([][]int{{2, 3}, {5, 5}, {2, 2}, {3, 4}, {3, 4}}))
	fmt.Printf("10. 合并区间：%v\n", merge([][]int{{2, 3}, {4, 5}, {6, 7}, {8, 9}, {1, 10}}))
	fmt.Printf("11. 合并区间：%v\n", merge([][]int{{4, 5}, {2, 4}, {4, 6}, {3, 4}, {0, 0}, {1, 1}, {3, 5}, {2, 2}}))
	fmt.Printf("12. 合并区间：%v\n", merge([][]int{{4, 5}, {2, 4}, {4, 6}, {3, 4}, {5, 5}, {1, 1}, {3, 5}, {2, 2}}))
	fmt.Printf("13. 合并区间：%v\n", merge([][]int{{0, 2}, {2, 3}, {4, 4}, {0, 1}, {5, 7}, {4, 5}, {0, 0}}))
	fmt.Printf("14. 合并区间：%v\n", merge([][]int{{0, 0}, {0, 0}, {4, 4}, {0, 0}, {1, 3}, {5, 5}, {4, 6}, {1, 1}, {0, 2}}))

	//两数之和
	fmt.Println(twoSum([]int{1, 2, 3}, 4))

}
