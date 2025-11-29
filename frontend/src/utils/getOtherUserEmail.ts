import axiosInstance from '../utils/axiosInstance';

const getOtherUserEmail = async (props: { username: string }) =>
  axiosInstance
    .post(
      `/user/get_email/other`,
      { username: props.username },
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      },
    )
    .then((response) => response.data)
    .catch((error) => console.log(error));

export default getOtherUserEmail;
