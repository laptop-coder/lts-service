import axiosInstanceUnauthorized from '../utils/axiosInstanceUnauthorized';

const fetchThingData = async (props: {
  thingType: 'lost' | 'found';
  thingId: string;
}) => {
  return axiosInstanceUnauthorized
    .get(
      `/thing/get_data?thing_type=${props.thingType}&thing_id=${props.thingId}`,
    )
    .then((response) => {
      return response.data;
    })
    .catch((error) => console.log(error));
};

export default fetchThingData;
