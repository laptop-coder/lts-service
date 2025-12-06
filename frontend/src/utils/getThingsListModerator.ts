import axiosInstance from '../utils/axiosInstance';
import { ThingType, NoticesVerification } from './consts';

const getThingsListModerator = async (props: {
  thingsType: ThingType;
  noticesVerification: NoticesVerification;
}) =>
  axiosInstance
    .get(
      `/things/get_list/moderator?things_type=${props.thingsType}&notices_verification=${props.noticesVerification}`,
    )
    .then((response) => response.data)
    .catch((error) =>
      // console.log(error)
      console.log('error'),
    );

export default getThingsListModerator;
