package transferbatch_test

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/transferbatch"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func ExampleTransferBatchApiService_GetTransferBatchByNo() {
	var (
		mchID                      string = "190000****"                               // 商户号
		mchCertificateSerialNumber string = "3775************************************" // 商户证书序列号
		mchAPIv3Key                string = "2ab9****************************"         // 商户APIv3密钥
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("/path/to/merchant/apiclient_key.pem")
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

	svc := transferbatch.TransferBatchApiService{Client: client}
	resp, result, err := svc.GetTransferBatchByNo(ctx,
		transferbatch.GetTransferBatchByNoRequest{
			BatchId:         core.String("1030000071100999991182020050700019480001"),
			NeedQueryDetail: core.Bool(true),
			Offset:          core.Int64(0),
			Limit:           core.Int64(20),
			DetailStatus:    core.String("FAIL"),
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call GetTransferBatchByNo err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	}
}

func ExampleTransferBatchApiService_GetTransferBatchByOutNo() {
	var (
		mchID                      string = "190000****"                               // 商户号
		mchCertificateSerialNumber string = "3775************************************" // 商户证书序列号
		mchAPIv3Key                string = "2ab9****************************"         // 商户APIv3密钥
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("/path/to/merchant/apiclient_key.pem")
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

	svc := transferbatch.TransferBatchApiService{Client: client}
	resp, result, err := svc.GetTransferBatchByOutNo(ctx,
		transferbatch.GetTransferBatchByOutNoRequest{
			OutBatchNo:      core.String("plfk2020042013"),
			NeedQueryDetail: core.Bool(true),
			Offset:          core.Int64(0),
			Limit:           core.Int64(20),
			DetailStatus:    core.String("FAIL"),
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call GetTransferBatchByOutNo err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	}
}

func ExampleTransferBatchApiService_InitiateBatchTransfer() {
	var (
		mchID                      string = "190000****"                               // 商户号
		mchCertificateSerialNumber string = "3775************************************" // 商户证书序列号
		mchAPIv3Key                string = "2ab9****************************"         // 商户APIv3密钥
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("/path/to/merchant/apiclient_key.pem")
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

	svc := transferbatch.TransferBatchApiService{Client: client}
	resp, result, err := svc.InitiateBatchTransfer(ctx,
		transferbatch.InitiateBatchTransferRequest{
			Appid:       core.String("wxf636efh567hg4356"),
			OutBatchNo:  core.String("plfk2020042013"),
			BatchName:   core.String("2019年1月深圳分部报销单"),
			BatchRemark: core.String("2019年1月深圳分部报销单"),
			TotalAmount: core.Int64(4000000),
			TotalNum:    core.Int64(200),
			TransferDetailList: []transferbatch.TransferDetailInput{transferbatch.TransferDetailInput{
				Openid:         core.String("o-MYE42l80oelYMDE34nYD456Xoy"),
				OutDetailNo:    core.String("x23zy545Bd5436"),
				TransferAmount: core.Int64(200000),
				TransferRemark: core.String("2020年4月报销"),
				UserIdCard:     core.String("UserIdCard_example"),
				UserName:       core.String("UserName_example"),
			}},
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call InitiateBatchTransfer err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	}
}

func InitiateBatchTransfer222(privPath string) {
	var (
		//登陆商户平台【账户中心】-【商户信息】->【微信支付商户号】
		mchID string = "1606746***" // 商户号

		//登陆商户平台【账户中心】-【API安全】->【API证书】->【查看证书】
		mchCertificateSerialNumber string = "5430464867C627F3FC939F94B9B63416&&&&&" // 商户证书序列号

		//登陆商户平台【账户中心】-【API安全】->【设置APIv3密钥】 仅在设置时能看到
		mchAPIv3Key string = "dhfICa404Ez0uCBwlxXEQkOA&&&&&" // 商户APIv3密钥
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(filepath.Join(privPath, "apiclient_key.pem"))
	if err != nil {
		log.Fatal("load merchant private key error")
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

	svc := transferbatch.TransferBatchApiService{Client: client}
	resp, result, err := svc.InitiateBatchTransfer(ctx,
		transferbatch.InitiateBatchTransferRequest{
			Appid:       core.String("wxd26fd1b5060cfd6b"), // 登陆商户平台【产品中心】->【AppID账号管理】
			OutBatchNo:  core.String("plfk202004201322211216"),
			BatchName:   core.String("测试批量转账到零钱"),
			BatchRemark: core.String("测试批量转账到零钱"),
			TotalAmount: core.Int64(30),
			TotalNum:    core.Int64(1),
			TransferDetailList: []transferbatch.TransferDetailInput{transferbatch.TransferDetailInput{
				// 获取openid的appid必须和上面的一致，否则会受理成功但转账失败
				Openid:         core.String("oM-AC6cmOClOpezp-EguzaGg7Dog"),
				OutDetailNo:    core.String("plfk202004201322211216"),
				TransferAmount: core.Int64(30),
				TransferRemark: core.String("TransferRemark"),
				//UserIdCard:     core.String(""),
				UserName: core.String("XX"), // 必须正确
			}},
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call InitiateBatchTransfer err:%s", err)
	} else {
		// 处理返回结果
		r, _ := result.Request.GetBody()
		b, _ := ioutil.ReadAll(r)
		b2, _ := ioutil.ReadAll(result.Response.Body)
		log.Printf("status=%d resp=%s；---2%+v\n---1%+v", result.Response.StatusCode, resp, string(b), string(b2))
	}
}
func Test_InitiateBatchTransfer(t *testing.T) {
	//b,_ := filepath.Abs(".")
	//print(b)
	p := "D:\\develop\\go\\s3\\server\\src\\tmp\\public_config\\conf\\beta\\wxpay"

	InitiateBatchTransfer222(p)
}
