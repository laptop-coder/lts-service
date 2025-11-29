import axiosInstance from '../utils/axiosInstance';

const getUsername = async () =>
  axiosInstance
    .get(`/user/get_username`)
    .then((response) => response.data)
    .catch((error) => console.log(error));

export default getUsername;
