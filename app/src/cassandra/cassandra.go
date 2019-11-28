package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

//Session Tạo một session
var Session *gocql.Session

//Hàm khởi tạo
func init() {
	cluster := gocql.NewCluster("cassandra.cassandra:9042")
	cluster.Keyspace = "shortenurl"
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "cassandra",
	}
	var err error
	Session, err = cluster.CreateSession()
	if err != nil {
		Session.Close()
		panic(err)
	}
	fmt.Println("cassandra init done")
	return
}

//Get Trả về URL tương ứng với shortcode trả về lỗi nếu có
func Get(shortcode string) (string, error) {
	m := map[string]interface{}{}
	query := "SELECT longurl FROM url WHERE shortcode = ?"
	iterable := Session.Query(query, shortcode).Iter()
	if iterable.MapScan(m) {
		return m["longurl"].(string), iterable.Close()
	}
	return "", iterable.Close()
}

//Insert Ghi vào database trả về lỗi nếu có
func Insert(shortcode, longurl string) error {
	query := "INSERT INTO url(shortcode,longurl) VALUES (?,?)"
	return Session.Query(query, shortcode, longurl).Exec()
}
