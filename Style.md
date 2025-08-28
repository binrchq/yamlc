### YAMLC Style

StyleTop
```yaml
# 用户姓名
name: 张三
# 用户年龄
age: 30
# 电子邮箱地址
email: "zhangsan@example.com"
# 电话号码
phone: 13800138000
# 用户地址
address:
  # 街道地址
  street: 中关村大街1号
  # 城市
  city: 北京
  # 国家
  country: 中国
  # 邮政编码
  zipcode: 100080

# 用户标签
tags:
  - 开发者
  - Go语言
  - 后端

# 任职经历
workExperience:
    # company:公司名
    # position:职位
    - company: co1
      position: golang开发
    - company: co2
      position: 高级golang开发

# 是否激活
active: true
# 用户评分
score: 95.5
```

StyleInline
```yaml
name: 张三                        # 用户姓名
age: 30                          # 用户年龄    
email: "zhangsan@example.com"    # 电子邮箱地址
phone: 13800138000               # 电话号码

address:                         # 用户地址
  street: 中关村大街1号               # 街道地址
  city: 北京                         # 城市
  country: 中国                      # 国家
  zipcode: 100080                    # 邮政编码

tags:                            # 用户标签
  - 开发者
  - Go语言
  - 后端

workExperience:                  # 任职经历
    - company: co1                   # company:公司名
      position: golang开发           # position:职位

    - company: co2                   # company:公司名
      position: 高级golang开发        # position:职位

active: true                     # 是否激活
score: 95.5                      # 用户评分
```

StyleSmart
```yaml
name: 张三                        # 用户姓名
age: 30                          # 用户年龄    
email: "zhangsan@example.com"    # 电子邮箱地址
phone: 13800138000               # 电话号码

# 用户地址
address:                         
  street: 中关村大街1号               # 街道地址
  city: 北京                         # 城市
  country: 中国                      # 国家
  zipcode: 100080                    # 邮政编码

 # 用户标签
tags:                           
  - 开发者
  - Go语言
  - 后端

# 任职经历
workExperience:                 
    - company: co1                   # company:公司名
      position: golang开发           # position:职位

    - company: co2                   # company:公司名
      position: 高级golang开发        # position:职位

active: true                     # 是否激活
score: 95.5                      # 用户评分
```


StyleCompact
```yaml
name: 张三 # 用户姓名
age: 30 # 用户年龄    
email: "zhangsan@example.com" # 电子邮箱地址
phone: 13800138000 # 电话号码

address: # 用户地址
  street: 中关村大街1号 # 街道地址
  city: 北京 # 城市
  country: 中国 # 国家
  zipcode: 100080 # 邮政编码

tags: # 用户标签
  - 开发者
  - Go语言
  - 后端

workExperience: # 任职经历
    - company: co1 # company:公司名
      position: golang开发 # position:职位

    - company: co2 # company:公司名
      position: 高级golang开发 # position:职位

active: true # 是否激活
score: 95.5 # 用户评分
```

StyleMinimal
```yaml
name: 张三
age: 30   
email: "zhangsan@example.com"
phone: 13800138000
address:
  street: 中关村大街1号
  city: 北京
  country: 中国
  zipcode: 100080
tags:
  - 开发者
  - Go语言
  - 后端
workExperience:
    - company: co1
      position: golang开发
    - company: co2
      position: 高级golang开发
active: true
score: 95.5
```

StyleVerbose
```yaml
# name(string):用户姓名
name: 张三
# age(int):用户年龄
age: 30
# email(string):电子邮箱地址
email: "zhangsan@example.com"
# phone(int64):电话号码
phone: 13800138000
# address(*Address):用户地址
address:
  # street(string):街道地址
  street: 中关村大街1号
  # city(string):城市
  city: 北京
  # country(string):国家
  country: 中国
  # zipcode(int):邮政编码
  zipcode: 100080

# tags([]string):用户标签
tags:
  - 开发者
  - Go语言
  - 后端

# workExperience([]*WorkExperience):任职经历
workExperience:
    # company(string):公司名
    # position(string):职位
    - company: co1
      position: golang开发
    - company: co2
      position: 高级golang开发

# active(bool):是否激活
active: true
# score(float):用户评分
score: 95.5
```

StyleSpaced
```yaml
# 用户姓名
name: 张三

# 用户年龄
age: 30

# 电子邮箱地址
email: "zhangsan@example.com"

# 电话号码
phone: 13800138000

# 用户地址
address:
  # 街道地址
  street: 中关村大街1号

  # 城市
  city: 北京
  
  # 国家
  country: 中国

  # 邮政编码
  zipcode: 100080

# 用户标签
tags:
  - 开发者
  - Go语言
  - 后端

# 任职经历
workExperience:
    # company:公司名
    # position:职位
    - company: co1
      position: golang开发

    - company: co2
      position: 高级golang开发

# 是否激活
active: true

# 用户评分
score: 95.5
```

StyleGrouped
```yaml
# 用户姓名
name: 张三
# 用户年龄
age: 30
# 电子邮箱地址
email: "zhangsan@example.com"
# 电话号码
phone: 13800138000

# 用户地址
address:
  # 街道地址
  street: 中关村大街1号
  # 城市
  city: 北京
  # 国家
  country: 中国
  # 邮政编码
  zipcode: 100080

# 用户标签
tags:
  - 开发者
  - Go语言
  - 后端

# 任职经历
workExperience:
    # company:公司名
    # position:职位
    - company: co1
      position: golang开发
    - company: co2
      position: 高级golang开发

# 是否激活
active: true
# 用户评分
score: 95.5
```

StyleSectioned
```yaml
# 用户姓名
# 用户年龄
# 电子邮箱地址
# 电话号码
# 用户地址

name: 张三
age: 30
email: "zhangsan@example.com"
phone: 13800138000
address:
  # 街道地址
  # 城市
  # 国家
  # 邮政编码
  street: 中关村大街1号
  city: 北京
  country: 中国
  zipcode: 100080

# 用户标签
tags:
  - 开发者
  - Go语言
  - 后端

# 任职经历
workExperience:
    # company:公司名
    # position:职位
    - company: co1
      position: golang开发
    - company: co2
      position: 高级golang开发

# 是否激活
# 用户评分
active: true
score: 95.5
```

StyleDoc
```yaml
############################################
# name(string):用户姓名
# age(int):用户年龄
# email(string):电子邮箱地址
# phone(int64):电话号码
# address(*Address):用户地址
###########################################

name: 张三
age: 30
email: "zhangsan@example.com"
phone: 13800138000
address:
  ###########################################
  # street(string):街道地址
  # city(string):城市
  # country(string):国家
  # zipcode(int):邮政编码
  # zipcode(int):邮政编码
  ############################################
  street: 中关村大街1号
  city: 北京
  country: 中国
  zipcode: 100080
###########################################
# tags([]string):用户标签
###########################################
tags:
  - 开发者
  - Go语言
  - 后端
###################################################
# workExperience([]*WorkExperience):任职经历
####################################################
workExperience:
    #############################################
    # company(string):公司名
    # position(string):职位
    #############################################
    - company: co1
      position: golang开发
    - company: co2
      position: 高级golang开发

#####################################
# active(bool):是否激活
####################################
active: true
####################################
# score(float):用户评分
#####################################
score: 95.5
```

StyleSeparate
```yaml
############################################
# name(string):用户姓名
# age(int):用户年龄
# email(string):电子邮箱地址
# phone(int64):电话号码
# address(*Address):用户地址
    # street(string):街道地址
    # city(string):城市
    # country(string):国家
    # zipcode(int):邮政编码
    # zipcode(int):邮政编码
# tags([]string):用户标签
# workExperience([]*WorkExperience):任职经历
    #- company(string):公司名
    #  position(string):职位
# active(bool):是否激活
# score(float):用户评分
###########################################

name: 张三
age: 30
email: "zhangsan@example.com"
phone: 13800138000
address:
  street: 中关村大街1号
  city: 北京
  country: 中国
  zipcode: 100080

tags:
  - 开发者
  - Go语言
  - 后端

workExperience:
    - company: co1
      position: golang开发
    - company: co2
      position: 高级golang开发
      
active: true
score: 95.5
```
