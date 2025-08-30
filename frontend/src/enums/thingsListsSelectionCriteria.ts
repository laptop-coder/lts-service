export enum ThingsListsSelectionCriteria {
  NotReceivedVerifiedThings,
  ReceivedVerifiedThings,
  NotVerifiedThings,
  AcceptedThings,
  RejectedThings,
}

export enum ThingsListsSelectionCriteriaStatusValues {
  NotReceived = '0',
  Received = '1',
}

export enum ThingsListsSelectionCriteriaVerifiedValues {
  NotVerified = '0',
  Accepted = '1',
  Rejected = '-1',
}
