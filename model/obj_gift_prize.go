package models

// 单独定义的奖品模型，原来的gift model有些字段不能给用户看
// 拆分后的中奖区间 PrizeCodeA， PrizeCodeB
type ObjGiftPrize struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	PrizeNum     int    `json:"-"`
	LeftNum      int    `json:"-"`
	PrizeCodeA   int    `json:"-"` //拆分后的中奖区间
	PrizeCodeB   int    `json:"-"` //拆分后的中奖区间
	Img          string `json:"img"`
	Displayorder int    `json:"displayorder"`
	Gtype        int    `json:"gtype"`
	Gdata        string `json:"gdata"`
}
