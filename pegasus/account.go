package pegasus

import (
	"github.com/HearthSim/hs-proto-go/pegasus/shared"
	"github.com/HearthSim/hs-proto-go/pegasus/util"
	"github.com/golang/protobuf/proto"
	"log"
	"time"
)

func (v *Account) Init(sess *Session) {
	// TODO: fetch the account using the bnet session
	
	db.Where("ID = ?", sess.host.GetbnetID()).Find(v)
	v.displayName = "FIXME"

	sess.RegisterPacket(util.GetAccountInfo_ID, OnGetAccountInfo)
	sess.RegisterPacket(util.UpdateLogin_ID, OnUpdateLogin)
	sess.RegisterPacket(util.CreateDeck_ID, OnCreateDeck)
	sess.RegisterPacket(util.DeleteDeck_ID, OnDeleteDeck)
	sess.RegisterPacket(util.RenameDeck_ID, OnRenameDeck)
	sess.RegisterPacket(util.GetDeck_ID, OnGetDeck)
	sess.RegisterPacket(util.DeckSetData_ID, OnDeckSetData)
	sess.RegisterPacket(util.AckCardSeen_ID, OnAckCardSeen)
	sess.RegisterPacket(util.OpenBooster_ID, OnOpenBooster)
	sess.RegisterPacket(util.SetOptions_ID, OnSetOptions)
	sess.RegisterPacket(util.GetOptions_ID, OnGetOptions)
	sess.RegisterPacket(util.AckAchieveProgress_ID, OnAckAchieveProgress)
	sess.RegisterPacket(util.GetAchieves_ID, OnGetAchieves)
	sess.RegisterPacket(util.CheckAccountLicenses_ID, OnCheckAccountLicenses)
	sess.RegisterPacket(util.CheckGameLicenses_ID, OnCheckGameLicenses)
	sess.RegisterPacket(util.CancelQuest_ID, OnCancelQuest)
	sess.RegisterPacket(util.ValidateAchieve_ID, OnValidateAchieve)
	sess.RegisterPacket(util.SetCardBack_ID, OnSetCardBack)
	sess.RegisterPacket(util.GetAdventureProgress_ID, OnGetAdventureProgress)
	sess.RegisterPacket(util.SetFavoriteHero_ID, OnSetFavoriteHero)
	sess.RegisterPacket(util.GenericRequestList_ID, OnGenericRequest)
}

func OnAckCardSeen(s *Session, body []byte) *Packet {
	req := util.AckCardSeen{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	log.Printf("FIXME: AckCardSeen = %s", req.String())
	return nil
}

func OnCheckAccountLicenses(s *Session, body []byte) *Packet {
	res := util.CheckAccountLicensesResponse{}
	res.Success = proto.Bool(true)
	return EncodePacket(util.CheckAccountLicensesResponse_ID, &res)
}

func OnCheckGameLicenses(s *Session, body []byte) *Packet {
	res := util.CheckGameLicensesResponse{}
	res.Success = proto.Bool(true)
	return EncodePacket(util.CheckGameLicensesResponse_ID, &res)
}

func OnUpdateLogin(s *Session, body []byte) *Packet {
	req := util.UpdateLogin{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	log.Printf("req = %s", req.String())
	res := util.UpdateLoginComplete{}
	return EncodePacket(util.UpdateLoginComplete_ID, &res)
}

func OnGetAccountInfo(s *Session, body []byte) *Packet {
	req := util.GetAccountInfo{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	log.Printf("req = %s", req.String())
	return s.HandleAccountInfoRequest(*req.Request)
}

const AccoutInfoRequestBoosters = int32(util.GetAccountInfo_BOOSTERS)

func (s *Session) HandleAccountInfoRequest(req util.GetAccountInfo_Request) *Packet {
	switch req {
	case util.GetAccountInfo_GAMES_PLAYED:
		res := util.GamesInfo{}
		res.GamesStarted = proto.Int32(0)
		res.GamesWon = proto.Int32(0)
		res.GamesLost = proto.Int32(0)
		res.FreeRewardProgress = proto.Int32(0)
		return EncodePacket(util.GamesInfo_ID, &res)
	case util.GetAccountInfo_CAMPAIGN_INFO:
		res := util.ProfileProgress{}
		res.Progress = proto.Int64(6)  // ILLIDAN_COMPLETE
		res.BestForge = proto.Int32(0) // Arena wins
		return EncodePacket(util.ProfileProgress_ID, &res)
	case util.GetAccountInfo_BOOSTERS:
		res := util.BoosterList{}
		classicPacks := s.GetBoosterInfo(1)
		gvgPacks := s.GetBoosterInfo(9)
		tgtPacks := s.GetBoosterInfo(10)
		log.Printf("Account: %d, Remain packs: %d %d %d\n", s.Account.ID, *classicPacks.Count, *gvgPacks.Count, *tgtPacks.Count)
		if *classicPacks.Count > 0 {
			res.List = append(res.List, classicPacks)
		}
		if *gvgPacks.Count > 0 {
			res.List = append(res.List, gvgPacks)
		}
		if *tgtPacks.Count > 0 {
			res.List = append(res.List, tgtPacks)
		}
		return EncodePacket(util.BoosterList_ID, &res)
	case util.GetAccountInfo_FEATURES:
		res := util.GuardianVars{}
		res.ShowUserUI = proto.Int32(1)
		return EncodePacket(util.GuardianVars_ID, &res)
	case util.GetAccountInfo_MEDAL_INFO:
		res := util.MedalInfo{}
		res.SeasonWins = proto.Int32(0)
		res.Stars = proto.Int32(2)
		res.Streak = proto.Int32(0)
		res.StarLevel = proto.Int32(1)
		res.LevelStart = proto.Int32(1)
		res.LevelEnd = proto.Int32(3)
		res.CanLoseLevel = proto.Bool(false)
		return EncodePacket(util.MedalInfo_ID, &res)
	case util.GetAccountInfo_MEDAL_HISTORY:
		res := util.MedalHistory{}
		for i := int32(1); i <= 3; i++ {
			info := &util.MedalHistoryInfo{}
			info.When = PegasusDate(time.Date(2015, 8, 1, 7, 0, 0, 0, time.UTC))
			info.Season = proto.Int32(i)
			info.Stars = proto.Int32(0)
			info.StarLevel = proto.Int32(0)
			info.LevelStart = proto.Int32(0)
			info.LevelEnd = proto.Int32(0)
			info.LegendRank = proto.Int32(1)
			res.Medals = append(res.Medals, info)
		}
		return EncodePacket(util.MedalHistory_ID, &res)
	case util.GetAccountInfo_NOTICES:
		res := util.ProfileNotices{}
		return EncodePacket(util.ProfileNotices_ID, &res)
	case util.GetAccountInfo_DECK_LIST:
		res := util.DeckList{}
		basicDecks := []Deck{}
		deckType := shared.DeckType_PRECON_DECK
		db.Where("deck_type = ?", deckType).Find(&basicDecks)
		for _, deck := range basicDecks {
			info := MakeDeckInfo(&deck)
			res.Decks = append(res.Decks, info)
		}
		decks := []Deck{}
		deckType = shared.DeckType_NORMAL_DECK
		db.Where("deck_type = ? and account_id = ?", deckType, s.Account.ID).Find(&decks)
		for _, deck := range decks {
			info := MakeDeckInfo(&deck)
			res.Decks = append(res.Decks, info)
		}
		return EncodePacket(util.DeckList_ID, &res)
	case util.GetAccountInfo_COLLECTION:
		res := util.Collection{}
		collectionCards := []CollectionCard{}
		db.Where("account_id = ?", s.Account.ID).Find(&collectionCards)
		for _, card := range collectionCards {
			stack1 := &shared.CardStack{}
			stack1.LatestInsertDate = PegasusDate(time.Now().UTC())
			stack1.NumSeen = proto.Int32(2)
			stack1.Count = proto.Int32(card.Num)
			carddef := &shared.CardDef{}
			carddef.Asset = proto.Int32(card.CardID)
			carddef.Premium = proto.Int32(card.Premium)
			stack1.CardDef = carddef
			res.Stacks = append(res.Stacks, stack1)
		}
		return EncodePacket(util.Collection_ID, &res)
	case util.GetAccountInfo_DECK_LIMIT:
		res := util.ProfileDeckLimit{}
		res.DeckLimit = proto.Int32(9)
		return EncodePacket(util.ProfileDeckLimit_ID, &res)
	case util.GetAccountInfo_CARD_VALUES:
		res := util.CardValues{}
		dbfCards := []DbfCard{}
		db.Where("is_collectible = ? AND buy_price is not ?", true, 0).Find(&dbfCards)
		for _, cards := range dbfCards {
			card := &util.CardValue{}
			cardPremium := &util.CardValue{}
			card.Card = MakeCardDef(cards.ID, 0)
			card.Buy = proto.Int32(cards.BuyPrice)
			card.Sell = proto.Int32(cards.SellPrice)
			card.Nerfed = proto.Bool(false)
			cardPremium.Card = MakeCardDef(cards.ID, 1)
			cardPremium.Buy = proto.Int32(cards.GoldBuyPrice)
			cardPremium.Sell = proto.Int32(cards.GoldSellPrice)
			cardPremium.Nerfed = proto.Bool(false)
			res.Cards = append(res.Cards, card)
			res.Cards = append(res.Cards, cardPremium)
		}
		res.CardNerfIndex = proto.Int32(0)
		return EncodePacket(util.CardValues_ID, &res)
	case util.GetAccountInfo_ARCANE_DUST_BALANCE:
		res := util.ArcaneDustBalance{}
		account := Account{}
		db.Where("id = ?", s.Account.ID).First(&account)
		res.Balance = proto.Int64(account.Dust)
		return EncodePacket(util.ArcaneDustBalance_ID, &res)
	case util.GetAccountInfo_GOLD_BALANCE:
		res := util.GoldBalance{}
		account := Account{}
		db.Where("id = ?", s.Account.ID).First(&account)
		res.Cap = proto.Int64(999999)
		res.CapWarning = proto.Int64(2000)
		res.CappedBalance = proto.Int64(account.Gold)
		res.BonusBalance = proto.Int64(0)
		return EncodePacket(util.GoldBalance_ID, &res)
	case util.GetAccountInfo_HERO_XP:
		res := util.HeroXP{}
		for i := 2; i <= 10; i++ {
			info := &util.HeroXPInfo{}
			level := 2*i + 5
			maxXp := 60 + level*10
			info.ClassId = proto.Int32(int32(i))
			info.Level = proto.Int32(int32(level))
			info.CurrXp = proto.Int64(int64(maxXp / 2))
			info.MaxXp = proto.Int64(int64(maxXp))
			res.XpInfos = append(res.XpInfos, info)
		}
		return EncodePacket(util.HeroXP_ID, &res)
	case util.GetAccountInfo_NOT_SO_MASSIVE_LOGIN:
		res := util.NotSoMassiveLoginReply{}
		// Client get event time including brawls from this request
		// Hardcode TGT until it defined in database
		events := []*util.SpecialEventTiming{}
		tgtEvent := util.SpecialEventTiming{
			Event: proto.String("tgt_normal_sale"),
		}
		events = append(events, &tgtEvent)
		res.SpecialEventTiming = events
		log.Printf("events res = %s", res.String())
		return EncodePacket(util.NotSoMassiveLoginReply_ID, &res)
	case util.GetAccountInfo_REWARD_PROGRESS:
		res := util.RewardProgress{}
		nextMonth := time.Date(2015, 8, 1, 7, 0, 0, 0, time.UTC)
		res.SeasonEnd = PegasusDate(nextMonth)
		res.WinsPerGold = proto.Int32(3)
		res.GoldPerReward = proto.Int32(10)
		res.MaxGoldPerDay = proto.Int32(100)
		res.SeasonNumber = proto.Int32(21)
		res.XpSoloLimit = proto.Int32(60)
		res.MaxHeroLevel = proto.Int32(60)
		res.NextQuestCancel = PegasusDate(time.Now().UTC())
		res.EventTimingMod = proto.Float32(0.291667)
		return EncodePacket(util.RewardProgress_ID, &res)
	case util.GetAccountInfo_PVP_QUEUE:
		res := util.PlayQueue{}
		queue := shared.PlayQueueInfo{}
		gametype := shared.BnetGameType_BGT_NORMAL
		queue.GameType = &gametype
		res.Queue = &queue
		return EncodePacket(util.PlayQueue_ID, &res)

	case util.GetAccountInfo_PLAYER_RECORD:
		res := util.PlayerRecords{}
		return EncodePacket(util.PlayerRecords_ID, &res)
	case util.GetAccountInfo_CARD_BACKS:
		res := util.CardBacks{}
		dbfCardBacks := []DbfCardBack{}
		res.DefaultCardBack = proto.Int32(0)
		db.Find(&dbfCardBacks)
		for _, cardBack := range dbfCardBacks {
			res.CardBacks = append(res.CardBacks, cardBack.ID)
		}
		return EncodePacket(util.CardBacks_ID, &res)
	case util.GetAccountInfo_FAVORITE_HEROES:
		res := util.FavoriteHeroesResponse{}
		favoriteHeros := []FavoriteHero{}
		db.Where("account_id = ?", s.Account.ID).Find(&favoriteHeros)
		for _, hero := range favoriteHeros {
			card := DbfCard{}
			db.Where("id = ?", hero.CardID).First(&card)
			carddef := &shared.CardDef{}
			carddef.Asset = proto.Int32(hero.CardID)
			carddef.Premium = proto.Int32(hero.Premium)
			fav := &shared.FavoriteHero{}
			fav.ClassId = proto.Int32(hero.ClassID)
			fav.Hero = carddef
			res.FavoriteHeroes = append(res.FavoriteHeroes, fav)
		}
		return EncodePacket(util.FavoriteHeroesResponse_ID, &res)
	case util.GetAccountInfo_ACCOUNT_LICENSES:
		res := util.AccountLicensesInfoResponse{}
		return EncodePacket(util.AccountLicensesInfoResponse_ID, &res)
	case util.GetAccountInfo_BOOSTER_TALLY:
		res := util.BoosterTallyList{}
		return EncodePacket(util.BoosterTallyList_ID, &res)
	case util.GetAccountInfo_CLIENT_OPTIONS:
		return MakeFakeOptions()
	default:
		log.Printf("Unhandled GetAccountInfo request type: %s", req.String())
		panic(nyi)
	}
}

func OnGenericRequest(s *Session, body []byte) *Packet {
	req := util.GenericRequestList{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	log.Printf("req = %s", req.String())
	for _, r := range req.Requests {
		switch *r.RequestId {
		case int32(util.GetAccountInfo_ID):
			p := s.HandleAccountInfoRequest(util.GetAccountInfo_Request(*r.RequestSubId))
			s.SendUtilPacket(p)
		default:
			panic(nyi)
		}
	}
	return nil
}

func OnGetAdventureProgress(s *Session, body []byte) *Packet {
	res := util.AdventureProgressResponse{}
	return EncodePacket(util.AdventureProgressResponse_ID, &res)
}

func OnSetOptions(s *Session, body []byte) *Packet {
	req := util.SetOptions{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	log.Printf("req = %s", req.String())
	return nil
}

func OnGetOptions(s *Session, body []byte) *Packet {
	req := util.GetOptions{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	log.Printf("req = %s", req.String())
	return MakeFakeOptions()
}

func MakeFakeOptions() *Packet {
	res := util.ClientOptions{}
	res.Options = append(res.Options, &util.ClientOption{
		Index:    proto.Int32(1),
		AsUint64: proto.Uint64(0x20FFFF3FFFCCFCFF),
	})
	res.Options = append(res.Options, &util.ClientOption{
		Index:    proto.Int32(2),
		AsUint64: proto.Uint64(0xF0BFFFEF3FFF),
	})
	res.Options = append(res.Options, &util.ClientOption{
		Index:   proto.Int32(18),
		AsInt64: proto.Int64(0xB765A8C),
	})
	return EncodePacket(util.ClientOptions_ID, &res)
}

func OnGetAchieves(s *Session, body []byte) *Packet {
	req := util.GetAchieves{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	res := util.Achieves{}
	dbfAchieves := []DbfAchieve{}
	db.Find(&dbfAchieves)
	for _, achieve := range dbfAchieves {
		achieveProgress := Achieve{}
		db.Where("achieve_id = ? and account_id = ?", achieve.ID, s.Account.ID).First(&achieveProgress)

		res.List = append(res.List, &util.Achieve{
			Id:              proto.Int32(achieve.ID),
			Progress:        proto.Int32(achieveProgress.Progress),
			AckProgress:     proto.Int32(achieveProgress.AckProgress),
			CompletionCount: proto.Int32(achieveProgress.CompletionCount),
			Active:          proto.Bool(achieveProgress.Active),
			DateGiven:       PegasusDate(achieveProgress.DateGiven),
			DateCompleted:   PegasusDate(achieveProgress.DateCompleted),
		})
	}
	return EncodePacket(util.Achieves_ID, &res)
}

func OnSetFavoriteHero(s *Session, body []byte) *Packet {
	req := util.SetFavoriteHero{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	// TODO: optionally deny changes and respond accordingly
	favoriteHero := FavoriteHero{}
	db.Where("class_id = ? and account_id = ?", req.GetFavoriteHero().ClassId, s.Account.ID).First(&favoriteHero)
	requestedHeroCard := req.GetFavoriteHero().GetHero()
	favoriteHero.CardID = requestedHeroCard.GetAsset()
	favoriteHero.Premium = requestedHeroCard.GetPremium()
	db.Save(&favoriteHero)

	res := util.SetFavoriteHeroResponse{
		Success:      proto.Bool(true),
		FavoriteHero: req.FavoriteHero,
	}
	return EncodePacket(util.SetFavoriteHeroResponse_ID, &res)
}

func OnAckAchieveProgress(s *Session, body []byte) *Packet {
	req := util.AckAchieveProgress{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	achieve := Achieve{}
	db.Where("achieve_id = ? and account_id = ?", req.Id, s.Account.ID).First(&achieve)
	achieve.AckProgress = req.GetAckProgress()
	db.Save(&achieve)

	return nil
}

func OnValidateAchieve(s *Session, body []byte) *Packet {
	req := util.ValidateAchieve{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	log.Printf("req = %s", req.String())
	res := util.ValidateAchieveResponse{}
	res.Achieve = proto.Int32(req.GetAchieve())
	return EncodePacket(util.ValidateAchieveResponse_ID, &res)
}

func OnCancelQuest(s *Session, body []byte) *Packet {
	req := util.CancelQuest{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	// TODO: check if the quest is a daily that can be canceled and if so cancel it.
	// Note that if CancelQuestResponse.Success is true the client will request
	//  a new set of achieves and expect a new daily
	log.Printf("TODO: OnCancelQuest stub = %s", req.String())

	res := util.CancelQuestResponse{
		QuestId:         proto.Int32(req.GetQuestId()),
		Success:         proto.Bool(true),
		NextQuestCancel: PegasusDate(time.Now().UTC()),
	}

	return EncodePacket(util.CancelQuestResponse_ID, &res)
}

func MakeDeckInfo(deck *Deck) *shared.DeckInfo {
	info := &shared.DeckInfo{}
	info.Id = proto.Int64(deck.ID)
	info.Name = proto.String(deck.Name)
	info.CardBack = proto.Int32(0)
	info.Hero = proto.Int32(deck.HeroID)
	deckType := shared.DeckType(deck.DeckType)
	info.DeckType = &deckType
	info.Validity = proto.Uint64(31)
	info.HeroPremium = proto.Int32(deck.HeroPremium)
	info.CardBackOverride = proto.Bool(false)
	info.HeroOverride = proto.Bool(false)

	return info
}

func MakeCardDef(id, premium int32) *shared.CardDef {
	res := &shared.CardDef{}
	res.Asset = proto.Int32(id)
	res.Premium = proto.Int32(premium)
	return res
}

func OnOpenBooster(s *Session, body []byte) *Packet {
	req := util.OpenBooster{}
	err := proto.Unmarshal(body, &req)
		
	if err != nil {
		panic(err)
	}

	res := util.BoosterContent{}
	booster := Booster{}
	db.Where("booster_type = ? and opened = ? and account_id = ?", req.GetBoosterType(), false, s.Account.ID).Preload("Cards").First(&booster)
	log.Println(booster)
	for _, card := range booster.Cards {
		boosterCard := &util.BoosterCard{
			CardDef:    MakeCardDef(card.CardID, card.Premium),
			InsertDate: PegasusDate(time.Now().UTC()),
		}
		cards := CollectionCard{}
		if !db.Where("account_id = ? AND card_id = ? AND premium = ?", s.Account.ID, card.CardID, card.Premium).First(&cards).RecordNotFound() {
			db.Model(&cards).Update("num", cards.Num+1)
		} else {
			cards.AccountID = s.Account.ID
			cards.CardID = card.CardID
			cards.Premium = card.Premium
			cards.Num = 1
			db.Save(&cards)
		}
		res.List = append(res.List, boosterCard)

	}
	booster.Opened = true
	db.Save(&booster)
	return EncodePacket(util.BoosterContent_ID, &res)
}

func OnGetDeck(s *Session, body []byte) *Packet {
	req := util.GetDeck{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	id := req.GetDeck()
	var deck Deck
	db.First(&deck, id)

	// TODO: does this also need to allow brawl/arena decks? what about AI decks?
	if deck.DeckType != int(shared.DeckType_PRECON_DECK) && deck.AccountID != s.Account.ID {
		log.Panicf("received OnGetDeck for non-precon deck not owned by account")
	}

	deckCards := []DeckCard{}
	db.Where("deck_id = ?", id).Find(&deckCards)

	res := util.DeckContents{
		Deck: proto.Int64(id),
	}

	for i, card := range deckCards {
		cardData := &shared.DeckCardData{
			Def:    MakeCardDef(card.CardID, card.Premium),
			Handle: proto.Int32(int32(i)),
			Qty:    proto.Int32(card.Num),
			Prev:   proto.Int32(int32(i) - 1),
		}
		res.Cards = append(res.Cards, cardData)
	}

	return EncodePacket(util.DeckContents_ID, &res)
}

func OnCreateDeck(s *Session, body []byte) *Packet {
	req := util.CreateDeck{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	deck := Deck{
		AccountID:    s.Account.ID,
		DeckType:     int(req.GetDeckType()),
		Name:         req.GetName(),
		HeroID:       int32(req.GetHero()),
		HeroPremium:  int32(req.GetHeroPremium()),
		CardBackID:   0,
		LastModified: time.Now().UTC(),
	}
	db.Create(&deck)

	res := util.DeckCreated{}

	info := shared.DeckInfo{}
	info.Id = proto.Int64(deck.ID)
	info.Name = req.Name
	info.DeckType = req.DeckType
	info.CardBack = proto.Int32(1)
	info.CardBackOverride = proto.Bool(false)
	info.Hero = req.Hero
	info.HeroPremium = req.HeroPremium
	info.HeroOverride = proto.Bool(false)
	info.Validity = proto.Uint64(1)
	res.Info = &info
	return EncodePacket(util.DeckCreated_ID, &res)
}

func OnDeckSetData(s *Session, body []byte) *Packet {
	req := util.DeckSetData{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	id := req.GetDeck()
	deck := Deck{}
	db.First(&deck, id)

	if deck.AccountID != s.Account.ID {
		log.Panicf("received DeckSetData for deck not owned by account")
	}

	// Clear the deck then re-populate it
	db.Where("deck_id = ?", id).Delete(DeckCard{})

	for _, card := range req.Cards {
		cardDef := card.GetDef()
		qty := int(card.GetQty())
		if qty == 0 {
			qty = 1
		}
		c := DeckCard{
			DeckID:  deck.ID,
			CardID:  int32(cardDef.GetAsset()),
			Premium: int32(cardDef.GetPremium()),
			Num:     int32(qty),
		}
		db.Create(&c)
	}

	hero := req.GetHero()
	if hero != nil {
		deck.HeroID = int32(hero.GetAsset())
		deck.HeroPremium = int32(hero.GetPremium())
	}

	cardBack := req.GetCardBack()
	if cardBack != 0 {
		deck.CardBackID = cardBack
	}

	deck.LastModified = time.Now().UTC()
	db.Save(&deck)

	res := util.DBAction{}
	action := shared.DatabaseAction(int32(5)) // DB_A_SET_DECK
	result := shared.DatabaseResult(int32(1)) // DB_E_SUCCESS
	res.Action = &action
	res.Result = &result
	res.MetaData = proto.Int64(id)

	return EncodePacket(util.DBAction_ID, &res)
}

func OnSetCardBack(s *Session, body []byte) *Packet {
	req := util.SetCardBack{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	log.Printf("FIXME: SetCardBack stub = %s", req.String())
	res := util.SetCardBackResponse{}
	cardback := req.GetCardBack()
	res.CardBack = &cardback
	res.Success = proto.Bool(false)
	return EncodePacket(util.SetCardBackResponse_ID, &res)
}

func OnRenameDeck(s *Session, body []byte) *Packet {
	req := util.RenameDeck{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	id := req.GetDeck()
	deck := Deck{}
	name := req.GetName()
	db.First(&deck, id)

	if deck.AccountID != s.Account.ID {
		log.Panicf("received RenameDeck for deck not owned by account")
	}

	deck.Name = name
	db.Save(&deck)

	res := util.DeckRenamed{
		Deck: proto.Int64(id),
		Name: proto.String(name),
	}
	return EncodePacket(util.DeckRenamed_ID, &res)
}

func OnDeleteDeck(s *Session, body []byte) *Packet {
	req := util.DeleteDeck{}
	err := proto.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	id := req.GetDeck()
	deck := Deck{}
	db.First(&deck, id)

	if deck.AccountID != s.Account.ID {
		log.Panicf("received DeleteDeck for deck not owned by account")
	}

	db.Where("deck_id = ?", id).Delete(DeckCard{})
	db.Delete(&deck)

	res := util.DeckDeleted{
		Deck: proto.Int64(id),
	}
	return EncodePacket(util.DeckDeleted_ID, &res)
}

func (s *Session) GetBoosterInfo(kind int32) *shared.BoosterInfo {
	var count int32
	db.Model(Booster{}).
		Where("booster_type = ? and opened = ? and account_id = ?", kind, false, s.Account.ID).
		Count(&count)
	log.Printf("kind: %d count: %d\n", kind, count)    
	res := &shared.BoosterInfo{}
	res.Count = proto.Int32(count)
	res.Type = proto.Int32(kind)
	return res
}

func PegasusDate(t time.Time) *shared.Date {
	return &shared.Date{
		Year:  proto.Int32(int32(t.Year())),
		Month: proto.Int32(int32(t.Month())),
		Day:   proto.Int32(int32(t.Day())),
		Hours: proto.Int32(int32(t.Hour())),
		Min:   proto.Int32(int32(t.Minute())),
		Sec:   proto.Int32(int32(t.Second())),
	}
}
