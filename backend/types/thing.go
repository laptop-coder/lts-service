package types

type LostThing struct {
	Id                  string
	PublicationDatetime string
	Name                string
	UserEmail           string
	UserMessage         string
	Verified            int
	Found               int
	AdvertisementOwner  string
}

type FoundThing struct {
	Id                  string
	PublicationDatetime string
	Name                string
	Location            string
	UserMessage         string
	Verified            int
	Found               int
	AdvertisementOwner  string
}
