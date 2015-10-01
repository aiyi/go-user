package base64

// 去掉 base64 编码最后面的 '=', dst 引用了 src 的空间.
func Trim(src []byte) (dst []byte) {
	for length := len(src); length > 0; length = len(src) {
		lastIndex := length - 1
		if src[lastIndex] == '=' {
			src = src[:lastIndex]
			continue
		}
		break
	}
	dst = src
	return
}

// 填充 base64 最后面的 '=', dst 可能引用了 src 的空间.
func Pad(src []byte) (dst []byte) {
	n := len(src) & 0x3
	if n == 0 {
		return src
	}
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
