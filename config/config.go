# 应用模式：debug开发模式 / release生产模式（影响日志、调试信息等）
AppMode = debug
# 服务监听地址（格式: ":端口号"，:8080 表示监听所有接口的8080端口）
HttpPort = :8080
# JWT令牌签名密钥（生产环境必须修改为复杂字符串！）
JwtKey = xieyangyuyue

[database]
# 数据库类型（代码中未使用该配置项）
Db = mariadb
# 数据库连接信息（示例为远程数据库配置）
DbHost = hnb1.wch1.top
DbPort = 10500
DbUser = xieyangyuyue
DbPassWord = xieyangyuyue  # 数据库密码（敏感信息！）
DbName =gin-gorm

[redis]
# 数据库连接信息（示例为远程数据库配置）
DbHost = hnb1.wch1.top
DbPort = 10600
DbPassWord = redis_ESGA4c  # 数据库密码（敏感信息！）
DbNumber =0


[mail]
MailPasswd = zjlrwlypyrjgbace

[qiniu]
# 存储区域（1:华东 2:华北 3:华南）
Zone = 1
# 以下为七牛云API凭证（需要自行申请配置）
AccessKey = ysd-CR0LBc7P2lo78aeoR00yQQEiAdzBnquj0lU7   # 留空需要代码中配置
SecretKey = zR42dEdMmoMrRSepW3B7aOHi7eSF2mSUBas-wwMV   # 留空需要代码中配置
ABucket =   xieyangyuyue    # 存储空间名称
QiniuSever =  http://svo0wzpvb.hd-bkt.clouddn.com/ # CDN加速域名或七牛云域名
