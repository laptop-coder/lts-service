import { JSX } from 'solid-js';

import Header from '../components/Header/Header';
import Page from '../ui/Page/Page';
import Content from '../components/Content/Content';

const ModeratorAccountSettingsPage = (): JSX.Element => {
  return (
    <Page>
      <Header moderator />
      <Content></Content>
    </Page>
  );
};

export default ModeratorAccountSettingsPage;
