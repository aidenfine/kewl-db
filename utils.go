package main


func IsAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || ( c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
