import { JSX, createSignal } from 'solid-js';

import Page from '../ui/Page/Page';
import Header from '../ui/Header/Header';
import Content from '../ui/Content/Content';
import Footer from '../ui/Footer/Footer';
import { Role } from '../utils/consts';
import getAuthorizedCookie from '../utils/getAuthorizedCookie';

const RegisterModeratorPage = (): JSX.Element => {
  const [authorized, setAuthorized] = createSignal(false);
  getAuthorizedCookie(setAuthorized);
  return (
    <Page
      role={Role.none}
      authorized={authorized()}
    >
      <Header
        role={Role.moderator}
        authorized={authorized()}
      />
      <Content></Content>
      <Footer />
    </Page>
  );
};

export default RegisterModeratorPage;
