import { JSX, createSignal } from 'solid-js';

import Page from '../ui/Page/Page';
import Header from '../ui/Header/Header';
import Content from '../ui/Content/Content';
import Footer from '../ui/Footer/Footer';
import { Role } from '../utils/consts';
import getAuthorizedCookie from '../utils/getAuthorizedCookie';
import RegisterUserForm from '../components/RegisterUserForm';

const RegisterUserPage = (): JSX.Element => {
  const [authorized, setAuthorized] = createSignal(false);
  getAuthorizedCookie(setAuthorized);
  return (
    <Page
      role={Role.none}
      authorized={authorized()}
    >
      <Header
        role={Role.user}
        authorized={authorized()}
      />
      <Content>
        <RegisterUserForm />
      </Content>
      <Footer />
    </Page>
  );
};

export default RegisterUserPage;
