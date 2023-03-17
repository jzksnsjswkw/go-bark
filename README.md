# go-bark

bark 推送 golang sdk，支持加密传输

## Example

```Go
package main

import "github.com/jzksnsjswkw/go-bark"

func main() {
	err := bark.Push(&bark.Option{
		Msg:   "test",
		Token: "xxxxxxxxxxxxxxxxxxxxxx",
	})
	if err != nil {
		panic(err)
	}

	err = bark.Push(&bark.Option{
		Msg:   "test",
		Token: "xxxxxxxxxxxxxxxxxxxxxx",
		Enc: &bark.EncOpt{
			Mode: bark.CBC,
			Key:  "1234567890abcdef",
			Iv:   "1111111111111111",
		},
	})
	if err != nil {
		panic(err)
	}
}

```

## Options

```Go
type Options struct {
	// 推送内容 (必填)
	Msg string `json:"body"`
	// token (必填)
	Token string `json:"-"`
	// 推送标题
	Title string `json:"title,omitempty"`
	// 消息分组
	Group string `json:"group,omitempty"`
	// 点击推送时，跳转的URL，支持URL Scheme 和 Universal Link
	Url string `json:"url,omitempty"`
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
```
