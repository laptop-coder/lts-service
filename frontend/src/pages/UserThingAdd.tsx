import { JSX, createSignal } from 'solid-js';

import Page from '../ui/Page/Page';
import Header from '../ui/Header/Header';
import Content from '../ui/Content/Content';
import Footer from '../ui/Footer/Footer';
import { Role, ThingType } from '../utils/consts';
import getAuthorizedCookie from '../utils/getAuthorizedCookie';
import UserThingAddForm from '../components/UserThingAddForm';

import { useSearchParams } from '@solidjs/router';

const UserThingAddPage = (): JSX.Element => {
  const [searchParams] = useSearchParams();
  const [authorized, setAuthorized] = createSignal(false);
  getAuthorizedCookie(setAuthorized);
  return (
    <Page
      role={Role.user}
      authorized={authorized()}
    >
      <Header
        role={Role.user}
        authorized={authorized()}
      />
      <Content>
        <UserThingAddForm
          defaultThingType={
            searchParams.default_thing_type === ThingType.found
              ? ThingType.found
              : ThingType.lost
          }
        />
      </Content>
      <Footer />
    </Page>
  );
};

export default UserThingAddPage;
