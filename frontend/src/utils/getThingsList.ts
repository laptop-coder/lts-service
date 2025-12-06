import axiosInstance from '../utils/axiosInstance';
import { ThingType, NoticesOwnership } from './consts';

const getThingsList = async (props: {
  thingsType: ThingType;
  noticesOwnership: NoticesOwnership;
}) =>
  axiosInstance
    .get(
      `/things/get_list?things_type=${props.thingsType}&notices_ownership=${props.noticesOwnership}`,
    )
    .then((response) => response.data)
    .catch((error) =>
      // console.log(error)
      console.log('error'),
    );

export default getThingsList;
