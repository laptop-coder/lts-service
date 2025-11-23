import getCookie from '../utils/getCookie';
import type { Setter } from 'solid-js';

const updateAuthorizedCookie = (setAuthorized: Setter<boolean>) => {
  var authorizedCookie = getCookie('authorized');
  if (authorizedCookie != undefined) {
    setAuthorized(JSON.parse(authorizedCookie));
  }
};

export default updateAuthorizedCookie;
