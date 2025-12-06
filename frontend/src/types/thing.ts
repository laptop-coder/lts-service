import { ThingType } from '../utils/consts';

export interface Thing {
  Id: string;
  Type: ThingType;
  PublicationDatetime: string;
  Name: string;
  Photo: string;
  UserMessage: string;
  Verified: number;
  Found: number;
  NoticeOwner: string;
}
