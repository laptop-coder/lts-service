import axios, { AxiosInstance } from 'axios';

import { BACKEND_URL } from './consts';

const axiosInstance: AxiosInstance = axios.create({
  baseURL: BACKEND_URL,
  withCredentials: true,
});

export default axiosInstance;
