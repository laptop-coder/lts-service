import { JSX, createSignal } from 'solid-js';

import Page from '../ui/Page/Page';
import Header from '../ui/Header/Header';
import Content from '../ui/Content/Content';
import Footer from '../ui/Footer/Footer';
import { Role, HeaderButton, NoticesVerification } from '../utils/consts';
import getAuthorizedCookie from '../utils/getAuthorizedCookie';
import NoticesVerificationToggle from '../components/NoticesVerificationToggle';

const ModeratorHomePage = (): JSX.Element => {
  const [authorized, setAuthorized] = createSignal(false);
  getAuthorizedCookie(setAuthorized);

  const [noticesVerification, setNoticesVerification] = createSignal(
    NoticesVerification.not_verified,
  );
  return (
    <Page
      role={Role.moderator}
      authorized={authorized()}
    >
      <Header
        role={Role.moderator}
        buttons={[authorized() ? HeaderButton.profile : HeaderButton.login]}
      />
      <Content>
        <NoticesVerificationToggle setter={setNoticesVerification} />
      </Content>
      <Footer />
    </Page>
  );
};

export default ModeratorHomePage;
