import axiosInstance from '../utils/axiosInstance';

const getUsernameModerator = async () =>
  axiosInstance
    .get(`/moderator/get_username`)
    .then((response) => response.data)
    .catch((error) => console.log(error));

export default getUsernameModerator;
