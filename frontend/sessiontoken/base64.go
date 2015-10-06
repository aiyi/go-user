package sessiontoken

// 去掉 base64 编码尾部的 '='.
//  dst 引用了 src 的空间.
func base64Trim(src []byte) (dst []byte) {
	for len(src) > 0 {
		lastIndex := len(src) - 1 // lastIndex >= 0
		if src[lastIndex] != '=' {
			break
		}
		src = src[:lastIndex]
	}
	return src
}

// 在编码尾部填充 '=', 使之符合 base64 编码规则.
//  dst 可能引用了 src 的空间, src 之后的空间可能被污染(当 cap(src) > len(src) 时).
func base64Pad(src []byte) (dst []byte) {
	n := len(src) & 0x3
	if n == 0 {
		return src
	}
	switch n = 4 - n; n {
	case 1:
		dst = append(src, '=')
	case 2:
		dst = append(src, "=="...)
	case 3:
		dst = append(src, "==="...)
	}
	return
}
