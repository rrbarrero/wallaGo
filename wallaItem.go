package main

const WALLA_ITEM_URL = ""

type WallaItems struct {
	Items []WallaItem
}

type WallaItem struct {
	ItemId             int64                   `json:"itemId"`
	ItemUUID           string                  `json:"itemUUID"`
	Title              string                  `json:"title"`
	Description        string                  `json:"description"`
	CategoryId         int                     `json:"categoryId"`
	Category           string                  `json:"category"`
	MainImage          WallaImage              `json:"mainImage"`
	Images             []WallaImage            `json:"images"`
	SalePrice          float64                 `json:"salePrice"`
	SellerUser         WallaUser               `json:"sellerUser"`
	URL                string                  `json:"itemURL"`
	PictureURL         string                  `json:"pictureURL"`
	Sold               bool                    `json:"sold"`
	Reserved           bool                    `json:"reserved"`
	PublishDate        int64                   `json:"publishDate"`
	ItemSaleConditions WallaItemSaleConditions `json:"itemSaleConditions"`
}

type WallaImage struct {
	ImageId           int    `json:"imageId"`
	PictureId         int64  `json:"pictureId"`
	Type              string `json:"type"`
	OriginalWitdh     int    `json:"originalWitdh"`
	OriginalHeight    int    `json:"originalHeight"`
	WebScaleHeight    int    `json:"webScaleHeight"`
	MobileScaledRatio int    `json:"mobileScaleRatio"`
	SmallURL          string `json:"smallURL"`
	MediumURL         string `json:"mediumURL"`
	BigURL            string `json:"bigURL"`
	XlargeURL         string `json:"xlargeURL"`
	AverageHexColor   string `json:"averageHexColor"`
}

type WallaUser struct {
	UserId     int64           `json:"userId"`
	MicroName  string          `json:"microName"`
	Image      WallaImage      `json:"image"`
	Banned     bool            `json:"banned"`
	Gender     string          `json:"gender"`
	ScreenName string          `json:"screenName"`
	URL        string          `json:"url"`
	Email      string          `json:"email"`
	StatUser   WallaStatsUser  `json:"statUser"`
	Validation WallaValidation `json:"validation"`
}

type WallaStatsUser struct {
	NotificationReadPendingCount int32 `json:"notificationReadPendingCount"`
	ConversationReadPendingCount int32 `json:"conversationReadPendingCount"`
	FavoritesCount               int32 `json:"favoritesCount"`
	SellingCount                 int32 `json:"sellingCount"`
	SelledCount                  int32 `json:"selledCount"`
	PurchasedCount               int32 `json:"purchasedCount"`
	SendReviewsCount             int32 `json:"sendReviewsCount"`
	ReceivedReviewsCount         int32 `json:"receivedReviewsCount"`
	ProdsSellingCount            int32 `json:"prodsSellingCount"`
	ProdsSelledCount             int32 `json:"prodsSelledCount"`
	UserReceivedReviewsCount     int32 `json:"userReceivedReviewsCount"`
	ProdsFavoritesCount          int32 `json:"prodsFavoritesCount"`
}

type WallaValidation struct {
	ScoringStarts int8 `json:"scoring_starts"`
}

type WallaItemSaleConditions struct {
	FixPrice        bool `json:"fix_price"`
	ExchangeAllowed bool `json:"exchange_allowed"`
	ShippingAllowed bool `json:"shipping_allowed"`
}
