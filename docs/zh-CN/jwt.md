## 简介

json web token



## 基础示例

使用Jwt必须先实现签发者接口

```go
JwtAudience interface {
		Audience() string
}
```

### 创建Jwt实例

```go
var	(
	jwtConfig = config.Item("jwt")
  memoryStore = jwt.NewMemoryStore()
)

// 创建jwt
jwt := New("secret", jwtConfig, memoryStore)
```

假设有`User`实现`JwtAudience`

```go
type User struct {
  Id string
}

func (u *User) Audience() string {
  return u.Id
}
```

### 生成token

```go
user := &User{
  Id: "uu3012"
}

token , err := jwt.Create(user)
if err != nil {
  fmt.Printf("create jwt error %s", err)
  return
}

// output
/*
Token struct {
  Lifetime int64  `json:"lifetime"`
  Token    string `json:"token"`
  Type     string `json:"type"`
}
*/
fmt.Printf("%v",token)
```

### 解析token

```go
var token = "token"

// 解析token
claims ,err := jwt.Parse(token)
if err != nil {
  fmt.Printf("parse token error %s", err)
}

// output
/*
type StandardClaims struct {
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
}
*/
fmt.Printf("%v", claims)

```

### 验证token

```go
var token = "token"

is,err := jwt.Valid(token)
if err != nil {
  // 验证是否过期
  if errors.Is(err, jwt.ValidationErrorExpired) {
    // 过期错误
  }
}

if is {
  // ok
}

```

### 其它方法

```go
// 刷新token
token := "token"
if err := jwt.Refresh(token); err != nil {
  fmt.Printf("refresh token error %s", err)
}

// 加入黑名单
if err := jwt.Invalidate(token); err != nil {
  fmt.Printf("invalidate token error", err)
}
```



## 自定义存储器

系统默认自带内存存储器，通常在分布式应用中你需要其它存储引擎，实现起来也非常简单，只需要实现`JwtStore` 接口，如下示例：

```go
JwtStore interface {
  Has(id string) bool

  Put(id string, audience JwtAudience, lifetime time.Time) error

  Forget(audience JwtAudience) error
}

type RedisStore struct {
  
}

jwt := New("secret", jwtConfig, new(RedisStore))
```



