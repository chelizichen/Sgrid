package saas

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"strings"
	"time"

	"fmt"
	"os"
)

const saasPermKey = "com.cloud.hellosgrid.server"

var UsesaasPerm = new(saasPerm)

type saasPerm struct {
	StartTime string
	EndTime   string
	Username  string
	Password  string
	_authstr  string `json:"-"`
}

func (s *saasPerm) CheckAuth() (bool, error) {
	if len(os.Args) != 2 {
		panic(`
		********* error **********
		鉴权失败！请联系作者进行软件购买
		phone: 13476973442
		email: leemulus21@gmail.com
		********* error **********
		`)
	}
	// 通过该 key 得到 开始时间，结束时间
	authString := os.Args[1]
	s.SetAuthString(authString)
	parseString, err := s.ParseAuthString()
	if err != nil {
		panic(`
		********* error **********
		鉴权Token失败！请联系作者进行软件购买
		phone: 13476973442
		email: leemulus21@gmail.com
		********* error **********
		`)
	}
	s.SetAuthBody(parseString)
	fmt.Println("authString", authString)
	isExpire, err := s.CheckExpireTime()
	if err != nil || isExpire {
		panic(`
		********* error **********
		凭证已过期或失效！请联系作者进行软件购买
		phone: 13476973442
		email: leemulus21@gmail.com
		********* error **********
		`)
	}
	return true, nil
}

func (s *saasPerm) SetAuthString(str string) {
	s._authstr = str
}

func (s *saasPerm) GenAuthString() (string, error) {
	plainText := fmt.Sprintf("%v\n%v\n%v\n%v", s.Username, s.Password, s.StartTime, s.EndTime)
	key := []byte(saasPermKey)[:16] // 确保密钥长度为16字节，适应AES-128
	ciphertextBytes, err := encrypt([]byte(plainText), key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertextBytes), nil
}

// ParseAuthString 解析并解密认证字符串
func (s *saasPerm) ParseAuthString() (string, error) {
	key := []byte(saasPermKey)[:16] // 确保密钥长度为16字节
	cipherTextBytes, err := base64.StdEncoding.DecodeString(s._authstr)
	if err != nil {
		return "", err
	}
	plainText, err := decrypt(cipherTextBytes, key)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

func (s *saasPerm) SetAuthBody(str string) {
	arguments := strings.Split(str, "\n")
	if len(arguments) != 4 {
		panic("error.length not equal 4")
	}
	s.Username = arguments[0]
	s.Password = arguments[1]
	s.StartTime = arguments[2]
	s.EndTime = arguments[3]

}

func (s *saasPerm) CheckExpireTime() (bool, error) {
	t, err := time.Parse(time.DateTime, s.EndTime)
	if err != nil {
		fmt.Println("sassPerm.CheckExpireTime", err.Error())
		return false, err
	}
	n := time.Now()

	isExpire := n.After(t)
	return isExpire, nil
}

func decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 如果密文是Base64编码的，先解码
	ciphertext, err = base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)
	cipher.NewCTR(block, iv).XORKeyStream(plaintext, ciphertext[aes.BlockSize:])
	return plaintext, nil
}

func encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	cipher.NewCTR(block, iv).XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	// 可选：返回Base64编码的密文，便于文本传输
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}
