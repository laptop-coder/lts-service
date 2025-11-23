import axiosInstanceUnauthorized from '../utils/axiosInstanceUnauthorized';
import {
  ThingsListsSelectionCriteria,
  ThingsListsSelectionCriteriaStatusValues,
  ThingsListsSelectionCriteriaVerifiedValues,
} from '../enums/thingsListsSelectionCriteria';

const fetchThingsList = async (props: {
  thingsType: 'lost' | 'found';
  selectBy: ThingsListsSelectionCriteria;
}) => {
  var thingsStatus = '';
  var thingsVerified = '';
  if (
    props.selectBy === ThingsListsSelectionCriteria.NotReceivedVerifiedThings
  ) {
    thingsStatus = ThingsListsSelectionCriteriaStatusValues.NotReceived;
    thingsVerified = ThingsListsSelectionCriteriaVerifiedValues.Accepted;
  } else if (
    props.selectBy === ThingsListsSelectionCriteria.ReceivedVerifiedThings
  ) {
    thingsStatus = ThingsListsSelectionCriteriaStatusValues.Received;
    thingsVerified = ThingsListsSelectionCriteriaVerifiedValues.Accepted;
  } else if (
    props.selectBy === ThingsListsSelectionCriteria.NotVerifiedThings
  ) {
    thingsVerified = ThingsListsSelectionCriteriaVerifiedValues.NotVerified;
  } else if (props.selectBy === ThingsListsSelectionCriteria.AcceptedThings) {
    thingsVerified = ThingsListsSelectionCriteriaVerifiedValues.Accepted;
  } else if (props.selectBy === ThingsListsSelectionCriteria.RejectedThings) {
    thingsVerified = ThingsListsSelectionCriteriaVerifiedValues.Rejected;
  }

  return axiosInstanceUnauthorized
    .get(
      `/things/get_list?things_type=${props.thingsType}&things_status=${thingsStatus}&things_verified=${thingsVerified}`,
    )
    .then((response) => {
      return response.data;
    })
    .catch((error) => console.log(error));
};

export default fetchThingsList;
