package sessiontoken

// 去掉 base64 编码尾部的 '='.
//  dst 引用了 src 的空间.
func base64Trim(src []byte) (dst []byte) {
	for len(src) > 0 {
		lastIndex := len(src) - 1 // lastIndex >= 0
		if src[lastIndex] == '=' {
			src = src[:lastIndex]
			continue
		}
		break
	}
	dst = src
	return
}

// 在编码尾部填充 '=', 使之符合 base64 编码规则.
//  dst 可能引用了 src 的空间.
func base64Pad(src []byte) (dst []byte) {
	srcLen := len(src)
	n := srcLen & 0x3
	if n == 0 {
		return src
	}
	src = src[:srcLen:srcLen]
	switch n = 4 - n; n {
	case 1:
		dst = append(src, "="...)
	case 2:
		dst = append(src, "=="...)
	case 3:
		dst = append(src, "==="...)
	}
	return
}
