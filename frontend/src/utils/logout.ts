import axiosInstance from './axiosInstance';
import {
  MODERATOR__HOME__ROUTE,
  HOME__ROUTE,
  Role,
  BACKEND__LOGOUT__ROUTE,
} from './consts';

const logout = async (props: { role: Role }) => {
  return axiosInstance
    .post(BACKEND__LOGOUT__ROUTE)
    .then(() => {
      window.location.href =
        props.role === Role.moderator ? MODERATOR__HOME__ROUTE : HOME__ROUTE;
    })
    .catch((error) =>
      // console.log(error)
      console.log('error'),
    );
};

export default logout;
