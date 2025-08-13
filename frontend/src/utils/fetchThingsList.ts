import axiosInstanceUnauthorized from '../utils/axiosInstanceUnauthorized';

const fetchThingsList = async (thingsType: 'lost' | 'found') => {
  return axiosInstanceUnauthorized
    .get(`/things/get_list?things_type=${thingsType}`)
    .then((response) => {
      return response.data;
    });
};

export default fetchThingsList;
