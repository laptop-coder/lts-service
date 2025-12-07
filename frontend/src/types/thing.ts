import { ThingType } from '../utils/consts';

export interface Thing {
  Id: string;
  Type: ThingType;
  PublicationDatetime: string;
  Name: string;
  Photo: string;
  UserMessage: string;
  Verified: string;
  Found: string;
  NoticeOwner: string;
}
