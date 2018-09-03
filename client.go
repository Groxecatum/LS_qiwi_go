package golang_commons

type Client struct {
	Id    int    `xml:"id"`
	Name  string `xml:"name"`
	Email string `xml:"email"`
	Phone string `xml:"cellPhone"`
}
