import { JSX, createSignal } from 'solid-js';

import Page from '../ui/Page/Page';
import Header from '../ui/Header/Header';
import Content from '../ui/Content/Content';
import Footer from '../ui/Footer/Footer';
import { Role, HeaderButton } from '../utils/consts';
import getAuthorizedCookie from '../utils/getAuthorizedCookie';
import getUsernameModerator from '../utils/getUsernameModerator';
import ModeratorProfileContent from '../components/ModeratorProfileContent/ModeratorProfileContent';

const ModeratorProfilePage = (): JSX.Element => {
  const [authorized, setAuthorized] = createSignal(false);
  getAuthorizedCookie(setAuthorized);

  const [username, setUsername] = createSignal('');

  getUsernameModerator().then((data) => setUsername(data));

  return (
    <Page
      role={Role.moderator}
      authorized={authorized()}
    >
      <Header
        role={Role.moderator}
        buttons={[authorized() ? HeaderButton.logout : HeaderButton.login]}
      />
      <Content>
        <ModeratorProfileContent username={username()} />
      </Content>
      <Footer />
    </Page>
  );
};

export default ModeratorProfilePage;
