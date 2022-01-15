package services

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"log"
	"path/filepath"
	"testing"
)

var p = "D:\\develop\\go\\s3\\server\\src\\tmp\\public_config\\conf\\beta\\wxpay"

// 查询商户余额
// https://pay.weixin.qq.com/wiki/doc/apiv3_partner/apis/chapter7_7_1.shtml
func Test_QueryBalance(t *testing.T) {
	var (
		//登陆商户平台【账户中心】-【商户信息】->【微信支付商户号】
		mchID string = "1606746840" // 商户号

		//登陆商户平台【账户中心】-【API安全】->【API证书】->【查看证书】
		mchCertificateSerialNumber string = "5430464867C627F3FC939F94B9B634162FB0A3F4" // 商户证书序列号

		//登陆商户平台【账户中心】-【API安全】->【设置APIv3密钥】 仅在设置时能看到
		mchAPIv3Key string = "dhfICa404Ez0uCBwlxXEQkOA3MUIlxVL" // 商户APIv3密钥
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(filepath.Join(p, "apiclient_key.pem"))
	if err != nil {
		log.Print("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalf("new wechat pay client err:%s", err)
	}
	//url := "https://api.mch.weixin.qq.com/v3/ecommerce/fund/balance/"+mchID+"?"+"account_type=OPERATION"
	url := "https://api.mch.weixin.qq.com/v3/merchant/fund/balance/OPERATION"
	ret, err := client.Get(context.TODO(), url)
	if err != nil {
		log.Fatalf("client.Get err:%s", err)
	}

	type BalanceRet struct {
		//SubMchid        string `json:"sub_mchid"`
		//AccountType     string `json:"account_type"`
		AvailableAmount int64 `json:"available_amount"`
		PendingAmount   int64 `json:"pending_amount"`
	}
	err = core.CheckResponse(ret.Response)
	if err != nil {
		log.Fatalf("CheckResponse err:%s", err)
	}
	data := new(BalanceRet)
	_ = core.UnMarshalResponse(ret.Response, data)
	log.Printf("%+v", data)
}
