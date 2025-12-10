import axiosInstance from '../utils/axiosInstance';

const getThingData = async (props: { thing: { id: string } }) =>
  axiosInstance
    .get(`/thing/get_data?thing_id=${props.thing.id}`)
    .then((response) => response.data)
    .catch((error) =>
      // console.log(error)
      console.log('error'),
    );

export default getThingData;
