import axiosInstance from '../utils/axiosInstance';
import { ThingType, NoticesOwnership } from './consts';

const getThingsListUser = async (props: {
  thingsType: ThingType;
  noticesOwnership: NoticesOwnership;
}) =>
  axiosInstance
    .get(
      `/things/get_list/user?things_type=${props.thingsType}&notices_ownership=${props.noticesOwnership}`,
    )
    .then((response) => response.data)
    .catch((error) =>
      // console.log(error)
      console.log('error'),
    );

export default getThingsListUser;
