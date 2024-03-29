// Copyright 2021 Tencent Inc. All rights reserved.
//
// 微信支付服务商基础支付
//
// 微信支付 API v3 服务商基础支付
//
// API version: 1.2.3

// Code generated by WechatPay APIv3 Generator based on [OpenAPI Generator](https://openapi-generator.tech); DO NOT EDIT.

package partnerpayments

import (
	"encoding/json"
	"fmt"
)

// PromotionDetail
type PromotionDetail struct {
	// 券ID
	CouponId *string `json:"coupon_id,omitempty"`
	// 优惠名称
	Name *string `json:"name,omitempty"`
	// GLOBAL：全场代金券；SINGLE：单品优惠
	Scope *string `json:"scope,omitempty"`
	// CASH：充值；NOCASH：预充值。
	Type *string `json:"type,omitempty"`
	// 优惠券面额
	Amount *int64 `json:"amount,omitempty"`
	// 活动ID，批次ID
	StockId *string `json:"stock_id,omitempty"`
	// 单位为分
	WechatpayContribute *int64 `json:"wechatpay_contribute,omitempty"`
	// 单位为分
	MerchantContribute *int64 `json:"merchant_contribute,omitempty"`
	// 单位为分
	OtherContribute *int64 `json:"other_contribute,omitempty"`
	// CNY：人民币，境内商户号仅支持人民币。
	Currency    *string                `json:"currency,omitempty"`
	GoodsDetail []PromotionGoodsDetail `json:"goods_detail,omitempty"`
}

func (o PromotionDetail) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.CouponId != nil {
		toSerialize["coupon_id"] = o.CouponId
	}

	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	if o.Scope != nil {
		toSerialize["scope"] = o.Scope
	}

	if o.Type != nil {
		toSerialize["type"] = o.Type
	}

	if o.Amount != nil {
		toSerialize["amount"] = o.Amount
	}

	if o.StockId != nil {
		toSerialize["stock_id"] = o.StockId
	}

	if o.WechatpayContribute != nil {
		toSerialize["wechatpay_contribute"] = o.WechatpayContribute
	}

	if o.MerchantContribute != nil {
		toSerialize["merchant_contribute"] = o.MerchantContribute
	}

	if o.OtherContribute != nil {
		toSerialize["other_contribute"] = o.OtherContribute
	}

	if o.Currency != nil {
		toSerialize["currency"] = o.Currency
	}

	if o.GoodsDetail != nil {
		toSerialize["goods_detail"] = o.GoodsDetail
	}
	return json.Marshal(toSerialize)
}

func (o PromotionDetail) String() string {
	var ret string
	if o.CouponId == nil {
		ret += "CouponId:<nil>, "
	} else {
		ret += fmt.Sprintf("CouponId:%v, ", *o.CouponId)
	}

	if o.Name == nil {
		ret += "Name:<nil>, "
	} else {
		ret += fmt.Sprintf("Name:%v, ", *o.Name)
	}

	if o.Scope == nil {
		ret += "Scope:<nil>, "
	} else {
		ret += fmt.Sprintf("Scope:%v, ", *o.Scope)
	}

	if o.Type == nil {
		ret += "Type:<nil>, "
	} else {
		ret += fmt.Sprintf("Type:%v, ", *o.Type)
	}

	if o.Amount == nil {
		ret += "Amount:<nil>, "
	} else {
		ret += fmt.Sprintf("Amount:%v, ", *o.Amount)
	}

	if o.StockId == nil {
		ret += "StockId:<nil>, "
	} else {
		ret += fmt.Sprintf("StockId:%v, ", *o.StockId)
	}

	if o.WechatpayContribute == nil {
		ret += "WechatpayContribute:<nil>, "
	} else {
		ret += fmt.Sprintf("WechatpayContribute:%v, ", *o.WechatpayContribute)
	}

	if o.MerchantContribute == nil {
		ret += "MerchantContribute:<nil>, "
	} else {
		ret += fmt.Sprintf("MerchantContribute:%v, ", *o.MerchantContribute)
	}

	if o.OtherContribute == nil {
		ret += "OtherContribute:<nil>, "
	} else {
		ret += fmt.Sprintf("OtherContribute:%v, ", *o.OtherContribute)
	}

	if o.Currency == nil {
		ret += "Currency:<nil>, "
	} else {
		ret += fmt.Sprintf("Currency:%v, ", *o.Currency)
	}

	ret += fmt.Sprintf("GoodsDetail:%v", o.GoodsDetail)

	return fmt.Sprintf("PromotionDetail{%s}", ret)
}

func (o PromotionDetail) Clone() *PromotionDetail {
	ret := PromotionDetail{}

	if o.CouponId != nil {
		ret.CouponId = new(string)
		*ret.CouponId = *o.CouponId
	}

	if o.Name != nil {
		ret.Name = new(string)
		*ret.Name = *o.Name
	}

	if o.Scope != nil {
		ret.Scope = new(string)
		*ret.Scope = *o.Scope
	}

	if o.Type != nil {
		ret.Type = new(string)
		*ret.Type = *o.Type
	}

	if o.Amount != nil {
		ret.Amount = new(int64)
		*ret.Amount = *o.Amount
	}

	if o.StockId != nil {
		ret.StockId = new(string)
		*ret.StockId = *o.StockId
	}

	if o.WechatpayContribute != nil {
		ret.WechatpayContribute = new(int64)
		*ret.WechatpayContribute = *o.WechatpayContribute
	}

	if o.MerchantContribute != nil {
		ret.MerchantContribute = new(int64)
		*ret.MerchantContribute = *o.MerchantContribute
	}

	if o.OtherContribute != nil {
		ret.OtherContribute = new(int64)
		*ret.OtherContribute = *o.OtherContribute
	}

	if o.Currency != nil {
		ret.Currency = new(string)
		*ret.Currency = *o.Currency
	}

	if o.GoodsDetail != nil {
		ret.GoodsDetail = make([]PromotionGoodsDetail, len(o.GoodsDetail))
		for i, item := range o.GoodsDetail {
			ret.GoodsDetail[i] = *item.Clone()
		}
	}

	return &ret
}

// PromotionGoodsDetail
type PromotionGoodsDetail struct {
	// 商品编码
	GoodsId *string `json:"goods_id"`
	// 商品数量
	Quantity *int64 `json:"quantity"`
	// 商品价格
	UnitPrice *int64 `json:"unit_price"`
	// 商品优惠金额
	DiscountAmount *int64 `json:"discount_amount"`
	// 商品备注
	GoodsRemark *string `json:"goods_remark,omitempty"`
}

func (o PromotionGoodsDetail) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.GoodsId == nil {
		return nil, fmt.Errorf("field `GoodsId` is required and must be specified in PromotionGoodsDetail")
	}
	toSerialize["goods_id"] = o.GoodsId

	if o.Quantity == nil {
		return nil, fmt.Errorf("field `Quantity` is required and must be specified in PromotionGoodsDetail")
	}
	toSerialize["quantity"] = o.Quantity

	if o.UnitPrice == nil {
		return nil, fmt.Errorf("field `UnitPrice` is required and must be specified in PromotionGoodsDetail")
	}
	toSerialize["unit_price"] = o.UnitPrice

	if o.DiscountAmount == nil {
		return nil, fmt.Errorf("field `DiscountAmount` is required and must be specified in PromotionGoodsDetail")
	}
	toSerialize["discount_amount"] = o.DiscountAmount

	if o.GoodsRemark != nil {
		toSerialize["goods_remark"] = o.GoodsRemark
	}
	return json.Marshal(toSerialize)
}

func (o PromotionGoodsDetail) String() string {
	var ret string
	if o.GoodsId == nil {
		ret += "GoodsId:<nil>, "
	} else {
		ret += fmt.Sprintf("GoodsId:%v, ", *o.GoodsId)
	}

	if o.Quantity == nil {
		ret += "Quantity:<nil>, "
	} else {
		ret += fmt.Sprintf("Quantity:%v, ", *o.Quantity)
	}

	if o.UnitPrice == nil {
		ret += "UnitPrice:<nil>, "
	} else {
		ret += fmt.Sprintf("UnitPrice:%v, ", *o.UnitPrice)
	}

	if o.DiscountAmount == nil {
		ret += "DiscountAmount:<nil>, "
	} else {
		ret += fmt.Sprintf("DiscountAmount:%v, ", *o.DiscountAmount)
	}

	if o.GoodsRemark == nil {
		ret += "GoodsRemark:<nil>"
	} else {
		ret += fmt.Sprintf("GoodsRemark:%v", *o.GoodsRemark)
	}

	return fmt.Sprintf("PromotionGoodsDetail{%s}", ret)
}

func (o PromotionGoodsDetail) Clone() *PromotionGoodsDetail {
	ret := PromotionGoodsDetail{}

	if o.GoodsId != nil {
		ret.GoodsId = new(string)
		*ret.GoodsId = *o.GoodsId
	}

	if o.Quantity != nil {
		ret.Quantity = new(int64)
		*ret.Quantity = *o.Quantity
	}

	if o.UnitPrice != nil {
		ret.UnitPrice = new(int64)
		*ret.UnitPrice = *o.UnitPrice
	}

	if o.DiscountAmount != nil {
		ret.DiscountAmount = new(int64)
		*ret.DiscountAmount = *o.DiscountAmount
	}

	if o.GoodsRemark != nil {
		ret.GoodsRemark = new(string)
		*ret.GoodsRemark = *o.GoodsRemark
	}

	return &ret
}

// Transaction
type Transaction struct {
	Amount          *TransactionAmount `json:"amount,omitempty"`
	SpAppid         *string            `json:"sp_appid,omitempty"`
	SubAppid        *string            `json:"sub_appid,omitempty"`
	SpMchid         *string            `json:"sp_mchid,omitempty"`
	SubMchid        *string            `json:"sub_mchid,omitempty"`
	Attach          *string            `json:"attach,omitempty"`
	BankType        *string            `json:"bank_type,omitempty"`
	OutTradeNo      *string            `json:"out_trade_no,omitempty"`
	Payer           *TransactionPayer  `json:"payer,omitempty"`
	PromotionDetail []PromotionDetail  `json:"promotion_detail,omitempty"`
	SuccessTime     *string            `json:"success_time,omitempty"`
	TradeState      *string            `json:"trade_state,omitempty"`
	TradeStateDesc  *string            `json:"trade_state_desc,omitempty"`
	TradeType       *string            `json:"trade_type,omitempty"`
	TransactionId   *string            `json:"transaction_id,omitempty"`
}

func (o Transaction) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Amount != nil {
		toSerialize["amount"] = o.Amount
	}

	if o.SpAppid != nil {
		toSerialize["sp_appid"] = o.SpAppid
	}

	if o.SubAppid != nil {
		toSerialize["sub_appid"] = o.SubAppid
	}

	if o.SpMchid != nil {
		toSerialize["sp_mchid"] = o.SpMchid
	}

	if o.SubMchid != nil {
		toSerialize["sub_mchid"] = o.SubMchid
	}

	if o.Attach != nil {
		toSerialize["attach"] = o.Attach
	}

	if o.BankType != nil {
		toSerialize["bank_type"] = o.BankType
	}

	if o.OutTradeNo != nil {
		toSerialize["out_trade_no"] = o.OutTradeNo
	}

	if o.Payer != nil {
		toSerialize["payer"] = o.Payer
	}

	if o.PromotionDetail != nil {
		toSerialize["promotion_detail"] = o.PromotionDetail
	}

	if o.SuccessTime != nil {
		toSerialize["success_time"] = o.SuccessTime
	}

	if o.TradeState != nil {
		toSerialize["trade_state"] = o.TradeState
	}

	if o.TradeStateDesc != nil {
		toSerialize["trade_state_desc"] = o.TradeStateDesc
	}

	if o.TradeType != nil {
		toSerialize["trade_type"] = o.TradeType
	}

	if o.TransactionId != nil {
		toSerialize["transaction_id"] = o.TransactionId
	}
	return json.Marshal(toSerialize)
}

func (o Transaction) String() string {
	var ret string
	ret += fmt.Sprintf("Amount:%v, ", o.Amount)

	if o.SpAppid == nil {
		ret += "SpAppid:<nil>, "
	} else {
		ret += fmt.Sprintf("SpAppid:%v, ", *o.SpAppid)
	}

	if o.SubAppid == nil {
		ret += "SubAppid:<nil>, "
	} else {
		ret += fmt.Sprintf("SubAppid:%v, ", *o.SubAppid)
	}

	if o.SpMchid == nil {
		ret += "SpMchid:<nil>, "
	} else {
		ret += fmt.Sprintf("SpMchid:%v, ", *o.SpMchid)
	}

	if o.SubMchid == nil {
		ret += "SubMchid:<nil>, "
	} else {
		ret += fmt.Sprintf("SubMchid:%v, ", *o.SubMchid)
	}

	if o.Attach == nil {
		ret += "Attach:<nil>, "
	} else {
		ret += fmt.Sprintf("Attach:%v, ", *o.Attach)
	}

	if o.BankType == nil {
		ret += "BankType:<nil>, "
	} else {
		ret += fmt.Sprintf("BankType:%v, ", *o.BankType)
	}

	if o.OutTradeNo == nil {
		ret += "OutBatchNo:<nil>, "
	} else {
		ret += fmt.Sprintf("OutBatchNo:%v, ", *o.OutTradeNo)
	}

	ret += fmt.Sprintf("Payer:%v, ", o.Payer)

	ret += fmt.Sprintf("PromotionDetail:%v, ", o.PromotionDetail)

	if o.SuccessTime == nil {
		ret += "SuccessTime:<nil>, "
	} else {
		ret += fmt.Sprintf("SuccessTime:%v, ", *o.SuccessTime)
	}

	if o.TradeState == nil {
		ret += "TradeState:<nil>, "
	} else {
		ret += fmt.Sprintf("TradeState:%v, ", *o.TradeState)
	}

	if o.TradeStateDesc == nil {
		ret += "TradeStateDesc:<nil>, "
	} else {
		ret += fmt.Sprintf("TradeStateDesc:%v, ", *o.TradeStateDesc)
	}

	if o.TradeType == nil {
		ret += "TradeType:<nil>, "
	} else {
		ret += fmt.Sprintf("TradeType:%v, ", *o.TradeType)
	}

	if o.TransactionId == nil {
		ret += "TransactionId:<nil>"
	} else {
		ret += fmt.Sprintf("TransactionId:%v", *o.TransactionId)
	}

	return fmt.Sprintf("Transaction{%s}", ret)
}

func (o Transaction) Clone() *Transaction {
	ret := Transaction{}

	if o.Amount != nil {
		ret.Amount = o.Amount.Clone()
	}

	if o.SpAppid != nil {
		ret.SpAppid = new(string)
		*ret.SpAppid = *o.SpAppid
	}

	if o.SubAppid != nil {
		ret.SubAppid = new(string)
		*ret.SubAppid = *o.SubAppid
	}

	if o.SpMchid != nil {
		ret.SpMchid = new(string)
		*ret.SpMchid = *o.SpMchid
	}

	if o.SubMchid != nil {
		ret.SubMchid = new(string)
		*ret.SubMchid = *o.SubMchid
	}

	if o.Attach != nil {
		ret.Attach = new(string)
		*ret.Attach = *o.Attach
	}

	if o.BankType != nil {
		ret.BankType = new(string)
		*ret.BankType = *o.BankType
	}

	if o.OutTradeNo != nil {
		ret.OutTradeNo = new(string)
		*ret.OutTradeNo = *o.OutTradeNo
	}

	if o.Payer != nil {
		ret.Payer = o.Payer.Clone()
	}

	if o.PromotionDetail != nil {
		ret.PromotionDetail = make([]PromotionDetail, len(o.PromotionDetail))
		for i, item := range o.PromotionDetail {
			ret.PromotionDetail[i] = *item.Clone()
		}
	}

	if o.SuccessTime != nil {
		ret.SuccessTime = new(string)
		*ret.SuccessTime = *o.SuccessTime
	}

	if o.TradeState != nil {
		ret.TradeState = new(string)
		*ret.TradeState = *o.TradeState
	}

	if o.TradeStateDesc != nil {
		ret.TradeStateDesc = new(string)
		*ret.TradeStateDesc = *o.TradeStateDesc
	}

	if o.TradeType != nil {
		ret.TradeType = new(string)
		*ret.TradeType = *o.TradeType
	}

	if o.TransactionId != nil {
		ret.TransactionId = new(string)
		*ret.TransactionId = *o.TransactionId
	}

	return &ret
}

// TransactionAmount
type TransactionAmount struct {
	Currency      *string `json:"currency,omitempty"`
	PayerCurrency *string `json:"payer_currency,omitempty"`
	PayerTotal    *int64  `json:"payer_total,omitempty"`
	Total         *int64  `json:"total,omitempty"`
}

func (o TransactionAmount) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Currency != nil {
		toSerialize["currency"] = o.Currency
	}

	if o.PayerCurrency != nil {
		toSerialize["payer_currency"] = o.PayerCurrency
	}

	if o.PayerTotal != nil {
		toSerialize["payer_total"] = o.PayerTotal
	}

	if o.Total != nil {
		toSerialize["total"] = o.Total
	}
	return json.Marshal(toSerialize)
}

func (o TransactionAmount) String() string {
	var ret string
	if o.Currency == nil {
		ret += "Currency:<nil>, "
	} else {
		ret += fmt.Sprintf("Currency:%v, ", *o.Currency)
	}

	if o.PayerCurrency == nil {
		ret += "PayerCurrency:<nil>, "
	} else {
		ret += fmt.Sprintf("PayerCurrency:%v, ", *o.PayerCurrency)
	}

	if o.PayerTotal == nil {
		ret += "PayerTotal:<nil>, "
	} else {
		ret += fmt.Sprintf("PayerTotal:%v, ", *o.PayerTotal)
	}

	if o.Total == nil {
		ret += "Total:<nil>"
	} else {
		ret += fmt.Sprintf("Total:%v", *o.Total)
	}

	return fmt.Sprintf("TransactionAmount{%s}", ret)
}

func (o TransactionAmount) Clone() *TransactionAmount {
	ret := TransactionAmount{}

	if o.Currency != nil {
		ret.Currency = new(string)
		*ret.Currency = *o.Currency
	}

	if o.PayerCurrency != nil {
		ret.PayerCurrency = new(string)
		*ret.PayerCurrency = *o.PayerCurrency
	}

	if o.PayerTotal != nil {
		ret.PayerTotal = new(int64)
		*ret.PayerTotal = *o.PayerTotal
	}

	if o.Total != nil {
		ret.Total = new(int64)
		*ret.Total = *o.Total
	}

	return &ret
}

// TransactionPayer
type TransactionPayer struct {
	SpOpenid  *string `json:"sp_openid,omitempty"`
	SubOpenid *string `json:"sub_openid,omitempty"`
}

func (o TransactionPayer) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.SpOpenid != nil {
		toSerialize["sp_openid"] = o.SpOpenid
	}

	if o.SubOpenid != nil {
		toSerialize["sub_openid"] = o.SubOpenid
	}
	return json.Marshal(toSerialize)
}

func (o TransactionPayer) String() string {
	var ret string
	if o.SpOpenid == nil {
		ret += "SpOpenid:<nil>, "
	} else {
		ret += fmt.Sprintf("SpOpenid:%v, ", *o.SpOpenid)
	}

	if o.SubOpenid == nil {
		ret += "SubOpenid:<nil>"
	} else {
		ret += fmt.Sprintf("SubOpenid:%v", *o.SubOpenid)
	}

	return fmt.Sprintf("TransactionPayer{%s}", ret)
}

func (o TransactionPayer) Clone() *TransactionPayer {
	ret := TransactionPayer{}

	if o.SpOpenid != nil {
		ret.SpOpenid = new(string)
		*ret.SpOpenid = *o.SpOpenid
	}

	if o.SubOpenid != nil {
		ret.SubOpenid = new(string)
		*ret.SubOpenid = *o.SubOpenid
	}

	return &ret
}
