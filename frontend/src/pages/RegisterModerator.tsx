import { JSX, createSignal } from 'solid-js';

import Page from '../ui/Page/Page';
import Header from '../ui/Header/Header';
import Content from '../ui/Content/Content';
import Footer from '../ui/Footer/Footer';
import { Role, HeaderButton } from '../utils/consts';
import getAuthorizedCookie from '../utils/getAuthorizedCookie';
import RegisterModeratorForm from '../components/RegisterModeratorForm';

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
        buttons={[authorized() ? HeaderButton.profile : HeaderButton.login]}
      />
      <Content>
        <RegisterModeratorForm />
      </Content>
      <Footer />
    </Page>
  );
};

export default RegisterModeratorPage;
