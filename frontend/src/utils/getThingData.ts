import axiosInstance from '../utils/axiosInstance';

const getThingData = async (props: { thingId: string }) =>
  axiosInstance
    .get(`/thing/get_data?thing_id=${props.thingId}`)
    .then((response) => response.data)
    .catch((error) =>
      // console.log(error)
      console.log('error'),
    );

export default getThingData;
