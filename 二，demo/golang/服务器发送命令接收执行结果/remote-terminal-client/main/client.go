package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"

	"remote-terminal-client/protocol"

	"github.com/gorilla/websocket"
	"github.com/kr/pty"
)

const (
	stateOpen uint32 = iota
	stateConnecting
	stateReady
	stateReconnecting
	stateClose
)

const (
	EventOpen      string = "open"
	EventConnect          = "connect"
	EventReconnect        = "reconnect"
	EventError            = "error"
)

type Listener func(args ...interface{})

type emitter struct {
	listeners map[string][]Listener
	m         sync.RWMutex
}

type option struct {
	AutoReconnect    bool
	MaxReconnections int32
}

var defaultOption = &option{
	AutoReconnect:    true,
	MaxReconnections: math.MaxInt32,
}

type socketClient struct {
	emitter
	state     uint32
	url       *url.URL
	option    *option
	transprot protocol.Transport
	outChan   chan *protocol.Packet
	closeChan chan bool
}

var connection *websocket.Conn

//---------------------虚拟终端--------------------------
type wsPty struct {
	Cmd *exec.Cmd // pty builds on os.exec
	Pty *os.File  // a pty is simply an os.File
}

func (wp *wsPty) Start() {
	var err error
	args := flag.Args()
	wp.Cmd = exec.Command(cmdFlag, args...)
	wp.Cmd.Env = append(os.Environ(), "TERM=xterm")
	wp.Pty, err = pty.Start(wp.Cmd)
	if err != nil {
		log.Fatalf("Failed to start command: %s\n", err)
	}
}

func (wp *wsPty) Stop() {
	wp.Pty.Close()
	wp.Cmd.Wait()
}

var cmdFlag string
var messageData interface{}

func init() {
	flag.StringVar(&cmdFlag, "cmd", "/bin/bash", "command to execute on slave side of the pty")
}

func main() {

	wp := wsPty{}
	wp.Start()

	var conHd = make(map[string]*websocket.Conn)

	fmt.Println(RsaEncrypt([]byte("aiyouwei")))
	var Header http.Header = map[string][]string{
		"moja":     {"ccccc, asdasdasdasd"},
		"terminal": {"en-esadasdasdwrw"},
		"success":  {"dasdadas", "wdsadaderew"},
		"ticket":   {RsaEncrypt([]byte("aiyouwei"))},
	}

	s, err := Socket("ws://127.0.0.1:3000")
	if err != nil {
		panic(err)
	}
	//s.Connect(Header)
	//建立主连接
	if atomic.CompareAndSwapUint32(&s.state, stateOpen, stateConnecting) {
		conn, c, err := s.transprot.Dial(s.url.String(), Header)
		connection = c
		if err != nil {
			s.emit(EventError, err)
			go s.reconnect(stateConnecting, Header)
			return
		}
		if atomic.CompareAndSwapUint32(&s.state, stateConnecting, stateReady) {
			go s.start(conn, Header)
			s.emit(EventConnect)
		} else {
			conn.Close()
		}
	}
	//建立子连接
	go func() {
		for {

			//每次轮训需要判断连接句柄是否存在

			//s, _ := ParseString(messageData)

			//fmt.Println("bbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
			// in := []byte(s)
			// var raw = make(map[string]interface{})
			// json.Unmarshal(in, &raw)
			// fmt.Println(raw["subconn"])
			if messageData == "subconn" {
				sub, err := Socket("ws://127.0.0.1:3000?a=sub")
				if err != nil {
					panic(err)
				}
				if atomic.CompareAndSwapUint32(&sub.state, stateOpen, stateConnecting) {
					subConn, c, err := sub.transprot.Dial(sub.url.String(), Header)
					conHd["1"] = c
					if err != nil {
						sub.emit(EventError, err)
						go sub.reconnect(stateConnecting, Header)
						return
					}
					if atomic.CompareAndSwapUint32(&sub.state, stateConnecting, stateReady) {
						go sub.start(subConn, Header)
						sub.emit(EventConnect)
					} else {
						subConn.Close()
					}
				}

				fmt.Println("pppppppppppppppppppppppppppp")
				fmt.Println(conHd["1"])
				sub.On("message", func(args ...interface{}) {
					enResult, _ := ParseString(args[0])
					messageData = DecryptWithAES("asdasdasdasdasd", enResult)
					//fmt.Println(cmd)
					//wp.Pty.Write([]byte(cmd))
				})
			} else if messageData == "cmd" {
				fmt.Println("wqeqweqwqw")
			} else {
				//	fmt.Println("qweqwerrrrtytyyyqwwetrtyutuiop")
				// decodeBytes, err := base64.StdEncoding.DecodeString(s)
				// if err != nil {
				// 	log.Fatalln(err)
				// }
				//	fmt.Println(string(decodeBytes))
			}
		}
	}()

	input := []byte("testtttt")
	// 演示base64编码
	encodeString := base64.StdEncoding.EncodeToString(input)
	s.Emit("messgae", encodeString)
	//主连接接收消息类型
	s.On("message", func(args ...interface{}) {
		enResult, _ := ParseString(args[0])
		messageData = DecryptWithAES("asdasdasdasdasd", enResult)
		//fmt.Println(cmd)
		//wp.Pty.Write([]byte(cmd))
	})

	go func() {
		resBuf := make([]byte, 1024)
		for {
			fmt.Println(string(resBuf))
			n, err := wp.Pty.Read(resBuf)
			if err != nil {
				log.Printf("Failed to read from pty master: %s", err)
				return
			}
			out := make([]byte, base64.StdEncoding.EncodedLen(n))
			base64.StdEncoding.Encode(out, resBuf[0:n])
			s.Emit("result", string(resBuf[0:n]))
		}

	}()

	for {

	}
}

func (e *emitter) On(event string, listener Listener) {
	e.m.Lock()
	defer e.m.Unlock()
	listeners, ok := e.listeners[event]
	if ok {
		listeners = append(listeners, listener)
	} else {
		listeners = []Listener{listener}
	}
	e.listeners[event] = listeners
}

func (e *emitter) emit(event string, args ...interface{}) bool {
	e.m.RLock()
	listeners, ok := e.listeners[event]
	if ok {
		for _, listener := range listeners {
			listener(args...)
		}
	}
	e.m.RUnlock()
	return ok
}

func Socket(urlstring string) (*socketClient, error) {
	u, err := url.Parse(urlstring)
	if err != nil {
		return nil, err
	}
	u.Path = "/socket.io/"
	q := u.Query()
	q.Add("EIO", "3")
	q.Add("transport", "websocket")
	u.RawQuery = q.Encode()
	return &socketClient{
		emitter:   emitter{listeners: make(map[string][]Listener)},
		url:       u,
		option:    defaultOption,
		transprot: protocol.NewWebSocketTransport(),
		outChan:   make(chan *protocol.Packet, 64),
		closeChan: make(chan bool),
	}, nil
}

func (s *socketClient) Connect(requestHeader http.Header) {
	if atomic.CompareAndSwapUint32(&s.state, stateOpen, stateConnecting) {
		conn, c, err := s.transprot.Dial(s.url.String(), requestHeader)
		connection = c
		if err != nil {
			s.emit(EventError, err)
			go s.reconnect(stateConnecting, requestHeader)
			return
		}
		if atomic.CompareAndSwapUint32(&s.state, stateConnecting, stateReady) {
			go s.start(conn, requestHeader)
			s.emit(EventConnect)
		} else {
			conn.Close()
		}
	}
}

func (s *socketClient) Disconnect() {
	atomic.StoreUint32(&s.state, stateClose)
	close(s.outChan)
	close(s.closeChan)
}

func (s *socketClient) Emit(event string, args ...interface{}) {
	if atomic.LoadUint32(&s.state) == stateReady && !s.emit(event, args) {
		m := &protocol.Message{
			Type:      protocol.MessageTypeEvent,
			Namespace: "/",
			ID:        -1,
			Event:     event,
			Payloads:  args,
		}
		p, err := m.Encode()
		if err != nil {
			s.emit(EventError, err)
		} else {
			s.outChan <- p
		}
	}
}

func (s *socketClient) reconnect(state uint32, requestHeader http.Header) {
	time.Sleep(time.Second)
	if atomic.CompareAndSwapUint32(&s.state, state, stateReconnecting) {
		conn, c, err := s.transprot.Dial(s.url.String(), requestHeader)
		connection = c
		if err != nil {
			s.emit(EventError, err)
			go s.reconnect(stateReconnecting, requestHeader)
			return
		}
		if atomic.CompareAndSwapUint32(&s.state, stateReconnecting, stateReady) {
			go s.start(conn, requestHeader)
			s.emit(EventReconnect)
		} else {
			conn.Close()
		}
	}
}

func (s *socketClient) start(conn protocol.Conn, requestHeader http.Header) {
	stopper := make(chan bool)
	go s.startRead(conn, stopper)
	go s.startWrite(conn, stopper)
	select {
	case <-stopper:
		go s.reconnect(stateReady, requestHeader)
		conn.Close()
	case <-s.closeChan:
		conn.Close()
	}
}

func (s *socketClient) startRead(conn protocol.Conn, stopper chan bool) {
	defer func() {
		recover()
	}()
	for atomic.LoadUint32(&s.state) == stateReady {
		p, err := conn.Read()
		if err != nil {
			s.emit(EventError, err)
			close(stopper)
			return
		}
		switch p.Type {
		case protocol.PacketTypeOpen:
			h, err := p.DecodeHandshake()
			if err != nil {
				s.emit(EventError, err)
			} else {
				go s.startPing(h, stopper)
			}
		case protocol.PacketTypePing:
			s.outChan <- protocol.NewPongPacket()
		case protocol.PacketTypeMessage:
			m, err := p.DecodeMessage()
			if err != nil {
				s.emit(EventError, err)
			} else {
				s.emit(m.Event, m.Payloads...)
			}
		}
	}
}

func (s *socketClient) startWrite(conn protocol.Conn, stopper chan bool) {
	defer func() {
		recover()
	}()
	for atomic.LoadUint32(&s.state) == stateReady {
		select {
		case <-stopper:
			return
		case p, ok := <-s.outChan:
			if !ok {
				return
			}
			err := conn.Write(p)
			if err != nil {
				s.emit(EventError, err)
				close(stopper)
				return
			}
		}

	}
}

func (s *socketClient) startPing(h *protocol.Handshake, stopper chan bool) {
	defer func() {
		recover()
	}()
	for {
		time.Sleep(time.Duration(h.PingInterval) * time.Millisecond)
		select {
		case <-stopper:
			return
		case <-s.closeChan:
			return
		default:
		}
		if atomic.LoadUint32(&s.state) != stateReady {
			return
		}
		s.outChan <- protocol.NewPingPacket()
	}
}

func EncryptWithAES(key, message string) string {

	hash := md5.New()
	hash.Write([]byte(key))
	keyData := hash.Sum(nil)

	block, err := aes.NewCipher(keyData)
	if err != nil {
		panic(err)
	}

	iv := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

	enc := cipher.NewCBCEncrypter(block, iv)
	content := PKCS5Padding([]byte(message), block.BlockSize())
	crypted := make([]byte, len(content))
	enc.CryptBlocks(crypted, content)
	return base64.StdEncoding.EncodeToString(crypted)
}

func DecryptWithAES(key, message string) string {

	hash := md5.New()
	hash.Write([]byte(key))
	keyData := hash.Sum(nil)

	block, err := aes.NewCipher(keyData)
	if err != nil {
		panic(err)
	}

	iv := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

	messageData, _ := base64.StdEncoding.DecodeString(message)
	dec := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(messageData))
	dec.CryptBlocks(decrypted, messageData)
	return string(PKCS5Unpadding(decrypted))
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Unpadding(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]

}
func ParseString(value interface{}) (string, error) {
	switch value.(type) {
	case string:
		return value.(string), nil
	default:
		return "", fmt.Errorf("unable to casting number %v (type %T)", value, value)
	}
}

func RsaEncrypt(data []byte) string {
	pubKey, err := ioutil.ReadFile("../public.pem")
	if err != nil {
		log.Fatal(err.Error())
	}

	block, _ := pem.Decode(pubKey) //将密钥解析成公钥实例
	if block == nil {
		fmt.Println("public key error")
		return ""
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	if err != nil {
		fmt.Println(err)
		return ""
	}
	pub := pubInterface.(*rsa.PublicKey)
	res, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pub, data, nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(res)
}
