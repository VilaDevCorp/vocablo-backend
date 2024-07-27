package utils

func GetStringPointer(s string) *string {
	return &s
}

func GetResponsePointer(r HttpResponse) *HttpResponse {
	return &r
}
