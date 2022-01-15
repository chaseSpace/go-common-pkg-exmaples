# TransferDetailInput

## 属性列表

名称 | 类型 | 描述 | 补充说明
------------ | ------------- | ------------- | -------------
**OutDetailNo** | **string** | 商户系统内部区分转账批次单下不同转账明细单的唯一标识  | 
**TransferAmount** | **int64** | 转账金额单位为“分”  | 
**TransferRemark** | **string** | 单条转账备注（微信用户会收到该备注），UTF8编码，最多允许32个字符  | 
**Openid** | **string** | 收款用户openid。如果转账特约商户授权类型是INFORMATION_AUTHORIZATION_TYPE，对应的是特约商户公众号下的openid。  | 
**UserName** | **string** | 收款用户姓名。采用标准RSA算法，公钥由微信侧提供  | 
**UserIdCard** | **string** | 收款方身份证号，可不用填（采用标准RSA算法，公钥由微信侧提供）  | [可选] 

[\[返回类型列表\]](README.md#类型列表)
[\[返回接口列表\]](README.md#接口列表)
[\[返回服务README\]](README.md)


