import { JSX, createSignal } from 'solid-js';

import Page from '../ui/Page/Page';
import Header from '../ui/Header/Header';
import Content from '../ui/Content/Content';
import Footer from '../ui/Footer/Footer';
import { Role, ThingType } from '../utils/consts';
import getAuthorizedCookie from '../utils/getAuthorizedCookie';
import ThingsTypeToggle from '../components/ThingsTypeToggle';
import ThingsList from '../components/ThingsList/ThingsList';

const HomePage = (): JSX.Element => {
  const [authorized, setAuthorized] = createSignal(false);
  getAuthorizedCookie(setAuthorized);

  const [thingsType, setThingsType] = createSignal(ThingType.lost);

  return (
    <Page
      role={authorized() ? Role.user : Role.none}
      authorized={authorized()}
    >
      <Header
        role={authorized() ? Role.user : Role.none}
        authorized={authorized()}
        addThingDefaultThingType={thingsType()}
      />
      <Content>
        <ThingsTypeToggle setter={setThingsType} />
        <ThingsList
          thingsType={thingsType}
          role={Role.user}
        />
      </Content>
      <Footer />
    </Page>
  );
};

export default HomePage;
