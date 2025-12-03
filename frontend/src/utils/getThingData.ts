import axiosInstance from '../utils/axiosInstance';
import { ThingType } from '../utils/consts';

const getThingData = async (props: { thingId: string; thingType: ThingType }) =>
  axiosInstance
    .get(
      `/thing/get_data?thing_type=${props.thingType}&thing_id=${props.thingId}`,
    )
    .then((response) => response.data)
    .catch((error) => console.log(error));

export default getThingData;
