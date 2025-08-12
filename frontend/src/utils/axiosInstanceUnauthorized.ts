import axios, { AxiosInstance } from 'axios';

import { BACKEND_URL } from './consts';

const axiosInstanceUnauthorized: AxiosInstance = axios.create({
  baseURL: BACKEND_URL,
});

export default axiosInstanceUnauthorized;
