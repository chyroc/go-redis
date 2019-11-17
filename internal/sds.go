package internal

// sds 简单动态字符串
type sdshdr struct {
	len  int
	free int
	buf  []byte
}
