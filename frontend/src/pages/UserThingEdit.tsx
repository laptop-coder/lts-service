import { JSX, createSignal } from 'solid-js';

import Page from '../ui/Page/Page';
import Header from '../ui/Header/Header';
import Content from '../ui/Content/Content';
import Footer from '../ui/Footer/Footer';
import { Role, HeaderButton } from '../utils/consts';
import getAuthorizedCookie from '../utils/getAuthorizedCookie';
import UserThingEditForm from '../components/UserThingEditForm';
import { useSearchParams } from '@solidjs/router';

const UserThingEditPage = (): JSX.Element => {
  const [searchParams] = useSearchParams();
  const thingId = (searchParams.thing_id || '').toString();

  const [authorized, setAuthorized] = createSignal(false);
  getAuthorizedCookie(setAuthorized);
  return (
    <Page
      role={Role.user}
      authorized={authorized()}
    >
      <Header
        role={Role.user}
        buttons={[authorized() ? HeaderButton.profile : HeaderButton.login]}
      />
      <Content>
        <UserThingEditForm thing={{ id: thingId }} />
      </Content>
      <Footer />
    </Page>
  );
};

export default UserThingEditPage;
