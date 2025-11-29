export interface LostThing {
  Id: string;
  PublicationDatetime: string;
  Name: string;
  Photo: string;
  UserMessage: string;
  Verified: number;
  Found: number;
  AdvertisementOwner: string;
}

export interface FoundThing {
  Id: string;
  PublicationDatetime: string;
  Name: string;
  Photo: string;
  Location: string;
  UserMessage: string;
  Verified: number;
  Found: number;
  AdvertisementOwner: string;
}
