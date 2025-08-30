package types

type LostThing struct {
	LostThingId         int
	PublicationDatetime string
	ThingName           string
	UserEmail           string
	CustomText          string
	Verified            int
	Status              int
}

type FoundThing struct {
	FoundThingId        int
	PublicationDatetime string
	ThingName           string
	ThingLocation       string
	CustomText          string
	Verified            int
	Status              int
}
