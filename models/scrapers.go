package models

// Ozbargain deal type
type OzBargainDeal struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Url      string `json:"url"`
	PostedOn string `json:"time"`
	Upvotes  string `json:"upvotes"`
	DealAge  string `json:"dealage"`
	DealType int    `json:"dealtype"`
}

// Camel Camel Camel deal type
type CamCamCamDeal struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	Published string `json:"time"`
	Image     string `json:"image"`
	DealType  string `json:"dealtype"`
}

// Setters and getters for OzBargainDeal
func (d *OzBargainDeal) SetId(id string) {
	d.Id = id
}
func (d *OzBargainDeal) SetTitle(title string) {
	d.Title = title
}
func (d *OzBargainDeal) SetUrl(url string) {
	d.Url = url
}
func (d *OzBargainDeal) SetPostedOn(postedOn string) {
	d.PostedOn = postedOn
}
func (d *OzBargainDeal) SetUpvotes(upvotes string) {
	d.Upvotes = upvotes
}
func (d *OzBargainDeal) SetDealAge(dealAge string) {
	d.DealAge = dealAge
}
func (d *OzBargainDeal) SetDealType(dealType int) {
	d.DealType = dealType
}
func (d *OzBargainDeal) GetId() string {
	return d.Id
}
func (d *OzBargainDeal) GetTitle() string {
	return d.Title
}
func (d *OzBargainDeal) GetUrl() string {
	return d.Url
}
func (d *OzBargainDeal) GetPostedOn() string {
	return d.PostedOn
}
func (d *OzBargainDeal) GetUpvotes() string {
	return d.Upvotes
}
func (d *OzBargainDeal) GetDealAge() string {
	return d.DealAge
}
func (d *OzBargainDeal) GetDealType() int {
	return d.DealType
}
