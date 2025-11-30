import { JSX, createSignal } from 'solid-js';

import Page from '../ui/Page/Page';
import Header from '../ui/Header/Header';
import Content from '../ui/Content/Content';
import Footer from '../ui/Footer/Footer';
import { Role, UserProfileSection, ThingType } from '../utils/consts';
import getAuthorizedCookie from '../utils/getAuthorizedCookie';
import ThingsTypeToggle from '../components/ThingsTypeToggle'; // TODO: refactor naming (plural)
import UserProfileSectionsToggle from '../components/UserProfileSectionsToggle';

const UserProfilePage = (): JSX.Element => {
  const [authorized, setAuthorized] = createSignal(false);
  getAuthorizedCookie(setAuthorized);

  const [section, setSection] = createSignal(UserProfileSection.advertisements);
  const [thingsType, setThingsType] = createSignal(ThingType.lost);

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
        <UserProfileSectionsToggle setter={setSection} />
        {section() === UserProfileSection.advertisements && (
          <ThingsTypeToggle setter={setThingsType} />
        )}
      </Content>
      <Footer />
    </Page>
  );
};

export default UserProfilePage;
