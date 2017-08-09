# yiban-homepage-backend-go
UPC易班后端
## 配置
首先配置config.json

将config_example.json更名为config.json 或者直接按照config_example.json 新建config.json文件
```json
{
         "server_port": 8080,//服务端口
         "admins":[
           {
             "username":"admin",
             "password":"admin"
           },{
             "username":"root",
             "password":"root"
           }
         ],
         "image_url_server":"E:\\projects\\yiban-homepage\\images\\",//后端保存上传图片文件的目录
         "session_secret":"my-secret",//session的加密token
         "data_url":"E:\\projects\\yiban-homepage\\backend\\data_copy.json",//前端访问的data.json位置
         "image_url_front":"http://localhost:8086/images/",//前端需要访问的图片目录（即是后端接收到上传图片之后存储的相对于前端wwwroot的目录）
         "cross_origin":"http://localhost:8086",//允许跨域的前端origin
         "account":"user",//用来扒取数据的公共账号
         "password":"pwd",//用来爬取数据的公共账号密码
         "interval_hours":12//爬取的时间间隔
       }

```
## 运行

```bash
# 假设 config.json CONFIG_JSON_PATH
chmod +x yiban-homepage-backend-go
./yiban-homepage-backend-go --config CONFIG_JSON_PATH
```