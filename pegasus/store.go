package pegasus

import (
	"github.com/HearthSim/hs-proto-go/pegasus/util"
	"github.com/golang/protobuf/proto"
	//"log"
	"math/rand"
	"time"
)

type Store struct{}

func (s *Store) Init(sess *Session) {
	sess.RegisterPacket(util.GetBattlePayConfig_ID, OnGetBattlePayConfig)
	sess.RegisterPacket(util.GetBattlePayStatus_ID, OnGetBattlePayStatus)
	sess.RegisterPacket(util.PurchaseWithGold_ID, OnPurchaseWithGold)
}

func OnGetBattlePayConfig(s *Session, body []byte) *Packet {
	res := util.BattlePayConfigResponse{}
	// Hardcode US Dollars until we setup the DB to handle other currencies
	res.Currency = proto.Int32(1)
	res.Unavailable = proto.Bool(false)
	res.SecsBeforeAutoCancel = proto.Int32(10)

	product := ProductGoldCost{}
	db.Where("product_type = ?", 2).Find(&product)
	res.GoldCostArena = proto.Int64(product.Cost)

	goldCostBoosters := []*util.GoldCostBooster{}
	cost := []ProductGoldCost{}
	db.Where("pack_type != ?", 0).Find(&cost)
	for _, costs := range cost {
		goldCostBoosters = append(goldCostBoosters, &util.GoldCostBooster{
			Cost:     proto.Int64(costs.Cost),
			PackType: proto.Int32(costs.PackType),
		})
	}
	res.GoldCostBoosters = goldCostBoosters

	bundles := []Bundle{}
	db.Find(&bundles)
	for _, bundle := range bundles {
		bundleItems := []*util.BundleItem{}
		products := []Product{}
		db.Model(&bundle).Association("Items").Find(&products)
		for _, items := range products {
			productType := util.ProductType(items.ProductType)
			bundleItems = append(bundleItems, &util.BundleItem{
				ProductType: &productType,
				Data:        proto.Int32(items.ProductData),
				Quantity:    proto.Int32(items.Quantity),
			})
		}
		res.Bundles = append(res.Bundles, &util.Bundle{
			Id: proto.String(bundle.ProductID),
			// Hardcode $1 until price data is implemented in DB
			Cost:         proto.Float64(1.00),
			AppleId:      proto.String(bundle.AppleID),
			AmazonId:     proto.String(bundle.AmazonID),
			GooglePlayId: proto.String(bundle.GoogleID),
			// Hardcode 100 until price data is implemented in DB
			GoldCost:         proto.Int64(100),
			ProductEventName: proto.String(bundle.EventName),
			Items:            bundleItems,
		})
	}
	return EncodePacket(util.BattlePayConfigResponse_ID, &res)
}

func OnGetBattlePayStatus(s *Session, body []byte) *Packet {
	res := util.BattlePayStatusResponse{}
	status := util.BattlePayStatusResponse_PS_READY
	res.Status = &status
	res.BattlePayAvailable = proto.Bool(true)
	return EncodePacket(util.BattlePayStatusResponse_ID, &res)
}

func BuyPacks(accountId int64, product ProductGoldCost, quantity int32) bool {
	if product.ProductType == 1 {
		//buy booster packs
		allCards := []DbfCard{}
		db.Where("is_collectible = ?", true).Find(&allCards)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		
		idPrefix := []string {}
		switch product.PackType {
			case 1:
				idPrefix = []string {"CS2", "DS1", "EX1", "NEW1", "tt"}
			case 9:
				idPrefix = []string {"GVG"}
			case 10:
				idPrefix = []string {"AT"}
		}
		runCards := allCards[:0]
		
		for _, prefix := range idPrefix {
			prefixLen := len(prefix)
			for _, x := range allCards {
				if (x.SellPrice > 0) && (x.NoteMiniGuid[0:prefixLen] == prefix) {
					runCards = append(runCards, x)
				}
			}
		}
		
		
		for j := int32(0); j < quantity; j++ {
			cards := []BoosterCard{}
			//A pack contains 5 cards
			whiteCount := 0
			for i := 0; i < 5; i++ {
				var cardId int32 = 0
				var premium int32 = 0
				if r.Intn(20) == 1 {
					premium = 1
				}
				
				for {
					chooseCard := runCards[r.Intn(len(runCards))]
					if whiteCount == 4 && chooseCard.Rarity == 1 {
						continue
					}
					if chooseCard.Rarity == 1 || r.Intn((int(chooseCard.Rarity)-1) * 25) == 1 {
						cardId = chooseCard.ID
						if chooseCard.Rarity == 1 {
							whiteCount++
						}
						break
					}
				}
				cards = append(cards, BoosterCard{
					CardID: cardId,
					Premium: premium,
				})
			}
			tmp := cards[4]
			tmpId := r.Intn(5)
			cards[4] = cards[tmpId]
			cards[tmpId] = tmp
			booster := Booster{
				AccountID: accountId, 
				BoosterType: int(product.PackType), 
				Opened: false,
				Cards: cards,
			}
			db.Create(&booster)
		}
		//TODO: Check if card gen
	}
	
	return true
}

func OnPurchaseWithGold(s *Session, body []byte) *Packet {
	req := util.PurchaseWithGold{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	product := ProductGoldCost{}
	productType := req.GetProduct()
	data := req.GetData()
	quantity := req.GetQuantity()
	productCount := 0
	// If data is > 0, we're buying a pack
	if data > 0 {
		db.Where("product_type = ? AND pack_type = ?", productType, data).Find(&product).Count(&productCount)
		if productCount == 1 {
			BuyPacks(s.Account.ID, product, quantity)
		}
	} else {
		db.Where("product_type = ?", productType).Find(&product)
	}

	res := util.PurchaseWithGoldResponse{}
	// TODO: Query the DB to ensure we have enough gold
	result := util.PurchaseWithGoldResponse_PR_SUCCESS
	res.Result = &result
	res.GoldUsed = proto.Int64(product.Cost * int64(req.GetQuantity()))
	return EncodePacket(util.PurchaseWithGoldResponse_ID, &res)
}
