import { JSX, createSignal } from 'solid-js';

import Page from '../ui/Page/Page';
import Header from '../ui/Header/Header';
import Content from '../ui/Content/Content';
import Footer from '../ui/Footer/Footer';
import {
  Role,
  ThingType,
  NoticesOwnership,
  HeaderButton,
} from '../utils/consts';
import getAuthorizedCookie from '../utils/getAuthorizedCookie';
import ThingsTypeToggle from '../components/ThingsTypeToggle';
import NoticesOwnershipToggle from '../components/NoticesOwnershipToggle';
import ThingsList from '../components/ThingsList/ThingsList';

const HomePage = (): JSX.Element => {
  const [authorized, setAuthorized] = createSignal(false);
  getAuthorizedCookie(setAuthorized);

  const [noticesOwnership, setNoticesOwnership] = createSignal(
    authorized() ? NoticesOwnership.not_my : NoticesOwnership.all,
  );
  const [thingsType, setThingsType] = createSignal(ThingType.lost);

  return (
    <Page
      role={authorized() ? Role.user : Role.none}
      authorized={authorized()}
    >
      <Header
        role={authorized() ? Role.user : Role.none}
        addThingDefaultThingType={thingsType()}
        buttons={[
          authorized() ? HeaderButton.add_thing : HeaderButton.none,
          authorized() ? HeaderButton.profile : HeaderButton.login,
        ]}
      />
      <Content>
        {authorized() && (
          <NoticesOwnershipToggle setter={setNoticesOwnership} />
        )}
        <ThingsTypeToggle setter={setThingsType} />
        <ThingsList
          thingsType={thingsType()}
          role={authorized() ? Role.user : Role.none}
          noticesOwnership={noticesOwnership()}
        />
      </Content>
      <Footer />
    </Page>
  );
};

export default HomePage;
