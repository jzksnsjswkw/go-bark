package bark

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"
)

type Client struct {
	// 服务器URL
	ServerURL string
}

type Options struct {
	Client `json:"-"`

	// 推送内容 (必填)
	Msg string `json:"body"`
	// token (必填)
	Token string `json:"-"`
	// 推送标题
	Title string `json:"title,omitempty"`
	// 消息分组
	Group string `json:"group,omitempty"`
	// 点击推送时，跳转的URL，支持URL Scheme 和 Universal Link
	URL string `json:"url,omitempty"`
	// 推送中断级别
	Level string `json:"level,omitempty"`
	// 指定复制的内容
	Copy string `json:"copy,omitempty"`
	// 自动复制
	AutoCopy bool `json:"autoCopy,omitempty"`
	// 推送铃声
	Sound string `json:"sound,omitempty"`
	// 自定义图标，传入URL
	Icon string `json:"icon,omitempty"`
	// 推送角标
	Badge int `json:"badge,omitempty"`
	// 传 1 保存推送，传其他的不保存推送，不传按APP内设置来决定是否保存。
	IsArchive int `json:"isArchive,omitempty"`

	// 加密传输
	Enc *EncOpt `json:"-"`
}

type EncOpt struct {
	Mode EncMode
	Key  string
	Iv   string
}

const DefaultDomain = "api.day.app"
const DefaultURL = "https://" + DefaultDomain

var DefaultClient = New(DefaultURL)

func New(url string) *Client {
	return &Client{
		ServerURL: strings.TrimSuffix(url, "/"),
	}
}

func Push(o *Options) error {
	return DefaultClient.Push(o)
}

func (c *Client) Push(o *Options) error {
	s, err := handleOpt(o)
	if err != nil {
		return err
	}
	var r *http.Response
	if o.ServerURL == "" {
		r, err = http.Post(c.ServerURL+"/"+o.Token, "application/json;charset:utf-8", strings.NewReader(s))
	} else {
		r, err = http.Post(o.ServerURL+"/"+o.Token, "application/json;charset:utf-8", strings.NewReader(s))
	}
	if err != nil {
		return err
	}
	defer r.Body.Close()

	s2 := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&s2); err != nil {
		return err
	}
	if s2.Code != 200 {
		return errors.New(s2.Message)
	}

	return nil
}

func barkEncrypt(e *EncOpt, s []byte) (string, error) {
	var c crypto.Cryptobin
	if strings.ToUpper(e.Mode) == ECB {
		c = crypto.FromBytes(s).SetKey(e.Key).Aes().ECB().PKCS7Padding().Encrypt()
	} else if strings.ToUpper(e.Mode) == CBC {
		c = crypto.FromBytes(s).SetKey(e.Key).SetIv(e.Iv).Aes().CBC().PKCS7Padding().Encrypt()
	} else {
		return "", errors.New("enc mode must be ECB or CBC")
	}
	if err := c.Error(); err != nil {
		return "", err
	}
	return c.ToBase64String(), nil
}

func handleOpt(o *Options) (string, error) {
	if o.Msg == "" {
		return "", errors.New("msg is empty")
	}
	if o.Token == "" {
		return "", errors.New("token is empty")
	}

	b, err := json.Marshal(o)
	if err != nil {
		return "", err
	}

	if o.Enc != nil {
		c, err := barkEncrypt(o.Enc, b)
		if err != nil {
			return "", err
		}
		b, err = json.Marshal(map[string]string{
			"ciphertext": c,
		})
		if err != nil {
			return "", err
		}
	}
	return string(b), nil
}
