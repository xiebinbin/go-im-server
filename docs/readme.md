### 系统相关

#### 获取系统公钥

**请求URL：**

- `/sys/pubKey`

**请求方式：**

- POST

**接口响应参数说明**

| 参数名     | 类型     | 说明     |
|:--------|:-------|:-------|
| pub_key | string | pubKey |

```
{
    "code": 200,
    "data": {
        "pub_key": "024297241d948c71dcc1b83f98563ceaca36150000fa389855aa84e79b88433e7b"
    },
    "msg": ""
}
```

#### 获取预签名URL

**请求URL：**

- `/sys/preSignURL`

**请求方式：**

- POST

**接口响应参数说明**

| 参数名 | 类型     | 说明  |
|:----|:-------|:----|
| url | string | 预签名 |

```
{
    "code": 200,
    "data": {
        "url": "https://bobobo-test.bb05e19bea2bba53858a8eeadb1c55f3.r2.cloudflarestorage.com/miyaya.txt?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=d71097891fd7fd38ccfe49720b37d20d%2F20231211%2F%2Fs3%2Faws4_request&X-Amz-Date=20231211T094304Z&X-Amz-Expires=900&X-Amz-SignedHeaders=expires%3Bhost&x-id=PutObject&X-Amz-Signature=bf9d7073492f14e88a492befd99d97b3fb48a340fbe8c31ff12f30b679b825ce"
    },
    "msg": ""
}
```

### PC扫码登录

#### 获取二维码生成信息

**请求URL：**

- `/auth/getQRCode`

**请求方式：**

- POST

**成功示例**

| 参数名         | 类型     | 说明                                                                     |
|:------------|:-------|:-----------------------------------------------------------------------|
| url_pre     | string | 二维码url前缀                                                               |
| value       | []byte | 二维码内容，客户端首先要进行base64encode 然后拼接在url_pre之后再生成二维码； 注意：第一位表示版本号，最后一位表示校验码 |
| expire      | int64  | 超时时间戳，毫秒                                                               |
| create_time | int64  | 创建时间                                                                   |

```
{
    "code": 200,
    "data": {
        "url_pre": "https://qr/",
        "value": [1,53,99,54,50,53,97,53,52,98,51,101,54,48,49,101,54,39],
        "expire": 0,
        "create_time": 1649294644208
    },
    "msg": ""
}
```

**错误示例**

```
{
    "code": 200,
    "data": {
        "err_code": 400,
        "err_msg": "fail"
    },
    "msg": ""
}
```

#### APP扫码

**请求URL：**

- `/auth/appScanLoginQrCode`

**请求方式：**

- POST

**请求参数：**

| 参数名  | 类型     | 说明      |
|:-----|:-------|:--------|
| code | []byte | 二维码code |

**成功或者失败的错误提示**

```
{
    "code": 200,
    "data": {
        "err_code": 400,
        "err_msg": "fail"
    },
    "msg": ""
}
```

**错误参数说明**

| 参数名    | 说明      |
|:-------|:--------|
| 200001 | 二维码已经过期 |
| 200003 | 无效的二维码  |
| 100010 | 用户未注册   |
| 100011 | 用户被禁用   |
| 100003 | 参数错误    |
| 400    | 失败      |

#### App确认登录

**请求URL：**

- `/auth/appConfirmLogin`

**请求方式：**

- POST

**请求参数：**

| 参数名  | 类型     | 说明     |
|:-----|:-------|:-------|
| code | []byte | 二维码    |
| sk   | string | 加密后的私钥 |

**成功或者失败的错误提示**

```
{
    "code": 200,
    "data": {
        "err_code": 400,
        "err_msg": "fail"
    },
    "msg": ""
}
```

**错误参数说明**

| 参数名    | 说明              |
|:-------|:----------------|
| 1      | 还未扫码            |
| 100010 | 用户未注册           |
| 100007 | 没有权限，禁止操作       |
| 200001 | 扫码后确认操作超时，默认1分钟 |
| 200003 | 无效的二维码          |
| 100011 | 用户被禁用           |
| 100003 | 参数错误            |
| 400    | 失败              |

#### PC端轮询查看扫码结果

**请求URL：**

- `/auth/getQRCodeScanRes`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名  | 类型     | 说明  |
|:-----|:-------|:----|
| code | []byte | 二维码 |

**接口返回参数说明**

| 参数名    | 类型     | 说明                                |
|:-------|:-------|:----------------------------------|
| status | string | 状态 0-待扫码 1-已扫码 2-已授权；当状态为0时只有此字段， |
| id     | string | 用户id（address） ，当状态为0时有此字段         |
| avatar | json   | 状态为1字段存在 ，当状态为1时有此字段              |
| sk     | string | 加密后的sk，当状态为1或2时有此字段               |

**示例**

```
{
    "code": 200,
    "data": {
        "err_code": 200001,
        "err_msg": "expired"
    },
    "msg": ""
}
或
{
    "code": 200,
    "data": {
        "id": "2tgmd3aiofyk",
        "avatar": "",
        "name":"",
        "status": 2,
        "sk":"",
    },
    "msg": ""
}
```

**错误参数说明**

| 参数名    | 说明              |
|:-------|:----------------|
| 402    | 发生异常            |
| 200001 | 扫码后确认操作超时，默认1分钟 |
| 200003 | 无效的二维码          |
| 100003 | 参数错误            |
| 400    | 失败              |

### 用户相关

#### 判断用户是否存在

**请求URL：**

- `/user/isRegister`

**请求方式：**

- POST

**接口响应参数说明**

| 参数名         | 类型   | 说明                  |
|:------------|:-----|:--------------------|
| is_register | bool | true-已经注册 false-未注册 |

```
{
    "code": 200,
    "data": {
        "is_register": ture/false
        },
    "msg": ""
}
```

#### 修改用户昵称

**请求URL：**

- `/user/updateName`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名  | 类型     | 说明   |
|:-----|:-------|:-----|
| name | string | 用户昵称 |

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 用户认证（注册）

**请求URL：**

- `/auth/register`

**请求方式：**

- POST

**接口响应参数说明**

| 参数名         | 类型     | 说明     |
|:------------|:-------|:-------|
| address     | string | 地址     |
| avatar      | string | 头像     |
| name        | string | 昵称     |
| pub_key     | string | pubKey |
| create_time | int64  | 创建时间   |

```
{
    "code": 200,
    "data": {
        "address": "0x0b84b2d122cb1c058b988d9f0291a6e25364c6f8d",
        "avatar": "https://img1.baidu.com/it/u=3709586903,1286591012&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=500",
        "name": "z4x56k",
        "pub_key": "0364d735fb01f4bfd2695b1805585de1aa1992243f5a685a9ead1045e4e10c5c41",
        "create_time": 1692531674664
    },
    "msg": ""
}
```

#### 批量获取用户信息

**请求URL：**

- `/user/getBatchInfo`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名  | 类型       | 说明    |
|:-----|:---------|:------|
| uids | []string | 用户ids |

**接口响应**

```
{
    "code": 200,
    "data": {
        "items":[
            {
            "id":"0xb1d3c24d3cd2ef52e6dc3ac6c06742a7dc17e041",
            "avatar":"https://img1.baidu.com/it/u=3709586903,1286591012\u0026fm=253\u0026fmt=auto\u0026app=138\u0026f=JPEG?w=500\u0026h=500",
            "name":"2hk3ki"
            }
        ]
    },
    "msg": ""
}
```

#### 注销账号-未定后续操作

**请求URL：**

- `/user/unsubscribe`

**请求方式：**

- POST

**接口请求参数说明**
无

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 获取用户关系

**请求URL：**

- `/friend/relationList`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名  | 类型    | 说明   |
|:-----|:------|:-----|
| uids | array | 用户地址 |

**接口返回参数说明**

| 参数名       | 类型     | 说明                                                         |
|:----------|:-------|:-----------------------------------------------------------|
| id        | string | 用户地址                                                       |
| is_friend | int    | 是否是好友 0-互为陌生人 1-互为好友 2-对方是我的好友/我是对方的陌生人 3-对方是我的陌生人/我是对方的好友 |
| name      | string | 用户名称                                                       |
| avatar    | string | 用户头像                                                       |

```
{
    "code": 200,
    "data": {
        "items":[
            {
            "id":"0xb1d3c24d3cd2ef52e6dc3ac6c06742a7dc17e041",
            "is_friend":1,
            "name":"",
            "avatar":"",
            }
        ]
    },
    "msg": ""
}
```

### 消息相关

#### 发送消息

**请求URL：**

- `/message/send`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名     | 类型         | 说明          |
|:--------|:-----------|:------------|
| mid     | string     | 消息id        |
| chat_id | string     | 会话id        |
| content | jsonString | 客户端定义，服务端不管 |

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### （双向）删除消息-根据消息IDs

**请求URL：**

- `/message/deleteBatch`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型       | 说明                  |
|:----|:---------|:--------------------|
| ids | []string | 消息id数组, 单个消息也是这个API |

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### （自己端-暂时是双向）删除消息-按消息Id

**请求URL：**

- `/message/deleteSelfAll`

**请求方式：**

- POST

**接口请求参数说明**
无

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### （双向）删除所有消息-根据会话IDs

**请求URL：**

- `/message/deleteAllByChatIds`

**请求方式：**

- POST

**接口请求参数说明**
| 参数名 | 类型 | 说明 |
|:----|:---------|:--------------------|
| chat_ids | []string | 会话id数组, 单个会话也是这个API |

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### （自己端-暂时是双向）删除所有消息-根据会话IDs

**请求URL：**

- `/message/deleteSelfByChatIds`

**请求方式：**

- POST

**接口请求参数说明**
| 参数名 | 类型 | 说明 |
|:----|:---------|:--------------------|
| ids | []string | 会话id数组, 单个会话也是这个API |

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 撤回消息-根据消息IDs

**请求URL：**

- `/message/revokeBatch`

**请求方式：**

- POST

**接口请求参数说明**
| 参数名 | 类型 | 说明 |
|:----|:---------|:--------------------|
| ids | []string | 消息id数组, 单个消息也是这个API |

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 撤回消息-根据会话IDs

**请求URL：**

- `/message/revokeByChatIds`

**请求方式：**

- POST

**接口请求参数说明**
| 参数名 | 类型 | 说明 |
|:----|:---------|:--------------------|
| ids | []string | 会话id数组, 单个会话也是这个API |

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

### 会话相关

#### 我的会话列表

**请求URL：**

- `/chat/list`

**请求方式：**

- POST

**接口响应参数说明**

| 参数名      | 类型     | 说明        |
|:---------|:-------|:----------|
| id       | string | 会话id      |
| group_id | string | 群聊的群id    |
| name     | string | 昵称        |
| avatar   | string | 头像        |
| type     | int    | 1-单聊 2-群聊 |

**接口响应**

```
{
    "code": 200,
    "data": {
        "items": [
            {
                "id": "s_a7be32fecc0b2015",
                "group_id": "",
                "name": "2hk3ki",
                "avatar": "https://img1.baidu.com/it/u=3709586903,1286591012&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=500",
                "type": 1,
                "last_read_sequence": 0,
                "last_sequence": 0,
                "last_time": 0
            }
        ]
    },
    "msg": ""
}
```

#### 删除会话

**请求URL：**

- `/chat/delete`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型       | 说明     |
|:----|:---------|:-------|
| ids | []string | 消息id数组 |

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

### 好友相关

#### 申请添加好友

**请求URL：**

- `/friend/inviteApply`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名     | 类型     | 说明    |
|:--------|:-------|:------|
| obj_uid | string | 对象uid |
| remark  | string | 备注    |

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 申请列表

**请求URL：**

- `/friend/inviteList`

**请求方式：**

- POST

**接口请求参数说明**
无

**接口响应**

| 参数名    | 类型     | 说明                  |
|:-------|:-------|:--------------------|
| id     | string | 数据库编号ID             |
| uid    | string | 请求的用户地址             |
| avatar | string | 请求的用户头像             |
| name   | string | 请求的用户名称             |
| remark | string | 请求的用户备注             |
| status | int8   | 0-等待验证 1-验证通过 2-已拒绝 |

```
{
    "code": 200,
    "data": {
        "items":[
            {
            "id":"fcad58e10337f612",
            "uid":"0xb1d3c24d3cd2ef52e6dc3ac6c06742a7dc17e041",
            "remark":"hello",
            "avatar":"https://img1.baidu.com/it/u=3709586903,1286591012\u0026fm=253\u0026fmt=auto\u0026app=138\u0026f=JPEG?w=500\u0026h=500",
            "name":"2hk3ki"
            }
        ]
    },
    "msg": ""
}
```

#### 同意

**请求URL：**

- `/friend/inviteAgree`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型     | 说明   |
|:----|:-------|:-----|
| id  | string | 主键id |

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 拒绝

**请求URL：**

- `/friend/inviteReject`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型     | 说明   |
|:----|:-------|:-----|
| id  | string | 主键id |

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 已读

**请求URL：**

- `/friend/inviteRead`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型       | 说明    |
|:----|:---------|:------|
| ids | []string | 主键ids |

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 好友列表

**请求URL：**

- `/friend/list`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名  | 类型       | 说明                |
|:-----|:---------|:------------------|
| uids | []string | 好友的用户id数组，为空的时候全量 |

**接口响应的参数说明**

| 参数名     | 类型     | 说明                 |
|:--------|:-------|:-------------------|
| uid     | string | 好友的用户              |
| chat_id | string | 和好友的会话id           |
| remark  | string | 好友的备注              |
| pub_key | string | 好友的公钥              |
| avatar  | string | 好友的头像              |
| name    | string | 好友的昵称              |
| gender  | int    | 好友的性别 0-未知 1-男 2-女 |
| sign    | string | 好友的签名              |

**接口响应**

```
{
    "code": 200,
    "data":   {
    "items":[
            {
                "uid":"0x4c3f6cb0cd7df2977ea98e006b61bb899637d1ca",
                "chat_id":"s_a7be32fecc0b2015",
                "remark":"",
                "pub_key":"",
                "avatar":"",
                "name":"",
                "gender":0,
                "sign":"",
            }
        ]
    },
    "msg": ""
}
```

#### 好友备注

**请求URL：**

- `/friend/updateRemark`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名     | 类型     | 说明      |
|:--------|:-------|:--------|
| obj_uid | string | 好友的用户id |
| remark  | string | 备注      |

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 删除好友（指定好友）

**请求URL：**

- `/friend/delete`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型       | 说明       |
|:----|:---------|:---------|
| ids | []string | 好友列表id数组 |

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 删除好友（单向）

**请求URL：**

- `/friend/deleteUnilateral`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名  | 类型       | 说明                  |
|:-----|:---------|:--------------------|
| uids | []string | 好友的id数组 如果为空-则为删除所有 |

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 删除所有好友（双向）

**请求URL：**

- `/friend/deleteBilateral`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名  | 类型       | 说明                  |
|:-----|:---------|:--------------------|
| uids | []string | 好友的id数组 如果为空-则为删除所有 |

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

### 群组相关

#### 创建群聊

**请求URL：**

- `/group/create`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名     | 类型       | 说明             |
|:--------|:---------|:---------------|
| id      | string   | 群ID，前端生成防止重复调用 |
| avatar  | string   | 头像             |
| name    | string   | 名称             |
| pub_key | string   | 公钥             |
| members | []string | 好友列表id数组       |

**请求Demo**

```
{
	"id": "test001",
	"pub_key":"",
	"avatar": "https://img1.baidu.com/it/u=3709586903,1286591012&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=500",
	"name": "miya-test001",
	"members": ["0x7fe4407b6de6b0ac3b0a02fe93ecd175e9b31aa8", "0x7f90fadd2e3fdbacfd3ffc0c554fcf5878cc1601"]
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 群聊用户

**请求URL：**

- `/group/members`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型     | 说明  |
|:----|:-------|:----|
| id  | string | 群ID |

**请求Demo**

```
{
	"id": "test001",
}
```

**接口响应参数说明**

| 参数名      | 类型     | 说明                   |
|:---------|:-------|:---------------------|
| id       | string | 主键ID                 |
| uid      | string | 用户id（地址）             |
| gid      | string | 群id                  |
| role     | int8   | 角色 1-群组 2-管理员 3-普通成员 |
| my_alias | string | 群昵称                  |

**接口响应**

```
{
    "code": 200,
    "data": {
        "items": [
            {
                "id": "2e2bd92a86eb9b61",
                "uid": "0x7f90fadd2e3fdbacfd3ffc0c554fcf5878cc1601",
                "gid": "test001",
                "role": 3,
                "my_alias": "",
                "admin_time": 0,
                "create_time": 1696839147569
            },
            {
                "id": "c3427c909968d7bc",
                "uid": "0x7fe4407b6de6b0ac3b0a02fe93ecd175e9b31aa8",
                "gid": "test001",
                "role": 1,
                "my_alias": "",
                "admin_time": 0,
                "create_time": 1696839147569
            }
        ],
        "status": 0
    },
    "msg": ""
}
```

#### 同意加入群聊

（只有群组和管理员有权限）

**请求URL：**

- `/group/agreeJoin`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名     | 类型       | 说明     |
|:--------|:---------|:-------|
| id      | string   | 群ID    |
| obj_uid | []string | 用户id数组 |

**请求Demo**

```
{
	"id": "test001"
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 邀请加入群聊

**请求URL：**

- `/group/inviteJoin`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型     | 说明         |
| :----- | :------- | :----------- |
| id     | string   | 群ID         |
| items  | []string | 用户相关信息 |

| 参数名        | 类型     | 说明        |
| :------------ | :------- | :---------- |
| items.uid     | string   | 用户ID      |
| items.enc_key | []string | 用户enc_key |

**请求Demo**

```
{
	"id": "test001",
	"items": [
		{
			"uid":"test002",
			"enc_key":"",
		}
	]
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 踢出群聊

（群组和管理员有权限）

**请求URL：**

- `/group/kickOut`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名     | 类型       | 说明     |
|:--------|:---------|:-------|
| id      | string   | 群ID    |
| obj_uid | []string | 用户id数组 |

**请求Demo**

```
{
	"id": "test001"
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 我的群聊

**请求URL：**

- `/group/list`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型 | 说明 |
|:----|:---|:---|

**请求Demo**

```
...
```

**接口响应**

```
{
    "code": 200,
    "data": {
        "items": [
            {
                "id": "test001",
                "uid": "0x7f90fadd2e3fdbacfd3ffc0c554fcf5878cc1601",
                "avatar": "https://img1.baidu.com/it/u=3709586903,1286591012&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=500",
                "name": "miya-test001",
                "create_time": 1697080826986
            }
        ]
    },
    "msg": ""
}
```

#### 修改群名称（群组和管理员有权限）

**请求URL：**

- `/group/updateName`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名  | 类型     | 说明  |
|:-----|:-------|:----|
| id   | string | 群ID |
| name | string | 群名称 |

**请求Demo**

```
{
	"id": "test001",
	"name": "test002"
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 修改群头像（群组和管理员有权限）

**请求URL：**

- `/group/updateAvatar`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名    | 类型     | 说明  |
|:-------|:-------|:----|
| id     | string | 群ID |
| avatar | string | 头像  |

**请求Demo**

```
{
	"id": "test001",
	"avatar": "",
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 修改我在群里面的昵称（只能修改自己的）

**请求URL：**

- `/group/updateAlias`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名   | 类型     | 说明    |
|:------|:-------|:------|
| id    | string | 群ID   |
| alias | string | 用户群昵称 |

**请求Demo**

```
{
	"id": "test001",
	"alias": "miyaya"
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 退出群聊

##### 指定单个群

**请求URL：**

- `/group/quit`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型     | 说明  |
|:----|:-------|:----|
| id  | string | 群ID |

**请求Demo**

```
{
	"id": "test001",
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

##### 指定多个群

**请求URL：**

- `/group/quitByIds`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型       | 说明   |
|:----|:---------|:-----|
| ids | []string | 群Ids |

**请求Demo**

```
{
	"id": ["test001","test002"]
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

##### 我的所有

**请求URL：**

- `/group/quitAll`

**请求方式：**

- POST

**接口请求参数说明**
无

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 修改群通告

**请求URL：**

- `/group/updateNotice`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名       | 类型     | 说明                |
|:----------|:-------|:------------------|
| id        | string | 群ID               |
| notice    | string | 群通告               |
| notice_id | string | 群通告Id md5(notice) |

**请求Demo**

```
{
	"id": "test001",
	"notice": "miyayayayayayay",
	"notice_md5": "aaaaaaaa",
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 获取群通告

**请求URL：**

- `/group/getNotice`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型     | 说明  |
|:----|:-------|:----|
| id  | string | 群ID |

**请求Demo**

```
{
	"id": "test001",
}
```

**接口响应**

```
{
	"id": "test001",
	"notice": "miyayayayayayay",
	"notice_md5": "aaaaaaaa",
}
```

#### 修改群简介

**请求URL：**

- `/group/updateDesc`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名      | 类型     | 说明              |
|:---------|:-------|:----------------|
| id       | string | 群ID             |
| desc     | string | 群简介             |
| desc_md5 | string | 群简介Id md5(desc) |

**请求Demo**

```
{
	"id": "test001",
	"desc": "miyayayayayayay",
	"desc_md5": "aaaaaaaa",
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 获取群简介

**请求URL：**

- `/group/getDesc`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型     | 说明  |
|:----|:-------|:----|
| id  | string | 群ID |

**请求Demo**

```
{
	"id": "test001",
}
```

**接口响应**

```
{
	"id": "test001",
	"desc": "miyayayayayayay",
	"desc_md5": "aaaaaaaa",
}
```

#### 解散群聊

（只有群组有权限）

**请求URL：**

- `/group/dismiss`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型     | 说明  |
|:----|:-------|:----|
| id  | string | 群ID |

**请求Demo**

```
{
	"id": "test001",
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 转移群组

（只有群主有权限）

**请求URL：**

- `/group/transfer`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名     | 类型     | 说明     |
|:--------|:-------|:-------|
| id      | string | 群ID    |
| obj_uid | string | 转移对象id |

**请求Demo**

```
{
	"id": "test001",
	"obj_uid": "test001",
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 添加管理员

（只有群主有权限）

**请求URL：**

- `/group/addAdmin`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名     | 类型          | 说明               |
|:--------|:------------|:-----------------|
| id      | string      | 群ID              |
| obj_uid | []string 数组 | 管理员的ids，一次可以添加多个 |

**请求Demo**

```
{
	"id": "test001",
	"obj_uid": ["test001","test002"],
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 移除管理员

（只有群主有权限）

**请求URL：**

- `/group/removeAdministrators`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名     | 类型          | 说明               |
|:--------|:------------|:-----------------|
| id      | string      | 群ID              |
| obj_uid | []string 数组 | 管理员的ids，一次可以移除多个 |

**请求Demo**

```
{
	"id": "test001",
	"obj_uid": ["test001","test002"],
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 清除群消息

（只有群主和管理员有权限）

**请求URL：**

- `/group/clearMessage`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名  | 类型          | 说明               |
|:-----|:------------|:-----------------|
| id   | string      | 群ID              |
| mids | []string 数组 | 消息的ids，为空则表示清楚全部 |

**请求Demo**

```
{
	"id": "test001",
	"mids": ["test001","test002"],
}
```

**接口响应**

```
{
    "code": 200,
    "data": null,
    "msg": ""
}
```

#### 待审核申请列表

---只有群主和管理员有权限

**请求URL：**

- `/group/applyList`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型     | 说明  |
|:----|:-------|:----|
| id  | string | 群ID |

**请求Demo**

```
{
	"id": "test001",
}
```

**接口响应参数说明**

| 参数名    | 类型     | 说明       |
|:-------|:-------|:---------|
| id     | string | 主键ID     |
| uid    | string | 用户id（地址） |
| gid    | string | 群id      |
| name   | int8   | 昵称       |
| avatar | string | 头像       |

**接口响应**

```
{
    "code": 200,
    "data": {
        "items": [
            {
                "id": "2e2bd92a86eb9b61",
                "gid": "2e2bd92a86eb9b61",
                "uid": "0x7f90fadd2e3fdbacfd3ffc0c554fcf5878cc1601",
                "name": "test001",
                "avatar": "",
                "create_time": 1696839147569
            }
        ],
        "status": 0
    },
    "msg": ""
}
```

#### 我的申请列表

**请求URL：**

`/group/myApplyList`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型     | 说明  |
| :----- | :------- | :---- |
| ids    | []string | 群Ids |

**请求Demo**

```
{
	"ids": ["test001","test002"],
}
```

**接口响应参数说明**

| 参数名       | 类型   | 说明           |
| :----------- | :----- | :------------- |
| id           | string | 主键ID         |
| uid          | string | 用户id（地址） |
| gid          | string | 群id           |
| name         | int8   | 群昵称         |
| avatar       | string | 群头像         |
| create_time  | int64  | 创建时间       |
| member_limit | int64  | 群人数限制     |
| total        | int64  | 当前人数       |

**接口响应**

```
{
    "code": 200,
    "data": {
        "items": [
            {
                "id": "2e2bd92a86eb9b61",
                "gid": "2e2bd92a86eb9b61",
                "uid": "0x7f90fadd2e3fdbacfd3ffc0c554fcf5878cc1601",
                "name": "test001",
                "avatar": "",
                "create_time": 1696839147569,
                "member_limit": 100,
                "total": 2
            }
        ],
        "status": 0
    },
    "msg": ""
}
```

#### 群详情

**请求URL：**

- `/group/detailByIds`

**请求方式：**

- POST

**接口请求参数说明**

| 参数名 | 类型       | 说明   |
|:----|:---------|:-----|
| ids | []string | 群Ids |

**请求Demo**

```
{
	"ids": ["test001","test002"],
}
```

**接口响应参数说明**

| 参数名          | 类型     | 说明       |
|:-------------|:-------|:---------|
| id           | string | 主键ID     |
| uid          | string | 用户id（地址） |
| gid          | string | 群id      |
| name         | int8   | 群昵称      |
| avatar       | string | 群头像      |
| create_time  | int64  | 创建时间     |
| member_limit | int64  | 群人数限制    |
| total        | int64  | 当前人数     |

**接口响应**

```
{
    "code": 200,
    "data": {
        "items": [
            {
                "id": "2e2bd92a86eb9b61",
                "gid": "2e2bd92a86eb9b61",
                "uid": "0x7f90fadd2e3fdbacfd3ffc0c554fcf5878cc1601",
                "name": "test001",
                "avatar": "",
                "create_time": 1696839147569,
                "member_limit": 100,
                "total": 2
            }
        ],
        "status": 0
    },
    "msg": ""
}
```

### WEBSOCKET相关

#### 初始化长链接

/wsConnect?address="0x7f90fadd2e3fdbacfd3ffc0c554fcf5878cc1601"&device_id="deviceId001"

#### 有新消息来长连接推送

**CMD类型说明**

| 参数名              | 说明    |
|:-----------------|:------|
| new_msg          | 新消息   |
| apply_friend_msg | 好友申请  |
| agree_add_friend | 同意加好友 |

```
{
"cmd":"new_msg",
"items":{
    "id":"ef4e009633ffde82", //msgId
    "sequence":13489,
    "chat_id":"1e567516ace3d6bebdbe6ea382227efb",
    }
}
```