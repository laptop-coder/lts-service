import axiosInstance from '../utils/axiosInstance';

const getUserEmail = async () =>
  axiosInstance
    .get(`/user/get_email`)
    .then((response) => response.data)
    .catch((error) => console.log(error));

export default getUserEmail;
