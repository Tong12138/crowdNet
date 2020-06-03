package main

import (
	"C"
	"bytes"
	//"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"encoding/hex"

	shell "github.com/ipfs/go-ipfs-api"
)

var sh *shell.Shell

//export UploadIPFS
func UploadIPFS(str *C.char, public_key *C.char) *C.char {
	sh = shell.NewShell("localhost:5001")
	arg := C.GoString(str)
	data, err := ioutil.ReadFile(arg)
	data_en := RSA_encrypter(C.GoString(public_key), []byte(data))
	//data_hex := hex.EncodeToString(data_en)
	hash, err := sh.Add(bytes.NewBufferString(string(data_en)))
	if err != nil {
		fmt.Println("上传ipfs时错误：", err)
	}
	return C.CString(hash)
}

//export CatIPFS
func CatIPFS(hash_c *C.char, filename_c *C.char) *C.char {
	sh = shell.NewShell("localhost:5001")
	hash := C.GoString(hash_c)
	filename := C.GoString(filename_c)
	read, err := sh.Cat(hash)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(read)
	data_de := RSA_decrypter("local_private.pem", body)
	err = ioutil.WriteFile(filename, data_de, 0666)
	if err != nil {
		panic(err)
	}
	result := string(data_de)
	return C.CString(result)
}

//使用公钥进行加密
func RSA_encrypter(path string,msg []byte)[]byte  {
	//首先从文件中提取公钥
	fp,_:=os.Open(path)
	defer fp.Close()
	//测量文件长度以便于保存
	fileinfo,_:=fp.Stat()
	buf:=make([]byte,fileinfo.Size())
	fp.Read(buf)
	//下面的操作是与创建秘钥保存时相反的
	//pem解码
	block,_:=pem.Decode(buf)
	//x509解码,得到一个interface类型的pub
	pub,_:=x509.ParsePKIXPublicKey(block.Bytes)
	//加密操作,需要将接口类型的pub进行类型断言得到公钥类型
	cipherText,_:=rsa.EncryptPKCS1v15(rand.Reader,pub.(*rsa.PublicKey),msg)
	return cipherText
}

//使用私钥进行解密
func RSA_decrypter(path string,cipherText []byte)[]byte  {
	//同加密时，先将私钥从文件中取出，进行二次解码
	fp,_:=os.Open(path)
	defer fp.Close()
	fileinfo,_:=fp.Stat()
	buf:=make([]byte,fileinfo.Size())
	fp.Read(buf)
	block,_:=pem.Decode(buf)
	PrivateKey,_:=x509.ParsePKCS1PrivateKey(block.Bytes)
	//二次解码完毕，调用解密函数
	afterDecrypter,_:=rsa.DecryptPKCS1v15(rand.Reader,PrivateKey,cipherText)
	return afterDecrypter
}

func main() {
	//尝试调用
	msg:=[]byte("RSA非对称加密很棒")
	//Getkeys()
	ciphertext:=RSA_encrypter("csdn_PublicKey.pem",msg)
	//转化为十六进制方便查看结果
	//fmt.Println(ciphertext)
	fmt.Println(hex.EncodeToString(ciphertext))
	result:=RSA_decrypter("csdn_private.pem",ciphertext)
	fmt.Println(string(result))
}

//func main() {
	//生成一个交易结构体(未来的通道)
	//transaction := Transaction{
	//	Person1:      "Aaron",
	//	Person2:      "Bob",
	//	Person1money: "100",
	//	Person2money: "200",
	//}
	//结构体序列化
	//data := marshalStruct(transaction)
	//上传到ipfs /Users/gaojiaxuan/Desktop/ipfs-sample.txt  ./result.txt
	//var filepath string
	/*var resultpath string
	fmt.Print("请输入上传文件：")
	fmt.Scanln(&filepath)
	hash := UploadIPFS(filepath)
	fmt.Println("文件hash是", hash)
	//从ipfs下载数据
	fmt.Print("请输入下载文件：")
	fmt.Scanln(&resultpath)
	str2 := CatIPFS(hash, resultpath)
	//数据反序列化
	//transaction2 := unmarshalStruct([]byte(str2))

	//验证下数据
	fmt.Println(str2)*/
//}