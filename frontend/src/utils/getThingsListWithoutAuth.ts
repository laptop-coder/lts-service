import axiosInstance from '../utils/axiosInstance';
import { ThingType } from './consts';

const getThingsListWithoutAuth = async (props: { thingsType: ThingType }) =>
  axiosInstance
    .get(`/things/get_list/without_auth?things_type=${props.thingsType}`)
    .then((response) => response.data)
    .catch((error) =>
      // console.log(error)
      console.log('error'),
    );

export default getThingsListWithoutAuth;
