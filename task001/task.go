// 包声明
package main

import (
	"fmt"
	"sort"
	"strconv"
)

// 引入包声明

// 136. 只出现一次的数字
func singleNumber(nums []int) int {

	m := make(map[int]int)

	for i := range nums {
		m[nums[i]]++
	}
	for key, value := range m {
		if value == 1 {
			return key
		}
	}

	return 0

}

// 回文数
func isPalindrome(x int) bool {
	s := strconv.Itoa(x)

	r := []rune(s)

	rLen := len(r)

	for i, v := range r {
		if v != r[rLen-i-1] {
			return false
		}
	}
	return true
}

// 20. 有效的括号
func isValid(s string) bool {
	n := len(s)
	if n < 2 && n%2 == 1 {
		return false
	}
	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}
	stack := []byte{}

	stack = append(stack, s[0])
	nextIndex := 1
	for i := 1; i < n; i++ {
		pv, pkf := pairs[s[i]]
		if pkf {
			if nextIndex < 1 {
				return false
			}
			if pv == stack[nextIndex-1] {
				nextIndex--
			} else {
				return false
			}
		} else {
			if nextIndex < len(stack) {
				stack[nextIndex] = s[i]
			} else {
				stack = append(stack, s[i])
			}
			nextIndex++
		}
	}
	return nextIndex == 0
}

// 14. 最长公共前缀
func longestCommonPrefix(strs []string) string {

	maxLen := getMinLen(strs)
	if maxLen == 0 {
		return ""
	}
	rowLen := len(strs)
	for i := 0; i < maxLen; i++ {
		for j := 0; j < rowLen; j++ {
			if strs[0][i] != strs[j][i] {
				return strs[0][0:i]
			}
		}
	}
	return strs[0][0:maxLen]

}
func getMinLen(strs []string) int {
	minLen := len(strs[0])
	for i := range strs {
		currLen := len(strs[i])
		if minLen > currLen {
			minLen = currLen
		}
	}
	return minLen
}

// 66. 加一
func plusOne(digits []int) []int {
	carry := true
	n := len(digits)
	for i := n - 1; i >= 0; i-- {
		if carry {
			temp := digits[i] + 1
			if temp > 9 {
				carry = true
				digits[i] = 0
			} else {
				carry = false
				digits[i]++
			}
		}
	}
	if carry {
		result := make([]int, n+1)
		result[0] = 1
		return result
	}
	return digits
}

// 26. 删除有序数组中的重复项
func removeDuplicates(nums []int) int {
	n := len(nums)
	if n < 2 {
		return n
	}
	i, j, k := 0, 1, 1
	for j < n {
		if nums[i] == nums[j] {
			j++
		} else {

			nums[i+1] = nums[j]
			i++
			k++

		}
	}
	return k

}

// 56. 合并区间
func merge(intervals [][]int) [][]int {
	sort.Slice(intervals, func(i, j int) bool {
		// 比较第i行和第j行的第一个元素
		return intervals[i][0] < intervals[j][0]
	})
	row := len(intervals)
	temp := make([]int, 2)

	haveData := false

	var result [][]int

	for i := 0; i < row; i++ {

		if !haveData {
			temp[0] = intervals[i][0]
			temp[1] = intervals[i][1]
			haveData = true
			continue
		}

		if temp[1] < intervals[i][0] {

			result = append(result, []int{temp[0], temp[1]})
			haveData = false
			i--
			continue
		}
		if temp[1] < intervals[i][1] {
			temp[1] = intervals[i][1]
		}
	}

	if haveData {
		result = append(result, []int{temp[0], temp[1]})
	}

	return result
}

// 1. 两数之和
func twoSum(nums []int, target int) []int {
	tempMap := make(map[int]int)
	for i, v := range nums {
		if p, ok := tempMap[target-v]; ok {
			return []int{p, i}
		}
		tempMap[v] = i
	}
	return nil
}

func main() {

	fmt.Println("result")
}
