import { JSX } from 'solid-js';

import { ASSETS_ROUTE } from '../utils/consts';
import Header from '../components/Header/Header';
import Page from '../ui/Page/Page';
import Content from '../components/Content/Content';
import SquareImageButton from '../ui/SquareImageButton/SquareImageButton';

const ModeratorAccountSettingsPage = (): JSX.Element => {
  const logout = () => {
    //TODO
    console.log('Logged out');
  };
  return (
    <Page>
      <Header moderator>
        <SquareImageButton onclick={logout}>
          <img src={`${ASSETS_ROUTE}/logout.svg`} />
        </SquareImageButton>
      </Header>
      <Content></Content>
    </Page>
  );
};

export default ModeratorAccountSettingsPage;
