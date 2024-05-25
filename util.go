/*
@author: sk
@date: 2024/5/16
*/
package main

import (
	"image/color"
	"strconv"
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MustFloat(val string) float64 {
	res, err := strconv.ParseFloat(val, 64)
	HandleErr(err)
	return res
}

// 只支持 6 位
func MustClr(val string) color.Color {
	num, err := strconv.ParseUint(val, 16, 64)
	HandleErr(err)
	return color.RGBA{
		R: uint8((num >> 16) & 0xFF),
		G: uint8((num >> 8) & 0xFF),
		B: uint8(num & 0xFF),
		A: 255,
	}
}
