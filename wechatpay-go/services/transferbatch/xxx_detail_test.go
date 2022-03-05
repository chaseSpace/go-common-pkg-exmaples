package transferbatch

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

func Test_GetTransferDetailByOutNo(t *testing.T) {

	var (
		//登陆商户平台【账户中心】-【商户信息】->【微信支付商户号】
		mchID string = "16067*****" // 商户号

		//登陆商户平台【账户中心】-【API安全】->【API证书】->【查看证书】
		mchCertificateSerialNumber string = "5430464867C627F3FC939F94B9B634162F******" // 商户证书序列号

		//登陆商户平台【账户中心】-【API安全】->【设置APIv3密钥】 仅在设置时能看到
		mchAPIv3Key string = "dhfICa404Ez0uCBwlxXEQkOA3M******" // 商户APIv3密钥
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
		log.Printf("new wechat pay client err:%s", err)
	}

	svc := TransferDetailApiService{Client: client}
	resp, result, err := svc.GetTransferDetailByOutNo(ctx,
		GetTransferDetailByOutNoRequest{
			OutDetailNo: core.String("plfk202004201322211215"),
			OutBatchNo:  core.String("plfk202004201322211215"),
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call GetTransferDetailByNo err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp.FailReason)
	}
}

func Test_GetTransferBatchByOutNo(t *testing.T) {

	var (
		//登陆商户平台【账户中心】-【商户信息】->【微信支付商户号】
		mchID string = "16067*****" // 商户号

		//登陆商户平台【账户中心】-【API安全】->【API证书】->【查看证书】
		mchCertificateSerialNumber string = "5430464867C627F3FC939F94B9B634162F******" // 商户证书序列号

		//登陆商户平台【账户中心】-【API安全】->【设置APIv3密钥】 仅在设置时能看到
		mchAPIv3Key string = "dhfICa404Ez0uCBwlxXEQkOA3M******" // 商户APIv3密钥
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
		log.Printf("new wechat pay client err:%s", err)
	}

	svc := TransferBatchApiService{Client: client}
	resp, result, err := svc.GetTransferBatchByOutNo(ctx,
		GetTransferBatchByOutNoRequest{
			OutBatchNo:      core.String("ARMB20220104152512200085666G6S"),
			NeedQueryDetail: core.Bool(true),
			Offset:          core.Int64(0),
			Limit:           core.Int64(20),
			DetailStatus:    core.String("FAIL"),
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call GetTransferDetailByNo err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	}
}

// 查批次
func Test_GetTransferBatchByNo(t *testing.T) {
	var (
		//登陆商户平台【账户中心】-【商户信息】->【微信支付商户号】
		mchID string = "16067*****" // 商户号

		//登陆商户平台【账户中心】-【API安全】->【API证书】->【查看证书】
		mchCertificateSerialNumber string = "5430464867C627F3FC939F94B9B634162F******" // 商户证书序列号

		//登陆商户平台【账户中心】-【API安全】->【设置APIv3密钥】 仅在设置时能看到
		mchAPIv3Key string = "dhfICa404Ez0uCBwlxXEQkOA3M******" // 商户APIv3密钥
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
		log.Printf("new wechat pay client err:%s", err)
	}

	svc := TransferBatchApiService{Client: client}
	resp, result, err := svc.GetTransferBatchByNo(ctx,
		GetTransferBatchByNoRequest{
			BatchId:         core.String("1030001017901275457032021123100595225416"),
			NeedQueryDetail: core.Bool(true),
			Offset:          core.Int64(0),
			Limit:           core.Int64(20),
			DetailStatus:    core.String("ALL"),
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call GetTransferDetailByNo err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	}
}
