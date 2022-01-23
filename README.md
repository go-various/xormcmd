**数据库对象创建工具**

**安装**
<pre>
go install github.com/go-various/xormcmd@v1.0.1
</pre>
**反向建模**
```
*创建的模型位于model目录下*

xormcmd reverse mysql "root:123456@tcp(localhost)/test" go/pkg/mod/github.com/go-various/xormcmd@v1.0.1/templates/goxorm model

*user.go 示例*
package model
import (
    "encoding/json"
    "time"
)

type User struct {
    Id             *string    `json:"id" name:"账户ID" xorm:"id not null pk comment('账户ID') VARCHAR(8)"`
    Mobile         *string    `json:"mobile" name:"" xorm:"mobile not null unique CHAR(11)"`
    ProductCode    *string    `json:"product_code" name:"产品编号" xorm:"product_code not null comment('产品编号') VARCHAR(16)"`
    NickName       *string    `json:"nick_name" name:"用户昵称" xorm:"nick_name default 'NULL' comment('用户昵称') VARCHAR(32)"`
    HeadImgUrl     *string    `json:"head_img_url" name:"用户头像" xorm:"head_img_url default 'NULL' comment('用户头像') VARCHAR(512)"`
    Gender         *string    `json:"gender" name:"性别(M,F)" xorm:"gender default 'NULL' comment('性别(M,F)') CHAR(1)"`
    Birthday       *string    `json:"birthday" name:"生日:1980-01-01" xorm:"birthday default 'NULL' comment('生日:1980-01-01') CHAR(10)"`
    WechatId       *string    `json:"wechat_id" name:"微信ID" xorm:"wechat_id default 'NULL' comment('微信ID') VARCHAR(64)"`
    WechatNickName *string    `json:"wechat_nick_name" name:"微信昵称" xorm:"wechat_nick_name default 'NULL' comment('微信昵称') VARCHAR(32)"`
    RegChannel     *string    `json:"reg_channel" name:"注册渠道" xorm:"reg_channel default 'NULL' comment('注册渠道') VARCHAR(16)"`
    CreatedTime    *time.Time `json:"created_time" name:"创建时间" xorm:"created_time not null default 'current_timestamp()' comment('创建时间') TIMESTAMP"`
}

func (o *User) TableName() string {
    return "t_user"
}
func (o *User) MarshalBinary() ([]byte, error) {
    return json.Marshal(o)
}
func (o *User) UnmarshalBinary(data []byte) error {
    return json.Unmarshal(data, &o)
}
func (o *User) ID() interface{} {
    return o.Id
}
```
**构建模型视图**
```
xormcmd view mysql "root:123456@tcp(localhost)/test" go/pkg/mod/github.com/go-various/xormcmd@v1.0.1/templates/goview views
```