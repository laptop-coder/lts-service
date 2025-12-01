import { JSX, createSignal } from 'solid-js';

import Page from '../ui/Page/Page';
import Header from '../ui/Header/Header';
import Content from '../ui/Content/Content';
import Footer from '../ui/Footer/Footer';
import { Role } from '../utils/consts';
import getAuthorizedCookie from '../utils/getAuthorizedCookie';
import getUserEmail from '../utils/getUserEmail';
import getUsername from '../utils/getUsername';
import UserProfileContent from '../components/UserProfileContent/UserProfileContent';

const UserProfilePage = (): JSX.Element => {
  const [authorized, setAuthorized] = createSignal(false);
  getAuthorizedCookie(setAuthorized);

  const [username, setUsername] = createSignal('');
  const [email, setEmail] = createSignal('');

  getUsername().then((data) => setUsername(data));
  getUserEmail().then((data) => setEmail(data));

  return (
    <Page
      role={Role.user}
      authorized={authorized()}
    >
      <Header
        role={Role.user}
        authorized={authorized()}
        showLogout
      />
      <Content>
        <UserProfileContent
          username={username()}
          email={email()}
        />
      </Content>
      <Footer />
    </Page>
  );
};

export default UserProfilePage;
