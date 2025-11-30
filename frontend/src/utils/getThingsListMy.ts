import axiosInstance from '../utils/axiosInstance';
import { ThingType } from './consts';

const getThingsListMy = async (props: { thingsType: ThingType }) =>
  axiosInstance
    .get(`/things/get_list/my?things_type=${props.thingsType}`)
    .then((response) => response.data)
    .catch((error) => console.log(error));

export default getThingsListMy;
