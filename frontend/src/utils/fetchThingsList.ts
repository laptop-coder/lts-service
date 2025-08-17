import axiosInstanceUnauthorized from '../utils/axiosInstanceUnauthorized';

const fetchThingsList = async (props: { thingsType: 'lost' | 'found' }) => {
  return axiosInstanceUnauthorized
    .get(`/things/get_list?things_type=${props.thingsType}`)
    .then((response) => {
      return response.data;
    })
    .catch((error) => console.log(error));
};

export default fetchThingsList;
